package middleware

import (
	"log/slog"
	"net/http"

	"github.com/sknv/protomock/pkg/log"
)

// ContextLogger injects a provided logger into request context.
func ContextLogger(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		handler := func(w http.ResponseWriter, r *http.Request) {
			ctxLog := log.ToContext(r.Context(), logger)

			next.ServeHTTP(w, r.WithContext(ctxLog))
		}

		return http.HandlerFunc(handler)
	}
}

// LogRequestID is a middleware that injects a request id into the context of each request.
func LogRequestID(next http.Handler) http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		requestID := GetRequestID(ctx)
		ctxLog := log.AppendCtx(ctx, slog.String("request_id", requestID))

		next.ServeHTTP(w, r.WithContext(ctxLog))
	}

	return http.HandlerFunc(handler)
}
