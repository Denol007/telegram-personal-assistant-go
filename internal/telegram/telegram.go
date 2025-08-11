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
	Message       *Message       `json:"message"`
	CallbackQuery *CallbackQuery `json:"callback_query"`
}

type Message struct {
	Text           string      `json:"text"`
	Caption        string      `json:"caption"`
	Photo          []PhotoSize `json:"photo"`
	Chat           Chat        `json:"chat"`
	ReplyToMessage *Message    `json:"reply_to_message"`
}

type PhotoSize struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	FileSize     int    `json:"file_size,omitempty"`
}

type Chat struct {
	ID int64 `json:"id"`
}

type CallbackQuery struct {
	ID      string  `json:"id"`
	From    User    `json:"from"`
	Message Message `json:"message"`
	Data    string  `json:"data"`
}

type User struct {
	ID int64 `json:"id"`
}

type InlineKeyboardButton struct {
	Text         string `json:"text"`
	CallbackData string `json:"callback_data"`
}

// InlineKeyboardMarkup представляет всю клавиатуру (набор кнопок).
type InlineKeyboardMarkup struct {
	InlineKeyboard [][]InlineKeyboardButton `json:"inline_keyboard"`
}

// ForceReply заставляет пользователя ответить на сообщение.
type ForceReply struct {
	ForceReply            bool   `json:"force_reply"`
	InputFieldPlaceholder string `json:"input_field_placeholder,omitempty"`
	Selective             bool   `json:"selective,omitempty"`
}

// sendMessagePayload представляет тело запроса для метода sendMessage.
type sendMessagePayload struct {
	ChatID      int64       `json:"chat_id"`
	Text        string      `json:"text"`
	// ReplyMarkup может быть либо клавиатурой, либо ForceReply.
	// `interface{}` позволяет нам использовать и то, и другое.
	ReplyMarkup interface{} `json:"reply_markup,omitempty"`
}


// Send отправляет сообщение в указанный чат Telegram, опционально с разметкой.
func Send(token string, chatID int64, text string, markup interface{}) {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)

	payload := sendMessagePayload{
		ChatID:      chatID,
		Text:        text,
		ReplyMarkup: markup,
	}

	requestBody, err := json.Marshal(payload)
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

// AnswerCallbackQuery отправляет ответ на callback query.
func AnswerCallbackQuery(token string, callbackQueryID string) {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/answerCallbackQuery", token)
	
	payload := map[string]string{
		"callback_query_id": callbackQueryID,
	}
	
	requestBody, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Telegram AnswerCallbackQuery: ошибка marshal json: %v", err)
		return
	}

	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Printf("Telegram AnswerCallbackQuery: ошибка отправки запроса: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		log.Printf("Telegram AnswerCallbackQuery: получен статус %d", resp.StatusCode)
	}
}

// sendPhotoPayload представляет тело запроса для метода sendPhoto.
type sendPhotoPayload struct {
	ChatID      int64       `json:"chat_id"`
	Photo       string      `json:"photo"`
	Caption     string      `json:"caption,omitempty"`
	ReplyMarkup interface{} `json:"reply_markup,omitempty"`
}

// SendPhoto отправляет фотографию в указанный чат Telegram.
func SendPhoto(token string, chatID int64, photoID string, caption string, markup interface{}) {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendPhoto", token)

	payload := sendPhotoPayload{
		ChatID:      chatID,
		Photo:       photoID,
		Caption:     caption,
		ReplyMarkup: markup,
	}

	requestBody, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Telegram SendPhoto: ошибка marshal json: %v", err)
		return
	}

	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Printf("Telegram SendPhoto: ошибка отправки запроса: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		log.Printf("Telegram SendPhoto: получен статус %d", resp.StatusCode)
	}
}
