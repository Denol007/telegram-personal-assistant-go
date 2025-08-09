package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// Структуры для разбора ответа от Telegram.
type Update struct {
	Message Message `json:"message"`
}
type Message struct {
	Text string `json:"text"`
	Chat Chat   `json:"chat"`
}
type Chat struct {
	ID int64 `json:"id"`
}

// Send отправляет сообщение в указанный чат Telegram.
func Send(chatID int64, text string) {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)
	requestBody, _ := json.Marshal(map[string]interface{}{
		"chat_id": chatID,
		"text":    text,
	})
	http.Post(apiURL, "application/json", bytes.NewBuffer(requestBody))
}