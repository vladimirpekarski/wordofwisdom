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

	var wg sync.WaitGroup

	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			c.Wisdom()
		}()
	}

	wg.Wait()
}
