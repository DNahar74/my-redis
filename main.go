package main

import (
	"fmt"

	"github.com/DNahar74/my-redis/server"
)

func main() {
	fmt.Println("This is the start of the main REDIS program")

	redisServer := server.NewServer("0.0.0.0:6379")
	err := redisServer.Start()
	if err != nil {
		fmt.Println("Error starting the REDIS server:", err)
		return
	}
}
