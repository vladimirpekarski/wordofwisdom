package client

import (
	"bytes"
	"encoding/gob"
	"github.com/vladimirpekarski/wordofwisdom/internal/pow"
	"golang.org/x/exp/slog"
	"net"
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

func (c *Client) Wisdom() string {
	conn, err := net.Dial("tcp", c.address)
	defer func() {
		_ = conn.Close()
	}()

	if err != nil {
		c.log.Error("failed to dial",
			slog.String("address", c.address),
			slog.String("error", err.Error()))

		return ""
	}

	dec := gob.NewDecoder(conn)

	var ch pow.Challenge

	if err := dec.Decode(&ch); err != nil {
		c.log.Error("failed to decode message", slog.String("error", err.Error()))
		return ""
	}

	c.log.Info("message received",
		slog.String("random_str", ch.RandomStr),
		slog.String("hash_prefix", ch.HashPrefix))

	sol := c.pow.Solve(ch)

	buf := new(bytes.Buffer)
	gobobj := gob.NewEncoder(buf)
	err = gobobj.Encode(sol)
	if err != nil {
		c.log.Error("failed to encode message", slog.String("error", err.Error()))
		return ""
	}

	_, err = conn.Write(buf.Bytes())
	if err != nil {
		c.log.Error("failed to send message", slog.String("error", err.Error()))
		return ""
	}

	var rec pow.Record
	if err := dec.Decode(&rec); err != nil {
		c.log.Error("failed to decode message", slog.String("error", err.Error()))
		return ""
	}

	c.log.Info("quote", slog.String("quote", rec.Quote), slog.String("author", rec.Author))

	return ""
}
