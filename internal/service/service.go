package service

import "fmt"

type repository interface {
	GetTranscriptions(string) ([]int64, error)
	SaveTranscriptions(string, int64) error
}

type transcriber interface {
	TranscribeAudio([]byte) (string, error)
}

type downloader interface {
	GetFileURL(fileID string) (string, error)
	DownloadFile(fileURL string) ([]byte, error)
}

type Service struct {
	repository  repository
	transcriber transcriber
	downloader  downloader
}

func New(
	repository repository,
	transcriber transcriber,
	downloader downloader,
) Service {
	return Service{
		repository:  repository,
		transcriber: transcriber,
		downloader:  downloader,
	}
}

func (s Service) TranscribeAndSave(fileID string, messageID int64) error {

	url, err := s.downloader.GetFileURL(fileID)
	if err != nil {
		return fmt.Errorf("faield to get file url:%w", err)
	}

	audio, err := s.downloader.DownloadFile(url)
	if err != nil {
		return fmt.Errorf("faield to download file:%w", err)
	}

	text, err := s.transcriber.TranscribeAudio(audio)
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
