// Файл: internal/bot/commands.go
package bot

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Denol007/telegram-personal-assistant-go/internal/note"
	"github.com/Denol007/telegram-personal-assistant-go/internal/telegram"
)

// handleListCommand обрабатывает команду /list.
func (h *Handler) handleListCommand(chatID int64) {
	notes, err := h.store.GetAllNotesByUser(context.Background(), chatID)
	if err != nil {
		log.Printf("ошибка получения заметок: %v", err)
		telegram.Send(h.token, chatID, "не удалось получить заметки :(")
		return
	}

	if len(notes) == 0 {
		telegram.Send(h.token, chatID, "У тебя пока нет заметок.")
		return
	}

	var sb strings.Builder
	sb.WriteString("Твои последние заметки:\n")
	for i, n := range notes {
		sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, n.Text))
	}

	telegram.Send(h.token, chatID, sb.String())
}

// handleSaveNote обрабатывает сохранение новой заметки.
func (h *Handler) handleSaveNote(chatID int64, text string) {
	newNote := note.Note{
		Text:      text,
		UserID:    chatID,
		CreatedAt: time.Now(),
	}

	if err := h.store.SaveNote(context.Background(), newNote); err != nil {
		log.Printf("Failed to save note: %v", err)
		telegram.Send(h.token, chatID, "Не удалось сохранить заметку. 😔")
		return
	}
	telegram.Send(h.token, chatID, "Заметка сохранена! 👍")
}
