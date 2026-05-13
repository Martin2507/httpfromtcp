package server

import (
	"fmt"
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"log"
	"net"
	"sync/atomic"
)

type Server struct {
	closed      atomic.Bool
	listener    net.Listener
	handlerFunc Handler
}

type Handler func(w *response.Writer, req *request.Request)

func Serve(port int, h Handler) (*Server, error) {

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))

	if err != nil {
		return nil, err
	}

	serv := Server{}
	serv.closed.Store(false)
	serv.listener = listener
	serv.handlerFunc = h

	go serv.listen()

	return &serv, nil
}

func (s *Server) Close() error {

	s.closed.Store(true)
	err := s.listener.Close()

	if err != nil {
		return err
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

			fmt.Printf("Error: %s", err)
			continue
		}

		go s.handle(conn)
	}

}

func (s *Server) handle(conn net.Conn) {

	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("Error closing connection: %s", err)
		}
	}()

	req, err := request.RequestFromReader(conn)
	if err != nil {
		log.Printf("Error parsing request : %s", err)
		return
	}

	writer := response.Writer{
		W: conn,
	}

	s.handlerFunc(&writer, req)

}
