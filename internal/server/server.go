package server

import (
	"fmt"
	"net"

	"github.com/DNahar74/my-redis/internal/command"
	"github.com/DNahar74/my-redis/internal/store"
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
	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			fmt.Println("Error closing the listener:", err)
		}
	}(listener)

	fmt.Println("Server listening on", s.address)

	var RedisStore = store.CreateStorage()
	command.InitStore(RedisStore)

	err = restoreStorage()
	if err != nil {
		return err
	}

	go handleAOF(RedisStore)
	go handleMemoryState(RedisStore)

	// Allow multiple connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error connecting the client", err)
			return err
		}

		// make a goroutine for handling R/W
		go handleConnection(conn, RedisStore)
	}
}
