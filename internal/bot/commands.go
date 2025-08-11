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
		telegram.Send(h.token, chatID, "не удалось получить заметки :(", nil)
		return
	}

	if len(notes) == 0 {
		telegram.Send(h.token, chatID, "У тебя пока нет заметок.", nil)
		return
	}

	var sb strings.Builder
	sb.WriteString("Твои последние заметки:\n")
	for i, n := range notes {
		if n.PhotoID != "" {
			sb.WriteString(fmt.Sprintf("%d. 📸 %s\n", i+1, n.Text))
		} else {
			sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, n.Text))
		}
	}

	telegram.Send(h.token, chatID, sb.String(), nil)
}

func (h *Handler) handleDeleteCommand(chatID int64, text string) {
	var noteNumber int 
	_, err := fmt.Sscanf(text, "/delete %d", &noteNumber)
	if err != nil {
		log.Printf("Неправильный формат команды: %s", text)
		telegram.Send(h.token, chatID, "Пожалуйста, укажи номер заметки, например: /delete 3", nil)
		return
	}

	notes, err := h.store.GetAllNotesByUser(context.Background(), chatID)
	if err != nil {
    	log.Printf("ошибка получения заметок: %v", err)
    	telegram.Send(h.token, chatID, "Не удалось получить заметки", nil)
    	return
	}

	if noteNumber < 1 || noteNumber > len(notes) {
		telegram.Send(h.token, chatID, "Заметки с таким номером не существует.", nil)
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
		telegram.Send(h.token, chatID, "Не удалось сохранить заметку. 😔", nil)
		return
	}
	telegram.Send(h.token, chatID, "Заметка сохранена! 👍", nil)
}

// handleSavePhoto обрабатывает сохранение фотографии.
func (h *Handler) handleSavePhoto(chatID int64, photo []telegram.PhotoSize, caption string) {
	// Выбираем фото наибольшего размера (последнее в массиве)
	var photoID string
	if len(photo) > 0 {
		photoID = photo[len(photo)-1].FileID
	}

	// Используем подпись как текст заметки, или "Фото" если подписи нет
	text := caption
	if text == "" {
		text = "Фото"
	}

	newNote := note.Note{
		Text:      text,
		UserID:    chatID,
		CreatedAt: time.Now(),
		PhotoID:   photoID,
	}

	if err := h.store.SaveNote(context.Background(), newNote); err != nil {
		log.Printf("Failed to save photo note: %v", err)
		telegram.Send(h.token, chatID, "Не удалось сохранить фото. 😔", nil)
		return
	}
	telegram.Send(h.token, chatID, "Фото сохранено! 📸", nil)
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
		telegram.Send(h.token, chatID, "Удаление отменено.", nil)
	} else if strings.HasPrefix(data, "edit_note:") {
		// Пользователь нажал кнопку редактирования
		noteID := strings.TrimPrefix(data, "edit_note:")
		h.handleEditFromCallback(chatID, noteID)
	} else if strings.HasPrefix(data, "delete_note:") {
		// Пользователь нажал кнопку удаления
		noteID := strings.TrimPrefix(data, "delete_note:")
		h.handleDeleteFromCallback(chatID, noteID)
	}
}

// handleConfirmDelete выполняет фактическое удаление заметки.
func (h *Handler) handleConfirmDelete(chatID int64, noteID string) {
	err := h.store.DeleteNote(context.Background(), noteID)
	if err != nil {
		log.Printf("ошибка удаления заметки %s: %v", noteID, err)
		telegram.Send(h.token, chatID, "Не удалось удалить заметку :(", nil)
		return
	}
	telegram.Send(h.token, chatID, "Заметка успешно удалена! 🗑️", nil)
}

// handleEditFromCallback обрабатывает нажатие кнопки "Редактировать".
func (h *Handler) handleEditFromCallback(chatID int64, noteID string) {
	// Получаем заметку по ID
	note, err := h.store.GetNoteByID(context.Background(), noteID)
	if err != nil {
		log.Printf("ошибка получения заметки %s: %v", noteID, err)
		telegram.Send(h.token, chatID, "Не удалось найти заметку", nil)
		return
	}

	// Формируем сообщение с текущим текстом заметки и скрытым ID
	messageText := fmt.Sprintf("Текущий текст заметки:\n\n\"%s\"\n\nОтправь новый текст заметки:\nedit_note:%s", note.Text, note.ID)
	
	// Создаем ForceReply
	forceReply := telegram.ForceReply{
		ForceReply:            true,
		InputFieldPlaceholder: "Введи новый текст заметки...",
		Selective:             false,
	}
	
	// Отправляем сообщение с ForceReply
	telegram.Send(h.token, chatID, messageText, &forceReply)
}

// handleDeleteFromCallback обрабатывает нажатие кнопки "Удалить".
func (h *Handler) handleDeleteFromCallback(chatID int64, noteID string) {
	// Получаем заметку по ID
	note, err := h.store.GetNoteByID(context.Background(), noteID)
	if err != nil {
		log.Printf("ошибка получения заметки %s: %v", noteID, err)
		telegram.Send(h.token, chatID, "Не удалось найти заметку", nil)
		return
	}

	// Формируем текст сообщения
	messageText := fmt.Sprintf("Вы уверены, что хотите удалить заметку?\n\n\"%s\"", note.Text)
	
	// Создаем клавиатуру с кнопками подтверждения
	keyboard := telegram.InlineKeyboardMarkup{
		InlineKeyboard: [][]telegram.InlineKeyboardButton{
			{
				{
					Text:         "Удалить",
					CallbackData: fmt.Sprintf("delete_note:%s", note.ID),
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


func (h *Handler) handleEditCommand(chatID int64, text string) {
	var noteNumber int
	_, err := fmt.Sscanf(text, "/edit %d", &noteNumber)
	if err != nil {
		log.Printf("Неправильный формат команды: %s", text)
		telegram.Send(h.token, chatID, "Пожалуйста, укажи номер заметки, например: /edit 3", nil)
		return
	}

	notes, err := h.store.GetAllNotesByUser(context.Background(), chatID)
	if err != nil {
		log.Printf("ошибка получения заметок: %v", err)
		telegram.Send(h.token, chatID, "Не удалось получить заметки", nil)
		return
	}

	if noteNumber < 1 || noteNumber > len(notes) {
		telegram.Send(h.token, chatID, "Заметки с таким номером не существует.", nil)
		return
	}

	noteToEdit := notes[noteNumber-1]
	
	// Формируем сообщение с текущим текстом заметки и скрытым ID
	messageText := fmt.Sprintf("Текущий текст заметки:\n\n\"%s\"\n\nОтправь новый текст заметки:\nedit_note:%s", noteToEdit.Text, noteToEdit.ID)
	
	// Создаем ForceReply
	forceReply := telegram.ForceReply{
		ForceReply:            true,
		InputFieldPlaceholder: "Введи новый текст заметки...",
		Selective:             false,
	}
	
	// Отправляем сообщение с ForceReply
	telegram.Send(h.token, chatID, messageText, &forceReply)
}

// handleEditResponse обрабатывает ответ пользователя на запрос редактирования заметки.
func (h *Handler) handleEditResponse(chatID int64, noteID string, newText string) {
	err := h.store.UpdateNote(context.Background(), noteID, newText)
	if err != nil {
		log.Printf("ошибка обновления заметки %s: %v", noteID, err)
		telegram.Send(h.token, chatID, "Не удалось обновить заметку :(", nil)
		return
	}
	telegram.Send(h.token, chatID, "Заметка успешно обновлена! ✏️", nil)
}

// handleReplyMessage обрабатывает ответы пользователя на ForceReply.
func (h *Handler) handleReplyMessage(chatID int64, text string, replyToMessage *telegram.Message) {
	// Проверяем, содержит ли исходное сообщение информацию о редактировании заметки
	if strings.Contains(replyToMessage.Text, "edit_note:") {
		// Извлекаем ID заметки из текста исходного сообщения
		// Ищем паттерн "edit_note:ID" в тексте
		lines := strings.Split(replyToMessage.Text, "\n")
		for _, line := range lines {
			if strings.Contains(line, "edit_note:") {
				// Извлекаем ID из строки вида "Отправь новый текст заметки: edit_note:abc123"
				parts := strings.Split(line, "edit_note:")
				if len(parts) >= 2 {
					noteID := strings.TrimSpace(parts[1])
					h.handleEditResponse(chatID, noteID, text)
					return
				}
			}
		}
	}
	
	// Если не удалось определить тип ответа, сохраняем как новую заметку
	h.handleSaveNote(chatID, text)
}

// handleShowCommand обрабатывает команду /show для отображения заметки с кнопками.
func (h *Handler) handleShowCommand(chatID int64, text string) {
	var noteNumber int
	_, err := fmt.Sscanf(text, "/show %d", &noteNumber)
	if err != nil {
		log.Printf("Неправильный формат команды: %s", text)
		telegram.Send(h.token, chatID, "Пожалуйста, укажи номер заметки, например: /show 3", nil)
		return
	}

	notes, err := h.store.GetAllNotesByUser(context.Background(), chatID)
	if err != nil {
		log.Printf("ошибка получения заметок: %v", err)
		telegram.Send(h.token, chatID, "Не удалось получить заметки", nil)
		return
	}

	if noteNumber < 1 || noteNumber > len(notes) {
		telegram.Send(h.token, chatID, "Заметки с таким номером не существует.", nil)
		return
	}

	note := notes[noteNumber-1]
	
	// Создаем клавиатуру с кнопками редактирования и удаления
	keyboard := telegram.InlineKeyboardMarkup{
		InlineKeyboard: [][]telegram.InlineKeyboardButton{
			{
				{
					Text:         "✏️ Редактировать",
					CallbackData: fmt.Sprintf("edit_note:%s", note.ID),
				},
				{
					Text:         "🗑️ Удалить",
					CallbackData: fmt.Sprintf("delete_note:%s", note.ID),
				},
			},
		},
	}

	// Если это заметка с фото, отправляем фото
	if note.PhotoID != "" {
		telegram.SendPhoto(h.token, chatID, note.PhotoID, note.Text, &keyboard)
	} else {
		// Если это текстовая заметка, отправляем текст
		telegram.Send(h.token, chatID, note.Text, &keyboard)
	}
}
