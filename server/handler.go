package server

import (
	"fmt"
	"net"

	"github.com/DNahar74/my-redis/command"
	"github.com/DNahar74/my-redis/resp"
	"github.com/DNahar74/my-redis/store"
	"github.com/DNahar74/my-redis/utils"
)

// RedisStore is a global variable to hold the Redis data store

// handleConnection takes the connecton request for a client and handles the input and output
func handleConnection(conn net.Conn) {
	fmt.Println("Client connected")
	fmt.Println("address:", conn.RemoteAddr().String())
	fmt.Println("")
	defer conn.Close()
	
	buf := make([]byte, 128)
	
	var RedisStore = store.CreateStorage()
	command.InitStore(RedisStore)

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
			m := resp.SimpleError{Value: err.Error()}
			utils.SendMessage(conn, m)
			continue
		}

		val, err := command.HandleCommands(commands)
		if err != nil {
			fmt.Println("Error handling commands: ", err.Error())
			m := resp.SimpleError{Value: err.Error()}
			utils.SendMessage(conn, m)
			continue
		}

		err = utils.SendMessage(conn, val)
		if err != nil {
			fmt.Println("Error sending response: ", err.Error())
			continue
		}
	}
}
