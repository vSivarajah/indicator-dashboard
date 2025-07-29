package handlers

import (
	"context"
	domainservices "crypto-indicator-dashboard/internal/domain/services"
	"crypto-indicator-dashboard/internal/infrastructure/config"
	"crypto-indicator-dashboard/pkg/logger"
	"math"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// IndicatorHandler handles HTTP requests for market indicators
type IndicatorHandler struct {
	mvrvService    domainservices.IndicatorService
	cache          domainservices.CacheService
	logger         logger.Logger
	dependencies   *config.Dependencies
}

// NewIndicatorHandler creates a new indicator handler
func NewIndicatorHandler(deps *config.Dependencies) *IndicatorHandler {
	return &IndicatorHandler{
		cache:        deps.Cache,
		logger:       deps.Logger,
		dependencies: deps,
	}
}

// RegisterRoutes registers all indicator routes
func (h *IndicatorHandler) RegisterRoutes(router *gin.RouterGroup) {
	indicators := router.Group("/indicators")
	{
		indicators.GET("/mvrv", h.GetMVRVIndicator)
		indicators.GET("/dominance", h.GetDominanceIndicator)
		indicators.GET("/fear-greed", h.GetFearGreedIndicator)
		indicators.GET("/bubble-risk", h.GetBubbleRiskIndicator)
	}

	// Chart data endpoints
	charts := router.Group("/charts")
	{
		charts.GET("/:indicator", h.GetChartData)
	}
}

// GetMVRVIndicator handles MVRV Z-Score indicator requests
func (h *IndicatorHandler) GetMVRVIndicator(c *gin.Context) {
	h.logger.Info("Processing MVRV indicator request")

	// Temporarily return mock data due to cache interface conflicts
	// TODO: Fix cache interface compatibility between old and new services
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"value":           "2.43",
			"change":          "+0.12", 
			"risk_level":      "medium",
			"status":          "Service temporarily unavailable - under maintenance",
			"last_updated":    time.Now(),
		},
	})
}

// GetDominanceIndicator handles Bitcoin dominance indicator requests
func (h *IndicatorHandler) GetDominanceIndicator(c *gin.Context) {
	h.logger.Info("Processing dominance indicator request")

	// Return mock data - use /api/v1/market/dominance for real data
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"value":           "56.8%",
			"change":          "-1.2%",
			"risk_level":      "low",
			"status":          "Use /api/v1/market/dominance for real data",
			"last_updated":    time.Now(),
		},
	})
}

// GetFearGreedIndicator handles Fear & Greed index requests
func (h *IndicatorHandler) GetFearGreedIndicator(c *gin.Context) {
	h.logger.Info("Processing Fear & Greed indicator request")

	// Return mock data
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"value":           "72",
			"change":          "+5",
			"risk_level":      "high",
			"status":          "Greed territory - Consider taking profits",
			"last_updated":    time.Now(),
		},
	})
}

// GetBubbleRiskIndicator handles bubble risk assessment requests
func (h *IndicatorHandler) GetBubbleRiskIndicator(c *gin.Context) {
	h.logger.Info("Processing bubble risk indicator request")

	// Return mock data
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"value":           "Medium",
			"change":          "Stable",
			"risk_level":      "medium",
			"status":          "Monitor closely for rapid changes",
			"last_updated":    time.Now(),
		},
	})
}

// GetChartData handles chart data requests for indicators
func (h *IndicatorHandler) GetChartData(c *gin.Context) {
	ctx := c.Request.Context()
	indicator := c.Param("indicator")
	h.logger.Info("Processing chart data request", "indicator", indicator)

	switch indicator {
	case "mvrv":
		chartData, err := h.getMVRVChartData(ctx)
		if err != nil {
			h.logger.Error("Failed to get MVRV chart data", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to fetch MVRV chart data",
			})
			return
		}
		c.JSON(http.StatusOK, chartData)

	case "dominance":
		chartData := h.generateDominanceChartData()
		c.JSON(http.StatusOK, chartData)

	case "fear-greed":
		chartData := h.generateFearGreedChartData()
		c.JSON(http.StatusOK, chartData)

	case "bubble-risk":
		chartData := h.generateBubbleRiskChartData()
		c.JSON(http.StatusOK, chartData)

	default:
		c.JSON(http.StatusOK, gin.H{
			"indicator": indicator,
			"message":   "Chart data coming soon",
			"mock_data": h.generateMockChartData(),
		})
	}

	h.logger.Info("Successfully processed chart data request", "indicator", indicator)
}

// Helper methods

// convertRiskLevel converts internal risk levels to frontend format
func (h *IndicatorHandler) convertRiskLevel(riskLevel string) string {
	switch riskLevel {
	case "extreme_high":
		return "high"
	case "high":
		return "high"
	case "medium":
		return "medium"
	case "low":
		return "low"
	case "extreme_low":
		return "low"
	default:
		return "medium"
	}
}

// getMVRVChartData retrieves MVRV chart data
func (h *IndicatorHandler) getMVRVChartData(ctx context.Context) (map[string]interface{}, error) {
	// Skip MVRV service initialization due to architecture migration
	// TODO: Complete migration of indicator services to new architecture
	
	// Return mock data since service is not available
	if h.mvrvService == nil {
		return h.generateMockMVRVChartData(), nil
	}

	// Get latest calculation which includes historical data
	indicator, err := h.mvrvService.GetLatest(ctx)
	if err != nil {
		return nil, err
	}

	// Convert to chart format
	var timestamps []int64
	var zScores []float64
	var prices []float64

	// For now, generate mock data based on the indicator
	// In production, this would extract and process historical_data from metadata
	for i := 0; i < 30; i++ {
		timestamp := time.Now().AddDate(0, 0, -30+i).Unix() * 1000
		timestamps = append(timestamps, timestamp)
		zScores = append(zScores, -2.0+float64(i)*0.15) // Mock z-score progression
		prices = append(prices, 30000+float64(i)*1000)  // Mock price progression
	}

	return map[string]interface{}{
		"timestamps":     timestamps,
		"zscore_data":    zScores,
		"price_data":     prices,
		"current_zscore": indicator.Value,
		"thresholds":     indicator.Metadata["zscore_thresholds"],
		"last_updated":   indicator.Timestamp,
	}, nil
}

// generateDominanceData creates mock dominance data
func (h *IndicatorHandler) generateDominanceData() map[string]interface{} {
	return gin.H{
		"value":      "54.2%",
		"change":     "-0.8%",
		"risk_level": "medium",
		"status":     "MEDIUM: Neutral dominance level - Monitor for trends",
		"timestamp":  time.Now().Format(time.RFC3339),
		"details": gin.H{
			"trend":             "declining",
			"trend_strength":    "moderate",
			"change_7d":         -2.1,
			"change_30d":        -5.4,
			"market_cycle":      "mid_bull",
			"alt_season":        false,
			"critical_levels": gin.H{
				"alt_season_trigger": 42.0,
				"strong_dominance":   65.0,
			},
		},
	}
}

// generateFearGreedData creates mock fear & greed data
func (h *IndicatorHandler) generateFearGreedData() map[string]interface{} {
	return gin.H{
		"value":      "72",
		"change":     "+5",
		"risk_level": "medium",
		"status":     "GREED: Market sentiment is greedy - Be cautious",
		"timestamp":  time.Now().Format(time.RFC3339),
		"details": gin.H{
			"classification":         "Greed",
			"change_7d":             8,
			"trading_recommendation": "Consider taking some profits",
			"data_source":           "Alternative.me API",
			"next_update":           time.Now().Add(24 * time.Hour).Format(time.RFC3339),
			"components": gin.H{
				"volatility": 75,
				"momentum":   80,
				"social":     65,
				"surveys":    70,
				"dominance":  68,
				"trends":     74,
			},
		},
	}
}

// generateBubbleRiskData creates mock bubble risk data
func (h *IndicatorHandler) generateBubbleRiskData() map[string]interface{} {
	return gin.H{
		"value":      "Medium Risk",
		"change":     "Real-time",
		"risk_level": "medium",
		"status":     "MEDIUM: Elevated risk levels - Monitor closely",
		"timestamp":  time.Now().Format(time.RFC3339),
		"details": gin.H{
			"risk_score":             45,
			"confidence_level":       82,
			"trading_recommendation": "Maintain current positions with tight stops",
			"data_source":           "Multi-factor analysis",
			"components": gin.H{
				"mvrv_score":    40,
				"nvt_score":     50,
				"social_score":  60,
				"flow_score":    35,
				"holder_score":  45,
			},
			"critical_levels": gin.H{
				"warning":  60,
				"danger":   80,
				"extreme":  90,
			},
		},
	}
}

// Chart data generators

func (h *IndicatorHandler) generateDominanceChartData() map[string]interface{} {
	timestamps := make([]int64, 30)
	values := make([]float64, 30)

	baseTime := time.Now().AddDate(0, 0, -30)
	for i := 0; i < 30; i++ {
		timestamps[i] = baseTime.AddDate(0, 0, i).Unix() * 1000
		// Simulate dominance oscillation between 45-65%
		values[i] = 55.0 + 10.0*math.Sin(float64(i)*0.2) + float64(i%3)*2.0
	}

	return map[string]interface{}{
		"timestamps":   timestamps,
		"values":       values,
		"last_updated": time.Now(),
		"current":      54.2,
		"levels": map[string]float64{
			"alt_season_trigger": 42.0,
			"strong_dominance":   65.0,
		},
	}
}

func (h *IndicatorHandler) generateFearGreedChartData() map[string]interface{} {
	timestamps := make([]int64, 30)
	values := make([]int, 30)

	baseTime := time.Now().AddDate(0, 0, -30)
	for i := 0; i < 30; i++ {
		timestamps[i] = baseTime.AddDate(0, 0, i).Unix() * 1000
		// Simulate fear & greed oscillation between 10-90
		values[i] = int(50.0 + 30.0*math.Sin(float64(i)*0.15) + float64(i%5)*3.0)
		if values[i] < 10 {
			values[i] = 10
		}
		if values[i] > 90 {
			values[i] = 90
		}
	}

	return map[string]interface{}{
		"timestamps":   timestamps,
		"values":       values,
		"last_updated": time.Now(),
		"current":      72,
		"levels": map[string]int{
			"extreme_fear":  25,
			"fear":          45,
			"greed":         75,
			"extreme_greed": 90,
		},
	}
}

func (h *IndicatorHandler) generateBubbleRiskChartData() map[string]interface{} {
	timestamps := make([]int64, 30)
	values := make([]int, 30)

	baseTime := time.Now().AddDate(0, 0, -30)
	for i := 0; i < 30; i++ {
		timestamps[i] = baseTime.AddDate(0, 0, i).Unix() * 1000
		// Simulate bubble risk progression
		values[i] = int(20.0 + float64(i)*1.2 + 10.0*math.Sin(float64(i)*0.1))
		if values[i] < 0 {
			values[i] = 0
		}
		if values[i] > 100 {
			values[i] = 100
		}
	}

	return map[string]interface{}{
		"timestamps":   timestamps,
		"values":       values,
		"last_updated": time.Now(),
		"current":      45,
		"levels": map[string]int{
			"low":      25,
			"medium":   50,
			"high":     75,
			"extreme":  90,
		},
	}
}

// generateMockChartData creates mock chart data for unknown indicators
func (h *IndicatorHandler) generateMockChartData() map[string]interface{} {
	timestamps := make([]int64, 30)
	values := make([]float64, 30)

	baseTime := time.Now().AddDate(0, 0, -30)
	for i := 0; i < 30; i++ {
		timestamps[i] = baseTime.AddDate(0, 0, i).Unix() * 1000
		values[i] = 50.0 + (float64(i%7) * 5.0) // Mock oscillating data
	}

	return map[string]interface{}{
		"timestamps":   timestamps,
		"values":       values,
		"last_updated": time.Now(),
	}
}

// generateMockMVRVChartData creates mock MVRV chart data
func (h *IndicatorHandler) generateMockMVRVChartData() map[string]interface{} {
	timestamps := make([]int64, 30)
	zScores := make([]float64, 30)
	prices := make([]float64, 30)

	baseTime := time.Now().AddDate(0, 0, -30)
	for i := 0; i < 30; i++ {
		timestamps[i] = baseTime.AddDate(0, 0, i).Unix() * 1000
		zScores[i] = -2.0 + float64(i)*0.15 // Mock z-score progression
		prices[i] = 30000 + float64(i)*1000  // Mock price progression
	}

	return map[string]interface{}{
		"timestamps":     timestamps,
		"zscore_data":    zScores,
		"price_data":     prices,
		"current_zscore": 2.43,
		"thresholds": map[string]float64{
			"extreme_low": -1.5,
			"low":        -0.5,
			"neutral":     0.5,
			"high":        3.0,
			"extreme_high": 7.0,
		},
		"last_updated": time.Now(),
	}
}