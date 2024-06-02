package mock_assembly

import (
	"context"
	"math/rand"
	"time"
)

type Mock struct{}

func New() *Mock {
	return &Mock{}
}

func (s *Mock) TranscribeAudio(ctx context.Context, url string) (string, error) {
	wordList := []string{
		"apple", "banana", "cherry", "date", "elderberry",
		"fig", "grape", "honeydew", "kiwi", "lemon",
		"mango", "nectarine", "orange", "papaya", "quince",
		"raspberry", "strawberry", "tangerine", "ugli", "violet",
		"watermelon", "xigua", "yellow", "zucchini",
	}
	return randomWord(wordList), nil
}

func randomWord(wordList []string) string {
	rand.Seed(time.Now().UnixNano())
	return wordList[rand.Intn(len(wordList))]
}
