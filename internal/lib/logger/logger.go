package logger

import (
	"os"

	"golang.org/x/exp/slog"

	"github.com/vladimirpekarski/wordofwisdom/internal/env"
)

func New(e string) *slog.Logger {
	var log *slog.Logger

	switch e {
	case env.Local:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
				Level: slog.LevelDebug,
			}))
	case env.DockerCompose:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
				Level: slog.LevelInfo,
			}))
	}

	return log
}
