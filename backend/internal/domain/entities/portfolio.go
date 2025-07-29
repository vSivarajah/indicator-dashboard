package entities

import (
	"time"
)

// Portfolio represents a user's investment portfolio
type Portfolio struct {
	ID          uint              `json:"id"`
	UserID      string            `json:"user_id"`
	Name        string            `json:"name"`
	Holdings    []PortfolioHolding `json:"holdings"`
	TotalValue  float64           `json:"total_value"`
	RiskLevel   string            `json:"risk_level"`
	LastUpdated time.Time         `json:"last_updated"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

// PortfolioHolding represents individual holdings in a portfolio
type PortfolioHolding struct {
	ID           uint      `json:"id"`
	PortfolioID  uint      `json:"portfolio_id"`
	Symbol       string    `json:"symbol"`
	Amount       float64   `json:"amount"`
	AveragePrice float64   `json:"average_price"`
	CurrentPrice float64   `json:"current_price"`
	Value        float64   `json:"value"`
	PnL          float64   `json:"pnl"`
	PnLPercent   float64   `json:"pnl_percent"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// PortfolioSummary represents aggregated portfolio data
type PortfolioSummary struct {
	TotalValue        float64                 `json:"total_value"`
	TotalPnL          float64                 `json:"total_pnl"`
	TotalPnLPercent   float64                 `json:"total_pnl_percent"`
	DayChange         float64                 `json:"day_change"`
	DayChangePercent  float64                 `json:"day_change_percent"`
	TopPerformer      *PortfolioHolding       `json:"top_performer"`
	WorstPerformer    *PortfolioHolding       `json:"worst_performer"`
	AllocationByAsset []AssetAllocation       `json:"allocation_by_asset"`
	RiskMetrics       PortfolioRiskMetrics    `json:"risk_metrics"`
}

// AssetAllocation represents asset allocation in portfolio
type AssetAllocation struct {
	Symbol     string  `json:"symbol"`
	Name       string  `json:"name"`
	Value      float64 `json:"value"`
	Percentage float64 `json:"percentage"`
	Color      string  `json:"color"`
}

// PortfolioRiskMetrics represents risk analysis
type PortfolioRiskMetrics struct {
	OverallRisk       string  `json:"overall_risk"`
	Volatility        float64 `json:"volatility"`
	SharpeRatio       float64 `json:"sharpe_ratio"`
	MaxDrawdown       float64 `json:"max_drawdown"`
	BetaToMarket      float64 `json:"beta_to_market"`
	ConcentrationRisk string  `json:"concentration_risk"`
}