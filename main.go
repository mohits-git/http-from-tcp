package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
  file, err := os.Open("message.txt")
  if err != nil {
    fmt.Println("Error while opening message.txt", err)
  }
  defer file.Close()

  var data = make([]byte, 8)
  currentLine := ""
  for {
    n, err := file.Read(data)
    if err == io.EOF && n == 0 {
      break
    }
    if err != nil {
      fmt.Println("Error while reading message.txt", err)
      break
    }
    parts := strings.Split(string(data[:n]), "\n")
    for i := 0; i < len(parts)-1; i++ {
      fmt.Printf("read: %s\n", currentLine + parts[i])
      currentLine = ""
    }
    currentLine += parts[len(parts)-1]
  }

  if len(currentLine) > 0 {
    fmt.Printf("read: %s\n", currentLine)
  }
}
