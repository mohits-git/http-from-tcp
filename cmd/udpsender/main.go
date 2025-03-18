package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

const adr = "localhost:42069"

func main() {
  // Resolve the udp address
  udpAddr, err := net.ResolveUDPAddr("udp", adr)
  if err != nil {
    fmt.Println("Error resolving UDP address", err)
    return
  }

  // Dial the remote udp address
  conn, err := net.DialUDP("udp", nil, udpAddr)
  if err != nil {
    fmt.Println("Error dialing UDP address", err)
    return
  }
  defer conn.Close()

  // Read from stdin and write to the udp connection
  rdr := bufio.NewReader(os.Stdin)
  for {
    fmt.Print("> ")
    // read line from stdin
    inputLine, err := rdr.ReadString('\n')
    if err != nil {
      fmt.Println("Error reading input", err)
      continue
    }
    // write line to udp connection
    _, err = conn.Write([]byte(inputLine))
    if err != nil {
      fmt.Println("Error writing to UDP", err)
    }
  }
}
