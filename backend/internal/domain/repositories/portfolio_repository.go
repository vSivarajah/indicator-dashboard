package repositories

import (
	"context"
	"crypto-indicator-dashboard/internal/domain/entities"
)

// PortfolioRepository defines the interface for portfolio data operations
type PortfolioRepository interface {
	// Portfolio CRUD operations
	Create(ctx context.Context, portfolio *entities.Portfolio) error
	GetByID(ctx context.Context, id uint) (*entities.Portfolio, error)
	GetByUserID(ctx context.Context, userID string) ([]entities.Portfolio, error)
	Update(ctx context.Context, portfolio *entities.Portfolio) error
	Delete(ctx context.Context, id uint) error
	
	// Portfolio Holdings operations
	AddHolding(ctx context.Context, portfolioID uint, holding *entities.PortfolioHolding) error
	UpdateHolding(ctx context.Context, holding *entities.PortfolioHolding) error
	RemoveHolding(ctx context.Context, holdingID uint) error
	GetHoldings(ctx context.Context, portfolioID uint) ([]entities.PortfolioHolding, error)
	
	// Portfolio analytics
	CalculateTotalValue(ctx context.Context, portfolioID uint) (float64, error)
	GetPortfolioSummary(ctx context.Context, portfolioID uint) (*entities.PortfolioSummary, error)
}