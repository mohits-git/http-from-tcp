package headers

import (
	"errors"
	"strings"
)

type Headers map[string]string

func NewHeaders() Headers {
	return Headers{}
}

// Parses a single key-value pair from the data and adds it to the Headers map
// Returns the number of bytes parsed, whether the parsing is done, and an error if one occurred
func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	parsedBytes := 0

	crlfIndex := strings.Index(string(data), "\r\n")

	if crlfIndex == -1 {
		return 0, false, nil // not enough data to parse
	}

	if crlfIndex == 0 {
		return 2, true, nil // done parsing
	}

	parsedBytes += crlfIndex + 2 // +1 for LF, +1 for 0-based-indexing

	// parse:
	fieldLine := strings.TrimSpace(string(data[:crlfIndex]))
	key, val, found := strings.Cut(fieldLine, ":") // split key and value
	if !found {
		return 0, false, errors.New("invalid format") // invalid format
	}

	if strings.TrimSpace(key) != key {
		return 0, false, errors.New("key has trailing spaces") // invalid (key) format
	}

  val = strings.TrimSpace(val)
	if len(val) == 0 {
		return 0, false, errors.New("value is empty") // invalid format (empty value)
	}

	h[key] = val // set key-value pair

	return parsedBytes, false, nil
}
