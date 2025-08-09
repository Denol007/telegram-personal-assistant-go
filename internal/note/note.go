package note

import "time"

// Note - это структура нашей "карточки" для сохранения в Firestore.
type Note struct {
	Text      string    `firestore:"text"`
	UserID    int64     `firestore:"userID"`
	CreatedAt time.Time `firestore:"createdAt"`
}