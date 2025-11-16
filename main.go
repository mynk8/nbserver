package main

import (
	"log"
	"net/http"

	"mynk8/nbserver/internal/websockets"
)

func main() {
	terminal := websockets.NewPTYTerminal()
	sessionManager := websockets.NewInMemorySessionManager()
	
	handler := websockets.NewTerminalHandler(
		terminal,
		sessionManager,
		"bash", // Can be changed to any shell command
	)

	// HTTP handler for websocket connections
	http.HandleFunc("/terminal", func(w http.ResponseWriter, r *http.Request) {
		if err := handler.Connect(w, r, nil); err != nil {
			log.Printf("Error handling terminal connection: %v", err)
		}
	})

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
