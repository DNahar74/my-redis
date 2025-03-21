package utils

import (
	"fmt"
	"net"

	"github.com/DNahar74/my-redis/resp"
)

// SendMessage takes connection instance of a client, an RESP string message and sends it to the client
// It prints the error and returns it if any error arises
func SendMessage(conn net.Conn, message resp.RESPType) error {
	response, err := message.Serialize()
	if err != nil {
		fmt.Printf("Error serializing [message: %s]: %v\n", message, err)
		return err
	}
	_, err = conn.Write([]byte(response))
	if err != nil {
		fmt.Printf("Error writing [message: %s] to client: %v\n", message, err)
		return err
	}

	fmt.Println("Message Sent:", response)

	return nil
}
