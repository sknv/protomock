package logger

import (
	"context"
	"log/slog"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"

	"github.com/sknv/protomock/pkg/log"
)

// LogUnaryRequest logs gRPC requests.
//
//nolint:nonamedreturns // used in defer
func LogUnaryRequest(
	ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
) (_ any, err error) {
	start := time.Now()

	defer func() {
		fields := []any{
			slog.String("method", info.FullMethod),
			slog.Uint64("code", uint64(status.Code(err))),
			slog.Int64("latency_ms", time.Since(start).Milliseconds()),
		}
		if err != nil {
			fields = append(fields, slog.Any("error", err))
		}

		log.FromContext(ctx).InfoContext(ctx, "gRPC request handled", fields...)
	}()

	return handler(ctx, req)
}
