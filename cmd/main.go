package main

import (
	"log/slog"
	"os"
	"os/signal"

	"transcripter_bot/internal/bot"
	"transcripter_bot/internal/client/assembly"
	"transcripter_bot/internal/repo"
	"transcripter_bot/internal/service"
	"transcripter_bot/pkg/config"
	"transcripter_bot/pkg/database"
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

	db, err := database.Connect(cfg)
	if err != nil {
		log.Error("error connecting to database", err)
		return
	}

	botClient, err := gotgbot.NewBot(cfg.TelegramToken, nil)
	if err != nil {
		log.Error("failed to connect to bot", err)
		return
	}
	defer botClient.Close(nil)
	log.Info("telegram bot connection established")

	transcriberService := assembly.NewAssembly(cfg.AssemblyKey)
	repo := repo.New(db, cfg.CollectionName)
	service := service.New(repo, transcriberService)
	botController := bot.NewBotController(service, log)

	if err := bot.RunTelegramBot(botClient, botController, log); err != nil {
		log.Error("failed to run the project", err)
		return
	}

	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt)

	<-ch

	log.Info("...project is shuting down")
}
