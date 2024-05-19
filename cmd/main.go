package main

import (
	"fmt"
	"log"
	"os"
	"transcripter_bot/internal/bot"

	"github.com/PaulSonOfLars/gotgbot/v2"
)

func main() {
	fmt.Println("Starting new project")

	token := os.Getenv("TELEGRAM_BOT_TOKEN")

	botClient, err := gotgbot.NewBot(token, nil)
	if err != nil {
		log.Printf("failed to connect to bot: %v", err)

		return
	}

	// TODO: implement transcriberService, searchService
	botController := bot.NewBotController(nil, nil)

	if err := bot.RunTelegramBot(botClient, botController); err != nil {
		log.Printf("failed to run the project: %v", err)

		return
	}
}
