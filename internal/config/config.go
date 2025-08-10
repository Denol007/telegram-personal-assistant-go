// Файл: internal/config/config.go
package config

import (
	"fmt"
	"os"
)

// Config хранит всю конфигурацию приложения.
type Config struct {
	TelegramToken string
	ProjectID     string
}

// Load загружает конфигурацию из переменных окружения.
func Load() (*Config, error) {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("переменная окружения TELEGRAM_BOT_TOKEN не установлена")
	}

	projectID := os.Getenv("GCP_PROJECT_ID")
	if projectID == "" {
		return nil, fmt.Errorf("переменная окружения GCP_PROJECT_ID не установлена")
	}

	return &Config{
		TelegramToken: token,
		ProjectID:     projectID,
	}, nil
}
