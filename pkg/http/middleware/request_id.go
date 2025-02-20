package middleware

import (
	"context"
	"crypto/rand"
	"net/http"
)

const _requestIDHeader = "X-Request-ID"

type ctxKey string

const _requestIDField ctxKey = "request_id"

// RequestID looks for header X-Request-ID and makes it as random id if not found,
// then populates it to the result's header and to request context.
func RequestID(next http.Handler) http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get(_requestIDHeader)
		if requestID == "" {
			requestID = rand.Text()
		}

		w.Header().Set(_requestIDHeader, requestID)
		ctxReqID := context.WithValue(r.Context(), _requestIDField, requestID)

		next.ServeHTTP(w, r.WithContext(ctxReqID))
	}

	return http.HandlerFunc(handler)
}

// GetRequestID returns request id from the context.
func GetRequestID(ctx context.Context) string {
	if id, ok := ctx.Value(_requestIDField).(string); ok {
		return id
	}

	return ""
}
