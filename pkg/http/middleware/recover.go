package middleware

import (
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/sknv/protomock/pkg/log"
)

// Recover is a middleware that recovers from panics, logs the panic and returns a HTTP 500 status if possible.
func Recover(next http.Handler) http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		defer func() { //nolint:contextcheck // panic recover
			if rvr := recover(); rvr != nil {
				ctx := r.Context()
				log.FromContext(ctx).ErrorContext(ctx, "Request panic",
					slog.String("url", r.URL.String()),
					slog.Any("reason", rvr),
					slog.String("stack", string(debug.Stack())),
				)

				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(handler)
}
