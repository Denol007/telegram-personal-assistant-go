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

func (h *Handler) handleDeleteCommand(chatID int64, text string) {
	var noteNumber int 
	_, err := fmt.Sscanf(text, "/delete %d", &noteNumber)
	if err != nil {
		log.Printf("Неправильный формат команды: %s", text)
		telegram.Send(h.token, chatID, "Пожалуйста, укажи номер заметки, например: /delete 3")
		return
	}

	notes, err := h.store.GetAllNotesByUser(context.Background(), chatID)
	if err != nil {
    	log.Printf("ошибка получения заметок: %v", err)
    	telegram.Send(h.token, chatID, "Не удалось получить заметки")
    	return
	}

	if noteNumber < 1 || noteNumber > len(notes) {
		telegram.Send(h.token, chatID, "Заметки с таким номером не существует.")
    	return
	} else {
		noteToDelete := notes[noteNumber-1]
		
		// Формируем текст сообщения
		messageText := fmt.Sprintf("Вы уверены, что хотите удалить заметку?\n\n\"%s\"", noteToDelete.Text)
		
		// Создаем клавиатуру с кнопками подтверждения
		keyboard := telegram.InlineKeyboardMarkup{
			InlineKeyboard: [][]telegram.InlineKeyboardButton{
				{
					{
						Text:         "Удалить",
						CallbackData: fmt.Sprintf("delete_note:%s", noteToDelete.ID),
					},
					{
						Text:         "Отмена",
						CallbackData: "cancel_delete",
					},
				},
			},
		}
		
		// Отправляем сообщение с клавиатурой
		telegram.Send(h.token, chatID, messageText, &keyboard)
	}

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

// handleCallbackQuery обрабатывает нажатия на инлайн-кнопки.
func (h *Handler) handleCallbackQuery(callbackQuery *telegram.CallbackQuery) {
	chatID := callbackQuery.Message.Chat.ID
	data := callbackQuery.Data

	// Отправляем ответ на callback query (убирает "часики" на кнопке)
	telegram.AnswerCallbackQuery(h.token, callbackQuery.ID)

	// Проверяем, какая кнопка была нажата
	if data == "cancel_delete" {
		// Пользователь отменил удаление
		telegram.Send(h.token, chatID, "Удаление отменено.")
	} else if strings.HasPrefix(data, "delete_note:") {
		// Пользователь подтвердил удаление
		// Извлекаем ID заметки из callback data
		noteID := strings.TrimPrefix(data, "delete_note:")
		h.handleConfirmDelete(chatID, noteID)
	}
}

// handleConfirmDelete выполняет фактическое удаление заметки.
func (h *Handler) handleConfirmDelete(chatID int64, noteID string) {
	err := h.store.DeleteNote(context.Background(), noteID)
	if err != nil {
		log.Printf("ошибка удаления заметки %s: %v", noteID, err)
		telegram.Send(h.token, chatID, "Не удалось удалить заметку :(")
		return
	}
	telegram.Send(h.token, chatID, "Заметка успешно удалена! 🗑️")
}
