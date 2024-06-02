package main

import (
	"log"

	"transcripter_bot/internal/bot"
	"transcripter_bot/internal/client/mock"
	"transcripter_bot/internal/repo"
	"transcripter_bot/internal/service"
	"transcripter_bot/pkg/config"
	"transcripter_bot/pkg/database"

	"github.com/PaulSonOfLars/gotgbot/v2"
)

func main() {

	log.Println("Starting a project...")

	cfg, err := config.LoadConfig("transcripter")
	if err != nil {
		log.Println(err)
		return
	}

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
		return
	}

	botClient, err := gotgbot.NewBot(cfg.TelegramToken, nil)
	if err != nil {
		log.Printf("failed to connect to bot: %v", err)
		return
	}

	transcriberService := mock.New()

	repo := repo.New(db, cfg.CollectionName)

	service := service.New(repo, transcriberService)

	botController := bot.NewBotController(service)

	if err := bot.RunTelegramBot(botClient, botController); err != nil {
		log.Printf("failed to run the project: %v", err)
		return
	}
}
