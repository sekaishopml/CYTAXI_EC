package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Port          int
	Env           string
	LogLevel      string
	AuthSecret    string
	RateLimitRPS  int
	BackendHosts  map[string]string
}

func Load() (*Config, error) {
	port, err := strconv.Atoi(env("GATEWAY_PORT", "8000"))
	if err != nil {
		return nil, fmt.Errorf("invalid GATEWAY_PORT: %w", err)
	}
	rps, err := strconv.Atoi(env("GATEWAY_RATE_LIMIT_RPS", "100"))
	if err != nil {
		return nil, fmt.Errorf("invalid GATEWAY_RATE_LIMIT_RPS: %w", err)
	}
	return &Config{
		Port:         port,
		Env:          env("APP_ENV", "development"),
		LogLevel:     env("LOG_LEVEL", "info"),
		AuthSecret:   env("JWT_SECRET", ""),
		RateLimitRPS: rps,
		BackendHosts: map[string]string{
			"trip":         env("BACKEND_TRIP", "http://localhost:8087"),
			"pricing":      env("BACKEND_PRICING", "http://localhost:8088"),
			"payment":      env("BACKEND_PAYMENT", "http://localhost:8091"),
			"customer":     env("BACKEND_CUSTOMER", "http://localhost:8085"),
			"driver":       env("BACKEND_DRIVER", "http://localhost:8086"),
			"notification": env("BACKEND_NOTIFICATION", "http://localhost:8090"),
			"admin":        env("BACKEND_ADMIN", "http://localhost:8094"),
			"analytics":    env("BACKEND_ANALYTICS", "http://localhost:8093"),
			"matching":     env("BACKEND_MATCHING", "http://localhost:8089"),
		},
	}, nil
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
