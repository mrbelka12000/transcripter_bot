package bot

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

type service interface {
	TranscribeAndSave(string, int64) error
	FindTranscriptions([]string) ([]int64, error)
}

type botController struct {
	service service
	l       *slog.Logger
}

func NewBotController(
	service service,
	l *slog.Logger,
) *botController {
	return &botController{
		service: service,
		l:       l,
	}
}

func (c *botController) listenToAudioAndVideo(b *gotgbot.Bot, ctx *ext.Context) error {
	_, err := ctx.EffectiveChat.SendMessage(b, "implement me", nil)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil

	msg := ctx.EffectiveMessage

	fileID, err := getFileID(msg)
	if err != nil {
		return nil
	}

	err = c.service.TranscribeAndSave(fileID, ctx.EffectiveMessage.MessageId)
	if err != nil {
		return fmt.Errorf("failed to transcrive and save: %w", err)
	}

	return nil
}

func (c *botController) findCommand(b *gotgbot.Bot, ctx *ext.Context) error {

	_, err := ctx.EffectiveChat.SendMessage(b, "implement me", nil)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil

	query := ctx.Args()
	matchingIDs, err := c.service.FindTranscriptions(query)
	if err != nil {
		return fmt.Errorf("failed to find transcriptions: %w", err)
	}

	var response string
	if len(matchingIDs) == 0 {
		response = "No matching messages("
	} else {
		_, err := b.ForwardMessages(
			ctx.EffectiveSender.ChatId,
			ctx.EffectiveSender.ChatId,
			matchingIDs,
			nil,
		)
		if err != nil {
			return fmt.Errorf("failed to forward messages: %w", err)
		}

		return nil
	}

	_, err = ctx.EffectiveChat.SendMessage(b, response, nil)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}

func (c *botController) ping(b *gotgbot.Bot, ctx *ext.Context) error {
	_, err := ctx.EffectiveMessage.Reply(b, "pong", nil)

	return err
}

func getFileID(msg *gotgbot.Message) (string, error) {
	var fileID string

	if msg.Audio != nil {
		fileID = msg.Audio.FileId
	} else if msg.Voice != nil {
		fileID = msg.Voice.FileId
	} else if msg.VideoNote != nil {
		fileID = msg.VideoNote.FileId
	} else {
		return "", errors.New("message is not audio or video type")
	}

	return fileID, nil
}
