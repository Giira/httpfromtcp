package server

import (
	"bytes"
	"fmt"
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"log"
	"net"
	"sync/atomic"
)

type Server struct {
	listener net.Listener
	closed   atomic.Bool
}

type Handler func(w *response.Writer, req *request.Request) HandlerError

type HandlerError struct {
	StatusCode int
	Message    string
}

func WriteError(w *response.Writer, h HandlerError) {
	w.WriteStatusLine(response.StatusCode(h.StatusCode))
	headers := response.GetDefaultHeaders(len(h.Message))
	w.WriteHeaders(headers)
	w.conn.Write([]byte(h.Message))
}

func Serve(port int, f Handler) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, fmt.Errorf("error listening on port: %v - %v", port, err)
	}
	s := &Server{
		listener: listener,
	}
	go s.listen(f)
	return s, nil
}

func (s *Server) Close() error {
	s.closed.Store(true)
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}

func (s *Server) listen(f Handler) {
	for {
		con, err := s.listener.Accept()
		if err != nil {
			if s.closed.Load() {
				return
			}
			log.Printf("Failed to accept connection: %v", err)
			continue
		}
		go s.handle(con, f)
	}
}

func (s *Server) handle(conn net.Conn, f Handler) {
	defer conn.Close()
	req, err := request.RequestFromReader(conn)
	if err != nil {
		log.Printf("error: failed to get request from reader: %v", err)
	}
	b := bytes.Buffer{}
	hErr := f(&b, req)
	if hErr.Message != "" {
		WriteError(conn, hErr)
		return
	} else {
		h := response.GetDefaultHeaders(len(b.Bytes()))
		err = response.WriteStatusLine(conn, response.CodeOK)
		err = response.WriteHeaders(conn, h)
		if err != nil {
			log.Printf("error: failed to write headers properly: %v", err)
		}
		conn.Write(b.Bytes())
		return
	}

}
