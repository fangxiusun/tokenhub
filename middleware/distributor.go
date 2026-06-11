package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/fangxiusun/tokenhub/model"
)

// ReadRequestBody reads the request body and stores it in context
func ReadRequestBody() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Only read body for POST/PUT/PATCH requests
		if c.Request.Method != "POST" && c.Request.Method != "PUT" && c.Request.Method != "PATCH" {
			c.Next()
			return
		}

		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Failed to read request body",
			})
			c.Abort()
			return
		}
		// Restore the body for downstream handlers
		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

		// Parse and store in context
		var bodyMap map[string]interface{}
		if err := json.Unmarshal(body, &bodyMap); err == nil {
			c.Set("requestBody", bodyMap)
		}

		c.Next()
	}
}

// Distribute is a middleware that selects a channel for the request
func Distribute() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get model from request body
		body, exists := c.Get("requestBody")
		if !exists {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Request body not found",
			})
			c.Abort()
			return
		}

		modelName := extractModelName(body)
		if modelName == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Model name not found in request",
			})
			c.Abort()
			return
		}

		// Get user group
		group := "default"
		if userGroup, groupExists := c.Get("userGroup"); groupExists {
			if g, ok := userGroup.(string); ok && g != "" {
				group = g
			}
		}

		// Find a suitable channel
		channel, err := model.GetRandomSatisfiedChannel(group, modelName)
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"success": false,
				"message": "No available channel for model: " + modelName,
			})
			c.Abort()
			return
		}

		// Set channel in context
		c.Set("channelId", channel.Id)
		c.Set("channelName", channel.Name)
		c.Set("channelType", channel.Type)
		c.Set("channelBaseURL", channel.BaseURL)
		c.Set("channelKey", channel.Key)

		c.Next()
	}
}

func extractModelName(body interface{}) string {
	if bodyMap, ok := body.(map[string]interface{}); ok {
		if modelVal, exists := bodyMap["model"]; exists {
			if modelName, ok := modelVal.(string); ok {
				return modelName
			}
		}
	}
	return ""
}

