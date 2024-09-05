package bot

import (
	"log/slog"

	"github.com/yanzay/tbot/v2"
)

func RunTelegramBot(bot *tbot.Server, botController *Controller, log *slog.Logger) error {

	bot.HandleMessage("/ping", botController.ping)
	bot.HandleMessage("/find", botController.findCommand)
	bot.HandleMessage(".*", botController.listenToAudioAndVideo)

	return bot.Start()
}
