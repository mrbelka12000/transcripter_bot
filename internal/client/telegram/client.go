package telegram

import (
	"fmt"
	"net/http"
	"time"
)

const (
	apiBaseURL = "https://api.telegram.org"

	eventSetReaction = "setMessageReaction"

	reactionNoted = "✍"
)

type (
	Client struct {
		hc  *http.Client
		url string
	}
)

func NewClient(token string) *Client {
	return &Client{
		hc: &http.Client{
			Timeout: 10 * time.Second,
		},
		url: fmt.Sprintf("%s/bot%s/", apiBaseURL, token) + "%s",
	}
}
