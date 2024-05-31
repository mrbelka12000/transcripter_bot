package service

import "fmt"

type repository interface {
	GetTranscriptions(string) ([]int64, error)
	SaveTranscriptions(string, int64) error
}

type transcriber interface {
	TranscribeAudio(string) (string, error)
}

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

func (s Service) TranscribeAndSave(fileURL string, messageID int64) error {
	text, err := s.transcriber.TranscribeAudio(fileURL)
	if err != nil {
		return fmt.Errorf("faield transcribe audio:%w", err)
	}

	if err := s.repository.SaveTranscriptions(text, messageID); err != nil {
		return fmt.Errorf("failed to save text: %w", err)
	}

	return nil
}

func (s Service) FindTranscriptions(target string) ([]int64, error) {
	messageIDs, err := s.repository.GetTranscriptions(target)
	if err != nil {
		return nil, fmt.Errorf("failed to find: %w", err)
	}

	return messageIDs, err
}
