package main

import (
	"github.com/vladimirpekarski/wordofwisdom/internal/pow"
	"golang.org/x/exp/slog"
	"sync"

	"github.com/vladimirpekarski/wordofwisdom/internal/client"
	"github.com/vladimirpekarski/wordofwisdom/internal/config"
	"github.com/vladimirpekarski/wordofwisdom/internal/lib/logger"
)

func main() {
	cfg := config.MustLoad()

	log := logger.New(cfg.Env)
	log.Info("starting client",
		slog.String("env", cfg.Env),
		slog.String("address", cfg.Address))

	p := pow.New(log)

	c := client.New(cfg.Client.Address, log, p)
	conns := cfg.Connections

	var wg sync.WaitGroup

	wg.Add(conns)

	for i := 0; i < conns; i++ {
		go func() {
			defer wg.Done()
			quote, author, err := c.Quote()
			if err != nil {
				log.Error("failed to get quote", slog.String("error", err.Error()))
				return
			}

			log.Info("received quote",
				slog.String("quote", quote),
				slog.String("author", author))
		}()
	}

	wg.Wait()
}
