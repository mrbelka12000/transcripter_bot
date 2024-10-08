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
	botName     string
}

func New(
	repository repository,
	transcriber transcriber,
	log *slog.Logger,
	botName string,
) Service {
	return Service{
		repository:  repository,
		transcriber: transcriber,
		log:         log,
		botName:     botName,
	}
}

// TranscribeAndSave ..
func (s Service) TranscribeAndSave(ctx context.Context, fileURL string, message models.Message) error {
	text, err := s.transcriber.TranscribeAudio(ctx, fileURL)
	if err != nil {
		return fmt.Errorf("failed to transcribe audio:%w", err)
	}

	if text == "" {
		return fmt.Errorf("empty audio text")
	}

	message.Text = text

	if err := s.repository.SaveMessage(ctx, message); err != nil {
		return fmt.Errorf("failed to save text: %w", err)
	}

	return nil
}

// FindMessages ..
func (s Service) FindMessages(ctx context.Context, target, chatID string) ([]int, error) {
	if isEmptyTarget(target) {
		return nil, fmt.Errorf("invalid target: %w", ErrEmptyTarget)
	}

	query := strings.Split(target, " ")
	if len(query) < 2 {
		return nil, fmt.Errorf("invalid target: %w", ErrEmptyTarget)
	}
	if query[0] == "/find" || query[0] == s.botName {
		target = strings.Join(query[1:], " ")
	}

	messages, err := s.repository.GetMessagesForFind(ctx, target, chatID)
	if err != nil {
		return nil, fmt.Errorf("failed to find: %w", err)
	}

	return messages, err
}

// GetMessageByMessageID ..
func (s Service) GetMessageByMessageID(ctx context.Context, messageID int) (models.Message, error) {
	return s.repository.GetMessageByMessageID(ctx, messageID)
}

func isEmptyTarget(s string) bool {
	s = strings.TrimSpace(s)
	return len(s) == 0
}
