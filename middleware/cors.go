package middleware

import (
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/fangxiusun/tokenhub/common"
)

// Cors is a middleware that handles CORS
func Cors() gin.HandlerFunc {
	allowOrigins := common.GetEnvOrDefault("CORS_ORIGINS", "*")
	allowAll := allowOrigins == "*"

	config := cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: !allowAll,
		MaxAge:           86400,
	}

	if allowAll {
		config.AllowAllOrigins = true
	} else {
		config.AllowOrigins = strings.Split(allowOrigins, ",")
	}

	return cors.New(config)
}

// CorsWithConfig is a middleware that handles CORS with custom config
func CorsWithConfig(allowOrigins []string) gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     allowOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           86400,
	})
}

