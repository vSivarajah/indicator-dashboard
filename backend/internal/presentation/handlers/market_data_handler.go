package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"crypto-indicator-dashboard/internal/domain/entities"
	"crypto-indicator-dashboard/internal/domain/services"
	"crypto-indicator-dashboard/internal/infrastructure/external"
	"crypto-indicator-dashboard/pkg/logger"
	"github.com/gin-gonic/gin"
)

// MarketDataHandler handles market data HTTP requests
type MarketDataHandler struct {
	marketDataService   services.MarketDataService
	coinMarketCapClient *external.CoinMarketCapClient
	tradingViewScraper  *external.TradingViewScraper
	logger              logger.Logger
}

// NewMarketDataHandler creates a new market data handler
func NewMarketDataHandler(
	marketDataService services.MarketDataService,
	coinMarketCapClient *external.CoinMarketCapClient,
	tradingViewScraper *external.TradingViewScraper,
	logger logger.Logger,
) *MarketDataHandler {
	return &MarketDataHandler{
		marketDataService:   marketDataService,
		coinMarketCapClient: coinMarketCapClient,
		tradingViewScraper:  tradingViewScraper,
		logger:              logger,
	}
}

// GetCryptoPrices handles GET /api/v1/market/prices
func (h *MarketDataHandler) GetCryptoPrices(c *gin.Context) {
	symbolsParam := c.Query("symbols")
	var symbols []string
	
	if symbolsParam != "" {
		symbols = strings.Split(symbolsParam, ",")
		// Clean up whitespace
		for i, symbol := range symbols {
			symbols[i] = strings.TrimSpace(strings.ToUpper(symbol))
		}
	} else {
		// Default symbols
		symbols = []string{"BTC", "ETH", "BNB", "SOL", "ADA", "XRP", "DOT", "AVAX", "MATIC", "LINK"}
	}

	h.logger.Info("Fetching crypto prices", "symbols", symbols)

	prices, err := h.marketDataService.GetCryptoPrices(c.Request.Context(), symbols)
	if err != nil {
		h.logger.Error("Failed to get crypto prices", "error", err, "symbols", symbols)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch crypto prices",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    prices,
		"count":   len(prices),
	})
}

// GetBitcoinDominance handles GET /api/v1/market/dominance
func (h *MarketDataHandler) GetBitcoinDominance(c *gin.Context) {
	h.logger.Info("Fetching Bitcoin dominance")

	dominance, err := h.marketDataService.GetBitcoinDominance(c.Request.Context())
	if err != nil {
		h.logger.Error("Failed to get Bitcoin dominance", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch Bitcoin dominance",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    dominance,
	})
}

// GetMarketSummary handles GET /api/v1/market/summary
func (h *MarketDataHandler) GetMarketSummary(c *gin.Context) {
	h.logger.Info("Fetching market summary")

	// Get top cryptocurrencies
	countParam := c.DefaultQuery("count", "10")
	count, err := strconv.Atoi(countParam)
	if err != nil || count <= 0 || count > 50 {
		count = 10
	}

	prices, err := h.marketDataService.GetTopCryptoPrices(c.Request.Context(), count)
	if err != nil {
		h.logger.Error("Failed to get crypto prices for summary", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch market summary",
			"message": err.Error(),
		})
		return
	}

	// Get Bitcoin dominance
	dominance, err := h.marketDataService.GetBitcoinDominance(c.Request.Context())
	if err != nil {
		h.logger.Warn("Failed to get Bitcoin dominance for summary", "error", err)
		// Continue without dominance data
	}

	// Calculate total market cap and volume from available data
	var totalMarketCap, totalVolume24h float64
	for _, price := range prices {
		totalMarketCap += price.MarketCap
		totalVolume24h += price.Volume24h
	}

	summary := map[string]interface{}{
		"total_market_cap":      totalMarketCap,
		"total_volume_24h":      totalVolume24h,
		"bitcoin_dominance":     dominance,
		"top_cryptocurrencies":  prices,
		"market_trend":          determineTrendFromPrices(prices),
		"crypto_count":          len(prices),
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    summary,
	})
}

// GetSinglePrice handles GET /api/v1/market/price/:symbol
func (h *MarketDataHandler) GetSinglePrice(c *gin.Context) {
	symbol := strings.ToUpper(c.Param("symbol"))
	
	h.logger.Info("Fetching single price", "symbol", symbol)

	prices, err := h.marketDataService.GetCryptoPrices(c.Request.Context(), []string{symbol})
	if err != nil {
		h.logger.Error("Failed to get single price", "error", err, "symbol", symbol)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch price",
			"message": err.Error(),
		})
		return
	}

	price, exists := prices[symbol]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Price not found",
			"message": "Price data not available for " + symbol,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    price,
	})
}

// RefreshMarketData handles POST /api/v1/market/refresh
func (h *MarketDataHandler) RefreshMarketData(c *gin.Context) {
	h.logger.Info("Refreshing market data")

	err := h.marketDataService.RefreshAllMarketData(c.Request.Context())
	if err != nil {
		h.logger.Error("Failed to refresh market data", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to refresh market data",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Market data refreshed successfully",
	})
}

// GetHealthCheck handles GET /api/v1/market/health
func (h *MarketDataHandler) GetHealthCheck(c *gin.Context) {
	h.logger.Info("Checking market data sources health")

	healthResults := h.marketDataService.HealthCheck(c.Request.Context())
	
	allHealthy := true
	for _, err := range healthResults {
		if err != nil {
			allHealthy = false
			break
		}
	}

	status := http.StatusOK
	if !allHealthy {
		status = http.StatusServiceUnavailable
	}

	response := gin.H{
		"success": allHealthy,
		"sources": make(map[string]interface{}),
	}

	for source, err := range healthResults {
		if err != nil {
			response["sources"].(map[string]interface{})[source] = map[string]interface{}{
				"healthy": false,
				"error":   err.Error(),
			}
		} else {
			response["sources"].(map[string]interface{})[source] = map[string]interface{}{
				"healthy": true,
			}
		}
	}

	c.JSON(status, response)
}

// RegisterRoutes registers all market data routes
func (h *MarketDataHandler) RegisterRoutes(router *gin.RouterGroup) {
	market := router.Group("/market")
	{
		market.GET("/prices", h.GetCryptoPrices)
		market.GET("/price/:symbol", h.GetSinglePrice)
		market.GET("/dominance", h.GetBitcoinDominance)
		market.GET("/summary", h.GetMarketSummary)
		market.POST("/refresh", h.RefreshMarketData)
		market.GET("/health", h.GetHealthCheck)
	}
}

// Helper function to determine market trend based on price changes
func determineTrendFromPrices(prices map[string]*entities.CryptoPrice) string {
	if len(prices) == 0 {
		return "unknown"
	}

	var totalChange24h float64
	count := 0

	for _, price := range prices {
		totalChange24h += price.PercentChange24h
		count++
	}

	if count == 0 {
		return "unknown"
	}

	avgChange := totalChange24h / float64(count)
	
	if avgChange > 3 {
		return "bullish"
	} else if avgChange < -3 {
		return "bearish"
	} else {
		return "sideways"
	}
}