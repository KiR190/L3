package config

import (
	"fmt"

	"github.com/DeanPDX/dotconfig"
)

type Config struct {
	HTTPPort         string `env:"APP_PORT"`
	DatabaseHost     string `env:"DATABASE_HOST"`
	DatabasePort     int    `env:"DATABASE_PORT"`
	DatabaseUser     string `env:"DATABASE_USER"`
	DatabasePassword string `env:"DATABASE_PASSWORD"`
	DatabaseName     string `env:"DATABASE_NAME"`
	DatabaseSSLMode  string `env:"DATABASE_SSLMODE"`
}

func Load() (Config, error) {
	cfg, err := dotconfig.FromFileName[Config](".env")
	if err != nil {
		fmt.Printf("Error: %v.", err)
	}

	fmt.Println(cfg)

	return cfg, nil
}

func (c *Config) DatabaseURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.DatabaseUser, c.DatabasePassword, c.DatabaseHost, c.DatabasePort, c.DatabaseName, c.DatabaseSSLMode)
}
