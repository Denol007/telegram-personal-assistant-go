// Файл: internal/bot/handler.go
package bot

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/Denol007/telegram-personal-assistant-go/internal/store"
	"github.com/Denol007/telegram-personal-assistant-go/internal/telegram"
)

// Handler - главный обработчик бота.
type Handler struct {
	token string
	store *store.Store
}

// NewHandler создает новый обработчик.
func NewHandler(token string, store *store.Store) *Handler {
	return &Handler{
		token: token,
		store: store,
	}
}

// HandleUpdate - главная функция, точка входа для всех обновлений от Telegram.
func (h *Handler) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	var update telegram.Update
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		log.Printf("could not decode request: %v", err)
		return
	}

	// Сначала проверяем, не пришел ли CallbackQuery
	if update.CallbackQuery != nil {
		h.handleCallbackQuery(update.CallbackQuery)
		return
	}

	// Если это не CallbackQuery, проверяем обычное сообщение
	if update.Message == nil || update.Message.Text == "" {
		return
	}

	// Разбор команд
	if strings.HasPrefix(update.Message.Text, "/list") {
		h.handleListCommand(update.Message.Chat.ID)

	} else if strings.HasPrefix(update.Message.Text, "/delete") {
		h.handleDeleteCommand(update.Message.Chat.ID, update.Message.Text)
	} else {
		// Действие по умолчанию - сохранить заметку
		h.handleSaveNote(update.Message.Chat.ID, update.Message.Text)
	}
}
