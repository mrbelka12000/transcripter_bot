package service

import (
	"context"
	"fmt"
)

type Service struct {
	repository  repository
	transcriber transcriber
}

func New(
	repository repository,
	transcriber transcriber,
) Service {
	return Service{
		repository:  repository,
		transcriber: transcriber,
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
