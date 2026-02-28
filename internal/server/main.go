package server

import (
	"fmt"
	"httpfromtcp/internal/request"
	"log"
	"net"
	"strconv"
)

type Server struct {
	listener net.Listener
	state    serverState
}

type serverState int

const (
	sOpen serverState = iota
	sClosed
)

func Serve(port int) (*Server, error) {
	f, err := net.Listen("tcp", strconv.Itoa(port))
	if err != nil {
		return nil, fmt.Errorf("error listening on port: %v - %v", port, err)
	}
	s := &Server{
		listener: f,
		state:    sOpen,
	}
	return s, nil
}

func (s *Server) Close() error {
	err := s.listener.Close()
	if err != nil {
		return err
	}
	s.state = sClosed
	return nil
}

func (s *Server) listen() {
	for {
		con, err := s.listener.Accept()
		if err != nil {
			log.Fatalf("error accepting connection: %v", err)
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
}
