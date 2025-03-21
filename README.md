## HTTP From TCP

HTTP server from scratch using Go's `net` package.

### Features Implemented

- [x] HTTP Request (`internal/request`)
  - [x] Parse Request Line
  - [x] Parse Headers (`internal/header`)
  - [x] Parse Body
- [x] HTTP Response (`internal/response`)
  - [x] Write Status Line
  - [x] Write Headers
  - [x] Write Body
  - [x] Write Chunked Encoding
- [x] HTTP Server (`internal/server`)
  - [x] Listen on TCP
  - [x] Accept Connections Asyncronously
  - [x] Custom User Handlers
    - [x] Read Requests
    - [x] Write Responses

## Get Started

### internal/server

- `server.Serve`
```go
func Serve(port int, handler Handler) (*Server, error)
```
- `server.Server`
```go
type Server struct {
  Close func() error
}
```
- `server.Handler`
```go
type Handler func(w response.Writer, req *request.Request)
```

### internal/request

- `request.Request`
```go
type Request struct {
	RequestLine RequestLine
	Headers     headers.Headers
	Body        []byte
}
type RequestLine struct {
	Method        string
	RequestTarget string
	HttpVersion   string
}
```

- `request.RequestFromReader`
```go
RequestFromReader(reader io.Reader) (*Request, error)
```

### internal/response

- `response.GetDefaultHeaders`
```go
func GetDefaultHeaders(contentLen int) headers.Headers
```
- `response.StatusCode`
```go
type StatusCode int
const (
	StatusCodeOK                  StatusCode = 200
	StatusCodeBadRequest          StatusCode = 400
	StatusCodeInternalServerError StatusCode = 500
)
```
- `response.Writer`
```go
type Writer interface {
  WriteStatusLine(statusCode StatusCode) error
  WriteHeaders(headers headers.Headers) error
  WriteBody(p []byte) error
  WriteChunked(body []byte)
  WriteChunkedBody(p []byte) (int, error)
  WriteChunkedBodyDone() (int, error)
  WriteTrailers(h headers.Headers) error
}
```

### internal/headers

- `headers.Headers`
```go
type Headers map[string][]string

func (h Headers) Get(key string) string
func (h Headers) Set(key, value string)
func (h Headers) Add(key, value string)
func (h Headers) Delete(key string)
func (h Headers) Parse(data []byte) (n int, done bool, err error)
```
- `headers.NewHeaders`
```go
func NewHeaders() Headers
```

### Example http server using `internal/server`
```go
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
```

### Run the demo server (`cmd/httpserver`)

```bash
go run ./cmd/httpserver
```
