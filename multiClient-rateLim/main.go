package main

import (
	"time"
	"user-onboarding/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Initialize Rate Limiter
	rateLimiter := middleware.NewRateLimiter(3, 30*time.Second, 1) // 3 tokens max, 1 token every 30 seconds

	// Apply Rate Limiter Middleware
	r.Use(rateLimiter.LimitMiddleware())

	// Define Routes
	r.POST("/api/v1/onboard", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "User onboarded successfully!",
		})
	})

	// Start Server
	r.Run(":8081")
}
