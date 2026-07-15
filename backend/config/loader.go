package config

import (
	"fmt"
	"os"
	"strconv"
)

func Load() (*Config, error) {
	port, err := strconv.Atoi(env("APP_PORT", "8080"))
	if err != nil {
		return nil, fmt.Errorf("invalid APP_PORT: %w", err)
	}

	dbPort, err := strconv.Atoi(env("DATABASE_PORT", "5432"))
	if err != nil {
		return nil, fmt.Errorf("invalid DATABASE_PORT: %w", err)
	}

	redisPort, err := strconv.Atoi(env("REDIS_PORT", "6379"))
	if err != nil {
		return nil, fmt.Errorf("invalid REDIS_PORT: %w", err)
	}

	return &Config{
		App: AppConfig{
			Env:   env("APP_ENV", "development"),
			Name:  env("APP_NAME", "cytaxi"),
			Port:  port,
			Debug: env("APP_DEBUG", "true") == "true",
		},
		DB: DBConfig{
			Host:     env("DATABASE_HOST", "localhost"),
			Port:     dbPort,
			User:     env("DATABASE_USER", "cytaxi"),
			Password: env("DATABASE_PASSWORD", ""),
			Name:     env("DATABASE_NAME", "cytaxi_dev"),
			SSLMode:  env("DATABASE_SSLMODE", "disable"),
		},
		Redis: RedisConfig{
			Host:     env("REDIS_HOST", "localhost"),
			Port:     redisPort,
			Password: env("REDIS_PASSWORD", ""),
		},
		Log: LogConfig{
			Level:  env("LOG_LEVEL", "debug"),
			Format: env("LOG_FORMAT", "json"),
		},
		Auth: AuthConfig{
			JWTSecret: env("JWT_SECRET", ""),
			JWTExpiry: env("JWT_EXPIRATION", "24h"),
		},
		Otel: OtelConfig{
			ServiceName:  env("OTEL_SERVICE_NAME", "cytaxi"),
			OTLPEndpoint: env("OTEL_EXPORTER_OTLP_ENDPOINT", ""),
		},
	}, nil
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

type envLoader struct{}

func NewLoader() Loader {
	return &envLoader{}
}

func (l *envLoader) Load() (*Config, error) {
	return Load()
}
