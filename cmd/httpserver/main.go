package main

import (
	"http-from-tcp/internal/request"
	"http-from-tcp/internal/response"
	"http-from-tcp/internal/server"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

const port = 42069

func requestHandler(w response.Writer, req *request.Request) {
	var resBody string
	var statusCode response.StatusCode
	headers := response.GetDefaultHeaders(len(resBody))
	headers.Set("Content-Type", "text/html")

	if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin/") {
		proxyHandler(w, req)
		return
	}

	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		statusCode = response.StatusCodeBadRequest
		resBody = `<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>`
	case "/myproblem":
		statusCode = response.StatusCodeInternalServerError
		resBody = `<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>`
	default:
		statusCode = response.StatusCodeOK
		resBody = `<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>`
	}

	// write response:
	err := w.WriteStatusLine(statusCode)
	if err != nil {
		log.Printf("Failed to write status line: %v\n", err)
		return
	}
	err = w.WriteHeaders(headers)
	if err != nil {
		log.Printf("Failed to write headers: %v\n", err)
		return
	}
	err = w.WriteBody([]byte(resBody))
	if err != nil {
		log.Printf("Failed to write body: %v\n", err)
		return
	}
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
