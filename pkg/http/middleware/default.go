package middleware

import (
	"log/slog"
	"net/http"
)

type Skipper func(r *http.Request) bool

type Config struct {
	Logger  *slog.Logger
	Skipper Skipper
}

func ApplyDefault(cfg Config, middlewares ...func(http.Handler) http.Handler) []func(http.Handler) http.Handler {
	defaultMiddlewares := []func(http.Handler) http.Handler{
		ContextLogger(cfg.Logger),
		RequestID,
		LogRequestID,
		LogRequest(cfg.Skipper),
		Recover,
	}

	newMiddlewares := make([]func(http.Handler) http.Handler,
		len(defaultMiddlewares), len(defaultMiddlewares)+len(middlewares))
	copy(newMiddlewares, defaultMiddlewares)
	newMiddlewares = append(newMiddlewares, middlewares...)

	return newMiddlewares
}
