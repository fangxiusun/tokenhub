package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RequestId is a middleware that generates a unique request ID
func RequestId() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if request ID already exists
		requestId := c.GetHeader("X-Request-Id")
		if requestId == "" {
			// Generate new request ID
			requestId = uuid.New().String()
		}

		// Set request ID in context
		c.Set("requestId", requestId)

		// Set response header
		c.Header("X-Request-Id", requestId)

		c.Next()
	}
}
