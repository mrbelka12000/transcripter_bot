package bot

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

type transcriberService interface {
	TranscribeAndSave(string, int64) error
}

type searchService interface {
	FindTranscriptions(string) ([]int64, error)
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

	url, err := getFileURL(b, fileID)
	if err != nil {
		return fmt.Errorf("failed to get file url: %w", err)
	}

	err = c.transcriberService.TranscribeAndSave(url, ctx.EffectiveMessage.MessageId)
	if err != nil {
		log.Println("failed to transcribe and save file:", err)
	}

	return nil
}

func (c *botController) findCommand(b *gotgbot.Bot, ctx *ext.Context) error {
	query := ctx.Args()
	matchingIDs, err := c.searchService.FindTranscriptions(strings.Join(query, " "))
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
