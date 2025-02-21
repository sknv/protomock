package container

import (
	"context"
	"errors"
	"fmt"
	stdlog "log"
	"log/slog"
	"net/http"
	"time"

	"github.com/sknv/protomock/pkg/http/routegroup"
	"github.com/sknv/protomock/pkg/option"
)

type httpServer struct {
	router *routegroup.Group
	server *http.Server
}

func (a *Application) RegisterHTTPServer(address string) *routegroup.Group {
	router := routegroup.New(http.NewServeMux())
	httpServer := &httpServer{
		router: router,
		server: newHTTPServer(address, router),
	}

	a.httpServer = option.Some(httpServer)

	return router
}

func (a *Application) Router() option.Option[*routegroup.Group] {
	if a.httpServer.IsSome() {
		return option.Some(
			a.httpServer.Unwrap().router,
		)
	}

	return option.None[*routegroup.Group]()
}

// ----------------------------------------------------------------------------

const _readHeaderTimeout = time.Second * 30

func newHTTPServer(address string, handler http.Handler) *http.Server {
	return &http.Server{ //nolint:exhaustruct // too many unused fields
		Addr:              address,
		Handler:           handler,
		ReadHeaderTimeout: _readHeaderTimeout,
	}
}

func (a *Application) runHTTPServer(ctx context.Context) error {
	if a.httpServer.IsNone() {
		return nil // No HTTP server registered.
	}

	logger := a.logger.UnwrapOrElse(slog.Default)
	server := a.httpServer.Unwrap().server

	logger.InfoContext(ctx, "Starting http server...", slog.String("address", server.Addr))
	defer logger.InfoContext(ctx, "Http server started")

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			stdlog.Fatalf("Can't start http server: %v", err)
		}
	}()

	// Remember to stop the server.
	a.closers.Add(func(closeCtx context.Context) error {
		logger.InfoContext(closeCtx, "Stopping http server...")

		if err := server.Shutdown(closeCtx); err != nil {
			return fmt.Errorf("shutdown http server: %w", err)
		}

		logger.InfoContext(closeCtx, "Http server stopped")

		return nil
	})

	return nil
}
