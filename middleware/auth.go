package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"github.com/your-username/your-project/common"
	"github.com/your-username/your-project/model"
)

// UserAuth is a middleware that requires user authentication
func UserAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		username := session.Get("username")
		role := session.Get("role")
		id := session.Get("id")
		status := session.Get("status")

		if username == nil {
			// Check access token
			accessToken := c.Request.Header.Get("Authorization")
			if accessToken == "" {
				c.JSON(http.StatusUnauthorized, gin.H{
					"success": false,
					"message": "Not logged in",
				})
				c.Abort()
				return
			}

			// Validate access token
			user, err := model.GetUserByUsername(accessToken)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{
					"success": false,
					"message": "Invalid access token",
				})
				c.Abort()
				return
			}

			username = user.Username
			role = user.Role
			id = user.Id
			status = user.Status
		}

		// Check user status
		if status != nil && status.(int) != 1 {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "User is disabled",
			})
			c.Abort()
			return
		}

		// Check role
		if role == nil || role.(int) < 1 {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Insufficient permissions",
			})
			c.Abort()
			return
		}

		// Set context values
		c.Set("username", username)
		c.Set("role", role)
		c.Set("id", id)
		c.Set("status", status)

		c.Next()
	}
}

// AdminAuth is a middleware that requires admin authentication
func AdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		role := session.Get("role")

		if role == nil || role.(int) < 10 {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Admin access required",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RootAuth is a middleware that requires root authentication
func RootAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		role := session.Get("role")

		if role == nil || role.(int) < 100 {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Root access required",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// TokenAuth is a middleware that validates API tokens
func TokenAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract token from request
		token := extractToken(c)
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "API key required",
			})
			c.Abort()
			return
		}

		// Validate token
		user, tokenModel, err := model.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Invalid API key: " + err.Error(),
			})
			c.Abort()
			return
		}

		// Set context values
		c.Set("userId", user.Id)
		c.Set("username", user.Username)
		c.Set("token", token)
		c.Set("tokenId", tokenModel.Id)
		c.Set("tokenName", tokenModel.Name)
		c.Set("tokenQuota", tokenModel.RemainQuota)
		c.Set("userGroup", user.Group)

		c.Next()
	}
}

func extractToken(c *gin.Context) string {
	// Check Authorization header
	auth := c.GetHeader("Authorization")
	if auth != "" {
		if strings.HasPrefix(auth, "Bearer ") {
			return strings.TrimPrefix(auth, "Bearer ")
		}
		return auth
	}

	// Check x-api-key header
	apiKey := c.GetHeader("x-api-key")
	if apiKey != "" {
		return apiKey
	}

	// Check query parameter
	return c.Query("api_key")
}
