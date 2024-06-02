package service

import (
	"context"
	"transcripter_bot/internal/models"
)

type repository interface {
	GetMessages(context.Context, string, int64) ([]int64, error)
	SaveMessage(context.Context, models.Message) error
}

type transcriber interface {
	TranscribeAudio(context.Context, string) (string, error)
}
