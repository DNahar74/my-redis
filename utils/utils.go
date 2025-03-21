package utils

import (
	"fmt"
	"net"
)

// SendMessage takes connection instance of a client, an RESP string message and sends it to the client
// It prints the error and returns it if any error arises
func SendMessage(conn net.Conn, message string) error {
	response := []byte(message)
	_, err := conn.Write(response)
	if err != nil {
		fmt.Printf("Error writing [message: %s] to client: %v\n", message, err)
		return err
	}

	return nil
}
