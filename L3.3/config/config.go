package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	HTTPPort    string
	DatabaseURL string
}

func Load() (*Config, error) {
	// Загружаем .env файл если он существует
	_ = godotenv.Load()

	fmt.Println("DATABASE_URL ENV =", os.Getenv("DATABASE_URL"))

	cfg := &Config{
		HTTPPort:    getEnv("APP_PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://customer:customer_pass@postgres:5432/orders_db?sslmode=disable"),
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
