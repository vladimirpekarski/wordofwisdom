package server

import (
	"net"
	"testing"

	"github.com/vladimirpekarski/wordofwisdom/internal/book"
	"github.com/vladimirpekarski/wordofwisdom/internal/env"
	"github.com/vladimirpekarski/wordofwisdom/internal/lib/logger"
	"github.com/vladimirpekarski/wordofwisdom/internal/message"
	"github.com/vladimirpekarski/wordofwisdom/internal/message/gob"
	"github.com/vladimirpekarski/wordofwisdom/internal/server/mocks"
)

func TestServer_handleConnections_validated(t *testing.T) {
	pow := mocks.NewPOWer(t)
	pow.On("GenerateChallenge", 5).
		Return(message.Challenge{RandomStr: "ab", HashPrefix: "00"}, nil).
		On("Validate", message.Challenge{RandomStr: "ab", HashPrefix: "00"},
			message.Solution{Hash: "ab", Nonce: 5}).
		Return(true)

	b := mocks.NewBooker(t)
	b.On("RandomQuote").
		Return(book.Record{
			Quote:  "some quote",
			Author: "some author",
		})

	srv := NewMock(Params{
		Log:           logger.New(env.Local),
		PowDifficulty: 5,
		Book:          b,
		Pow:           pow,
	})

	serv, client := net.Pipe()
	defer func() {
		_ = serv.Close()
		_ = client.Close()
	}()

	go func() {
		var ch message.Challenge
		_ = gob.ReceiveMessage(client, &ch)
		_ = gob.SendMessage(client, message.Solution{Hash: "ab", Nonce: 5})

		var rec message.BookRecord
		_ = gob.ReceiveMessage(client, &rec)
	}()

	srv.handleConnection(serv)
}

func TestServer_handleConnections_not_validated(t *testing.T) {
	pow := mocks.NewPOWer(t)
	pow.On("GenerateChallenge", 5).
		Return(message.Challenge{RandomStr: "ab", HashPrefix: "00"}, nil).
		On("Validate", message.Challenge{RandomStr: "ab", HashPrefix: "00"},
			message.Solution{Hash: "ab", Nonce: 5}).
		Return(false)

	b := mocks.NewBooker(t)

	srv := NewMock(Params{
		Log:           logger.New(env.Local),
		PowDifficulty: 5,
		Book:          b,
		Pow:           pow,
	})

	serv, client := net.Pipe()
	defer func() {
		_ = serv.Close()
		_ = client.Close()
	}()

	go func() {
		var ch message.Challenge
		_ = gob.ReceiveMessage(client, &ch)
		_ = gob.SendMessage(client, message.Solution{Hash: "ab", Nonce: 5})

		var rec message.BookRecord
		_ = gob.ReceiveMessage(client, &rec)
	}()

	srv.handleConnection(serv)
}
