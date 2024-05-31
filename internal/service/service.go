package service

import (
	"context"
	"fmt"
	"transcripter_bot/internal/service/decipher"
)

type repository interface {
	GetTranscriptions(context.Context, string) ([]int64, error)
	SaveTranscriptions(context.Context, string, int64) error
}

type transcriber interface {
	TranscribeAudio(context.Context, string) (string, error)
}

type Service struct {
	repository  repository
	transcriber transcriber
}

func New(
	repository repository,
	assemblyApiKey string,
) Service {
	return Service{
		repository:  repository,
		transcriber: decipher.NewTranscribeService(assemblyApiKey),
	}
}

func (s Service) TranscribeAndSave(ctx context.Context, fileURL string, messageID int64) error {
	text, err := s.transcriber.TranscribeAudio(ctx, fileURL)
	if err != nil {
		return fmt.Errorf("faield transcribe audio:%w", err)
	}

	fmt.Println(text)

	// if err := s.repository.SaveTranscriptions(ctx, text, messageID); err != nil {
	// 	return fmt.Errorf("failed to save text: %w", err)
	// }

	return nil
}

func (s Service) FindTranscriptions(ctx context.Context, target string) ([]int64, error) {
	messageIDs, err := s.repository.GetTranscriptions(ctx, target)
	if err != nil {
		return nil, fmt.Errorf("failed to find: %w", err)
	}

	return messageIDs, err
}
