package server

import (
	"fmt"
	"net"

	"github.com/DNahar74/my-redis/command"
	"github.com/DNahar74/my-redis/store"
)

// Server represents a Redis server configurations
type Server struct {
	address string
}

// NewServer creates a new Server object
func NewServer(address string) *Server {
	return &Server{address: address}
}

// Start starts the Redis server
func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		fmt.Println("Error starting the listener:", err)
		return err
	}
	defer listener.Close()

	fmt.Println("Server listening on port 6379")

	var RedisStore = store.CreateStorage()
	command.InitStore(RedisStore)

	// Allow multiple connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error connecting the client", err)
			return err
		}

		// make a goroutine for handling R/W
		go handleConnection(conn)
	}
}
