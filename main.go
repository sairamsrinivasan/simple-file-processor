package main

import (
	"simple-file-processor/internal/server"
)

func main() {
	// Initialize the server with the default configurations
	server := server.NewServer()

	// Start the server
	server.Start()
}
