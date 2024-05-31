package converter

import "context"

type Converter struct {
}

func (c Converter) TranscribeAudio(ctx context.Context, fileURL string) (string, error) {
	return "", nil
}
