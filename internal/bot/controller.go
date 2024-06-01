package bot

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

type service interface {
	TranscribeAndSave(context.Context, string, int64) error
	FindTranscriptions(context.Context, string) ([]int64, error)
}

type botController struct {
	service service
}

func NewBotController(
	service service,
) *botController {
	return &botController{
		service: service,
	}
}

func (c *botController) listenToAudioAndVideo(b *gotgbot.Bot, ctx *ext.Context) error {
	cont := context.Background()

	msg := ctx.EffectiveMessage

	fileID, err := getFileID(msg)
	if err != nil {
		return nil
	}

	url, err := getFileURL(b, fileID)
	if err != nil {
		return fmt.Errorf("failed to get file url: %w", err)
	}

	if err = c.service.TranscribeAndSave(cont, url, ctx.EffectiveMessage.MessageId); err != nil {
		log.Println("failed to transcribe and save file:", err)
	}

	return nil
}

func (c *botController) findCommand(b *gotgbot.Bot, ctx *ext.Context) error {
	query := ctx.Args()
	matchingIDs, err := c.service.FindTranscriptions(context.TODO(), strings.Join(query, " "))
	if err != nil {
		log.Println("failed to find transcriptions:", err)

		return nil
	}

	var response string
	if len(matchingIDs) == 0 {
		response = "No matching messages("
	} else {
		_, err := b.ForwardMessages(ctx.EffectiveSender.ChatId, ctx.EffectiveSender.ChatId, matchingIDs, nil)
		if err != nil {
			log.Println("failed to forward messages")
		}

		return nil
	}

	_, err = ctx.EffectiveChat.SendMessage(b, response, nil)
	if err != nil {
		log.Println("failed to send message:", err)
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
