// Command server runs the books REST API.
package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"restapi/internal/config"
	"restapi/internal/database"
	"restapi/internal/server"
)

func main() {
	if err := run(); err != nil {
		slog.Error("startup failed", "error", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	logger := newLogger(cfg.LogLevel)
	slog.SetDefault(logger)

	db, err := database.Open(cfg)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := database.Close(db); cerr != nil {
			logger.Error("closing database", "error", cerr)
		}
	}()

	app := server.New(cfg, db, logger)

	// Start the server in the background.
	errCh := make(chan error, 1)
	go func() {
		logger.Info("server listening", "port", cfg.Port)
		if lerr := app.Listen(":" + cfg.Port); lerr != nil {
			errCh <- lerr
		}
	}()

	// Wait for an interrupt or a fatal server error.
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-errCh:
		return err
	case sig := <-stop:
		logger.Info("shutdown signal received", "signal", sig.String())
	}

	// Graceful shutdown with a bounded timeout.
	ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()
	if err := app.ShutdownWithContext(ctx); err != nil {
		return err
	}
	logger.Info("server stopped cleanly")
	return nil
}

func newLogger(level string) *slog.Logger {
	var lvl slog.Level
	switch level {
	case "debug":
		lvl = slog.LevelDebug
	case "warn":
		lvl = slog.LevelWarn
	case "error":
		lvl = slog.LevelError
	default:
		lvl = slog.LevelInfo
	}
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: lvl}))
}
