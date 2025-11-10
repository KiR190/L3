package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"delayed-notifier/config"
	"delayed-notifier/internal/cache"
	"delayed-notifier/internal/handler"
	"delayed-notifier/internal/models"
	"delayed-notifier/internal/queue"
	"delayed-notifier/internal/repository"
	"delayed-notifier/internal/sender"
	"delayed-notifier/internal/server"
	"delayed-notifier/internal/service"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/wb-go/wbf/dbpg"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Загрузка конфигурации
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Подключение к БД
	dbOpts := &dbpg.Options{MaxOpenConns: 5, MaxIdleConns: 1}
	dbConn, err := dbpg.New(cfg.DatabaseURL, []string{}, dbOpts)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Инициализация репозитория
	notifRepo := repository.NewPostgresNotificationRepo(dbConn)

	// Telegram
	token := cfg.TG_BOT_TOKEN
	telegramSender, err := sender.NewTelegramSender(token)
	if err != nil {
		log.Fatal(err)
	}
	go telegramSender.ListenAndServe()

	emailSender := sender.NewEmailSender(
		cfg.SMTP_HOST,
		cfg.SMTP_PORT,
		cfg.SMTP_USERNAME,
		cfg.SMTP_PASSWORD,
		cfg.SMTP_FROM,
	)

	consoleSender := &sender.NativeSender{}

	multiSender := sender.NewMultiSender(consoleSender, emailSender, telegramSender)

	notificationQueue, err := queue.NewQueue(cfg.RabbitMQURL)
	if err != nil {
		log.Fatalf("Failed to init RabbitMQ: %v", err)
	}
	defer notificationQueue.Close()
	notificationQueue.StartMainConsumer()

	// Инициализация кеша
	redisCache := cache.NewCache(cfg.REDIS_ADDR, cfg.REDIS_PASSWORD, 0)

	// Сервис
	notifService := service.NewNotificationService(notifRepo, redisCache, notificationQueue, multiSender)

	// ctx := context.Background()
	notifService.RestoreCacheFromDB(ctx)

	notificationQueue.SetHandler(func(ctx context.Context, notif models.Notification) error {
		return notifService.ProcessNotification(ctx, &notif)
	})

	// Handler
	notifHandler := handler.NewNotificationHandler(notifService)
	router := server.NewRouter(notifHandler)
	httpServer := server.NewHTTPServer(cfg, router)

	log.Printf("Server started at http://localhost:%s", cfg.HTTPPort)
	if err := httpServer.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("Server started on port %s", cfg.HTTPPort)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server failed: %v", err)
		}
	}()

	// Ждём сигнал на остановку
	<-stop
	log.Println("Shutting down gracefully...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	// Закрываем HTTP сервер
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
	}

	// Останавливаем очередь
	if err := notificationQueue.Close(); err != nil {
		log.Printf("Queue shutdown error: %v", err)
	}

	// Останавливаем Telegram
	telegramSender.Stop()

	log.Println("Application stopped")
}
