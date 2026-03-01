package server

import (
	"fmt"
	"httpfromtcp/internal/request"
	"log"
	"net"
	"sync/atomic"
)

type Server struct {
	listener net.Listener
	state    atomic.Bool
}

func Serve(port int) (*Server, error) {
	f, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, fmt.Errorf("error listening on port: %v - %v", port, err)
	}
	s := &Server{
		listener: f,
	}
	go s.listen()
	return s, nil
}

func (s *Server) Close() error {
	err := s.listener.Close()
	if err != nil {
		return err
	}
	s.state.Store(false)
	return nil
}

func (s *Server) listen() {
	for {
		con, err := s.listener.Accept()
		if err != nil {
			if s.state.Load() {
				return
			}
			log.Printf("Failed to accept connection: %v", err)
			continue
		}
		go s.handle(con)
	}
}

func (s *Server) handle(conn net.Conn) {
	req, err := request.RequestFromReader(conn)
	if err != nil {
		log.Fatalf("error parsing request: %v", err)
	}
	request.PrintData(req)
	conn.Close()
}
