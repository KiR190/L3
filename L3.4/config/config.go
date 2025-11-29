package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	HTTPPort    string
	KafkaURL    string
	KafkaTopic  string
	MinioURL    string
	MinioKey    string
	MinioSecret string
	BucketName  string
}

func Load() (*Config, error) {
	// Загружаем .env файл если он существует
	_ = godotenv.Load()

	cfg := &Config{
		HTTPPort:    getEnv("APP_PORT", "8080"),
		KafkaURL:    getEnv("KAFKA_URL", "kafka:9092"),
		KafkaTopic:  getEnv("KAFKA_TOPIC", "images"),
		MinioURL:    getEnv("MINIO_URL", "minio:9000"),
		MinioKey:    getEnv("MINIO_ACCESS_KEY", "minio"),
		MinioSecret: getEnv("MINIO_SECRET_KEY", "minio123"),
		BucketName:  getEnv("MINIO_BUCKET", "images"),
	}

	fmt.Println("Loaded config:", cfg)

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
