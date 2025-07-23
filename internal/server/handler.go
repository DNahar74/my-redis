package server

import (
	"fmt"
	"net"

	"github.com/DNahar74/my-redis/internal/command"
	"github.com/DNahar74/my-redis/internal/resp"
	"github.com/DNahar74/my-redis/internal/store"
	"github.com/DNahar74/my-redis/internal/utils"
)

// RedisStore is a global variable to hold the Redis data store

// handleConnection takes the connection request for a client and handles the input and output
func handleConnection(conn net.Conn, s *store.Store) {
	fmt.Println("Client connected")
	fmt.Println("address:", conn.RemoteAddr().String())
	fmt.Println("")
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			fmt.Println("Error closing the connection:", err)
		}
	}(conn)

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
			m := resp.SimpleError{Value: err.Error()}
			err := utils.SendMessage(conn, m)
			if err != nil {
				fmt.Println("Error sending response: ", err.Error())
				return
			}
			continue
		}

		val, err := command.HandleCommands(commands)
		if err != nil {
			fmt.Println("Error handling commands: ", err.Error())
			m := resp.SimpleError{Value: err.Error()}
			err := utils.SendMessage(conn, m)
			if err != nil {
				fmt.Println("Error sending response: ", err.Error())
				return
			}
			continue
		}

		s.AOFChan <- cmd

		err = utils.SendMessage(conn, val)
		if err != nil {
			fmt.Println("Error sending response: ", err.Error())
			continue
		}
	}
}
