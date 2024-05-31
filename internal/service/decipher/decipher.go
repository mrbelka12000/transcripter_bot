package decipher

import (
	"context"
	"transcripter_bot/pkg/assembly"
)

type TranscriberService interface {
	TranscribeAudio(context.Context, string) (string, error)
}

type ServiceImpl struct {
	Assembly assembly.Assembly
}

func NewTranscribeService(transciberApiKey string) TranscriberService {
	return &ServiceImpl{
		Assembly: assembly.NewAssembly(transciberApiKey),
	}
}

func (s *ServiceImpl) TranscribeAudio(ctx context.Context, url string) (string, error) {
	return s.Assembly.TranscribeFromUrl(ctx, url)
}
