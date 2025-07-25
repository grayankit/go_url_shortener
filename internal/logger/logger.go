package logger

import (
	"log/slog"
	"os"
)

func NewLogger() *slog.Logger {
	replaceAttr := func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.TimeKey || a.Key == slog.MessageKey {
			return slog.Attr{}
		}
		return a
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		ReplaceAttr: replaceAttr,
	}))

	return logger
}
