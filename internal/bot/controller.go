package bot

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/yanzay/tbot/v2"

	"transcripter_bot/internal/client/telegram"
	"transcripter_bot/internal/models"
	srv "transcripter_bot/internal/service"
)

type (
	// Controller ..
	Controller struct {
		client   *tbot.Client
		service  service
		reaction reaction
		log      *slog.Logger
		name     string
	}

	service interface {
		TranscribeAndSave(ctx context.Context, text string, msg models.Message) error
		FindMessages(ctx context.Context, target string, chatID string) ([]int, error)
		GetMessageByMessageID(ctx context.Context, messageID int) (models.Message, error)
	}

	reaction interface {
		SetReaction(reaction telegram.MessageData, emojis ...string) error
	}
)

func New(
	client *tbot.Client,
	service service,
	reaction reaction,
	log *slog.Logger,
	name string,
) *Controller {
	return &Controller{
		client:   client,
		service:  service,
		reaction: reaction,
		log:      log,
		name:     name,
	}
}

func (c *Controller) listenToAudioAndVideo(msg *tbot.Message) {
	if isInline(c.name, msg.Text) {
		c.findCommand(msg)
		return
	}

	fileID := getFileID(msg)
	if fileID == "" {
		return
	}

	url, err := c.getFileURL(fileID)
	if err != nil {
		c.log.With("error", err).Error("get file url")
		return
	}

	ctx := context.Background()

	err = c.service.TranscribeAndSave(ctx, url, models.Message{
		MessageID: msg.MessageID,
		ChatID:    msg.Chat.ID,
	})
	if err != nil {
		c.log.With("error", err).Error("transcribe and save")
		return
	}

	err = c.reaction.SetReaction(telegram.MessageData{
		ChatID:    msg.Chat.ID,
		MessageID: msg.MessageID,
	}, "✍")
	if err != nil {
		c.log.With("error", err).Error("set reaction")
		return
	}

	c.log.With("fileID", fileID).Info("successfully saved")
}

func (c *Controller) findCommand(msg *tbot.Message) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	matchingIDs, err := c.service.FindMessages(ctx, msg.Text, msg.Chat.ID)
	if err != nil {
		if errors.Is(err, srv.ErrEmptyTarget) {
			c.client.SendMessage(msg.Chat.ID, "Please specify what to find")
			return
		}

		c.client.SendMessage(msg.Chat.ID, "Something went wrong")
		c.log.With("error", err).Error("find messages")
		return
	}

	if len(matchingIDs) == 0 {
		c.client.SendMessage(msg.Chat.ID, "No matching messages(")
		return
	}

	for _, id := range matchingIDs {
		c.client.SendMessage(msg.Chat.ID, "Found", tbot.OptReplyToMessageID(id))
	}
}

func (c *Controller) textCommand(msg *tbot.Message) {
	if msg.ReplyToMessage == nil {
		c.log.Debug("no reply in text command")
		return
	}

	message, err := c.service.GetMessageByMessageID(context.Background(), msg.ReplyToMessage.MessageID)
	if err != nil {
		c.log.With("error", err).Error("get reply to message")
		return
	}

	c.client.SendMessage(msg.Chat.ID, message.Text, tbot.OptReplyToMessageID(msg.ReplyToMessage.MessageID))
}

func (c *Controller) ping(msg *tbot.Message) {
	c.client.SendMessage(msg.Chat.ID, "pong")
}

func (c *Controller) getFileURL(fileID string) (string, error) {
	file, err := c.client.GetFile(fileID)
	if err != nil {
		return "", fmt.Errorf("can't get file: %w", err)
	}

	return c.client.FileURL(file), nil
}

func getFileID(msg *tbot.Message) string {
	switch {
	case msg.Audio != nil:
		return msg.Audio.FileID
	case msg.Voice != nil:
		return msg.Voice.FileID
	case msg.Video != nil:
		return msg.Video.FileID
	case msg.VideoNote != nil:
		return msg.VideoNote.FileID
	default:
		return ""
	}
}

func isInline(botName, text string) bool {
	words := strings.Split(text, " ")
	if len(words) == 0 {
		return false
	}

	return botName == words[0]
}
