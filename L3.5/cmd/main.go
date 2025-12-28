package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"event-booker/config"
	"event-booker/internal/auth"
	"event-booker/internal/handler"
	"event-booker/internal/queue"
	"event-booker/internal/repository"
	"event-booker/internal/sender"
	"event-booker/internal/service"
	"event-booker/internal/worker"

	"github.com/wb-go/wbf/dbpg"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to database
	dbOpts := &dbpg.Options{MaxOpenConns: 10, MaxIdleConns: 5}
	dbConn, err := dbpg.New(cfg.DatabaseURL(), []string{}, dbOpts)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("Connected to database")

	// Initialize repositories
	repo := repository.NewRepository(dbConn)
	userRepo := repository.NewUserRepository(dbConn)

	// Initialize auth service
	authService := auth.NewAuthService(cfg.JWTSecret, cfg.JWTExpirationHours)
	log.Println("Auth service initialized")

	// Initialize Email sender
	var emailSender *sender.EmailSender
	if cfg.SMTPHost != "" {
		emailSender = sender.NewEmailSender(
			cfg.SMTPHost,
			cfg.SMTPPort,
			cfg.SMTPUsername,
			cfg.SMTPPassword,
			cfg.SMTPFrom,
		)
		log.Println("Email sender initialized")
	} else {
		log.Println("Warning: Email sender not configured (SMTP_HOST not set)")
	}

	// Initialize Telegram sender
	var telegramSender *sender.TelegramSender
	if cfg.TelegramBotToken != "" {
		telegramSender, err = sender.NewTelegramSender(cfg.TelegramBotToken)
		if err != nil {
			log.Printf("Warning: Failed to initialize Telegram bot: %v", err)
		} else {
			log.Printf("Telegram bot initialized")
			
			// Start Telegram bot in background
			go telegramSender.ListenAndServe()
		}
	} else {
		log.Println("Warning: Telegram bot not configured (TG_BOT_TOKEN not set)")
	}

	// Initialize multi sender (with priority fallback)
	multiSender := sender.NewMultiSender(emailSender, telegramSender)

	// Initialize queue
	queue := queue.NewQueue([]string{cfg.KafkaURL}, cfg.KafkaTopic, "event-booker-group")
	log.Println("Connected to Kafka")

	// Initialize service
	svc := service.NewService(repo, userRepo, queue, authService, multiSender)

	// Initialize handler
	h := handler.NewHandler(svc, cfg.DefaultPaymentTimeout)

	// Create admin middleware
	adminMiddleware := auth.NewAdminMiddleware(userRepo)

	// Create router and HTTP server
	router := handler.NewRouter(h, authService, adminMiddleware)
	httpServer := handler.NewHTTPServer(cfg.HTTPPort, router)

	// Start worker
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	w := worker.NewWorker(repo, queue, svc)
	w.Start(ctx)

	log.Printf("Server starting on http://localhost:%s", cfg.HTTPPort)

	// Start HTTP server in a goroutine
	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			if err.Error() == "http: Server closed" {
				return
			}
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Println("Shutting down server...")

	// Cancel worker context
	cancel()

	// Shutdown HTTP server
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	// Stop Telegram bot if initialized
	if telegramSender != nil {
		telegramSender.Stop()
	}

	// Close queue connections
	if err := queue.Close(); err != nil {
		log.Printf("Error closing queue: %v", err)
	}

	log.Println("Server exited gracefully")
}

