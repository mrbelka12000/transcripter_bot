package bot

import (
	"log"

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

	var fileID string

	// TODO: do we need to consider videos ?
	if msg.Audio != nil {
		fileID = msg.Audio.FileId
	} else if msg.Voice != nil {
		fileID = msg.Voice.FileId
	} else if msg.VideoNote != nil {
		fileID = msg.VideoNote.FileId
	} else {
		return nil
	}

	err := c.transcriberService.TranscribeAndSave(fileID, ctx.EffectiveMessage.MessageId)
	if err != nil {
		log.Println("failed to transcribe and save file:", err)
	}

	return nil
}

func (c *botController) findCommand(b *gotgbot.Bot, ctx *ext.Context) error {
	query := ctx.Args()
	matchingIDs, err := c.searchService.FindTranscriptions(query)
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
