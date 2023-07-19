package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/vladimirpekarski/wordofwisdom/internal/book"
	"github.com/vladimirpekarski/wordofwisdom/internal/config"
	"github.com/vladimirpekarski/wordofwisdom/internal/lib/logger"
	"github.com/vladimirpekarski/wordofwisdom/internal/pow"
	"github.com/vladimirpekarski/wordofwisdom/internal/server"
	"golang.org/x/exp/slog"
)

func main() {
	cfg := config.MustLoad()

	log := logger.New(cfg.Env)
	log.Info("starting wordofwisdom", slog.String("env", cfg.Env),
		slog.String("host", cfg.Host),
		slog.String("port", cfg.Port),
		slog.Int("difficulty", cfg.Difficulty))

	b, err := book.New()
	if err != nil {
		panic(err)
	}

	p := pow.New(log)

	srv := server.New(server.Params{
		Host:          cfg.Host,
		Port:          cfg.Port,
		Log:           log,
		Book:          b,
		Pow:           p,
		PowDifficulty: cfg.Difficulty,
	})

	srv.Run()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Info("Shutting down server...")
	srv.Stop()
	log.Info("Server stopped.")
}
