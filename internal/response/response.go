package response

import (
	"fmt"
	"http-from-tcp/internal/headers"
)

type StatusCode int

const (
	StatusCodeOK                  StatusCode = 200
	StatusCodeBadRequest          StatusCode = 400
	StatusCodeInternalServerError StatusCode = 500
)

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
