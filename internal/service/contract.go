package service

import "context"

type repository interface {
	GetTranscriptions(context.Context, string) ([]int64, error)
	SaveTranscriptions(context.Context, string, int64) error
}

type transcriber interface {
	TranscribeAudio(context.Context, string) (string, error)
}
