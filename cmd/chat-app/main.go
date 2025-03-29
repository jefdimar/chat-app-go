package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/jefdimar/go-chat-app/internal/handlers"
	"github.com/jefdimar/go-chat-app/internal/server"
)

func main() {
	// Initialize and start the server
	s := server.NewServer()

	// Start handling messages in a goroutine
	go handlers.HandleMessages()

	// Start the server
	fmt.Println("Starting chat server on :9090")
	log.Fatal(http.ListenAndServe(":9090", s.Router))
}
