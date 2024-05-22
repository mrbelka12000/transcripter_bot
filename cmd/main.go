package main

import (
	"fmt"
	"log"

	"transcripter_bot/internal/bot"
	"transcripter_bot/pkg/config"

	"github.com/PaulSonOfLars/gotgbot/v2"
)

func main() {
	fmt.Println("Starting new project")

	cfg, err := config.LoadConfig("transcripter")
	if err != nil {
		log.Println(err)

		return
	}

	botClient, err := gotgbot.NewBot(cfg.TelegramToken, nil)
	if err != nil {
		log.Printf("failed to connect to bot: %v", err)

		return
	}

	_ = bot.NewFileDownloader(botClient)

	botController := bot.NewBotController(nil, nil)

	if err := bot.RunTelegramBot(botClient, botController); err != nil {
		log.Printf("failed to run the project: %v", err)

		return
	}
}
