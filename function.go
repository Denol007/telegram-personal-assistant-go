package functions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

func init() {
	functions.HTTP("TelegramWebhookHandler", telegramWebhookHandler)
}

// --- Структуры для "расшифровки" JSON от Telegram ---

// Update - это самая внешняя структура, которую присылает Telegram.
type Update struct {
	Message Message `json:"message"`
}

// Message содержит информацию о самом сообщении.
type Message struct {
	Text string `json:"text"`
	Chat Chat   `json:"chat"`
}

// Chat содержит информацию о чате, из которого пришло сообщение.
type Chat struct {
	ID int64 `json:"id"`
}

// --- Логика нашего бота ---

// telegramWebhookHandler - главная функция, которая обрабатывает входящий запрос.
func telegramWebhookHandler(w http.ResponseWriter, r *http.Request) {
	// 1. "Расшифровываем" входящее обновление от Telegram.
	var update Update
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		log.Printf("could not decode request: %v", err)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	// Проверяем, есть ли в сообщении текст.
	if update.Message.Text == "" {
		// Если текста нет (например, прислали стикер), ничего не делаем.
		return
	}

	// 2. Отправляем эхо-сообщение обратно пользователю.
	// Мы передаем ID чата и текст, который нужно отправить.
	if err := sendMessage(update.Message.Chat.ID, update.Message.Text); err != nil {
		log.Printf("could not send message: %v", err)
		// Не отправляем ошибку пользователю, просто логируем у себя.
	}

	// Отвечаем Telegram "OK", чтобы он понял, что мы получили обновление.
	fmt.Fprint(w, "OK")
}

// sendMessage формирует и отправляет запрос к API Telegram.
func sendMessage(chatID int64, text string) error {
	// 3. Получаем токен из переменных окружения.
	// Это безопасный способ хранить секреты.
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		return fmt.Errorf("TELEGRAM_BOT_TOKEN not set")
	}

	// Формируем URL для отправки сообщения.
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)

	// Готовим тело запроса в формате JSON: кому и что отправить.
	requestBody, err := json.Marshal(map[string]interface{}{
		"chat_id": chatID,
		"text":    text,
	})
	if err != nil {
		return err
	}

	// Отправляем HTTP POST запрос в Telegram.
	_, err = http.Post(apiURL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}

	return nil
}
