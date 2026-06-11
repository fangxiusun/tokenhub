package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/your-username/your-project/model"
)

// Distribute is a middleware that selects a channel for the request
func Distribute() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get model from request body
		body, err := c.Get("requestBody")
		if err != nil {
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
		if userGroup, exists := c.Get("userGroup"); exists {
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
		if model, exists := bodyMap["model"]; exists {
			if modelName, ok := model.(string); ok {
				return modelName
			}
		}
	}
	return ""
}
