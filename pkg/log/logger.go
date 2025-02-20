package log

import (
	"log/slog"
	"os"
)

type Config struct {
	Level slog.Level
}

func New(cfg Config) *slog.Logger {
	opts := slog.HandlerOptions{
		AddSource:   false,
		Level:       cfg.Level,
		ReplaceAttr: nil,
	}
	textHandler := slog.NewTextHandler(os.Stderr, &opts)
	ctxHandler := NewContextHandler(textHandler)

	return slog.New(ctxHandler)
}
