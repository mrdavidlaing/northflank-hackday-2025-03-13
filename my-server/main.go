package main

import (
	"fmt"
	"log"
	"net/http"
)

const (
	defaultPort = "8080"
)

func main() {
	// Register handlers
	http.HandleFunc("/info", infoHandler)

	// Start the server
	port := defaultPort
	fmt.Printf("Server starting on port %s...\n", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
} 