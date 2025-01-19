package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// Person struct for extracting request body data
type Person struct {
	Name         string `json:"name"`
	MobileNumber string `json:"mobileNumber"`
}

// TokenBucket struct for rate-limiting
type TokenBucket struct {
	Tokens         int       // Current number of tokens
	LastRefillTime time.Time // Last time tokens were refilled
}

// RateLimiter struct to manage rate-limiting
type RateLimiter struct {
	UserBuckets map[string]*TokenBucket
	mu          sync.Mutex
	Capacity    int           // Maximum tokens in the bucket
	RefillRate  time.Duration // Time duration to add tokens
	RefillCount int           // Tokens to refill per interval
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(capacity int, refillRate time.Duration, refillCount int) *RateLimiter {
	return &RateLimiter{
		UserBuckets: make(map[string]*TokenBucket),
		Capacity:    capacity,
		RefillRate:  refillRate,
		RefillCount: refillCount,
	}
}

// LimitMiddleware applies rate-limiting
func (rl *RateLimiter) LimitMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {
		if c.Request.Method != http.MethodPost {
			c.Next() // Only limit POST requests
			return
		}

		// Read the request body
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			fmt.Printf("%v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid request body.",
			})
			c.Abort()
			return
		}

		// Rewind the body so it can be read by the handler
		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

		// Proceed with rate-limiting logic
		var payload Person
		if err := json.Unmarshal(body, &payload); err != nil {
			fmt.Printf("Error parsing JSON: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid or missing mobile number.",
			})
			c.Abort()
			return
		}

		if payload.MobileNumber == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Mobile number is required.",
			})
			c.Abort()
			return
		}

		// Use the mobile number as the user identifier for rate-limiting
		userID := payload.MobileNumber

		rl.mu.Lock()
		defer rl.mu.Unlock()

		// Get or initialize the user's token bucket
		bucket, exists := rl.UserBuckets[userID]
		if !exists {
			bucket = &TokenBucket{
				Tokens:         rl.Capacity,
				LastRefillTime: time.Now(),
			}
			rl.UserBuckets[userID] = bucket
		}

		// Refill tokens based on time elapsed
		now := time.Now()
		elapsed := now.Sub(bucket.LastRefillTime)
		if elapsed >= rl.RefillRate {
			refillTokens := int(elapsed/rl.RefillRate) * rl.RefillCount
			bucket.Tokens = min(bucket.Tokens+refillTokens, rl.Capacity)
			bucket.LastRefillTime = now
		}

		// Check if tokens are available
		if bucket.Tokens > 0 {
			bucket.Tokens-- // Consume one token
			c.Next()        // Allow the request
		} else {
			// Reject request with 429 Too Many Requests
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded. Try again later.",
			})
		}
	}
}

// Helper function to return the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
