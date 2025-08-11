// Файл: internal/note/note.go
package note

import "time"

// Note - это структура нашей "карточки" для сохранения в Firestore.
type Note struct {
	Text      string    `firestore:"text"`
	UserID    int64     `firestore:"userID"`
	CreatedAt time.Time `firestore:"createdAt"`
	PhotoID   string    `firestore:"photoID,omitempty"` // ID фото в Telegram
	ID        string    `firestore:"-"`
}
