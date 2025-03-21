package main

import (
	"http-from-tcp/internal/request"
	"http-from-tcp/internal/response"
	"http-from-tcp/internal/server"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const port = 42069

func requestHandler(w io.Writer, req *request.Request) *server.HandlerError {
	if req.RequestLine.RequestTarget == "/yourproblem" {
		return &server.HandlerError{
			StatusCode: response.StatusCodeBadRequest,
			Message:    "Your problem is not my problem\n",
		}
	}
	if req.RequestLine.RequestTarget == "/myproblem" {
		return &server.HandlerError{
			StatusCode: response.StatusCodeInternalServerError,
			Message:    "Woopsie, my bad\n",
		}
	}
	_, err := w.Write([]byte("All good, frfr\n"))
	if err != nil {
		log.Printf("Failed to write response: %v\n", err)
		return &server.HandlerError{
			StatusCode: response.StatusCodeInternalServerError,
			Message:    "Failed to write response\n",
		}
	}
	return nil
}

func main() {
	server, err := server.Serve(port, requestHandler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
