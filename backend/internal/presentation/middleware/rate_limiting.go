package middleware

import (
	"net/http"
	"sync"
	"time"
	"crypto-indicator-dashboard/pkg/logger"
	"github.com/gin-gonic/gin"
)

// RateLimiter implements a simple in-memory rate limiter
type RateLimiter struct {
	clients map[string]*clientInfo
	mutex   sync.RWMutex
	rate    int           // requests per minute
	window  time.Duration // time window
	logger  logger.Logger
}

type clientInfo struct {
	requests  int
	resetTime time.Time
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(requestsPerMinute int, logger logger.Logger) *RateLimiter {
	rl := &RateLimiter{
		clients: make(map[string]*clientInfo),
		rate:    requestsPerMinute,
		window:  time.Minute,
		logger:  logger,
	}
	
	// Start cleanup goroutine
	go rl.cleanupLoop()
	
	return rl
}

// RateLimit returns a rate limiting middleware
func (rl *RateLimiter) RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		
		if !rl.allow(clientIP) {
			rl.logger.Warn("Rate limit exceeded", "client_ip", clientIP, "path", c.Request.URL.Path)
			
			c.JSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"error": gin.H{
					"type":    "RATE_LIMIT_ERROR",
					"message": "Rate limit exceeded. Please try again later.",
				},
			})
			c.Abort()
			return
		}
		
		c.Next()
	}
}

// allow checks if a client is allowed to make a request
func (rl *RateLimiter) allow(clientIP string) bool {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()
	
	now := time.Now()
	client, exists := rl.clients[clientIP]
	
	if !exists || now.After(client.resetTime) {
		// New client or window expired
		rl.clients[clientIP] = &clientInfo{
			requests:  1,
			resetTime: now.Add(rl.window),
		}
		return true
	}
	
	if client.requests >= rl.rate {
		return false
	}
	
	client.requests++
	return true
}

// cleanupLoop periodically removes expired entries
func (rl *RateLimiter) cleanupLoop() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	
	for range ticker.C {
		rl.cleanup()
	}
}

// cleanup removes expired entries
func (rl *RateLimiter) cleanup() {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()
	
	now := time.Now()
	for clientIP, client := range rl.clients {
		if now.After(client.resetTime) {
			delete(rl.clients, clientIP)
		}
	}
}