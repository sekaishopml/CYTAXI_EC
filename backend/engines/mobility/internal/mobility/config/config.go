package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Port            int
	Env             string
	LogLevel        string
	DefaultStrategy string
	MaxCandidates   int
	DispatchTimeout int
}

func Load() (*Config, error) {
	port, err := strconv.Atoi(env("MOBILITY_PORT", "8084"))
	if err != nil {
		return nil, fmt.Errorf("invalid MOBILITY_PORT: %w", err)
	}
	maxCandidates, err := strconv.Atoi(env("MOBILITY_MAX_CANDIDATES", "20"))
	if err != nil {
		return nil, fmt.Errorf("invalid MOBILITY_MAX_CANDIDATES: %w", err)
	}
	timeout, err := strconv.Atoi(env("MOBILITY_DISPATCH_TIMEOUT", "30"))
	if err != nil {
		return nil, fmt.Errorf("invalid MOBILITY_DISPATCH_TIMEOUT: %w", err)
	}
	return &Config{
		Port:            port,
		Env:             env("APP_ENV", "development"),
		LogLevel:        env("LOG_LEVEL", "info"),
		DefaultStrategy: env("MOBILITY_DEFAULT_STRATEGY", "balanced_score"),
		MaxCandidates:   maxCandidates,
		DispatchTimeout: timeout,
	}, nil
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
