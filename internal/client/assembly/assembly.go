package assembly

import (
	"context"

	"github.com/AssemblyAI/assemblyai-go-sdk"
)

type Assembly struct {
	client *assemblyai.Client
}

func NewAssembly(apiKey string) *Assembly {
	return &Assembly{
		client: assemblyai.NewClient(apiKey),
	}
}

func (s *Assembly) TranscribeAudio(ctx context.Context, url string) (string, error) {
	params := assemblyai.TranscriptOptionalParams{
		LanguageCode: "ru",
	}

	transcript, err := s.client.Transcripts.TranscribeFromURL(ctx, url, &params)
	if err != nil {
		return "", err
	}

	return assemblyai.ToString(transcript.Text), nil
}
