//! (1) When a client disconnects it shows :: {Error reading from connection:  EOF}
//TODO: (1) Allow multiple clients to connect

package main

import (
	"fmt"
	"net"
	"strings"
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

	conn, err := listener.Accept()
	if err != nil {
		fmt.Println("Error connecting the client", err)
		return
	}

	fmt.Println("Client connected")

	buf := make([]byte, 128)

	for {
		_, err = conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading from connection: ", err.Error())
			return
		}
		defer conn.Close()
	
		fmt.Printf("read command:\n%s", buf)
	
		commands := getCommands(buf)
	
		if commands[2] == "PING" {
			message := []byte("+PONG\r\n")
			_, err := conn.Write(message)
			if err != nil {
				fmt.Printf("Error writing to client: %v\n", err)
			}
		}
	}
}

func getCommands(buf []byte) []string {
	val := string(buf)
	val = strings.ToUpper(val)
	
	inputs := strings.Split(val, "\r\n")

	return inputs
}