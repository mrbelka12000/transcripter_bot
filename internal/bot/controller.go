package bot

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"transcripter_bot/internal/models"
	srv "transcripter_bot/internal/service"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

type service interface {
	TranscribeAndSave(context.Context, string, models.Message) error
	FindMessages(context.Context, string, int64) ([]int64, error)
}

type botController struct {
	service service
	log     *slog.Logger
}

func NewBotController(
	service service,
	log *slog.Logger,
) *botController {
	return &botController{
		service: service,
		log:     log,
	}
}

func (c *botController) listenToAudioAndVideo(b *gotgbot.Bot, ctx *ext.Context) error {
	fileID, err := getFileID(ctx.EffectiveMessage)
	if err != nil {
		return nil
	}

	url, err := getFileURL(b, fileID)
	if err != nil {
		return fmt.Errorf("failed to get file url: %w", err)
	}

	message := models.Message{
		MessageID: ctx.EffectiveMessage.MessageId,
		ChatID:    ctx.EffectiveChat.Id,
	}

	if err = c.service.TranscribeAndSave(context.TODO(), url, message); err != nil {
		return fmt.Errorf("failed to transcribe and save: %w", err)
	}

	return nil
}

func (c *botController) findCommand(b *gotgbot.Bot, ctx *ext.Context) error {
	query := ctx.Args()

	matchingIDs, err := c.service.FindMessages(context.TODO(), strings.Join(query[1:], " "), ctx.EffectiveSender.ChatId)
	if err != nil {
		if errors.Is(err, srv.ErrEmptyTarget) {
			b.SendMessage(ctx.EffectiveSender.ChatId, "Please specify what to find", nil)
		}
		return fmt.Errorf("failed to find transcriptions: %w", err)
	}

	var response string
	if len(matchingIDs) == 0 {
		response = "No matching messages("
	} else {
		for _, id := range matchingIDs {
			_, err = b.SendMessage(ctx.EffectiveSender.ChatId, "found", &gotgbot.SendMessageOpts{
				ReplyParameters: &gotgbot.ReplyParameters{
					MessageId: id,
				},
			})
			if err != nil {
				return fmt.Errorf("failed to reply to message: %w", err)
			}
		}

		return nil
	}

	_, err = ctx.EffectiveChat.SendMessage(b, response, &gotgbot.SendMessageOpts{})
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

func getFileURL(b *gotgbot.Bot, fileID string) (string, error) {
	file, err := b.GetFile(fileID, nil)
	if err != nil {
		return "", fmt.Errorf("failed to get file: %w", err)
	}

	return file.URL(b, nil), nil
}
