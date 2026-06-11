package common

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// GetEnvOrDefault returns the value of an environment variable or a default value
func GetEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GetEnvOrDefaultInt returns the value of an environment variable as int or a default value
func GetEnvOrDefaultInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// GetEnvOrDefaultBool returns the value of an environment variable as bool or a default value
func GetEnvOrDefaultBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

// PasswordLoginEnabled indicates if password login is enabled
var PasswordLoginEnabled = true

// RegisterEnabled indicates if registration is enabled
var RegisterEnabled = true

// PasswordRegisterEnabled indicates if password registration is enabled
var PasswordRegisterEnabled = true

// EmailVerificationEnabled indicates if email verification is enabled
var EmailVerificationEnabled = false

// QuotaForNewUser is the quota given to new users
var QuotaForNewUser = 100

// BatchUpdateEnabled indicates if batch updates are enabled
var BatchUpdateEnabled = false

// IsValidateRole checks if the role is valid
func IsValidateRole(role int) bool {
	return role >= 0 && role <= 100
}

// StringsContains checks if a slice contains a string
func StringsContains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// GetTimestamp returns the current timestamp in seconds
func GetTimestamp() int64 {
	return time.Now().Unix()
}

// TranslateMessage translates a message key
func TranslateMessage(c interface{}, key string) string {
	return key
}

// SysLog logs a system message
func SysLog(msg string) {
	log.Printf("[SYS] %s", msg)
}

// SysError logs a system error
func SysError(msg string) {
	log.Printf("[ERROR] %s", msg)
}

// ApiError returns an error response
func ApiError(c *gin.Context, err error) {
	c.JSON(200, gin.H{
		"success": false,
		"message": err.Error(),
	})
}

// ApiSuccess returns a success response
func ApiSuccess(c *gin.Context, data interface{}) {
	c.JSON(200, gin.H{
		"success": true,
		"message": "",
		"data":    data,
	})
}

// SetContextKey sets a value in the context
func SetContextKey(c *gin.Context, key string, value interface{}) {
	c.Set(key, value)
}

