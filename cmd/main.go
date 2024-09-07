package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/yanzay/tbot/v2"

	"transcripter_bot/internal/bot"
	"transcripter_bot/internal/client/assembly"
	"transcripter_bot/internal/client/telegram"
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

	transcriberService := assembly.NewAssembly(cfg.AssemblyKey)
	repo := repository.New(db, cfg.TableName)
	svc := service.New(repo, transcriberService, log, cfg.BotName)
	telBot := tbot.New(cfg.TelegramToken)
	botController := bot.New(telBot.Client(), svc, telegram.NewClient(cfg.TelegramToken), log, cfg.BotName)

	go func() {
		//health check
		http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("ok"))
		})

		// metrics
		http.Handle("/metrics", promhttp.Handler())

		err := http.ListenAndServe(fmt.Sprintf(":%s", cfg.HTTPPort), nil)
		if err != nil {
			log.Error("error starting http server", "error", err)
			os.Exit(1)
		}
	}()

	if err := bot.RunTelegramBot(telBot, botController, log); err != nil {
		log.Error("failed to run the project", err)
		return
	}

	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt)

	<-ch

	log.Info("...project is shuting down")
}
