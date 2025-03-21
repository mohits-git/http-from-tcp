package response

import (
	"http-from-tcp/internal/headers"
	"strconv"
)

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	n := len(p)
	chunksize := strconv.FormatInt(int64(n), 16)
	_, err := w.Write([]byte(chunksize + "\r\n" + string(p) + "\r\n"))
	if err != nil {
		return 0, err
	}
	return n, nil
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	_, err := w.Write([]byte("0\r\n"))
	return 0, err
}

// WriteHeaders writes headers to the response writer (with last CRLF).
func (w *Writer) WriteTrailers(h headers.Headers) error {
  headersStr := ""
  for k, v := range h {
    headersStr += k + ": " + v + "\r\n"
  }
  _, err := w.Write([]byte(headersStr + "\r\n")) // write with last CRLF
  return err
}
