package headers

import (
	"errors"
	"regexp"
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

  // parse key:
  key = strings.ToLower(key)
	if len(key) == 0 {
		return 0, false, errors.New("key is empty") // invalid format (empty key)
	}
	if strings.TrimSpace(key) != key {
		return 0, false, errors.New("key has trailing spaces") // invalid format (space in between key and colon)
	}
  // regex to check if key is valid
  rgx := regexp.MustCompile(`^[a-zA-Z0-9!#$%&'*+-.^_` + "`" + `|~]+$`)
  rgxMatch := rgx.FindString(key)
  if rgxMatch != key {
    return 0, false, errors.New("key is invalid") // invalid format (invalid key)
  }

  // parse value:
	val = strings.TrimSpace(val)
	if len(val) == 0 {
		return 0, false, errors.New("value is empty") // invalid format (empty value)
	}

	h[key] = val // set key-value pair

	return parsedBytes, false, nil
}
