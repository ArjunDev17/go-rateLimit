package main

import (
	"encoding/json"
	"net/http"

	"github.com/ArjunDev17/go-rateLimit/internal/handlers"
	"github.com/ArjunDev17/go-rateLimit/internal/middleware"
)

// Define the Message struct
type Message struct {
	Status string `json:"status"`
	Body   string `json:"body"`
}

// Endpoint handler function
func endPointHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	message := Message{
		Status: "Successful",
		Body:   "What can I assist you with?",
	}
	err := json.NewEncoder(w).Encode(&message)
	if err != nil {
		return
	}
}

func main() {
	// Create a new message handler instance
	messageHandler := handlers.NewMessageHandler()
	// Set up routes with middleware
	http.HandleFunc("/hello", middleware.RateLimiter(messageHandler.HandleMessage))
	http.HandleFunc("/hi", messageHandler.HandleMessage)

	// Start the server on port 8080
	http.ListenAndServe(":8080", nil)
}
