package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiter is a simple in-memory rate limiter
type RateLimiter struct {
	visitors map[string][]time.Time
	mu       sync.RWMutex
	limit    int
	window   time.Duration
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}

	// Cleanup old entries
	go func() {
		for {
			time.Sleep(time.Minute)
			rl.cleanup()
		}
	}()

	return rl
}

// isAllowed checks if a request is allowed
func (rl *RateLimiter) isAllowed(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-rl.window)

	// Get visitor's requests
	requests := rl.visitors[key]

	// Filter out old requests
	var validRequests []time.Time
	for _, req := range requests {
		if req.After(windowStart) {
			validRequests = append(validRequests, req)
		}
	}

	// Check if limit exceeded
	if len(validRequests) >= rl.limit {
		return false
	}

	// Add current request
	rl.visitors[key] = append(validRequests, now)
	return true
}

// cleanup removes old entries
func (rl *RateLimiter) cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-rl.window)

	for key, requests := range rl.visitors {
		var validRequests []time.Time
		for _, req := range requests {
			if req.After(windowStart) {
				validRequests = append(validRequests, req)
			}
		}
		if len(validRequests) == 0 {
			delete(rl.visitors, key)
		} else {
			rl.visitors[key] = validRequests
		}
	}
}

// Global rate limiter
var globalLimiter = NewRateLimiter(100, time.Minute)

// GlobalAPIRateLimit is a middleware that limits API requests
func GlobalAPIRateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.ClientIP()
		if !globalLimiter.isAllowed(key) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"message": "Rate limit exceeded",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// CriticalRateLimit is a middleware that limits critical operations
func CriticalRateLimit() gin.HandlerFunc {
	criticalLimiter := NewRateLimiter(20, 20*time.Minute)
	return func(c *gin.Context) {
		key := c.ClientIP()
		if !criticalLimiter.isAllowed(key) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"message": "Too many requests, please try again later",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

