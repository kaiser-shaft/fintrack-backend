package logger

import (
	"log/slog"
	"os"
)

type Config struct {
	Level string `envconfig:"LOG_LEVEL"`
}

func Init(cfg Config) *slog.Logger {
	level := parseLevel(cfg.Level)

	opts := &slog.HandlerOptions{
		Level: level,
	}

	handler := slog.NewJSONHandler(os.Stdout, opts)
	logger := slog.New(handler)

	slog.SetDefault(logger)

	return logger
}

func parseLevel(lvl string) slog.Level {
	var level slog.Level
	switch lvl {
	case "debug":
		level = slog.LevelDebug
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}
	return level
}
