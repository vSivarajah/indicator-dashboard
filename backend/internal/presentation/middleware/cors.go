package middleware

import (
	"crypto-indicator-dashboard/internal/infrastructure/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)

// CORS returns a CORS middleware configured for the application
func CORS(cfg *config.Config) gin.HandlerFunc {
	config := cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000",
			"http://localhost:5173",
			"http://localhost:5174",
			"http://localhost:5175",
		},
		AllowMethods: []string{
			"GET",
			"POST",
			"PUT",
			"PATCH",
			"DELETE",
			"OPTIONS",
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Accept",
			"Authorization",
			"X-Requested-With",
			"X-Request-ID",
		},
		ExposeHeaders: []string{
			"Content-Length",
			"X-Request-ID",
		},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	
	// In production, use specific origins from config
	if cfg.Server.IsProduction() {
		config.AllowOrigins = []string{
			// Add production URLs from config
		}
	}
	
	return cors.New(config)
}