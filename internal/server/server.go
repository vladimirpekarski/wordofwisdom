package server

import (
	"fmt"
	"net"
	"sync"
	"time"

	"golang.org/x/exp/slog"

	"github.com/vladimirpekarski/wordofwisdom/internal/book"
)

type Server struct {
	host       string
	port       string
	log        *slog.Logger
	book       book.Book
	wg         sync.WaitGroup
	listener   net.Listener
	shutdown   chan struct{}
	connection chan net.Conn
}

type Params struct {
	Host string
	Port string
	Log  *slog.Logger
	Book book.Book
}

func New(p Params) *Server {
	l, err := net.Listen("tcp", p.Host+":"+p.Port)
	if err != nil {
		panic(fmt.Sprintf("failed to start: %s", err.Error()))
	}
	return &Server{
		host:       p.Host,
		port:       p.Port,
		log:        p.Log,
		book:       p.Book,
		listener:   l,
		shutdown:   make(chan struct{}),
		connection: make(chan net.Conn),
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

	responseStr := s.book.RandomQuote()
	n, err := conn.Write([]byte(responseStr))
	s.log.Info("", slog.Int("n", n))
	if err != nil {
		panic(err)
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
