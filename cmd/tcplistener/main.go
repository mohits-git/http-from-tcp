package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
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
		fmt.Println("Reading data...")
		lines := getLinesChannel(conn)
		for line := range lines {
			fmt.Println(line)
		}
		fmt.Println("Connection Closed.")
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	lines := make(chan string)

	go func() {
		defer f.Close()
		defer close(lines)
		currentLine := ""
		for {
			var data = make([]byte, 8)
			n, err := f.Read(data)
			if err != nil {
				if currentLine > "" {
					lines <- currentLine
				}
				if errors.Is(err, io.EOF) {
					return
				}
				fmt.Println("Error while reading file: ", err)
				return
			}
			parts := strings.Split(string(data[:n]), "\n")
			for i := 0; i < len(parts)-1; i++ {
				lines <- currentLine + parts[i]
				currentLine = ""
			}
			currentLine += parts[len(parts)-1]
		}
	}()

	return lines
}
