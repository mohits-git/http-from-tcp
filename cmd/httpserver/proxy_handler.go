package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"http-from-tcp/internal/headers"
	"http-from-tcp/internal/request"
	"http-from-tcp/internal/response"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

const bufferSize = 1024

func proxyHandler(w response.Writer, req *request.Request) {
	// proxy request to httpbin.org:
	reqPath := strings.TrimPrefix(req.RequestLine.RequestTarget, "/httpbin/")
	url := "https://httpbin.org/" + reqPath
	res, err := http.Get(url)
	if err != nil {
		w.WriteStatusLine(response.StatusCodeBadRequest)
		w.WriteHeaders(response.GetDefaultHeaders(0))
		return
	}
	defer res.Body.Close()

	// write status line:
	err = w.WriteStatusLine(response.StatusCodeOK)
	if err != nil {
		return
	}

	// write headers:
	h := response.GetDefaultHeaders(0)
	for k, v := range res.Header { // copy headers from res to headers
		h.Set(k, v[0])
	}
	h.Delete("Content-Length")            // remove Content-Length header
	h.Set("Transfer-Encoding", "chunked") // add Transfer-Encoding header
	h.Add("Trailer", "X-Content-SHA256")  // add Trailer headers:
	h.Add("Trailer", "X-Content-Length")
	err = w.WriteHeaders(h)
	if err != nil {
		return
	}

	// write chunked body:
	body := make([]byte, 0)
	buf := make([]byte, bufferSize)
	for {
		n, err := res.Body.Read(buf)
		if err != nil {
			if errors.Is(err, io.EOF) {
				w.WriteChunkedBodyDone() // write last chunked body 0\r\n
			}
			break
		}
		log.Printf("Read: %d bytes\n", n)
		// write chunked body:
		_, err = w.WriteChunkedBody(buf[:n])
		if err != nil {
			break
		}
		body = append(body, buf[:n]...)
	}

	// write trailers:
	contentLength := strconv.Itoa(len(body))
	hash := sha256.Sum256(body)
	hashStr := hex.EncodeToString(hash[:])

	trailers := headers.NewHeaders()
	trailers.Set("X-Content-SHA256", hashStr)
	trailers.Set("X-Content-Length", contentLength)
	w.WriteTrailers(trailers)
}
