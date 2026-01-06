package main

import (
	"fmt"
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"io"
	"log"
	"net"
	"sync/atomic"
)

type HandlerError struct {
	code    response.StatusCode
	message string
}

type Handler func(w *response.Writer, r *request.Request)

type Server struct {
	listener net.Listener
	handler  Handler
	closed   atomic.Bool
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.closed.Load() {
				return
			}
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		go s.handle(conn)

	}
}
func (s *Server) handle(conn io.ReadWriteCloser) {
	defer conn.Close()
	rw := &response.Writer{Conn: conn}
	req, err := request.RequestFromReader(conn)
	if err != nil {
		h := response.GetDefaultHeaders(0)
		body := []byte("Bad Request: " + err.Error())
		h.Replace("content-length", fmt.Sprintf("%d", len(body)))
		rw.WriteStatusLine(response.StatusBadRequest)
		rw.WriteHeaders(h)
		rw.WriteBody(body)
		return
	}
	s.handler(rw, req)
}

func Serve(port uint16, handler Handler) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	s := &Server{
		handler:  handler,
		listener: listener,
	}
	s.closed.Store(false)
	go s.listen()
	return s, nil
}

func (s *Server) Close() error {
	s.closed.Store(true)
	return nil
}
