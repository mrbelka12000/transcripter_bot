package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"

	"github.com/PaulSonOfLars/gotgbot/v2"

	"transcripter_bot/internal/bot"
	"transcripter_bot/internal/client/assembly"
	"transcripter_bot/internal/repo"
	"transcripter_bot/internal/service"
	"transcripter_bot/pkg/config"
	"transcripter_bot/pkg/database"
)

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	log.Info("starting project...")
	fmt.Println("suka")
	cfg, err := config.LoadConfig("transcripter")
	if err != nil {
		log.Error("failed to load config", "error", err)
		return
	}

	db, err := database.Connect(cfg)
	if err != nil {
		log.Error("error connecting to database", "error", err)
		return
	}

	botClient, err := gotgbot.NewBot(cfg.TelegramToken, nil)
	if err != nil {
		log.Error("failed to connect to bot", "error", err)
		return
	}
	defer botClient.Close(nil)
	log.Info("telegram bot connection established")

	transcriberService := assembly.NewAssembly(cfg.AssemblyKey)
	repo := repo.New(db, cfg.CollectionName, log)
	service := service.New(repo, transcriberService, log)
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
