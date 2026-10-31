// Package config reads configuration from environment variables. Secrets
// (bot token, database credentials) are never written in code or Git.
package config

import (
	"errors"
	"fmt"
	"os"
)

type Config struct {
	DatabaseURL    string
	TelegramToken  string
	TelegramAPIURL string
	TelegramMode   string // "polling" (local) or "webhook"
	WebhookURL     string
	WebhookSecret  string
	HTTPAddr       string
	LogLevel       string
}

func Load() (Config, error) {
	c := Config{
		DatabaseURL:    os.Getenv("DATABASE_URL"),
		TelegramToken:  os.Getenv("TELEGRAM_BOT_TOKEN"),
		TelegramAPIURL: getenv("TELEGRAM_API_URL", "https://api.telegram.org"),
		TelegramMode:   getenv("TELEGRAM_MODE", "polling"),
		WebhookURL:     os.Getenv("TELEGRAM_WEBHOOK_URL"),
		WebhookSecret:  os.Getenv("TELEGRAM_WEBHOOK_SECRET"),
		HTTPAddr:       getenv("HTTP_ADDR", ":8080"),
		LogLevel:       getenv("LOG_LEVEL", "info"),
	}
	var errs []error
	if c.DatabaseURL == "" {
		errs = append(errs, errors.New("DATABASE_URL is required"))
	}
	if c.TelegramToken == "" {
		errs = append(errs, errors.New("TELEGRAM_BOT_TOKEN is required"))
	}
	switch c.TelegramMode {
	case "polling":
	case "webhook":
		if c.WebhookURL == "" || c.WebhookSecret == "" {
			errs = append(errs, errors.New("webhook mode requires TELEGRAM_WEBHOOK_URL and TELEGRAM_WEBHOOK_SECRET"))
		}
	default:
		errs = append(errs, fmt.Errorf("unknown TELEGRAM_MODE %q (polling|webhook)", c.TelegramMode))
	}
	return c, errors.Join(errs...)
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
