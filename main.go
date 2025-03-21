package main

import (
	"fmt"
	"net"
)

func main() {
	fmt.Println("This is the start of the main REDIS program")

	listener, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
    fmt.Println("Error starting the listener:", err)
    return
  }
	defer listener.Close()

	fmt.Println("Server listening on port 6379")
}