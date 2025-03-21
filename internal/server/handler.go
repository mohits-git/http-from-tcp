package server

import (
	"http-from-tcp/internal/request"
	"http-from-tcp/internal/response"
	"io"
)

type HandlerError struct {
  StatusCode response.StatusCode
  Message string
}

type Handler func(w io.Writer, req *request.Request) *HandlerError

func (hErr *HandlerError) Write(w io.Writer) {
  response.WriteStatusLine(w, hErr.StatusCode)
  response.WriteHeaders(w, response.GetDefaultHeaders(len(hErr.Message)))
  w.Write([]byte(hErr.Message))
}
