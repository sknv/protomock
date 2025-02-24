package container

import (
	"context"
	"fmt"
	"log/slog"

	"golang.org/x/sync/errgroup"

	"github.com/sknv/protomock/pkg/closer"
	"github.com/sknv/protomock/pkg/option"
)

type Application struct {
	closers    *closer.Closers
	logger     option.Option[*slog.Logger]
	httpServer option.Option[*httpServer]
	grpcServer option.Option[*grpcServer]
}

func NewApplication() *Application {
	return &Application{
		closers:    closer.New(),
		logger:     option.None[*slog.Logger](),
		httpServer: option.None[*httpServer](),
		grpcServer: option.None[*grpcServer](),
	}
}

func (a *Application) Run(ctx context.Context) error {
	logger := a.logger.UnwrapOrElse(slog.Default)
	logger.InfoContext(ctx, "Starting application...")

	if err := runParallel(ctx,
		a.runHTTPServer,
		a.runGRPCServer,
	); err != nil {
		return fmt.Errorf("run components in parallel: %w", err)
	}

	logger.InfoContext(ctx, "Application started")

	return nil
}

func (a *Application) Stop(ctx context.Context) error {
	logger := a.logger.UnwrapOrElse(slog.Default)
	logger.InfoContext(ctx, "Stopping application...")

	if err := a.closers.Close(ctx); err != nil {
		return fmt.Errorf("close component: %w", err)
	}

	logger.InfoContext(ctx, "Application stopped")

	return nil
}

// ----------------------------------------------------------------------------

type runner func(ctx context.Context) error

func runParallel(ctx context.Context, runners ...runner) error {
	var group errgroup.Group

	for _, run := range runners {
		group.Go(func() error {
			return run(ctx)
		})
	}

	if err := group.Wait(); err != nil {
		return fmt.Errorf("wait group: %w", err)
	}

	return nil
}
