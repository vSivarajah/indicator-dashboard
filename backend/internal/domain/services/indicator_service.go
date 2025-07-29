package services

import (
	"context"
	"crypto-indicator-dashboard/internal/domain/entities"
)

// IndicatorService defines the general interface for indicator calculations
type IndicatorService interface {
	Calculate(ctx context.Context, params map[string]interface{}) (*entities.Indicator, error)
	GetHistoricalData(ctx context.Context, period string) ([]entities.Indicator, error)
	GetLatest(ctx context.Context) (*entities.Indicator, error)
}

// MVRVService defines the interface for MVRV analysis
type MVRVService interface {
	GetMVRVZScore(ctx context.Context) (*entities.MVRVResult, error)
	GetHistoricalMVRVChart(ctx context.Context) (map[string]interface{}, error)
	CalculateMVRVRisk(ctx context.Context, zScore float64) (string, string)
}

// DominanceService defines the interface for Bitcoin dominance analysis
type DominanceService interface {
	GetDominanceAnalysis(ctx context.Context) (*entities.DominanceResult, error)
	GetDominanceChart(ctx context.Context) (map[string]interface{}, error)
	DetectAltSeason(ctx context.Context, dominance float64) bool
}

// FearGreedService defines the interface for Fear & Greed index analysis
type FearGreedService interface {
	GetFearGreedAnalysis(ctx context.Context) (*entities.FearGreedResult, error)
	GetFearGreedChart(ctx context.Context) (map[string]interface{}, error)
	AnalyzeSentiment(ctx context.Context, value int) string
}

// BubbleRiskService defines the interface for bubble risk analysis
type BubbleRiskService interface {
	GetBubbleRiskAnalysis(ctx context.Context) (*entities.BubbleRiskResult, error)
	GetBubbleRiskChart(ctx context.Context) (map[string]interface{}, error)
	CalculateRiskScore(ctx context.Context) (float64, error)
}

// MacroService defines the interface for macroeconomic analysis
type MacroService interface {
	GetInflationAnalysis(ctx context.Context) (*entities.InflationResult, error)
	GetInterestRateAnalysis(ctx context.Context) (*entities.InterestRateResult, error)
	AnalyzeMacroImpact(ctx context.Context) (map[string]interface{}, error)
}

// MarketCycleService defines the interface for market cycle analysis
type MarketCycleService interface {
	GetCurrentCycle(ctx context.Context) (*entities.MarketCycle, error)
	PredictCycleStage(ctx context.Context) (string, float64, error)
	EstimateCycleDuration(ctx context.Context) (int, error)
}