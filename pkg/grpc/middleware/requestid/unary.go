package requestid

import (
	"context"
	"crypto/rand"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const _requestIDHeader = "X-Request-ID"

type ctxKey string

const _requestIDField ctxKey = "request_id"

// ProvideUnaryRequestID looks for metadata key X-Request-ID and makes it as random id if not found,
// then populates it to the request context.
func ProvideUnaryRequestID(
	ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
) (any, error) {
	var requestID string

	meta, _ := metadata.FromIncomingContext(ctx)

	headerValues := meta.Get(_requestIDHeader)
	if len(headerValues) > 0 {
		requestID = headerValues[0]
	} else {
		requestID = rand.Text()
	}

	ctxReqID := context.WithValue(ctx, _requestIDField, requestID)

	return handler(ctxReqID, req)
}

// GetRequestID returns request id from the context.
func GetRequestID(ctx context.Context) string {
	if id, ok := ctx.Value(_requestIDField).(string); ok {
		return id
	}

	return ""
}
