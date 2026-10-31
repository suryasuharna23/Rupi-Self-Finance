// Package config membaca konfigurasi dari environment variable. Rahasia
// (token bot, kredensial database) tidak pernah ditulis di kode atau Git.
package config

import (
	"errors"
	"fmt"
	"os"
)

type Config struct {
	DatabaseURL   string
	TelegramToken string
	TelegramMode  string // "polling" (lokal) atau "webhook"
	WebhookURL    string
	WebhookSecret string
	HTTPAddr      string
	LogLevel      string
}

func Load() (Config, error) {
	c := Config{
		DatabaseURL:   os.Getenv("DATABASE_URL"),
		TelegramToken: os.Getenv("TELEGRAM_BOT_TOKEN"),
		TelegramMode:  getenv("TELEGRAM_MODE", "polling"),
		WebhookURL:    os.Getenv("TELEGRAM_WEBHOOK_URL"),
		WebhookSecret: os.Getenv("TELEGRAM_WEBHOOK_SECRET"),
		HTTPAddr:      getenv("HTTP_ADDR", ":8080"),
		LogLevel:      getenv("LOG_LEVEL", "info"),
	}
	var errs []error
	if c.DatabaseURL == "" {
		errs = append(errs, errors.New("DATABASE_URL wajib diisi"))
	}
	if c.TelegramToken == "" {
		errs = append(errs, errors.New("TELEGRAM_BOT_TOKEN wajib diisi"))
	}
	switch c.TelegramMode {
	case "polling":
	case "webhook":
		if c.WebhookURL == "" || c.WebhookSecret == "" {
			errs = append(errs, errors.New("mode webhook butuh TELEGRAM_WEBHOOK_URL dan TELEGRAM_WEBHOOK_SECRET"))
		}
	default:
		errs = append(errs, fmt.Errorf("TELEGRAM_MODE %q tidak dikenal (polling|webhook)", c.TelegramMode))
	}
	return c, errors.Join(errs...)
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
