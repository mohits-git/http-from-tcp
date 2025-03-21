package server

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"
)

type Server struct {
	port     int
	listener net.Listener
	closed   atomic.Bool
}

func Serve(port int) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	s := &Server{
		port:     port,
		listener: listener,
		closed:   atomic.Bool{},
	}

	go s.listen()

	return s, nil
}

func (s *Server) Close() error {
	s.closed.Store(true)
	return s.listener.Close()
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.closed.Load() {
				return // ignore error if server is closed
			}
			log.Println("Error while accepting connection: ", err)
			continue
		}

		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
  _, err := conn.Write([]byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n\r\nHello World!"))
  if err != nil {
    log.Println("Error while writing response: ", err)
  }
  // TODO: Implement request handling
}
