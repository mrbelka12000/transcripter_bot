package bot

import (
	"errors"
	"fmt"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

type transcriberService interface {
	TranscribeAndSave(string, int64) error
}

type searchService interface {
	FindTranscriptions([]string) ([]int64, error)
}

type botController struct {
	transcriberService transcriberService
	searchService      searchService
}

func NewBotController(
	transcriberService transcriberService,
	searchService searchService,
) *botController {
	return &botController{
		transcriberService: transcriberService,
		searchService:      searchService,
	}
}

func (c *botController) listenToAudioAndVideo(b *gotgbot.Bot, ctx *ext.Context) error {
	msg := ctx.EffectiveMessage

	fileID, err := getFileID(msg)
	if err != nil {
		return nil
	}

	err = c.transcriberService.TranscribeAndSave(fileID, ctx.EffectiveMessage.MessageId)
	if err != nil {
		return fmt.Errorf("failed to transcrive and save: %w", err)
	}

	return nil
}

func (c *botController) findCommand(b *gotgbot.Bot, ctx *ext.Context) error {
	query := ctx.Args()
	matchingIDs, err := c.searchService.FindTranscriptions(query)
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
	if err == nil {
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
