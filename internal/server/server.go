package server

import (
	"bytes"
	"fmt"
	"http-from-tcp/internal/request"
	"http-from-tcp/internal/response"
	"log"
	"net"
	"sync/atomic"
)

type Server struct {
	port     int
	listener net.Listener
	closed   atomic.Bool
	handler  Handler
}

func Serve(port int, handler Handler) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	s := &Server{
		port:     port,
		listener: listener,
		closed:   atomic.Bool{},
		handler:  handler,
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
	// parse req
	req, err := request.RequestFromReader(conn)
	if err != nil {
		hErr := HandlerError{
			StatusCode: response.StatusCodeBadRequest,
			Message:    "Bad Request",
		}
		hErr.Write(conn)
		return
	}

	// call user handler
	buf := bytes.NewBuffer([]byte{})
	handlerErr := s.handler(buf, req)
	if handlerErr != nil {
		handlerErr.Write(conn)
		return
	}

	// write response (ok, if no handler error)
	err = response.WriteStatusLine(conn, response.StatusCodeOK)
	headers := response.GetDefaultHeaders(buf.Len())
	err = response.WriteHeaders(conn, headers)
	_, err = conn.Write(buf.Bytes())
	if err != nil {
		log.Println("Error while writing body: ", err)
	}
}
