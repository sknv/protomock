package container

import (
	"log/slog"

	"github.com/sknv/protomock/pkg/log"
	"github.com/sknv/protomock/pkg/option"
)

func (a *Application) RegisterLogger(cfg log.Config) *slog.Logger {
	logger := log.New(cfg)
	a.logger = option.Some(logger)

	return logger
}

func (a *Application) Logger() option.Option[*slog.Logger] {
	return a.logger
}
