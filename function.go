package functions

import (
    "context"
    "encoding/json"
    "log"
    "net/http"
    "time"

    "github.com/GoogleCloudPlatform/functions-framework-go/functions"

    // Импортируем наши новые пакеты
    "github.com/Denol007/telegram-personal-assistant-go/internal/note"
    "github.com/Denol007/telegram-personal-assistant-go/internal/store"
    "github.com/Denol007/telegram-personal-assistant-go/internal/telegram"
)

func init() {
    store.Init() // Инициализируем хранилище
    functions.HTTP("TelegramWebhookHandler", TelegramWebhookHandler)
}

// TelegramWebhookHandler - главная функция, точка входа.
func TelegramWebhookHandler(w http.ResponseWriter, r *http.Request) {
    // 1. Декодируем запрос от Telegram
    var update telegram.Update
    if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
        log.Printf("could not decode request: %v", err)
        return
    }

    if update.Message.Text == "" {
        return
    }

    // 2. Создаем заметку
    newNote := note.Note{
        Text:      update.Message.Text,
        UserID:    update.Message.Chat.ID,
        CreatedAt: time.Now(),
    }

    // 3. Сохраняем заметку
    if err := store.SaveNote(context.Background(), newNote); err != nil {
        log.Printf("Failed to save note: %v", err)
        telegram.Send(update.Message.Chat.ID, "Не удалось сохранить заметку. 😔")
        return
    }

    // 4. Отправляем подтверждение
    telegram.Send(update.Message.Chat.ID, "Заметка сохранена! 👍")
}
