//TODO: (1) Handle errors by continuing to next loop and returning valid error to client

package server

import (
	"fmt"
	"net"

	"github.com/DNahar74/my-redis/command"
	"github.com/DNahar74/my-redis/resp"
	"github.com/DNahar74/my-redis/utils"
)

// handleConnection takes the connecton request for a client and handles the input and output
func handleConnection(conn net.Conn) {
	fmt.Println("Client connected")
	fmt.Println("address:", conn.RemoteAddr().String())
	fmt.Println("")
	defer conn.Close()

	buf := make([]byte, 128)

	// keep reading the input until the client disconnects
	for {
		n, err := conn.Read(buf)
		if err != nil {
			// EOF can be used to find if the user disconnected
			if err.Error() == "EOF" {
				fmt.Println("Client Disconnected:", conn.RemoteAddr().String())
				message := resp.BulkString{Value: "DISCONNECTED"}
				err = utils.SendMessage(conn, message)
				if err != nil {
					fmt.Println("Error sending response: ", err.Error())
					return
				}
				return
			}

			fmt.Println("Error reading from connection: ", err.Error())
		}

		cmd := string(buf[:n])

		fmt.Println("Input:", cmd)

		commands, err := resp.Deserialize(cmd)
		if err != nil {
			fmt.Println("Error deserializing commands: ", err.Error())
			return
		}

		val, err := command.HandleCommands(commands)
		if err != nil {
			fmt.Println("Error handling commands: ", err.Error())
			return
		}

		err = utils.SendMessage(conn, val)
		if err != nil {
			fmt.Println("Error sending response: ", err.Error())
			return
		}
	}
}
