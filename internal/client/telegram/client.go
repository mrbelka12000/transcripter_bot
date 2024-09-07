package telegram

import (
	"fmt"
)

const (
	apiBaseURL = "https://api.telegram.org"

	eventSetReaction = "setMessageReaction"

	reactionNoted = "✍"
)

type (
	Client struct {
		url string
	}
)

func NewClient(token string) *Client {
	return &Client{
		url: fmt.Sprintf("%s/bot%s/", apiBaseURL, token) + "%s",
	}
}
