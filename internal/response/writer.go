package response

import (
	"fmt"
	"http-from-tcp/internal/headers"
	"io"
)

type writerState int

const (
	writerStateStatusLine writerState = iota
	writerStateHeaders
	writerStateBody
	writerStateDone
)

type Writer struct {
	io.Writer
	writerState writerState
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.writerState != writerStateStatusLine {
		return fmt.Errorf("status code can only be written once")
	}
	statusLine := "HTTP/1.1 "
	switch statusCode {
	case StatusCodeOK:
		statusLine += "200 OK\r\n"
	case StatusCodeBadRequest:
		statusLine += "400 Bad Request\r\n"
	case StatusCodeInternalServerError:
		statusLine += "500 Internal Server Error\r\n"
	default:
		statusLine += fmt.Sprintf("%d \r\n", statusCode)
	}
	_, err := w.Write([]byte(statusLine))
	if err != nil {
		return err
	}

	w.writerState = writerStateHeaders
	return nil
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.writerState != writerStateHeaders {
		if w.writerState == writerStateStatusLine {
			return fmt.Errorf("status line must be written before headers")
		}
		return fmt.Errorf("headers can only be written once")
	}
	for key, val := range headers {
		_, err := w.Write([]byte(fmt.Sprintf("%s: %s\r\n", key, val)))
		if err != nil {
			return err
		}
	}
	_, err := w.Write([]byte("\r\n")) // the final CRLF
	if err != nil {
		return err
	}

	w.writerState = writerStateBody
	return nil
}

func (w *Writer) WriteBody(p []byte) error {
	if w.writerState != writerStateBody {
		return fmt.Errorf("status line and headers must be written before body")
	}
	_, err := w.Write(p)
	if err != nil {
		return err
	}
	w.writerState = writerStateDone
	return nil
}
