package middleware

import (
	"crypto-indicator-dashboard/pkg/logger"
	"github.com/gin-gonic/gin"
	"time"
)

// RequestLogging creates a logging middleware
func RequestLogging(logger logger.Logger) gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// Custom log format
		logger.Info("HTTP Request",
			"timestamp", param.TimeStamp.Format(time.RFC3339),
			"status", param.StatusCode,
			"latency", param.Latency,
			"client_ip", param.ClientIP,
			"method", param.Method,
			"path", param.Path,
			"user_agent", param.Request.UserAgent(),
			"error_message", param.ErrorMessage,
		)
		return ""
	})
}

// ErrorLogging creates an error logging middleware
func ErrorLogging(logger logger.Logger) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		logger.Error("Panic recovered",
			"error", recovered,
			"path", c.Request.URL.Path,
			"method", c.Request.Method,
			"client_ip", c.ClientIP(),
		)
		
		c.JSON(500, gin.H{
			"success": false,
			"error": gin.H{
				"type":    "INTERNAL_ERROR",
				"message": "An internal error occurred",
			},
		})
	})
}