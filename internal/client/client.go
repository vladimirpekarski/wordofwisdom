package client

import (
	"bufio"
	"golang.org/x/exp/slog"
	"net"
)

type Client struct {
	address string
	log     *slog.Logger
}

func New(address string, log *slog.Logger) *Client {
	return &Client{
		address: address,
		log:     log,
	}
}

func (c *Client) Wisdom() string {
	conn, err := net.Dial("tcp", c.address)
	if err != nil {
		// TODO return error
		c.log.Error("failed to dial", slog.String("address", c.address))
		return ""
	}

	message, _ := bufio.NewReader(conn).ReadString('\n')

	c.log.Info(message)
	return message
}
