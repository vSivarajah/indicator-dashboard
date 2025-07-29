package services

import (
	"context"
	"crypto-indicator-dashboard/internal/domain/entities"
)

// DCAService defines the interface for DCA strategy business logic
type DCAService interface {
	// Strategy management
	CreateStrategy(ctx context.Context, userID string, strategy *entities.DCAStrategy) error
	GetStrategy(ctx context.Context, strategyID uint) (*entities.DCAStrategy, error)
	GetUserStrategies(ctx context.Context, userID string) ([]entities.DCAStrategy, error)
	UpdateStrategy(ctx context.Context, strategy *entities.DCAStrategy) error
	DeleteStrategy(ctx context.Context, strategyID uint) error
	
	// DCA simulation and backtesting
	SimulateDCA(ctx context.Context, request entities.DCARequest) (map[string]interface{}, error)
	BacktestStrategy(ctx context.Context, strategy *entities.DCAStrategy) (*entities.DCASimulation, error)
	
	// Purchase execution
	ExecutePurchase(ctx context.Context, strategyID uint) (*entities.DCAPurchase, error)
	GetPurchaseHistory(ctx context.Context, strategyID uint) ([]entities.DCAPurchase, error)
	
	// Analytics
	CalculateStrategyPerformance(ctx context.Context, strategyID uint) (map[string]interface{}, error)
	GetOptimalDCAFrequency(ctx context.Context, symbol string) (string, error)
}