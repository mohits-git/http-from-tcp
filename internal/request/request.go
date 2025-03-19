package request

import (
	"io"
)

type parserState int

const (
	initialized parserState = iota
	done
)

type RequestLine struct {
	Method        string
	RequestTarget string
	HttpVersion   string
}

type Request struct {
	state       parserState
	RequestLine RequestLine
	// Headers     map[string]string
	// Body        string
}

const bufferSize = 8

func RequestFromReader(reader io.Reader) (*Request, error) {
	req := &Request{
		RequestLine: RequestLine{},
		state:       initialized,
	}

	buff := make([]byte, bufferSize, bufferSize)
	readToIndex := 0
	// read buffSize bytes
	for req.state != done {
		if len(buff) <= readToIndex {
			t := make([]byte, (readToIndex+1)*2, (readToIndex+1)*2)
			copy(t, buff)
			buff = t
		}

		n, err := reader.Read(buff[readToIndex:])
		if err == io.EOF {
			req.state = done
			break
		}
		if err != nil {
			return nil, err
		}
		readToIndex += n

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

	return req, nil
}
