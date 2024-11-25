package main

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// Message struct for JSON responses
type Message struct {
	Status string `json:"status"`
	Body   string `json:"body"`
}

// perClientRateLimt is a middleware that applies rate limiting per client IP
func perClientRateLimt(next http.Handler) http.Handler {
	// Define a client struct to store limiter and last seen time
	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}

	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)

	// Clean up old clients every minute
	go func() {
		for {
			time.Sleep(time.Minute)
			mu.Lock()
			for ip, client := range clients {
				// Remove clients who haven't been seen in the last 3 minutes
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}
	}()

	// Return the rate-limiting middleware
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the client's IP address
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			// If there's an error parsing the IP, return an internal error
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Check if the client already exists, if not, create a new entry
		mu.Lock()
		if _, found := clients[ip]; !found {
			clients[ip] = &client{
				limiter: rate.NewLimiter(2, 4), // 2 requests per second, burst capacity of 4
			}
		}
		clients[ip].lastSeen = time.Now()
		mu.Unlock()

		// Check if the rate limit is exceeded for this client
		if !clients[ip].limiter.Allow() {
			// If rate limit exceeded, return a 429 response with message
			message := Message{
				Status: "Request Failed",
				Body:   "The API is at capacity, try again later",
			}
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(&message)
			return
		}

		// If the request is allowed, pass it to the next handler
		next.ServeHTTP(w, r)
	})
}

// endpointHandler responds to requests with a JSON message
func endpointHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	message := Message{
		Status: "Successful",
		Body:   "Hi! You have reached the API. How may I help you?",
	}

	// Encode the message to JSON and send it in the response
	err := json.NewEncoder(w).Encode(&message)
	if err != nil {
		// Log the error and return a 500 internal server error if encoding fails
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func main() {
	// Set up the server with the rate limiter middleware
	http.Handle("/ping", perClientRateLimt(http.HandlerFunc(endpointHandler)))

	// Start the server
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
