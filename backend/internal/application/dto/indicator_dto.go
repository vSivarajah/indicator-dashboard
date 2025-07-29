package dto

import (
	"fmt"
	"time"
	"crypto-indicator-dashboard/internal/domain/entities"
)

// BaseIndicatorResponse represents common indicator response fields
type BaseIndicatorResponse struct {
	Value     string    `json:"value"`
	Change    string    `json:"change"`
	RiskLevel string    `json:"risk_level"`
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
}

// MVRVResponse represents MVRV indicator response
type MVRVResponse struct {
	BaseIndicatorResponse
	Details map[string]interface{} `json:"details"`
}

// NewMVRVResponse creates a new MVRV response from entity
func NewMVRVResponse(result *entities.MVRVResult) *MVRVResponse {
	return &MVRVResponse{
		BaseIndicatorResponse: BaseIndicatorResponse{
			Value:     fmt.Sprintf("%.2f", result.CurrentZScore),
			Change:    "+0.12", // This would be calculated from previous value
			RiskLevel: result.RiskLevel,
			Status:    result.Status,
			Timestamp: result.LastUpdated,
		},
		Details: map[string]interface{}{
			"mvrv_ratio":    result.MVRVRatio,
			"market_cap":    result.MarketCap,
			"realized_cap":  result.RealizedCap,
			"price":         result.Price,
			"thresholds":    result.ZScoreThresholds,
		},
	}
}

// DominanceResponse represents Bitcoin dominance response
type DominanceResponse struct {
	BaseIndicatorResponse
	Details map[string]interface{} `json:"details"`
}

// NewDominanceResponse creates a new dominance response from entity
func NewDominanceResponse(result *entities.DominanceResult) *DominanceResponse {
	changeStr := "0%"
	if result.Change24h > 0 {
		changeStr = fmt.Sprintf("+%.1f%%", result.Change24h)
	} else if result.Change24h < 0 {
		changeStr = fmt.Sprintf("%.1f%%", result.Change24h)
	}
	
	return &DominanceResponse{
		BaseIndicatorResponse: BaseIndicatorResponse{
			Value:     fmt.Sprintf("%.1f%%", result.CurrentDominance),
			Change:    changeStr,
			RiskLevel: result.RiskLevel,
			Status:    result.Status,
			Timestamp: result.LastUpdated,
		},
		Details: map[string]interface{}{
			"trend":             result.Trend,
			"trend_strength":    result.TrendStrength,
			"change_7d":         result.Change7d,
			"change_30d":        result.Change30d,
			"market_cycle":      result.MarketCycleStage,
			"alt_season":        result.AltSeasonSignal,
			"critical_levels":   result.CriticalLevels,
		},
	}
}

// FearGreedResponse represents Fear & Greed index response
type FearGreedResponse struct {
	BaseIndicatorResponse
	Details map[string]interface{} `json:"details"`
}

// NewFearGreedResponse creates a new Fear & Greed response from entity
func NewFearGreedResponse(result *entities.FearGreedResult) *FearGreedResponse {
	changeStr := "0"
	if result.Change24h > 0 {
		changeStr = fmt.Sprintf("+%d", result.Change24h)
	} else if result.Change24h < 0 {
		changeStr = fmt.Sprintf("%d", result.Change24h)
	}
	
	return &FearGreedResponse{
		BaseIndicatorResponse: BaseIndicatorResponse{
			Value:     fmt.Sprintf("%d", result.CurrentValue),
			Change:    changeStr,
			RiskLevel: result.RiskLevel,
			Status:    result.Status,
			Timestamp: result.LastUpdated,
		},
		Details: map[string]interface{}{
			"classification":         result.Classification,
			"change_7d":             result.Change7d,
			"components":            result.Components,
			"trading_recommendation": result.TradingRecommendation,
			"data_source":           result.DataSource,
			"next_update":           result.NextUpdate,
		},
	}
}

// BubbleRiskResponse represents bubble risk response
type BubbleRiskResponse struct {
	BaseIndicatorResponse
	Details map[string]interface{} `json:"details"`
}

// NewBubbleRiskResponse creates a new bubble risk response from entity
func NewBubbleRiskResponse(result *entities.BubbleRiskResult) *BubbleRiskResponse {
	return &BubbleRiskResponse{
		BaseIndicatorResponse: BaseIndicatorResponse{
			Value:     result.RiskCategory,
			Change:    "Real-time",
			RiskLevel: result.RiskLevel,
			Status:    result.Status,
			Timestamp: result.LastUpdated,
		},
		Details: map[string]interface{}{
			"risk_score":             result.CurrentRiskScore,
			"confidence_level":       result.ConfidenceLevel,
			"components":             result.Components,
			"trading_recommendation": result.TradingRecommendation,
			"data_source":            result.DataSource,
			"critical_levels":        result.CriticalLevels,
		},
	}
}

// InflationResponse represents inflation indicator response
type InflationResponse struct {
	BaseIndicatorResponse
	Details map[string]interface{} `json:"details"`
}

// NewInflationResponse creates a new inflation response from entity
func NewInflationResponse(result *entities.InflationResult) *InflationResponse {
	changeStr := "0%"
	if result.ChangePercent > 0 {
		changeStr = fmt.Sprintf("+%.1f%%", result.ChangePercent)
	} else if result.ChangePercent < 0 {
		changeStr = fmt.Sprintf("%.1f%%", result.ChangePercent)
	}
	
	return &InflationResponse{
		BaseIndicatorResponse: BaseIndicatorResponse{
			Value:     fmt.Sprintf("%.1f%%", result.CurrentRate),
			Change:    changeStr,
			RiskLevel: result.Trend,
			Status:    result.ImpactOnCrypto,
			Timestamp: result.LastUpdated,
		},
		Details: map[string]interface{}{
			"current_rate":       result.CurrentRate,
			"previous_rate":      result.PreviousRate,
			"change":            result.Change,
			"change_percent":    result.ChangePercent,
			"trend":             result.Trend,
			"impact_on_crypto":  result.ImpactOnCrypto,
			"data_source":       result.DataSource,
			"confidence_level":  result.ConfidenceLevel,
		},
	}
}

// InterestRateResponse represents interest rate indicator response
type InterestRateResponse struct {
	BaseIndicatorResponse
	Details map[string]interface{} `json:"details"`
}

// NewInterestRateResponse creates a new interest rate response from entity
func NewInterestRateResponse(result *entities.InterestRateResult) *InterestRateResponse {
	changeStr := "0%"
	if result.ChangePercent > 0 {
		changeStr = fmt.Sprintf("+%.2f%%", result.ChangePercent)
	} else if result.ChangePercent < 0 {
		changeStr = fmt.Sprintf("%.2f%%", result.ChangePercent)
	}
	
	return &InterestRateResponse{
		BaseIndicatorResponse: BaseIndicatorResponse{
			Value:     fmt.Sprintf("%.2f%%", result.CurrentRate),
			Change:    changeStr,
			RiskLevel: result.Trend,
			Status:    result.ImpactOnCrypto,
			Timestamp: result.LastUpdated,
		},
		Details: map[string]interface{}{
			"current_rate":       result.CurrentRate,
			"previous_rate":      result.PreviousRate,
			"change":            result.Change,
			"change_percent":    result.ChangePercent,
			"trend":             result.Trend,
			"expected_change":   result.ExpectedChange,
			"impact_on_crypto":  result.ImpactOnCrypto,
			"data_source":       result.DataSource,
			"confidence_level":  result.ConfidenceLevel,
		},
	}
}

// MarketCycleResponse represents market cycle response
type MarketCycleResponse struct {
	CycleStage           string  `json:"cycle_stage"`
	Confidence           string  `json:"confidence"`
	EstimatedTimeToPeak  string  `json:"estimated_time_to_peak"`
	Timestamp            time.Time `json:"timestamp"`
}

// NewMarketCycleResponse creates a new market cycle response from entity
func NewMarketCycleResponse(cycle *entities.MarketCycle) *MarketCycleResponse {
	return &MarketCycleResponse{
		CycleStage:          cycle.Stage,
		Confidence:          fmt.Sprintf("%.0f%%", cycle.Confidence),
		EstimatedTimeToPeak: fmt.Sprintf("%d months", cycle.EstimatedDuration),
		Timestamp:           cycle.Timestamp,
	}
}

// ChartDataResponse represents chart data response
type ChartDataResponse struct {
	Indicator string                 `json:"indicator"`
	Data      map[string]interface{} `json:"data"`
	Message   string                 `json:"message,omitempty"`
}

// NewChartDataResponse creates a new chart data response
func NewChartDataResponse(indicator string, data map[string]interface{}) *ChartDataResponse {
	return &ChartDataResponse{
		Indicator: indicator,
		Data:      data,
	}
}