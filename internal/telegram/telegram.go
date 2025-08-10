// Файл: internal/telegram/telegram.go
package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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
func Send(token string, chatID int64, text string) {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)

	requestBody, err := json.Marshal(map[string]interface{}{
		"chat_id": chatID,
		"text":    text,
	})
	if err != nil {
		log.Printf("Telegram Send: ошибка marshal json: %v", err)
		return
	}

	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Printf("Telegram Send: ошибка отправки запроса: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		log.Printf("Telegram Send: получен статус %d", resp.StatusCode)
	}
}
