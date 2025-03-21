//TODO: (✔️fixed) Understand this :: (The issue is that we are not handling this command properly bcoz. of which the client never moves forward)
//? read command:
//? *2
//? $7
//? COMMAND
//? $4
//? DOCS

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

	// Allow multiple connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error connecting the client", err)
			return
		}

		// make a goroutine for handling R/W
		go handleConnection(conn)
	}
}

// handleConnection for each client in a seperate goroutine
func handleConnection(conn net.Conn) {
	fmt.Println("Client connected")
	fmt.Println("address:", conn.RemoteAddr().String())
	defer conn.Close()

	buf := make([]byte, 128)

	// keep reading the input until the client disconnects
	for {
		_, err := conn.Read(buf)
		if err != nil {
			// EOF can be used to find if the user disconnected
			if err.Error() == "EOF" {
				fmt.Println("Client Disconnected:", conn.RemoteAddr().String())
				err = sendMessage(conn, "+DISCONNECTED\r\n")
				return
			}

			fmt.Println("Error reading from connection: ", err.Error())
		}

		fmt.Printf("read command:\n%s", buf)

		commands := getCommands(buf)

		if commands[2] == "PING" {
			sendMessage(conn, "+PONG\r\n")
		} else if commands[2] == "COMMAND" {
			// Respond to the COMMAND DOCS request
			sendMessage(conn, "-ERR Unknown command\r\n")
		} else {
			// Handle unknown commands
			sendMessage(conn, "-ERR Unknown command\r\n")
		}
	}
}

func getCommands(buf []byte) []string {
	val := string(buf)
	val = strings.ToUpper(val)

	inputs := strings.Split(val, "\r\n")

	return inputs
}

func sendMessage(conn net.Conn, message string) error {
	response := []byte(message)
	_, err := conn.Write(response)
	if err != nil {
		fmt.Printf("Error writing [message: %s] to client: %v\n", message, err)
		return err
	}

	return nil
}
