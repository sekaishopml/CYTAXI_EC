package whatsapp

import (
	"fmt"
	"os"
)

type Config struct {
	Provider      ProviderKind
	SessionID     string
	Reconnect     bool
	AutoLoadQR    bool
	WebhookURL    string
	MaxRetries    int
	RetryDelaySec int
}

func LoadConfig() (*Config, error) {
	provider := ProviderKind(env("WHATSAPP_PROVIDER", "whatsmeow"))

	switch provider {
	case ProviderWhatsMeow, ProviderWAWebJS, ProviderBusinessAPI:
	default:
		return nil, fmt.Errorf("invalid WHATSAPP_PROVIDER: %s", provider)
	}

	return &Config{
		Provider:      provider,
		SessionID:     env("WHATSAPP_SESSION_ID", "cytaxi-main"),
		Reconnect:     env("WHATSAPP_RECONNECT", "true") == "true",
		AutoLoadQR:    env("WHATSAPP_AUTO_LOAD_QR", "true") == "true",
		WebhookURL:    env("WHATSAPP_WEBHOOK_URL", ""),
		MaxRetries:    3,
		RetryDelaySec: 5,
	}, nil
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
