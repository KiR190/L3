package config

import (
	"fmt"

	"github.com/DeanPDX/dotconfig"
)

type Config struct {
	// Application Settings
	HTTPPort string `env:"APP_PORT" default:"8080"`

	// Database Settings
	DatabaseHost     string `env:"DATABASE_HOST" default:"localhost"`
	DatabasePort     int    `env:"DATABASE_PORT" default:"5432"`
	DatabaseUser     string `env:"DATABASE_USER" default:"eventbooker"`
	DatabasePassword string `env:"DATABASE_PASSWORD" default:"password"`
	DatabaseName     string `env:"DATABASE_NAME" default:"eventbooker"`
	DatabaseSSLMode  string `env:"DATABASE_SSLMODE" default:"disable"`

	// Kafka Settings
	KafkaURL   string `env:"KAFKA_URL" default:"kafka:9092"`
	KafkaTopic string `env:"KAFKA_TOPIC" default:"booking-expirations"`

	// Business Logic Settings
	DefaultPaymentTimeout int `env:"DEFAULT_PAYMENT_TIMEOUT_MINUTES" default:"30"`

	// JWT Settings
	JWTSecret          string `env:"JWT_SECRET" default:"change-this-secret-key-in-production"`
	JWTExpirationHours int    `env:"JWT_EXPIRATION_HOURS" default:"24"`

	// SMTP Settings
	SMTPHost     string `env:"SMTP_HOST" default:"smtp.gmail.com"`
	SMTPPort     int    `env:"SMTP_PORT" default:"587"`
	SMTPUsername string `env:"SMTP_USERNAME"`
	SMTPPassword string `env:"SMTP_PASSWORD"`
	SMTPFrom     string `env:"SMTP_FROM"`

	// Telegram Bot
	TelegramBotToken string `env:"TG_BOT_TOKEN"`
}

func Load() (Config, error) {
	cfg, err := dotconfig.FromFileName[Config](".env")
	if err != nil {
		fmt.Printf("Error loading config: %v. Using defaults...\n", err)
	}

	fmt.Printf("Loaded config: %+v\n", cfg)

	return cfg, nil
}

func (c *Config) DatabaseURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.DatabaseUser, c.DatabasePassword, c.DatabaseHost, c.DatabasePort, c.DatabaseName, c.DatabaseSSLMode)
}
