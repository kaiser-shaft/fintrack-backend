package app

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/kaiser-shaft/fintrack-backend/config"
)

func Run(ctx context.Context, cfg *config.Config) error {
	container := NewContainer(ctx, cfg)
	defer container.Close()

	srv, err := container.HTTPServer()
	if err != nil {
		return fmt.Errorf("app.Run: %w", err)
	}

	sigQuit := make(chan os.Signal, 1)
	signal.Notify(sigQuit, syscall.SIGTERM, syscall.SIGINT)

	logger := container.Logger()

	serverErrors := make(chan error, 1)
	go func() {
		logger.Info("HTTP server starting", slog.String("port", cfg.HTTP.Port))
		serverErrors <- srv.Start()
	}()

	select {
	case err := <-serverErrors:
		return fmt.Errorf("critical server error: %v", err)
	case <-sigQuit:
		logger.Warn("App got signal to quit")
		container.Close()
	}

	return nil
}
