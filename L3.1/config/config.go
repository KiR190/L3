package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	HTTPPort       string
	DatabaseURL    string
	RabbitMQURL    string
	TG_BOT_TOKEN   string
	SMTP_HOST      string
	SMTP_PORT      int
	SMTP_USERNAME  string
	SMTP_PASSWORD  string
	SMTP_FROM      string
	REDIS_ADDR     string
	REDIS_PASSWORD string
}

func Load() (*Config, error) {
	// Загружаем .env файл если он существует
	_ = godotenv.Load()

	port, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))

	cfg := &Config{
		HTTPPort:       getEnv("APP_PORT", ":8080"),
		DatabaseURL:    getEnv("DATABASE_URL", "postgres://customer:customer_pass@postgres:5432/orders_db?sslmode=disable"),
		RabbitMQURL:    getEnv("RabbitMQ_URL", "amqp://guest:guest@rabbitmq:5672/"),
		TG_BOT_TOKEN:   getEnv("TG_BOT_TOKEN", ""),
		SMTP_HOST:      getEnv("SMTP_HOST", ""),
		SMTP_PORT:      port,
		SMTP_USERNAME:  getEnv("SMTP_USERNAME", ""),
		SMTP_PASSWORD:  getEnv("SMTP_PASSWORD", ""),
		SMTP_FROM:      getEnv("SMTP_FROM", ""),
		REDIS_ADDR:     getEnv("REDIS_ADDR", "redis:6379"),
		REDIS_PASSWORD: getEnv("REDIS_PASSWORD", ""),
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
