package request

import (
	"errors"
	"strings"
)

// Parse the request from the data buffer
func (r *Request) parse(data []byte) (int, error) {
	totalBytesParsed := 0
	for r.state != requestStateDone {
    n, err := r.parseSingle(data[totalBytesParsed:])
		if err != nil {
			return totalBytesParsed, err
		}
		if n == 0 {
			break
		}
		totalBytesParsed += n
	}
	return totalBytesParsed, nil
}

// parseSingle parses a single peice of data (request line, single header, etc)
func (r *Request) parseSingle(data []byte) (int, error) {
	switch r.state {
	case requestStateInitialized:
		rl, n, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		if n > 0 {
			r.RequestLine = rl
			r.state = requestStateParsingHeaders
		}
		return n, nil
	case requestStateParsingHeaders:
		n, d, err := r.Headers.Parse(data)
		if err != nil {
			return 0, err
		}
		if d {
			r.state = requestStateDone
		}
		return n, nil
	case requestStateDone:
		return 0, errors.New("error: trying to read data in a done state")
	default:
		return 0, errors.New("error: unknown state")
	}
}

// parse the request line from the data buffer
// returns RequestLine struct (if parsed), number of bytes parsed (0 if not enough data yet), and an error if one occurred
func parseRequestLine(data []byte) (RequestLine, int, error) {
	reqChunkStr := string(data)
	// Split the request data into lines
	lines := strings.Split(reqChunkStr, "\r\n")
	if len(lines) < 2 {
		// No lines to parse yet
		return RequestLine{}, 0, nil
	}
	reqLineStr := lines[0]

	// Split the request line into (3) parts
	parts := strings.Split(reqLineStr, " ")
	if len(parts) != 3 {
		return RequestLine{}, 0, errors.New("Invalid number of parts in request line")
	}

	// Verify method part is all CAPS
	if strings.ToUpper(parts[0]) != parts[0] {
		return RequestLine{}, 0, errors.New("Invalid Method")
	}

	// Verify the HTTP version
	if parts[2] != "HTTP/1.1" {
		return RequestLine{}, 0, errors.New("Only HTTP/1.1 is supported")
	}

	reqLine := RequestLine{
		Method:        parts[0],
		RequestTarget: parts[1],
		HttpVersion:   "1.1",
	}
  return reqLine, len(reqLineStr) + 2, nil // +2 for CRLF :) :) :) :) :)
}
