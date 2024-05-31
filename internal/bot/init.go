package bot

import (
	"time"

	"transcripter_bot/pkg/logger"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

func RunTelegramBot(bot *gotgbot.Bot, botController *botController) error {
	dispatcher := ext.NewDispatcher(&ext.DispatcherOpts{
		Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
			logger.Log.Error("an error occurred while handling update", err)
			return ext.DispatcherActionNoop
		},
		MaxRoutines: ext.DefaultMaxRoutines,
	})

	filter := func(msg *gotgbot.Message) bool {

		if msg.Audio != nil || msg.Voice != nil || msg.VideoNote != nil {
			return true
		}

		return false
	}

	dispatcher.AddHandler(handlers.NewCommand("find", botController.findCommand))
	dispatcher.AddHandler(handlers.NewCommand("ping", botController.ping))
	dispatcher.AddHandler(handlers.NewMessage(filter, botController.listenToAudioAndVideo))

	updater := ext.NewUpdater(dispatcher, nil)

	err := updater.StartPolling(bot, &ext.PollingOpts{
		DropPendingUpdates: true,
		GetUpdatesOpts: &gotgbot.GetUpdatesOpts{
			Timeout: 9,
			RequestOpts: &gotgbot.RequestOpts{
				Timeout: time.Second * 10,
			},
		},
	})
	if err != nil {
		return err
	}
	logger.Log.Info("bot has been started", "username", bot.User.Username)

	go updater.Idle()

	return nil
}
