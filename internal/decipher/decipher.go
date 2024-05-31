package decipher

import (
	"context"
	"transcripter_bot/pkg/assembly"

	aai "github.com/AssemblyAI/assemblyai-go-sdk"
)

type TranscriberService interface {
	TranscribeAudio(context.Context, string) (string, error)
}

type ServiceImpl struct {
	Assembly *assembly.ServiceImpl
}

func NewTranscribeService(cli *aai.Client) TranscriberService {
	return &ServiceImpl{
		Assembly: assembly.NewAssembly(cli),
	}
}

func (s *ServiceImpl) TranscribeAudio(ctx context.Context, url string) (string, error) {
	return s.Assembly.TranscribeFromUrl(ctx, url)
}
