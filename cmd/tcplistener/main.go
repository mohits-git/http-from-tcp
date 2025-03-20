package main

import (
	"fmt"
	"http-from-tcp/internal/request"
	"net"
)

const port = ":42069"

func main() {
	ln, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println("Error while listening: ", err)
		return
	}
	defer ln.Close()
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error while accepting connection: ", err)
			continue
		}
		fmt.Printf("New Connection Accepted.\n * Remote Addr: %s\n", conn.RemoteAddr())

		req, err := request.RequestFromReader(conn)
		if err != nil {
			fmt.Println("Error while reading request: ", err)
		}

		fmt.Printf("Request Line: \n")
		fmt.Printf(" - Method: %s\n", req.RequestLine.Method)
		fmt.Printf(" - Target: %s\n", req.RequestLine.RequestTarget)
		fmt.Printf(" - Version: %s\n", req.RequestLine.HttpVersion)

		// TODO: conn close?
		conn.Close()
		fmt.Println("Connection Closed.")
	}
}
