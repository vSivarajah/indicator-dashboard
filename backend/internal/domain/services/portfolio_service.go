package services

import (
	"context"
	"crypto-indicator-dashboard/internal/domain/entities"
)

// PortfolioService defines the interface for portfolio business logic
type PortfolioService interface {
	// Portfolio management
	CreatePortfolio(ctx context.Context, userID, name string) (*entities.Portfolio, error)
	GetPortfolio(ctx context.Context, portfolioID uint) (*entities.Portfolio, error)
	GetUserPortfolios(ctx context.Context, userID string) ([]entities.Portfolio, error)
	UpdatePortfolio(ctx context.Context, portfolio *entities.Portfolio) error
	DeletePortfolio(ctx context.Context, portfolioID uint) error
	
	// Holdings management
	AddHolding(ctx context.Context, portfolioID uint, symbol string, amount, averagePrice float64) (*entities.PortfolioHolding, error)
	UpdateHolding(ctx context.Context, holdingID uint, amount, averagePrice float64) error
	RemoveHolding(ctx context.Context, holdingID uint) error
	
	// Portfolio analytics
	GetPortfolioSummary(ctx context.Context, portfolioID uint) (*entities.PortfolioSummary, error)
	CalculateRiskMetrics(ctx context.Context, portfolioID uint) (*entities.PortfolioRiskMetrics, error)
	GetAssetAllocation(ctx context.Context, portfolioID uint) ([]entities.AssetAllocation, error)
	UpdatePortfolioValues(ctx context.Context, portfolioID uint) error
}

// RiskAnalysisService defines the interface for portfolio risk analysis
type RiskAnalysisService interface {
	AnalyzePortfolioRisk(ctx context.Context, portfolio *entities.Portfolio) (*entities.PortfolioRiskMetrics, error)
	CalculateVaR(ctx context.Context, portfolio *entities.Portfolio, confidence float64) (float64, error)
	RunMonteCarloSimulation(ctx context.Context, portfolio *entities.Portfolio, simulations, timeHorizon int) (map[string]interface{}, error)
	GetPositionSizingRecommendations(ctx context.Context, portfolio *entities.Portfolio) (map[string]interface{}, error)
	AnalyzeCorrelations(ctx context.Context, portfolio *entities.Portfolio) (map[string]interface{}, error)
}