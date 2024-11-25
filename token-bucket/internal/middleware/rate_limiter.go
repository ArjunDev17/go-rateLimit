package middleware

import (
	"encoding/json"
	"net/http"

	models "github.com/ArjunDev17/go-rateLimit/pkg"
	"golang.org/x/time/rate"
)

// RateLimiter is the middleware function for rate limiting
func RateLimiter(next http.HandlerFunc) http.HandlerFunc {
	limiter := rate.NewLimiter(2, 4) // 2 requests per second, burst capacity of 4
	return func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			// If rate limit exceeded, send 429 response
			message := models.Message{
				Status: "Request Failed",
				Body:   "The API is at capacity, try again later.",
			}
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(message)
			return
		}
		// Call the next handler if allowed
		next(w, r)
	}
}
