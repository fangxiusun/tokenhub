package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RelayNotFound handles 404 for relay routes
func RelayNotFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{
		"error": gin.H{
			"message": "The API endpoint you requested does not exist",
			"type":    "not_found",
		},
	})
}
