package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"sales-tracker/config"
	"sales-tracker/internal/handler"
	"sales-tracker/internal/repository"
	"sales-tracker/internal/service"

	"github.com/wb-go/wbf/dbpg"
)

func main() {
	// Загрузка конфигурации
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Подключение к БД
	dbOpts := &dbpg.Options{MaxOpenConns: 5, MaxIdleConns: 1}
	dbConn, err := dbpg.New(cfg.DatabaseURL(), []string{}, dbOpts)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Инициализация репозитория
	itemsRepo := repository.NewRepository(dbConn)

	// Сервис
	itemsService := service.NewService(itemsRepo)

	// Хендлер
	itemsHandler := handler.NewHandler(itemsService)

	// Роутер
	router := handler.NewRouter(itemsHandler)
	httpServer := handler.NewHTTPServer(cfg.HTTPPort, router)

	log.Printf("Server started at http://localhost:%s", cfg.HTTPPort)

	// Запуск сервера
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully")
}
