package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeadersParse(t *testing.T) {
	// Test: Valid single header
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["Host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

  // Test: Valid single header with extra spaces
	headers = NewHeaders()
	data = []byte("         Host:       localhost:42069  \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["Host"])
	assert.Equal(t, 40, n)
	assert.False(t, done)

  // Test: Valid 2 headers with existing headers
  headers = NewHeaders()
  headers["Host"] = "localhost:42069"
  data = []byte("User-Agent: curl/7.68.0\r\n\r\n")
  n, done, err = headers.Parse(data)
  require.NoError(t, err)
  require.NotNil(t, headers)
  assert.Equal(t, "localhost:42069", headers["Host"])
  assert.Equal(t, "curl/7.68.0", headers["User-Agent"])
  assert.Equal(t, 25, n)
  assert.False(t, done)

  // Test: Valid done
  headers = NewHeaders()
  data = []byte("\r\n")
  n, done, err = headers.Parse(data)
  require.NoError(t, err)
  require.NotNil(t, headers)
  assert.Equal(t, 2, n)
  assert.True(t, done)

  // Test: Invalid space header
  headers = NewHeaders()
  data = []byte("Host : localhost:42069\r\n\r\n")
  n, done, err = headers.Parse(data)
  require.Error(t, err)
  assert.Equal(t, 0, n)
  assert.False(t, done)

  // Test: Invalid empty value header
  headers = NewHeaders()
  data = []byte("Host: \r\n\r\n")
  n, done, err = headers.Parse(data)
  require.Error(t, err)
  assert.Equal(t, 0, n)
  assert.False(t, done)
}
