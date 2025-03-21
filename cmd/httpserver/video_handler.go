package main

import (
	"http-from-tcp/internal/request"
	"http-from-tcp/internal/response"
	"log"
	"os"
)

func handleVideoResponse(w response.Writer, req *request.Request) {
	if req.RequestLine.RequestTarget != "/video" {
		w.WriteStatusLine(response.StatusCodeBadRequest)
		w.WriteHeaders(response.GetDefaultHeaders(0))
		w.WriteBody([]byte{})
		return
	}

  // load video file intom memory:
	b, err := os.ReadFile("./assets/vim.mp4")
	if err != nil {
		log.Printf("Failed to read video file: %v\n", err)
		w.WriteStatusLine(response.StatusCodeInternalServerError)
		w.WriteHeaders(response.GetDefaultHeaders(0))
		w.WriteBody([]byte{})
		return
	}

  // write response:
	w.WriteStatusLine(response.StatusCodeOK)

	headers := response.GetDefaultHeaders(len(b))
	headers.Set("Content-Type", "video/mp4")
	w.WriteHeaders(headers)

	w.WriteBody(b)
}
