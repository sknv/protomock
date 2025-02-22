package middleware

import (
	"cmp"
	"log/slog"
	"net/http"
	"time"

	"github.com/uptrace/bunrouter"

	"github.com/sknv/protomock/pkg/log"
)

// LogRequest is a slightly modified version of the provided logger middleware.
func LogRequest(next bunrouter.HandlerFunc) bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, r bunrouter.Request) error {
		start := time.Now()

		respWriter := newCustomResponseWriter(w) // Save a response status.
		err := next(&respWriter, r)

		// Log data.
		ctx := r.Context()

		fields := []any{
			slog.Int("status", respWriter.status),
			slog.String("method", r.Method),
			slog.String("uri", r.RequestURI),
			slog.String("host", r.Host),
			slog.String("remote_ip", r.RemoteAddr),
			slog.String("user_agent", r.UserAgent()),
			slog.Int64("latency_ms", time.Since(start).Milliseconds()),
			slog.String("bytes_in", cmp.Or(r.Header.Get("Content-Length"), "0")),
			slog.Int("bytes_out", respWriter.size),
		}
		if err != nil {
			fields = append(fields, slog.Any("error", err))
		}

		log.FromContext(ctx).InfoContext(ctx, "HTTP request handled", fields...)

		return err
	}
}

// customResponseWriter is an HTTP response logger that keeps HTTP status code and the number of bytes written.
type customResponseWriter struct {
	http.ResponseWriter

	status int
	size   int
}

func newCustomResponseWriter(w http.ResponseWriter) customResponseWriter {
	return customResponseWriter{
		ResponseWriter: w,
		status:         http.StatusOK,
		size:           0,
	}
}

// WriteHeader implements http.ResponseWriter and saves status.
func (c *customResponseWriter) WriteHeader(status int) {
	c.status = status
	c.ResponseWriter.WriteHeader(status)
}

// Write implements http.ResponseWriter and tracks number of bytes written.
func (c *customResponseWriter) Write(b []byte) (int, error) {
	size, err := c.ResponseWriter.Write(b)
	c.size += size

	return size, err //nolint:wrapcheck // proxy
}
