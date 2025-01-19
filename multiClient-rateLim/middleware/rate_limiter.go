package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// TokenBucket represents a user's rate-limiting bucket
type TokenBucket struct {
	Tokens         int       // Current number of tokens
	LastRefillTime time.Time // Last time tokens were refilled
}

// RateLimiter manages token buckets for users
type RateLimiter struct {
	UserBuckets map[string]*TokenBucket
	mu          sync.Mutex
	Capacity    int           // Maximum tokens in the bucket
	RefillRate  time.Duration // Time duration to add tokens
	RefillCount int           // Tokens to refill per interval
}

// NewRateLimiter initializes the rate limiter
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

		userID := c.ClientIP() // Use IP as user identifier (replace with Auth ID in production)

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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
