package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"github.com/fangxiusun/tokenhub/common"
	"github.com/fangxiusun/tokenhub/constant"
	"github.com/fangxiusun/tokenhub/model"
)

// UserAuth is a middleware that requires user authentication
func UserAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Try JWT token first
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader != "" {
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			claims, err := common.ValidateJWT(tokenString)
			if err == nil {
				// JWT is valid
				c.Set("id", claims.UserId)
				c.Set("username", claims.Username)
				c.Set("role", claims.Role)
				c.Set("status", claims.Status)
				c.Set("group", claims.GroupId)
				c.Next()
				return
			}
		}

		// Fallback to session
		session := sessions.Default(c)
		username := session.Get("username")
		role := session.Get("role")
		id := session.Get("id")
		status := session.Get("status")
		group := session.Get("group")

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

			// Strip Bearer prefix before validating access token
			accessToken = strings.TrimPrefix(accessToken, "Bearer ")

			// Validate access token
			user, err := model.ValidateAccessToken(accessToken)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{
					"success": false,
					"message": "Invalid access token",
				})
				c.Abort()
				return
			}

			if user != nil && user.Username != "" {
				username = user.Username
				role = user.Role
				id = user.Id
				status = user.Status
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{
					"success": false,
					"message": "Invalid access token",
				})
				c.Abort()
				return
			}
		}

		// Check user status
		if status != nil && status.(int) != constant.UserStatusEnabled {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "User is disabled",
			})
			c.Abort()
			return
		}

		// Check role
		if role == nil || role.(int) < constant.RoleCommonUser {
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
		c.Set("group", group)

		c.Next()
	}
}

// AdminAuth is a middleware that requires admin authentication
func AdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Try JWT token first
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader != "" {
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			claims, err := common.ValidateJWT(tokenString)
			if err == nil {
				if claims.Role < 10 {
					c.JSON(http.StatusForbidden, gin.H{
						"success": false,
						"message": "Admin access required",
					})
					c.Abort()
					return
				}
				if claims.Status != constant.UserStatusEnabled {
					c.JSON(http.StatusForbidden, gin.H{
						"success": false,
						"message": "User is disabled",
					})
					c.Abort()
					return
				}
				c.Set("id", claims.UserId)
				c.Set("username", claims.Username)
				c.Set("role", claims.Role)
				c.Set("status", claims.Status)
				c.Set("group", claims.GroupId)
				c.Next()
				return
			}
		}

		// Fallback to session
		session := sessions.Default(c)
		role := session.Get("role")
		id := session.Get("id")
		username := session.Get("username")
		status := session.Get("status")
		group := session.Get("group")

		if role == nil || role.(int) < constant.RoleAdminUser {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Admin access required",
			})
			c.Abort()
			return
		}

		// Check user status
		if status != nil && status.(int) != constant.UserStatusEnabled {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "User is disabled",
			})
			c.Abort()
			return
		}

		c.Set("id", id)
		c.Set("username", username)
		c.Set("role", role)
		c.Set("status", status)
		c.Set("group", group)
		c.Next()
	}
}

// RootAuth is a middleware that requires root authentication
func RootAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Try JWT token first
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader != "" {
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			claims, err := common.ValidateJWT(tokenString)
			if err == nil {
				if claims.Role < 100 {
					c.JSON(http.StatusForbidden, gin.H{
						"success": false,
						"message": "Root access required",
					})
					c.Abort()
					return
				}
				if claims.Status != constant.UserStatusEnabled {
					c.JSON(http.StatusForbidden, gin.H{
						"success": false,
						"message": "User is disabled",
					})
					c.Abort()
					return
				}
				c.Set("id", claims.UserId)
				c.Set("username", claims.Username)
				c.Set("role", claims.Role)
				c.Set("status", claims.Status)
				c.Set("group", claims.GroupId)
				c.Next()
				return
			}
		}

		// Fallback to session
		session := sessions.Default(c)
		role := session.Get("role")
		id := session.Get("id")
		username := session.Get("username")
		status := session.Get("status")
		group := session.Get("group")

		if role == nil || role.(int) < constant.RoleRootUser {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Root access required",
			})
			c.Abort()
			return
		}

		// Check user status
		if status != nil && status.(int) != constant.UserStatusEnabled {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "User is disabled",
			})
			c.Abort()
			return
		}

		c.Set("id", id)
		c.Set("username", username)
		c.Set("role", role)
		c.Set("status", status)
		c.Set("group", group)
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

