package assembly

import (
	"bytes"
	"context"

	aai "github.com/AssemblyAI/assemblyai-go-sdk"
)

type ServiceImpl struct {
	Client *aai.Client
}

func NewAssembly(client *aai.Client) *ServiceImpl {
	return &ServiceImpl{
		Client: client,
	}
}

func (s *ServiceImpl) TranscribeFromUrl(ctx context.Context, audioUrl string) (string, error) {
	params := aai.TranscriptOptionalParams{
		LanguageCode: "ru",
	}

	transcript, err := s.Client.Transcripts.TranscribeFromURL(ctx, audioUrl, &params)
	if err != nil {
		return "", err
	}

	return aai.ToString(transcript.Text), nil
}

func (s *ServiceImpl) TranscribeFromBytes(ctx context.Context, info []byte) (string, error) {
	r := bytes.NewReader(info)

	transcript, err := s.Client.Transcripts.TranscribeFromReader(ctx, r, nil)
	if err != nil {
		return "", err
	}

	return aai.ToString(transcript.Text), nil
}
