package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"

	"github.com/PaulSonOfLars/gotgbot/v2"

	"transcripter_bot/internal/bot"
	"transcripter_bot/internal/client/assembly"
	"transcripter_bot/internal/repository"
	"transcripter_bot/internal/service"
	"transcripter_bot/pkg/config"
	"transcripter_bot/pkg/database"
)

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	log.Info("starting project...")

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

	transcriberService := assembly.NewAssembly(cfg.AssemblyKey)
	repo := repository.New(db, cfg.CollectionName, log)
	svc := service.New(repo, transcriberService, log)
	botController := bot.New(svc, log, cfg.BotName)

	go func() {
		//health check
		http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("ok"))
		})
		err := http.ListenAndServe(fmt.Sprintf(":%s", cfg.HTTPPort), nil)
		if err != nil {
			log.Error("error starting http server", "error", err)
			os.Exit(1)
		}
	}()

	if err := bot.RunTelegramBot(botClient, botController, log); err != nil {
		log.Error("failed to run the project", err)
		return
	}

	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt)

	<-ch

	log.Info("...project is shuting down")
}
