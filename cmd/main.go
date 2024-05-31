package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"transcripter_bot/internal/bot"
	"transcripter_bot/pkg/config"
	"transcripter_bot/pkg/logger"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"golang.org/x/sync/errgroup"
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

	// TODO: implement searchService and transcriberService
	botController := bot.NewBotController(nil, nil)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create a channel to listen for OS signals.
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	errGroup, ctx := errgroup.WithContext(ctx)

	errGroup.Go(func() error {
		select {
		case sig := <-sigChan:
			logger.Log.Info("received shutdown signal", "signal", sig)
			cancel()
		case <-ctx.Done():
		}
		return nil
	})

	errGroup.Go(func() error {
		if err := bot.RunTelegramBot(botClient, botController); err != nil {
			logger.Log.Error("failed to run the bot", err)
			return fmt.Errorf("failed to run the bot: %w", err)
		}
		return nil
	})

	errGroup.Go(func() error {
		if err := runDiagnostics(cfg.PromtPort); err != nil {
			logger.Log.Error("failed to run diagnostics", err)
			return fmt.Errorf("failed to run diagnostics: %w", err)
		}
		return nil
	})

	if err := errGroup.Wait(); err != nil {
		logger.Log.Error("failed to start services", err)
	}

	logger.Log.Info("...project is shutting down")
}

func runDiagnostics(port string) error {
	if err := http.ListenAndServe(port); err != nil {
		return fmt.Errorf("")
	}
	return nil
}
