package main

import (
	"context"
	"flag"
	"fmt"
	stdlog "log"
	"log/slog"
	"time"

	"github.com/uptrace/bunrouter"
	_ "go.uber.org/automaxprocs"

	"github.com/sknv/protomock/internal/config"
	"github.com/sknv/protomock/internal/container"
	"github.com/sknv/protomock/internal/transport/http"
	"github.com/sknv/protomock/pkg/http/middleware"
	"github.com/sknv/protomock/pkg/log"
	"github.com/sknv/protomock/pkg/os"
)

const _stopTimeout = time.Second * 10

func main() {
	configPath := config.FilePathFlag()
	flag.Parse() //nolint:wsl // process a variable above

	cfg, err := config.Parse(*configPath)
	fatalIfError(err)

	err = run(cfg)
	fatalIfError(err)
}

func run(cfg *config.Config) error {
	appCtx, cancelApp := context.WithCancel(context.Background())
	defer cancelApp()

	app, err := buildApp(cfg)
	if err != nil {
		return fmt.Errorf("build application: %w", err)
	}

	// Start the application and wait for the signal to shutdown.
	if err = app.Run(appCtx); err != nil {
		return fmt.Errorf("run apllcation: %w", err)
	}

	<-os.NotifyAboutExit()

	// Stop the application.
	cancelApp()

	if err = stopApp(app, _stopTimeout); err != nil {
		app.Logger().Unwrap().
			ErrorContext(appCtx, "Can't stop application properly", slog.Any("error", err))
	}

	return nil
}

func buildApp(cfg *config.Config) (*container.Application, error) {
	app := container.NewApplication()

	// Logger.
	logger := app.RegisterLogger(log.Config{Level: cfg.Log.Level})
	slog.SetDefault(logger) // Sets the global default logger.

	// HTTP server.
	if cfg.HTTPServer.Enabled {
		if err := buildHTTPServer(app, cfg); err != nil {
			return nil, fmt.Errorf("build http server: %w", err)
		}
	}

	return app, nil
}

func buildHTTPServer(app *container.Application, cfg *config.Config) error {
	mocks, err := http.BuildMocks(cfg.HTTPServer.MocksDir)
	if err != nil {
		return fmt.Errorf("build http mocks: %w", err)
	}

	defaultMiddlewares := []bunrouter.MiddlewareFunc{
		middleware.ProvideContextLogger(app.Logger().Unwrap()),
		middleware.ProvideRequestID,
		middleware.ProvideLogRequestID,
		middleware.LogRequest,
		middleware.HandleError,
		middleware.Recover,
	}
	router := app.RegisterHTTPServer(
		fmt.Sprintf(":%d", cfg.HTTPServer.Port),
		bunrouter.Use(defaultMiddlewares...),
	)

	handlers := http.NewHandlers(mocks)
	handlers.Route(router)

	return nil
}

// stopApp tryes to stop the app gracefully.
func stopApp(app *container.Application, timeout time.Duration) error {
	stopCtx, cancelStop := context.WithTimeout(context.Background(), timeout)
	defer cancelStop()

	if err := app.Stop(stopCtx); err != nil {
		return fmt.Errorf("stop application in timeout: %w", err)
	}

	return nil
}

func fatalIfError(err error) {
	if err != nil {
		stdlog.Fatal(err)
	}
}
