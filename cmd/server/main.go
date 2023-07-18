package main

import (
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/exp/slog"

	"github.com/vladimirpekarski/wordofwisdom/internal/book"
	"github.com/vladimirpekarski/wordofwisdom/internal/config"
	"github.com/vladimirpekarski/wordofwisdom/internal/lib/logger"
	"github.com/vladimirpekarski/wordofwisdom/internal/server"
)

func main() {
	cfg := config.MustLoad()

	log := logger.New(cfg.Env)
	log.Info("starting _wordofwisdom", slog.String("env", cfg.Env),
		slog.String("host", cfg.Host),
		slog.String("port", cfg.Port))

	b, err := book.New()
	if err != nil {
		panic(err)
	}

	srv := server.New(server.Params{
		Host: cfg.Host,
		Port: cfg.Port,
		Log:  log,
		Book: b,
	})

	srv.Run()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Info("Shutting down server...")
	srv.Stop()
	log.Info("Server stopped.")
}
