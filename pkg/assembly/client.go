package assembly

import (
	aai "github.com/AssemblyAI/assemblyai-go-sdk"
)

// TODO : how can i get error from client ???
func NewClient(apikey string) *aai.Client {
	client := aai.NewClient(apikey)

	return client
}
