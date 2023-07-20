package client

import (
	"context"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/vladimirpekarski/wordofwisdom/internal/client/mocks"
	"github.com/vladimirpekarski/wordofwisdom/internal/env"
	"github.com/vladimirpekarski/wordofwisdom/internal/lib/logger"
	"github.com/vladimirpekarski/wordofwisdom/internal/message"
	"github.com/vladimirpekarski/wordofwisdom/internal/message/gob"
)

func TestClient_Quote_Solved(t *testing.T) {
	ctx := context.Background()

	pow := mocks.NewSolver(t)
	pow.On("Solve", ctx, message.Challenge{
		RandomStr:  "str",
		HashPrefix: "000",
	}).
		Return(message.Solution{Hash: "az", Nonce: 1}, nil)

	client := NewMock("addr", logger.New(env.Local), pow)

	cl, serv := net.Pipe()
	defer func() {
		_ = serv.Close()
		_ = cl.Close()
	}()

	client.connFunc = func(network, address string) (net.Conn, error) {
		return cl, nil
	}

	go func() {
		ch := message.Challenge{
			RandomStr:  "str",
			HashPrefix: "000",
		}
		_ = gob.SendMessage(serv, ch)

		var sl message.Solution
		_ = gob.ReceiveMessage(serv, &sl)

		rec := message.BookRecord{
			Quote:            "quote",
			Author:           "author",
			PassedValidation: true,
		}
		_ = gob.SendMessage(serv, rec)
	}()

	quote, author, err := client.GetQuote(ctx)

	assert.NoError(t, err)
	assert.Equal(t, "quote", quote)
	assert.Equal(t, "author", author)
}

func TestClient_Quote_NotSolved(t *testing.T) {
	ctx := context.Background()

	pow := mocks.NewSolver(t)
	pow.On("Solve", ctx, message.Challenge{
		RandomStr:  "str",
		HashPrefix: "000",
	}).
		Return(message.Solution{Hash: "az", Nonce: 1}, nil)

	client := NewMock("addr", logger.New(env.Local), pow)

	cl, serv := net.Pipe()
	defer func() {
		_ = serv.Close()
		_ = cl.Close()
	}()

	client.connFunc = func(network, address string) (net.Conn, error) {
		return cl, nil
	}

	go func() {
		ch := message.Challenge{
			RandomStr:  "str",
			HashPrefix: "000",
		}
		_ = gob.SendMessage(serv, ch)

		var sl message.Solution
		_ = gob.ReceiveMessage(serv, &sl)

		rec := message.BookRecord{
			PassedValidation: false,
		}
		_ = gob.SendMessage(serv, rec)
	}()

	quote, author, err := client.GetQuote(ctx)

	assert.Error(t, err)
	assert.Equal(t, "", quote)
	assert.Equal(t, "", author)
}
