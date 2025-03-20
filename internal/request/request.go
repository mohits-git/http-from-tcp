package request

import (
	"errors"
	"http-from-tcp/internal/headers"
	"io"
)

type requestState int

const (
	requestStateInitialized requestState = iota
	requestStateParsingHeaders
	requestStateParsingBody
	requestStateDone
)

type RequestLine struct {
	Method        string
	RequestTarget string
	HttpVersion   string
}

type Request struct {
	state       requestState
	RequestLine RequestLine
	Headers     headers.Headers
	Body        []byte
}

const bufferSize = 8

func RequestFromReader(reader io.Reader) (*Request, error) {
	req := &Request{
		state:       requestStateInitialized,
		RequestLine: RequestLine{},
		Headers:     headers.NewHeaders(),
		Body:        []byte{},
	}

	buff := make([]byte, bufferSize, bufferSize)
	readToIndex := 0

	for req.state != requestStateDone {
		if len(buff) <= readToIndex {
			t := make([]byte, (readToIndex+1)*2, (readToIndex+1)*2)
			copy(t, buff)
			buff = t
		}

		// read chunk from buffer
		n, err := reader.Read(buff[readToIndex:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, err
		}
		readToIndex += n

		// parse the chunk
		parsedN, err := req.parse(buff[:readToIndex])
		if err != nil {
			return nil, err
		}

		if parsedN > 0 {
			t := make([]byte, len(buff)-parsedN, len(buff)-parsedN)
			copy(t, buff[parsedN:])
			buff = t
			readToIndex -= parsedN
		}
	}

  // check if the request is done
	if req.state != requestStateDone {
		return nil, errors.New("error: incomplete request")
	}

  // reverfiy body parsing
  contentLength, ok, err := getContentLength(req.Headers)
  if err != nil {
    return nil, err
  }
  if ok {
    _, err := reader.Read(buff)
    if !errors.Is(err, io.EOF) {
      return nil, errors.New("error: body is too long")
    }
    if len(req.Body) != contentLength {
      return nil, errors.New("error: body is too long")
    }
  }

	return req, nil
}
