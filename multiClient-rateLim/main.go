package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
	"user-onboarding/middleware"

	"github.com/gin-gonic/gin"
)

type Person struct {
	Name         string `json:"name"`
	MobileNumber string `json:"mobileNumber"`
}

func main() {
	r := gin.Default()

	// Initialize Rate Limiter
	rateLimiter := middleware.NewRateLimiter(3, 30*time.Second, 1) // 3 tokens max, 1 token every 30 seconds

	// Define Routes
	r.POST("/api/v1/onboard", rateLimiter.LimitMiddleware(), msg) // Apply middleware only to this route

	// Start Server
	r.Run(":8081")
}

func msg(c *gin.Context) {
	// Log the body before parsing
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		fmt.Println("Error reading body in handler:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error reading body.",
		})
		return
	}

	var payload Person
	if err := json.Unmarshal(body, &payload); err != nil {
		fmt.Println("Error parsing JSON in handler:", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid or missing mobile number.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User onboarded successfully!",
		"mobile":  payload.MobileNumber,
	})
}
