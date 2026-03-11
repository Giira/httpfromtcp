package server

import (
	"fmt"
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"io"
	"log"
	"net"
	"sync/atomic"
)

type Server struct {
	listener net.Listener
	closed   atomic.Bool
}

type Handler func(w io.Writer, req *request.Request) *HandlerError

type HandlerError struct {
	StatusCode int
	Message    string
}

func Serve(port int, f Handler) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, fmt.Errorf("error listening on port: %v - %v", port, err)
	}
	s := &Server{
		listener: listener,
	}
	go s.listen()
	return s, nil
}

func (s *Server) Close() error {
	s.closed.Store(true)
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}

func (s *Server) listen() {
	for {
		con, err := s.listener.Accept()
		if err != nil {
			if s.closed.Load() {
				return
			}
			log.Printf("Failed to accept connection: %v", err)
			continue
		}
		go s.handle(con)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	err := response.WriteStatusLine(conn, response.CodeOK)
	if err != nil {
		log.Printf("%v", err)
	}
	h := response.GetDefaultHeaders(0)
	err = response.WriteHeaders(conn, h)
	if err != nil {
		log.Printf("error: failed to write headers properly: %v", err)
	}
	conn.Write([]byte("\r\n"))
}
