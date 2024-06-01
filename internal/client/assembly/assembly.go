package assembly

import (
	"context"

	"github.com/AssemblyAI/assemblyai-go-sdk"
)

type assembly struct {
	client *assemblyai.Client
}

func NewAssembly(client *assemblyai.Client) assembly {
	return assembly{
		client: client,
	}
}

func (s assembly) TranscribeAudio(ctx context.Context, url string) (string, error) {
	params := assemblyai.TranscriptOptionalParams{
		LanguageCode: "ru",
	}

	transcript, err := s.client.Transcripts.TranscribeFromURL(ctx, url, &params)
	if err != nil {
		return "", err
	}

	return assemblyai.ToString(transcript.Text), nil
}
