package main

import (
	"context"
	"crypto-indicator-dashboard/internal/infrastructure/config"
	"crypto-indicator-dashboard/internal/presentation/handlers"
	"crypto-indicator-dashboard/internal/presentation/middleware"
	"crypto-indicator-dashboard/models"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	
	"github.com/gin-gonic/gin"
)


func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		panic("Failed to load configuration: " + err.Error())
	}

	// Initialize dependencies
	deps, err := config.NewDependencies(cfg)
	if err != nil {
		panic("Failed to initialize dependencies: " + err.Error())
	}
	defer deps.Cleanup()

	// Run database migrations if database is available
	if deps.DB != nil {
		if err := models.AutoMigrate(deps.DB); err != nil {
			deps.Logger.Error("Database migration failed", "error", err)
		} else {
			deps.Logger.Info("Database migrations completed successfully")
		}
	}

	// Set Gin mode based on environment
	if cfg.Server.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create Gin router
	router := gin.New()

	// Add middleware
	router.Use(middleware.ErrorLogging(deps.Logger))
	router.Use(middleware.RequestLogging(deps.Logger))
	router.Use(middleware.CORS(cfg))
	
	// Rate limiting (100 requests per minute)
	rateLimiter := middleware.NewRateLimiter(100, deps.Logger)
	router.Use(rateLimiter.RateLimit())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"message":   "Crypto Indicator Dashboard API",
			"timestamp": time.Now().Format(time.RFC3339),
			"version":   "2.0.0",
		})
	})

	// Initialize handlers
	portfolioHandler := handlers.NewPortfolioHandler(deps.PortfolioUseCase, deps.Logger)
	indicatorHandler := handlers.NewIndicatorHandler(deps)
	marketDataHandler := handlers.NewMarketDataHandler(
		deps.MarketDataService,
		deps.CoinMarketCapClient,
		deps.TradingViewScraper,
		deps.Logger,
	)

	// API routes
	apiV1 := router.Group("/api/v1")
	{
		// Portfolio routes
		portfolios := apiV1.Group("/portfolios")
		{
			portfolios.POST("", portfolioHandler.CreatePortfolio)
			portfolios.GET("", portfolioHandler.GetUserPortfolios)
			portfolios.GET("/:id", portfolioHandler.GetPortfolio)
			portfolios.GET("/:id/summary", portfolioHandler.GetPortfolioSummary)
			portfolios.POST("/:id/holdings", portfolioHandler.AddHolding)
			portfolios.PUT("/:id/holdings/:holdingId", portfolioHandler.UpdateHolding)
			portfolios.DELETE("/:id/holdings/:holdingId", portfolioHandler.RemoveHolding)
		}

		// Register indicator routes using the new handler
		indicatorHandler.RegisterRoutes(apiV1)

		// Register market data routes using proper handler
		marketDataHandler.RegisterRoutes(apiV1)

		// Market cycle
		apiV1.GET("/market/cycle", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "Market cycle endpoint - new implementation coming soon",
			})
		})
	}

	// Create HTTP server
	server := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// Start server in a goroutine
	go func() {
		deps.Logger.Info("Starting HTTP server", "port", cfg.Server.Port, "environment", cfg.Server.Environment)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			deps.Logger.Error("Failed to start server", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	deps.Logger.Info("Shutting down server...")

	// Create a deadline for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	// Gracefully shutdown the server
	if err := server.Shutdown(ctx); err != nil {
		deps.Logger.Error("Server forced to shutdown", "error", err)
		os.Exit(1)
	}

	deps.Logger.Info("Server gracefully stopped")
}