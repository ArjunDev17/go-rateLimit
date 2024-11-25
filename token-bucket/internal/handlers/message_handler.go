package handlers

import (
	"encoding/json"
	"net/http"

	models "github.com/ArjunDev17/go-rateLimit/pkg"
)

// MessageHandler holds logic for the message endpoints
type MessageHandler struct{}

// NewMessageHandler creates a new instance of MessageHandler
func NewMessageHandler() *MessageHandler {
	return &MessageHandler{}
}

// HandleMessage is the handler function for the message endpoint
func (m *MessageHandler) HandleMessage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	message := models.Message{
		Status: "Successful",
		Body:   "What can I assist you with?",
	}

	err := json.NewEncoder(w).Encode(&message)
	if err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}
