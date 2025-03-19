package request

import (
	"errors"
	"io"
	"strings"
)

type RequestLine struct {
	Method        string
	RequestTarget string
	HttpVersion   string
}

type Request struct {
	RequestLine RequestLine
	// Headers     map[string]string
	// Body        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
  // read entire bytes from the reader
  reqBytes, err := io.ReadAll(reader)
  if err != nil {
    return nil, err
  }
  reqStr := string(reqBytes)

  // Split the request data into lines
  lines := strings.Split(reqStr, "\r\n")
  if len(lines) < 1 {
    return nil, errors.New("Request is empty")
  }
  reqLineStr := lines[0]

  // Split the request line into (3) parts
  parts := strings.Split(reqLineStr, " ")
  if len(parts) != 3 {
    return nil, errors.New("Invalid number of parts in request line")
  }

  // verify method part is all CAPS
  if strings.ToUpper(parts[0]) != parts[0] {
    return nil, errors.New("Invalid Method")
  }

  // verify the HTTP version
  if parts[2] != "HTTP/1.1" {
    return nil, errors.New("Only HTTP/1.1 is supported")
  }

  reqLine := RequestLine{
    Method: parts[0],
    RequestTarget: parts[1],
    HttpVersion: "1.1",
  }
  req := &Request{
    RequestLine: reqLine,
  }

  return req, nil
}
