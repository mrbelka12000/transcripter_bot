package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type (
	MessageData struct {
		ChatID    string `json:"chat_id"`
		MessageID int    `json:"message_id"`
	}
)

type (
	apiReq struct {
		MessageId int        `json:"message_id"`
		ChatId    string     `json:"chat_id"`
		Reaction  []reaction `json:"reaction"`
	}
	reaction struct {
		Type  string `json:"type"`
		Emoji string `json:"emoji"`
	}
)

func (c *Client) SetReaction(ctx context.Context, data MessageData) error {
	apiReq := apiReq{
		MessageId: data.MessageID,
		ChatId:    data.ChatID,
		Reaction: []reaction{
			{
				Type:  "emoji",
				Emoji: reactionNoted,
			},
		},
	}

	url := fmt.Sprintf(c.url, eventSetReaction)

	jsonData, err := json.Marshal(apiReq)
	if err != nil {
		return fmt.Errorf("failed to marshal reaction: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.hc.Do(req)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("invalid status code: %d", resp.StatusCode)
	}

	return nil
}
