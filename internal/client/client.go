package client

import (
	"encoding/gob"
	"errors"
	"fmt"
	"net"

	"golang.org/x/exp/slog"

	"github.com/vladimirpekarski/wordofwisdom/internal/message"
	"github.com/vladimirpekarski/wordofwisdom/internal/pow"
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

func (c *Client) Quote() (string, string, error) {
	conn, err := net.Dial("tcp", c.address)

	if err != nil {
		return "", "", fmt.Errorf("failed to connect: %w", err)
	}

	defer func() {
		_ = conn.Close()
	}()

	enc := gob.NewEncoder(conn)
	dec := gob.NewDecoder(conn)

	var ch message.Challenge
	if err := dec.Decode(&ch); err != nil {
		return "", "", fmt.Errorf("failed to decode: %w", err)
	}

	c.log.Info("message received",
		slog.String("random_str", ch.RandomStr),
		slog.String("hash_prefix", ch.HashPrefix))

	sol := c.pow.Solve(ch)
	err = enc.Encode(sol)
	if err != nil {
		return "", "", fmt.Errorf("failed to encode: %w", err)
	}

	var rec message.BookRecord
	if err := dec.Decode(&rec); err != nil {
		return "", "", fmt.Errorf("failed to decode: %w", err)
	}

	if rec.PassedValidation {
		return rec.Quote, rec.Author, nil
	}

	return "", "", errors.New("wrong solution")
}
