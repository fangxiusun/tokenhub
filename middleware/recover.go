package middleware

import (
	"log"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

// Recover is a middleware that recovers from panics
func Recover() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Log the error
				log.Printf("Panic recovered: %v\n%s", err, debug.Stack())

				// Return 500 error
				c.JSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"message": "Internal server error",
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}

