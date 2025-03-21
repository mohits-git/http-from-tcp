package server

import (
	"http-from-tcp/internal/request"
	"http-from-tcp/internal/response"
	"log"
)

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

type Handler func(w response.Writer, req *request.Request)

func (hErr *HandlerError) Write(w response.Writer) {
	err := w.WriteStatusLine(hErr.StatusCode)
	err = w.WriteHeaders(response.GetDefaultHeaders(len(hErr.Message)))
	err = w.WriteBody([]byte(hErr.Message))
	if err != nil {
		log.Println("Error while writing error response: ", err)
	}
}
