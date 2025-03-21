package response

import (
	"fmt"
	"http-from-tcp/internal/headers"
	"io"
)

type StatusCode int

const (
	StatusCodeOK                  StatusCode = 200
	StatusCodeBadRequest          StatusCode = 400
	StatusCodeInternalServerError StatusCode = 500
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
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
	return err
}

func GetDefaultHeaders(contentLen int) headers.Headers {
  h := headers.NewHeaders()
  h.Set("Content-Length", fmt.Sprintf("%d", contentLen))
  h.Set("Connection", "close")
  h.Set("Content-Type", "text/plain")
  return h
  // other default headers can be:
  // h.Set("Content-Encoding", "gzip") // or "deflate"
  // h.Set("Date", "Mon, 27 Jul 2020 12:28:53 GMT") // or time.Now().UTC().Format(time.RFC1123)
  // h.Set("Cache-Control", "no-cache") // or "max-age=3600"
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
  for key, val := range headers {
    _, err := w.Write([]byte(fmt.Sprintf("%s: %s\r\n", key, val)))
    if err != nil {
      return err
    }
  }
  _, err := w.Write([]byte("\r\n")) // the final CRLF
  return err
}
