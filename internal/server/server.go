package server

import (
	"encoding/gob"
	"fmt"
	"net"
	"sync"
	"time"

	"golang.org/x/exp/slog"

	"github.com/vladimirpekarski/wordofwisdom/internal/book"
	"github.com/vladimirpekarski/wordofwisdom/internal/message"
	"github.com/vladimirpekarski/wordofwisdom/internal/pow"
)

type Server struct {
	host          string
	port          string
	log           *slog.Logger
	book          book.Book
	pow           pow.Pow
	powDifficulty int
	wg            sync.WaitGroup
	listener      net.Listener
	shutdown      chan struct{}
	connection    chan net.Conn
}

type Params struct {
	Host          string
	Port          string
	Log           *slog.Logger
	Book          book.Book
	Pow           pow.Pow
	PowDifficulty int
}

func New(p Params) *Server {
	l, err := net.Listen("tcp", p.Host+":"+p.Port)
	if err != nil {
		panic(fmt.Sprintf("failed to start: %s", err.Error()))
	}
	return &Server{
		host:          p.Host,
		port:          p.Port,
		log:           p.Log,
		book:          p.Book,
		pow:           p.Pow,
		powDifficulty: p.PowDifficulty,
		listener:      l,
		shutdown:      make(chan struct{}),
		connection:    make(chan net.Conn),
	}
}

func (s *Server) Run() {
	s.wg.Add(2)
	go s.acceptConnections()
	go s.handleConnections()

}

func (s *Server) acceptConnections() {
	defer s.wg.Done()

	for {
		select {
		case <-s.shutdown:
			return
		default:
			conn, err := s.listener.Accept()
			if err != nil {
				continue
			}
			s.connection <- conn
		}
	}
}

func (s *Server) handleConnections() {
	defer s.wg.Done()

	for {
		select {
		case <-s.shutdown:
			return
		case conn := <-s.connection:
			go s.handleConnection(conn)
		}
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer func() {
		_ = conn.Close()
	}()

	ch, err := s.pow.GenerateChallenge(16, s.powDifficulty)
	if err != nil {
		s.log.Error("failed to get challenge", slog.String("error", err.Error()))
		return
	}

	enc := gob.NewEncoder(conn)
	dec := gob.NewDecoder(conn)

	err = enc.Encode(ch)
	if err != nil {
		s.log.Error("failed to encode message", slog.String("error", err.Error()))
		return
	}

	var sl message.Solution
	if err := dec.Decode(&sl); err != nil {
		s.log.Error("failed to decode message", slog.String("error", err.Error()))
	}

	if s.pow.Validate(ch, sl) {
		s.log.Debug("validation passed")

		q := s.book.RandomQuote()

		rec := message.BookRecord{
			Quote:            q.Quote,
			Author:           q.Author,
			PassedValidation: true,
		}

		err = enc.Encode(rec)
		if err != nil {
			s.log.Error("failed to encode message", slog.String("error", err.Error()))
			return
		}
	} else {
		s.log.Debug("validation failed")

		rec := message.BookRecord{}

		err = enc.Encode(rec)
		if err != nil {
			s.log.Error("failed to encode message", slog.String("error", err.Error()))
			return
		}
	}
}

func (s *Server) Stop() {
	close(s.shutdown)
	_ = s.listener.Close()

	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return
	case <-time.After(time.Second):
		s.log.Warn("Timed out waiting for connections to finish.")
		return
	}
}
