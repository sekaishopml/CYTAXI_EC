package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Port      int
	Env       string
	LogLevel  string
	Provider  string
	APIKey    string
	CacheTTL  int
}

func Load() (*Config, error) {
	port, err := strconv.Atoi(env("GEOSPATIAL_PORT", "8082"))
	if err != nil {
		return nil, fmt.Errorf("invalid GEOSPATIAL_PORT: %w", err)
	}

	cacheTTL, err := strconv.Atoi(env("GEOSPATIAL_CACHE_TTL", "300"))
	if err != nil {
		return nil, fmt.Errorf("invalid GEOSPATIAL_CACHE_TTL: %w", err)
	}

	return &Config{
		Port:     port,
		Env:      env("APP_ENV", "development"),
		LogLevel: env("LOG_LEVEL", "info"),
		Provider: env("GEOSPATIAL_PROVIDER", "openstreetmap"),
		APIKey:   env("GEOSPATIAL_API_KEY", ""),
		CacheTTL: cacheTTL,
	}, nil
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
