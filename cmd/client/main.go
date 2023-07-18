package main

import (
	"golang.org/x/exp/slog"

	"github.com/vladimirpekarski/wordofwisdom/internal/client"
	"github.com/vladimirpekarski/wordofwisdom/internal/config"
	"github.com/vladimirpekarski/wordofwisdom/internal/lib/logger"
)

func main() {
	cfg := config.MustLoad()

	log := logger.New(cfg.Env)
	log.Info("starting client", slog.String("env", cfg.Env))

	c := client.New(cfg.Client.Address, log)

	c.Wisdom()
}
