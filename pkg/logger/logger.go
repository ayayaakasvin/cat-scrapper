package logger

import (
	"log/slog"
	"os"

	"github.com/ayayaakasvin/cat-scrapper/internal/config"
)

func New(cfg config.LoggerConfig) *slog.Logger {
	opts := &slog.HandlerOptions{
		Level: slog.Level(cfg.Level),
	}

	var handler slog.Handler

	if cfg.JSON {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	logger := slog.New(handler).With(
		slog.String("service", cfg.Service),
	)

	return logger
}