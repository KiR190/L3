package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"shortener/config"
	"shortener/internal/cache"
	"shortener/internal/generator"
	"shortener/internal/handler"
	"shortener/internal/repository"
	"shortener/internal/service"

	"github.com/wb-go/wbf/dbpg"
	"github.com/wb-go/wbf/ginext"
)

func main() {
	// Загрузка конфигурации
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Настройка базы данных
	dbOpts := &dbpg.Options{MaxOpenConns: 5, MaxIdleConns: 1}
	dbConn, err := dbpg.New(cfg.DatabaseURL, []string{}, dbOpts)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Репозитории
	shortRepo := repository.NewShortURLRepo(dbConn)
	analyticsRepo := repository.NewAnalyticsRepo(dbConn)

	// Генератор коротких кодов (sqids-go)
	sqidsGen := generator.NewShortCodeGenerator() // твой конструктор

	// Инициализация кеша
	redisCache := cache.NewCache(cfg.REDIS_ADDR, cfg.REDIS_PASSWORD, 0)

	// Сервис
	shortService := service.NewShortenerService(shortRepo, analyticsRepo, redisCache, *sqidsGen)
	ctx := context.Background()
	shortService.RestoreCacheFromDB(ctx)

	// Хендлер
	h := handler.NewURLHandler(shortService)

	// Gin Engine
	r := ginext.New("")
	r.Use(ginext.Logger())
	r.Use(ginext.Recovery())

	// Роуты
	handler.SetupRouter(r, h)

	srv := &http.Server{
		Addr:    cfg.HTTPPort,
		Handler: r,
	}

	// Запуск
	go func() {
		log.Printf("Starting server at %s", cfg.HTTPPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server failed: %v", err)
		}
	}()

	// ==== GRACEFUL SHUTDOWN ====

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Println("Shutting down server...")

	ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Завершаем HTTP сервер
	if err := srv.Shutdown(ctxShutdown); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting gracefully")
}
