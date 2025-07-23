package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/DNahar74/my-redis/internal/server"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	var (
		addr    = flag.String("addr", "0.0.0.0:6379", "Server address to bind to")
		help    = flag.Bool("help", false, "Show help information")
		ver     = flag.Bool("version", false, "Show version information")
		verbose = flag.Bool("verbose", false, "Enable verbose logging")
	)
	flag.Parse()

	if *help {
		flag.Usage()
		return
	}

	if *ver {
		fmt.Printf("my-redis %s (commit: %s, built: %s)\n", version, commit, date)
		return
	}

	if *verbose {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
	}

	fmt.Printf("Starting my-redis server on %s\n", *addr)
	fmt.Printf("Version: %s\n", version)

	redisServer := server.NewServer(*addr)

	// Graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := redisServer.Start(); err != nil {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	<-c
	fmt.Println("\nShutting down server...")
	// Add graceful shutdown logic here
}
