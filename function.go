// Файл: function.go
package functions

import (
	"log"
	"net/http"

	"github.com/Denol007/telegram-personal-assistant-go/internal/bot"
	"github.com/Denol007/telegram-personal-assistant-go/internal/config"
	"github.com/Denol007/telegram-personal-assistant-go/internal/store"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

var botHandler *bot.Handler

func init() {
	// 1. Загружаем конфигурацию (токены, ID проекта)
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	// 2. Инициализируем хранилище (подключаемся к Firestore)
	noteStore, err := store.New(cfg.ProjectID)
	if err != nil {
		log.Fatalf("Ошибка инициализации хранилища: %v", err)
	}

	// 3. Создаем наш главный обработчик, передавая ему все зависимости
	botHandler = bot.NewHandler(cfg.TelegramToken, noteStore)

	// 4. Регистрируем HTTP-функцию
	functions.HTTP("TelegramWebhookHandler", TelegramWebhookHandler)
}

// TelegramWebhookHandler просто вызывает наш основной обработчик.
func TelegramWebhookHandler(w http.ResponseWriter, r *http.Request) {
	botHandler.HandleUpdate(w, r)
}
