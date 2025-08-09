package store

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	"github.com/Denol007/telegram-personal-assistant-go/internal/note" // Импортируем наш пакет note
)

var client *firestore.Client


func Init() {
	ctx := context.Background()
	projectID := os.Getenv("GCP_PROJECT_ID")
	if projectID == "" {
		log.Printf("GCP_PROJECT_ID environment variable must be set.")
		return
	}
	var err error
	client, err = firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Printf("Failed to create Firestore client: %v", err)
		return
	}
}


func SaveNote(ctx context.Context, n note.Note) error {
	if client == nil {
		return nil // no-op if Firestore not initialized; caller may log success in dev
	}
	_, _, err := client.Collection("notes").Add(ctx, n)
	return err
}