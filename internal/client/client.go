package client

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/vladimirpekarski/wordofwisdom/internal/message"
	"github.com/vladimirpekarski/wordofwisdom/internal/message/gob"
	"github.com/vladimirpekarski/wordofwisdom/internal/pow"
	"golang.org/x/exp/slog"
)

type Client struct {
	address string
	log     *slog.Logger
	pow     pow.Pow
}

func New(address string, log *slog.Logger, pow pow.Pow) *Client {
	return &Client{
		address: address,
		log:     log,
		pow:     pow,
	}
}

func (c *Client) Quote(ctx context.Context) (string, string, error) {
	conn, err := net.Dial("tcp", c.address)
	if err != nil {
		return "", "", fmt.Errorf("failed to connect: %w", err)
	}
	defer func() {
		_ = conn.Close()
	}()

	var ch message.Challenge
	if err := gob.ReceiveMessage(conn, &ch); err != nil {
		return "", "", fmt.Errorf("failed to receive challenge: %w", err)
	}

	c.log.Info("challenge received",
		slog.String("random_str", ch.RandomStr),
		slog.String("hash_prefix", ch.HashPrefix))

	sol, err := c.pow.Solve(ctx, ch)
	if err != nil {
		return "", "", fmt.Errorf("failed to solve challenge: %w", err)
	}

	if err := gob.SendMessage(conn, sol); err != nil {
		return "", "", fmt.Errorf("failed to send solution: %w", err)
	}

	var rec message.BookRecord
	if err := gob.ReceiveMessage(conn, &rec); err != nil {
		return "", "", fmt.Errorf("failed to receive book record: %w", err)
	}

	if rec.PassedValidation {
		return rec.Quote, rec.Author, nil
	}

	return "", "", errors.New("wrong solution")
}
