package repositories

import (
	"context"
	"crypto-indicator-dashboard/internal/domain/entities"
)

// DCARepository defines the interface for DCA strategy data operations
type DCARepository interface {
	// DCA Strategy CRUD operations
	CreateStrategy(ctx context.Context, strategy *entities.DCAStrategy) error
	GetStrategyByID(ctx context.Context, id uint) (*entities.DCAStrategy, error)
	GetStrategiesByUserID(ctx context.Context, userID string) ([]entities.DCAStrategy, error)
	UpdateStrategy(ctx context.Context, strategy *entities.DCAStrategy) error
	DeleteStrategy(ctx context.Context, id uint) error
	
	// DCA Purchase operations
	CreatePurchase(ctx context.Context, purchase *entities.DCAPurchase) error
	GetPurchasesByStrategy(ctx context.Context, strategyID uint) ([]entities.DCAPurchase, error)
	
	// DCA Simulation operations
	SaveSimulation(ctx context.Context, simulation *entities.DCASimulation) error
	GetSimulationsByUser(ctx context.Context, userID string) ([]entities.DCASimulation, error)
	GetSimulationByID(ctx context.Context, id uint) (*entities.DCASimulation, error)
}