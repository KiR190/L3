package config

import (
	"fmt"

	"github.com/DeanPDX/dotconfig"
)

type Config struct {
	// Application Settings
	HTTPPort string `env:"APP_PORT" default:"8080"`

	// Database Settings
	DatabaseHost     string `env:"DATABASE_HOST" default:"postgres"`
	DatabasePort     int    `env:"DATABASE_PORT" default:"5432"`
	DatabaseUser     string `env:"DATABASE_USER" default:"warehouse"`
	DatabasePassword string `env:"DATABASE_PASSWORD" default:"warehouse"`
	DatabaseName     string `env:"DATABASE_NAME" default:"warehouse"`
	DatabaseSSLMode  string `env:"DATABASE_SSLMODE" default:"disable"`

	// JWT Settings
	JWTSecret          string `env:"JWT_SECRET"`
	JWTExpirationHours int    `env:"JWT_EXPIRATION_HOURS"`
}

func Load() (Config, error) {
	cfg, err := dotconfig.FromFileName[Config](".env")
	if err != nil {
		fmt.Printf("Error loading config: %v. Using defaults...\n", err)
	}

	return cfg, nil
}

func (c *Config) DatabaseURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.DatabaseUser, c.DatabasePassword, c.DatabaseHost, c.DatabasePort, c.DatabaseName, c.DatabaseSSLMode)
}
