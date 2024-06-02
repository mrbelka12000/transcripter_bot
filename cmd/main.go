package main

import (
	"log/slog"
	"os"
	"os/signal"

	"transcripter_bot/internal/bot"
	"transcripter_bot/pkg/config"
	timeformat "transcripter_bot/pkg/time-format"

	"github.com/PaulSonOfLars/gotgbot/v2"
)

func main() {
	log := slog.New(&timeformat.CustomHandler{})

	log.Info("starting project...")

	cfg, err := config.LoadConfig("transcripter")
	if err != nil {
		log.Error("failed to load config", err)
		return
	}
	log.Info("config loaded")

	botClient, err := gotgbot.NewBot(cfg.TelegramToken, nil)
	if err != nil {
		log.Error("failed to connect to bot", err)
		return
	}
	defer botClient.Close(nil)
	log.Info("telegram bot connection established")

	// TODO: implement searchService and transriberService
	botController := bot.NewBotController(nil, nil)

	if err := bot.RunTelegramBot(botClient, botController, log); err != nil {
		log.Error("failed to run the project", err)
		return
	}

	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt)

	<-ch

	log.Info("...project is shuting down")
}
