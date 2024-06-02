package timeformat

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"
)

type CustomHandler struct{}

func (h *CustomHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return true
}

func (h *CustomHandler) Handle(ctx context.Context, record slog.Record) error {
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	logMessage := fmt.Sprintf("%s [%s] %s\n", timestamp, record.Level.String(), record.Message)

	_, err := os.Stdout.WriteString(logMessage)
	return err
}

func (h *CustomHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

func (h *CustomHandler) WithGroup(name string) slog.Handler {
	return h
}
