package middleware

import (
	"log/slog"
	"net/http"

	"github.com/uptrace/bunrouter"

	"github.com/sknv/protomock/pkg/log"
)

// ProvideContextLogger injects a provided logger into request context.
func ProvideContextLogger(logger *slog.Logger) bunrouter.MiddlewareFunc {
	return func(next bunrouter.HandlerFunc) bunrouter.HandlerFunc {
		return func(w http.ResponseWriter, r bunrouter.Request) error {
			ctxLog := log.ToContext(r.Context(), logger)

			return next(w, r.WithContext(ctxLog))
		}
	}
}

// ProvideLogRequestID is a middleware that injects a request id into the context of each request.
func ProvideLogRequestID(next bunrouter.HandlerFunc) bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, r bunrouter.Request) error {
		ctx := r.Context()

		requestID := GetRequestID(ctx)
		ctxLog := log.AppendCtx(ctx, slog.String("request_id", requestID))

		return next(w, r.WithContext(ctxLog))
	}
}
