package models

// Message represents a response message structure
type Message struct {
	Status string `json:"status"`
	Body   string `json:"body"`
}
