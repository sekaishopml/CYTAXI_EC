package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Port        int
	Env         string
	LogLevel    string
	PoliciesDir string
	AutoReload  bool
	ReloadSec   int
}

func Load() (*Config, error) {
	port, err := strconv.Atoi(env("POLICY_PORT", "8083"))
	if err != nil {
		return nil, fmt.Errorf("invalid POLICY_PORT: %w", err)
	}

	reloadSec, err := strconv.Atoi(env("POLICY_RELOAD_SEC", "60"))
	if err != nil {
		return nil, fmt.Errorf("invalid POLICY_RELOAD_SEC: %w", err)
	}

	return &Config{
		Port:        port,
		Env:         env("APP_ENV", "development"),
		LogLevel:    env("LOG_LEVEL", "info"),
		PoliciesDir: env("POLICIES_DIR", "./policies"),
		AutoReload:  env("POLICY_AUTO_RELOAD", "true") == "true",
		ReloadSec:   reloadSec,
	}, nil
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
