// Файл: internal/store/firestore.go
package store

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
	"github.com/Denol007/telegram-personal-assistant-go/internal/note"
)

// Store - это наше хранилище.
type Store struct {
	client *firestore.Client
}

// New создает новое подключение к хранилищу.
func New(projectID string) (*Store, error) {
	client, err := firestore.NewClient(context.Background(), projectID)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания клиента Firestore: %w", err)
	}
	return &Store{client: client}, nil
}

// SaveNote сохраняет заметку в Firestore.
func (s *Store) SaveNote(ctx context.Context, n note.Note) error {
	_, _, err := s.client.Collection("notes").Add(ctx, n)
	return err
}

// GetAllNotesByUser находит все заметки для указанного пользователя.
func (s *Store) GetAllNotesByUser(ctx context.Context, userID int64) ([]note.Note, error) {
	var notes []note.Note

	query := s.client.Collection("notes").
		Where("userID", "==", userID).
		OrderBy("createdAt", firestore.Desc)

	iter := query.Documents(ctx)
	defer iter.Stop()
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("ошибка итерации: %w", err)
		}

		var n note.Note
		if err := doc.DataTo(&n); err != nil {
			log.Printf("ошибка преобразования документа %s: %v", doc.Ref.ID, err)
			continue
		}
		 n.ID = doc.Ref.ID

		notes = append(notes, n)
	}
	return notes, nil
}

// DeleteNote удаляет заметку по ID.
func (s *Store) DeleteNote(ctx context.Context, noteID string) error {
	_, err := s.client.Collection("notes").Doc(noteID).Delete(ctx)
	if err != nil {
		return fmt.Errorf("ошибка удаления заметки %s: %w", noteID, err)
	}
	return nil
}

// UpdateNote обновляет текст заметки по ID.
func (s *Store) UpdateNote(ctx context.Context, noteID string, newText string) error {
	_, err := s.client.Collection("notes").Doc(noteID).Update(ctx, []firestore.Update{
		{Path: "text", Value: newText},
	})
	if err != nil {
		return fmt.Errorf("ошибка обновления заметки %s: %w", noteID, err)
	}
	return nil
}

// GetNoteByID получает заметку по ID.
func (s *Store) GetNoteByID(ctx context.Context, noteID string) (*note.Note, error) {
	doc, err := s.client.Collection("notes").Doc(noteID).Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения заметки %s: %w", noteID, err)
	}

	var n note.Note
	if err := doc.DataTo(&n); err != nil {
		return nil, fmt.Errorf("ошибка преобразования заметки %s: %w", noteID, err)
	}
	
	n.ID = doc.Ref.ID
	return &n, nil
}
