package server

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"

	"github.com/JA50N14/httpfromtcp/internal/request"
	"github.com/JA50N14/httpfromtcp/internal/response"
)

type Server struct {
	listener net.Listener
	handler  Handler
	closed   atomic.Bool
}

type Handler func(w *response.Writer, req *request.Request)


func Serve(port int, handler Handler) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	s := &Server{
		listener: listener,
		handler:  handler,
	}
	s.closed.Store(false)

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
		conn, err := s.listener.Accept()
		if err != nil {
			if s.closed.Load() {
				return
			}
			log.Printf("error accepting connection: %v", err)
			continue
		}

		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	
	w := &response.Writer{
		Out:         conn,
		WriterState: response.StateWriteStatusLine,
	}


	req, err := request.RequestFromReader(conn)
	if err != nil {
		w.StatusCode = response.StatusCodeBadRequest
		
	}
	s.handler(w, req)
	return
}
