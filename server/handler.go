package server

import (
	"fmt"
	"net"

	"github.com/DNahar74/my-redis/resp"
	"github.com/DNahar74/my-redis/utils"
)

// handleConnection takes the connecton request for a client and handles the input and output
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
				err = utils.SendMessage(conn, "+DISCONNECTED\r\n")
				return
			}

			fmt.Println("Error reading from connection: ", err.Error())
		}

		commands := resp.GetCommands(buf)

		if commands[2] == "PING" {
			utils.SendMessage(conn, "+PONG\r\n")
		} else if commands[2] == "COMMAND" {
			// Respond to the COMMAND DOCS request
			utils.SendMessage(conn, "-ERR Unknown command\r\n")
		} else {
			// Handle unknown commands
			utils.SendMessage(conn, "-ERR Unknown command\r\n")
		}
	}
}
