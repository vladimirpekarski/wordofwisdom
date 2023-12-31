package server

import (
	"errors"
	"io"
	"net"

	"golang.org/x/exp/slog"

	"github.com/vladimirpekarski/wordofwisdom/internal/book"
	"github.com/vladimirpekarski/wordofwisdom/internal/message"
	"github.com/vladimirpekarski/wordofwisdom/internal/message/gob"
)

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=RandomRecordGetter
type RandomRecordGetter interface {
	GetRandomRecord() book.Record
}

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=ValidateGenerateChallenger
type ValidateGenerateChallenger interface {
	GenerateChallenge(difficulty int) (message.Challenge, error)
	Validate(ch message.Challenge, sl message.Solution) bool
}

func (s *Server) handleConnection(conn net.Conn) {
	defer func() {
		_ = conn.Close()
	}()

	ch, err := s.pow.GenerateChallenge(s.powDifficulty)
	if err != nil {
		s.log.Error("failed to get challenge", slog.String("error", err.Error()))
		return
	}

	if err := gob.SendMessage(conn, ch); err != nil {
		s.log.Error("failed to send message", slog.String("error", err.Error()))
		return
	}

	var sl message.Solution
	if err := gob.ReceiveMessage(conn, &sl); err != nil {
		if !errors.Is(err, io.EOF) {
			s.log.Error("failed to receive solution", slog.String("error", err.Error()))
		}
		return
	}

	if s.pow.Validate(ch, sl) {
		s.log.Debug("validation passed")

		q := s.book.GetRandomRecord()

		rec := message.BookRecord{
			Quote:            q.Quote,
			Author:           q.Author,
			PassedValidation: true,
		}

		if err := gob.SendMessage(conn, rec); err != nil {
			if !errors.Is(err, io.EOF) {
				s.log.Error("failed to send quote", slog.String("error", err.Error()))
			}
			return
		}
	} else {
		s.log.Debug("validation failed")

		rec := message.BookRecord{}
		if err := gob.SendMessage(conn, rec); err != nil {
			if !errors.Is(err, io.EOF) {
				s.log.Error("failed to send empty quote", slog.String("error", err.Error()))
			}
			return
		}
	}
}
