package telegram

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type (
	MessageData struct {
		ChatID    string `json:"chat_id"`
		MessageID int    `json:"message_id"`
	}

	reaction struct {
		Type  string `json:"type"`
		Emoji string `json:"emoji"`
	}
)

func (c *Client) SetReaction(data MessageData, emojis ...string) error {
	if len(emojis) == 0 {
		return errors.New("no emojis provided")
	}

	var reactions []reaction
	for _, emoji := range emojis {
		reactions = append(reactions, reaction{
			Type:  "emoji",
			Emoji: emoji,
		})
	}

	reactionsData, _ := json.Marshal(reactions)
	req := url.Values{}
	req.Set("chat_id", data.ChatID)
	req.Set("message_id", fmt.Sprint(data.MessageID))
	req.Set("is_big", strconv.FormatBool(false))
	req.Set("reaction", string(reactionsData))

	setReactionURL := fmt.Sprintf(c.url, eventSetReaction)
	resp, err := http.PostForm(setReactionURL, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("invalid status code: %d", resp.StatusCode)
	}

	return nil
}
