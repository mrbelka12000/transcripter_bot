package service

import (
	"context"

	"transcripter_bot/internal/models"
)

type repository interface {
	GetMessagesForFind(ctx context.Context, target string, chatID string) ([]int, error)
	SaveMessage(ctx context.Context, msg models.Message) error
	GetMessageByMessageID(ctx context.Context, messageID int) (models.Message, error)
}

type transcriber interface {
	TranscribeAudio(context.Context, string) (string, error)
}
