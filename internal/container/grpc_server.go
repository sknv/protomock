package container

import (
	"context"
	"fmt"
	stdlog "log"
	"log/slog"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/sknv/protomock/pkg/closer"
	"github.com/sknv/protomock/pkg/option"
)

type grpcServer struct {
	address string
	server  *grpc.Server
}

func (a *Application) RegisterGRPCServer(address string, opts ...grpc.ServerOption) *grpc.Server {
	server := grpc.NewServer(opts...)
	reflection.Register(server) // Register reflection service on gRPC server.

	grpcServer := &grpcServer{
		address: address,
		server:  server,
	}

	a.grpcServer = option.Some(grpcServer)

	return server
}

func (a *Application) GRPCServer() option.Option[*grpc.Server] {
	if a.grpcServer.IsSome() {
		return option.Some(
			a.grpcServer.Unwrap().server,
		)
	}

	return option.None[*grpc.Server]()
}

// ----------------------------------------------------------------------------

func (a *Application) runGRPCServer(ctx context.Context) error {
	if a.grpcServer.IsNone() {
		return nil // No gRPC server registered.
	}

	logger := a.logger.UnwrapOrElse(slog.Default)
	grpcServer := a.grpcServer.Unwrap()

	logger.InfoContext(ctx, "Starting grpc server...", slog.String("address", grpcServer.address))
	defer logger.InfoContext(ctx, "Grpc server started")

	lis, err := net.Listen("tcp", grpcServer.address)
	if err != nil {
		return fmt.Errorf("listen tcp address: %w", err)
	}

	go func() {
		if err := grpcServer.server.Serve(lis); err != nil {
			stdlog.Fatalf("Can't start grpc server: %v", err)
		}
	}()

	// Remember to stop the server.
	a.closers.Add(func(closeCtx context.Context) error {
		logger.InfoContext(closeCtx, "Stopping grpc server...")

		if err := closer.CloseWithContext(closeCtx, func() error {
			grpcServer.server.GracefulStop()

			return nil
		}); err != nil {
			return fmt.Errorf("shutdown grpc server: %w", err)
		}

		logger.InfoContext(closeCtx, "Grpc server stopped")

		return nil
	})

	return nil
}
