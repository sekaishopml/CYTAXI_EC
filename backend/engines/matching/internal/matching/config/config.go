package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Port     int
	Env      string
	LogLevel string
}

func Load() (*Config, error) {
	port, err := strconv.Atoi(env("MATCHING_PORT", "8089"))
	if err != nil {
		return nil, fmt.Errorf("invalid MATCHING_PORT: %w", err)
	}
	return &Config{
		Port:     port,
		Env:      env("APP_ENV", "development"),
		LogLevel: env("LOG_LEVEL", "info"),
	}, nil
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
