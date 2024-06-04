package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"transcripter_bot/internal/models"
)

var (
	ErrEmptyTarget = errors.New("empty target")
)

type Service struct {
	repository  repository
	transcriber transcriber
	log         *slog.Logger
}

func New(
	repository repository,
	transcriber transcriber,
	log *slog.Logger,
) Service {
	return Service{
		repository:  repository,
		transcriber: transcriber,
		log:         log,
	}
}

func (s Service) TranscribeAndSave(ctx context.Context, fileURL string, message models.Message) error {
	text, err := s.transcriber.TranscribeAudio(ctx, fileURL)
	if err != nil {
		return fmt.Errorf("failed to transcribe audio:%w", err)
	}

	message.Text = text

	if err := s.repository.SaveMessage(ctx, message); err != nil {
		return fmt.Errorf("failed to save text: %w", err)
	}

	return nil
}

func (s Service) FindMessages(ctx context.Context, target string, chatID int64) ([]int64, error) {
	if isEmptyTarget(target) {
		return nil, fmt.Errorf("invalid target: %w", ErrEmptyTarget)
	}

	messages, err := s.repository.GetMessages(ctx, target, chatID)
	if err != nil {
		return nil, fmt.Errorf("failed to find: %w", err)
	}

	return messages, err
}

func isEmptyTarget(s string) bool {
	s = strings.TrimSpace(s)
	return len(s) == 0
}
