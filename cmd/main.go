package main

import (
	"log/slog"
	"os"
	"os/signal"

	"transcripter_bot/internal/bot"
	"transcripter_bot/pkg/config"
	"transcripter_bot/pkg/logger"

	"github.com/PaulSonOfLars/gotgbot/v2"
)

func main() {
	logger.Init(slog.NewTextHandler(os.Stdout, nil))

	logger.Log.Info("starting project...")

	cfg, err := config.LoadConfig("transcripter")
	if err != nil {
		logger.Log.Error("failed to load config", err)
		return
	}
	logger.Log.Info("config loaded")

	botClient, err := gotgbot.NewBot(cfg.TelegramToken, nil)
	if err != nil {
		logger.Log.Error("failed to connect to bot", err)
		return
	}
	defer botClient.Close(nil)
	logger.Log.Info("telegram bot connection established")

	// TODO: implement searchService and transriberService
	botController := bot.NewBotController(nil, nil)

	if err := bot.RunTelegramBot(botClient, botController); err != nil {
		logger.Log.Error("failed to run the project", err)
		return
	}

	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt)

	<-ch

	logger.Log.Info("...project is shuting down")
}
