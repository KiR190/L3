package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"image-processor/config"
	"image-processor/internal/handler"
	"image-processor/internal/processor"
	"image-processor/internal/queue"
	"image-processor/internal/service"
	"image-processor/internal/storage"
	"image-processor/internal/worker"
)

func main() {
	// Load config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	storage, err := storage.NewMinioStorage(
		cfg.MinioURL,
		cfg.MinioKey,
		cfg.MinioSecret,
		cfg.BucketName,
		false,
	)
	if err != nil {
		log.Fatalf("storage init failed: %v", err)
	}

	queue := queue.NewQueue([]string{cfg.KafkaURL}, cfg.KafkaTopic, "images-group")

	svc := service.NewImageService(storage, queue)
	h := handler.NewHandler(svc)

	// Воркер
	processor := processor.NewProcessor()
	worker := worker.NewWorker(storage, processor, queue)
	worker.Start()

	// Роутер
	router := handler.NewRouter(h)
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
