package logger

import (
	"log/slog"
)

type log struct {
	log *slog.Logger
}

var Log log

func Init(h slog.Handler) {
	Log.log = slog.New(h)
}

func (l *log) Error(msg string, err error, keysAndValues ...interface{}) {
	l.log.Error(msg, append(keysAndValues, "error", err)...)
}

func (l *log) Info(msg string, keysAndValues ...interface{}) {
	l.log.Info(msg, keysAndValues...)
}

func (l *log) Debug(msg string, keysAndValues ...interface{}) {
	l.log.Debug(msg, keysAndValues...)
}
