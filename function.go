package functions

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"encoding/json"
	"bytes"

	"cloud.google.com/go/firestore"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

// firestoreClient - это наш "пропуск" в картотеку.
// Мы создаем его один раз, чтобы не делать это при каждом сообщении.
var firestoreClient *firestore.Client

func init() {
	// Эта функция выполняется один раз при "прогреве" облачной функции.
	// Здесь мы настраиваем подключение к Firestore.
	ctx := context.Background()
	projectID := os.Getenv("GCP_PROJECT_ID") // Получаем ID проекта из переменной окружения
	if projectID == "" {
		log.Fatalf("GCP_PROJECT_ID environment variable must be set.")
	}

	var err error
	firestoreClient, err = firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create Firestore client: %v", err)
	}

	functions.HTTP("TelegramWebhookHandler", telegramWebhookHandler)
}

// Update, Message, Chat - структуры для разбора ответа от Telegram.
// Они остались такими же, как и в эхо-боте.
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

// Note - это структура нашей "карточки" для сохранения в Firestore.
type Note struct {
	Text      string    `firestore:"text"`
	UserID    int64     `firestore:"userID"`
	CreatedAt time.Time `firestore:"createdAt"`
}

func telegramWebhookHandler(w http.ResponseWriter, r *http.Request) {
	var update Update
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		log.Printf("could not decode request: %v", err)
		return
	}

	if update.Message.Text == "" {
		return // Не сохраняем пустые сообщения
	}

	// Создаем новую "карточку"-заметку
	newNote := Note{
		Text:      update.Message.Text,
		UserID:    update.Message.Chat.ID,
		CreatedAt: time.Now(),
	}

	// Добавляем нашу "карточку" в ящик (коллекцию) "notes".
	// Firestore сам присвоит ей уникальный ID.
	_, _, err := firestoreClient.Collection("notes").Add(context.Background(), newNote)
	if err != nil {
		log.Printf("Failed to add note to Firestore: %v", err)
		// Если не удалось сохранить, сообщим об этом пользователю.
		sendMessage(update.Message.Chat.ID, "Не удалось сохранить заметку. 😔")
		return
	}

	// Отправляем пользователю подтверждение.
	sendMessage(update.Message.Chat.ID, "Заметка сохранена! 👍")
}

// sendMessage осталась такой же, как и раньше.
// Мы вынесли ее, чтобы не дублировать код.
func sendMessage(chatID int64, text string) {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)
	requestBody, _ := json.Marshal(map[string]interface{}{
		"chat_id": chatID,
		"text":    text,
	})
	http.Post(apiURL, "application/json", bytes.NewBuffer(requestBody))
}