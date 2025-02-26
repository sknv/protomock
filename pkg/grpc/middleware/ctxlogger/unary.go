package ctxlogger

import (
	"context"
	"log/slog"

	"google.golang.org/grpc"

	"github.com/sknv/protomock/pkg/grpc/middleware/requestid"
	"github.com/sknv/protomock/pkg/log"
)

// ProvideUnaryContextLogger returns a new unary server interceptor that adds slog.Logger to the context.
func ProvideUnaryContextLogger(logger *slog.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
	) (any, error) {
		ctxLog := log.ToContext(ctx, logger)

		return handler(ctxLog, req)
	}
}

// ProvideUnaryLogRequestID returns a new unary server interceptor that injects a request id
// into the context of each request.
func ProvideUnaryLogRequestID(
	ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
) (any, error) {
	requestID := requestid.GetRequestID(ctx)
	ctxLog := log.AppendCtx(ctx, slog.String("request_id", requestID))

	return handler(ctxLog, req)
}
