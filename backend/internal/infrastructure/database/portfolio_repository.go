package database

import (
	"context"
	"fmt"
	"crypto-indicator-dashboard/internal/domain/entities"
	"crypto-indicator-dashboard/internal/domain/repositories"
	"crypto-indicator-dashboard/models"
	"gorm.io/gorm"
)

// portfolioRepository implements the PortfolioRepository interface
type portfolioRepository struct {
	db *gorm.DB
}

// NewPortfolioRepository creates a new portfolio repository
func NewPortfolioRepository(db *gorm.DB) repositories.PortfolioRepository {
	return &portfolioRepository{
		db: db,
	}
}

// Create creates a new portfolio
func (r *portfolioRepository) Create(ctx context.Context, portfolio *entities.Portfolio) error {
	dbPortfolio := &models.Portfolio{
		UserID:     portfolio.UserID,
		Name:       portfolio.Name,
		TotalValue: portfolio.TotalValue,
		RiskLevel:  portfolio.RiskLevel,
	}
	
	if err := r.db.WithContext(ctx).Create(dbPortfolio).Error; err != nil {
		return fmt.Errorf("failed to create portfolio: %w", err)
	}
	
	// Update entity with generated ID
	portfolio.ID = dbPortfolio.ID
	portfolio.CreatedAt = dbPortfolio.CreatedAt
	portfolio.UpdatedAt = dbPortfolio.UpdatedAt
	
	return nil
}

// GetByID retrieves a portfolio by ID
func (r *portfolioRepository) GetByID(ctx context.Context, id uint) (*entities.Portfolio, error) {
	var dbPortfolio models.Portfolio
	
	if err := r.db.WithContext(ctx).Preload("Holdings").First(&dbPortfolio, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("portfolio not found")
		}
		return nil, fmt.Errorf("failed to get portfolio: %w", err)
	}
	
	return r.mapToEntity(&dbPortfolio), nil
}

// GetByUserID retrieves all portfolios for a user
func (r *portfolioRepository) GetByUserID(ctx context.Context, userID string) ([]entities.Portfolio, error) {
	var dbPortfolios []models.Portfolio
	
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Preload("Holdings").Find(&dbPortfolios).Error; err != nil {
		return nil, fmt.Errorf("failed to get user portfolios: %w", err)
	}
	
	portfolios := make([]entities.Portfolio, len(dbPortfolios))
	for i, dbPortfolio := range dbPortfolios {
		portfolios[i] = *r.mapToEntity(&dbPortfolio)
	}
	
	return portfolios, nil
}

// Update updates a portfolio
func (r *portfolioRepository) Update(ctx context.Context, portfolio *entities.Portfolio) error {
	dbPortfolio := r.mapToModel(portfolio)
	
	if err := r.db.WithContext(ctx).Save(dbPortfolio).Error; err != nil {
		return fmt.Errorf("failed to update portfolio: %w", err)
	}
	
	return nil
}

// Delete deletes a portfolio
func (r *portfolioRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&models.Portfolio{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete portfolio: %w", err)
	}
	
	return nil
}

// AddHolding adds a holding to a portfolio
func (r *portfolioRepository) AddHolding(ctx context.Context, portfolioID uint, holding *entities.PortfolioHolding) error {
	dbHolding := &models.PortfolioHolding{
		PortfolioID:  portfolioID,
		Symbol:       holding.Symbol,
		Amount:       holding.Amount,
		AveragePrice: holding.AveragePrice,
		CurrentPrice: holding.CurrentPrice,
		Value:        holding.Value,
		PnL:          holding.PnL,
		PnLPercent:   holding.PnLPercent,
	}
	
	if err := r.db.WithContext(ctx).Create(dbHolding).Error; err != nil {
		return fmt.Errorf("failed to add holding: %w", err)
	}
	
	// Update entity with generated ID
	holding.ID = dbHolding.ID
	holding.CreatedAt = dbHolding.CreatedAt
	holding.UpdatedAt = dbHolding.UpdatedAt
	
	return nil
}

// UpdateHolding updates a holding
func (r *portfolioRepository) UpdateHolding(ctx context.Context, holding *entities.PortfolioHolding) error {
	dbHolding := &models.PortfolioHolding{
		ID:           holding.ID,
		PortfolioID:  holding.PortfolioID,
		Symbol:       holding.Symbol,
		Amount:       holding.Amount,
		AveragePrice: holding.AveragePrice,
		CurrentPrice: holding.CurrentPrice,
		Value:        holding.Value,
		PnL:          holding.PnL,
		PnLPercent:   holding.PnLPercent,
	}
	
	if err := r.db.WithContext(ctx).Save(dbHolding).Error; err != nil {
		return fmt.Errorf("failed to update holding: %w", err)
	}
	
	return nil
}

// RemoveHolding removes a holding
func (r *portfolioRepository) RemoveHolding(ctx context.Context, holdingID uint) error {
	if err := r.db.WithContext(ctx).Delete(&models.PortfolioHolding{}, holdingID).Error; err != nil {
		return fmt.Errorf("failed to remove holding: %w", err)
	}
	
	return nil
}

// GetHoldings retrieves all holdings for a portfolio
func (r *portfolioRepository) GetHoldings(ctx context.Context, portfolioID uint) ([]entities.PortfolioHolding, error) {
	var dbHoldings []models.PortfolioHolding
	
	if err := r.db.WithContext(ctx).Where("portfolio_id = ?", portfolioID).Find(&dbHoldings).Error; err != nil {
		return nil, fmt.Errorf("failed to get holdings: %w", err)
	}
	
	holdings := make([]entities.PortfolioHolding, len(dbHoldings))
	for i, dbHolding := range dbHoldings {
		holdings[i] = entities.PortfolioHolding{
			ID:           dbHolding.ID,
			PortfolioID:  dbHolding.PortfolioID,
			Symbol:       dbHolding.Symbol,
			Amount:       dbHolding.Amount,
			AveragePrice: dbHolding.AveragePrice,
			CurrentPrice: dbHolding.CurrentPrice,
			Value:        dbHolding.Value,
			PnL:          dbHolding.PnL,
			PnLPercent:   dbHolding.PnLPercent,
			CreatedAt:    dbHolding.CreatedAt,
			UpdatedAt:    dbHolding.UpdatedAt,
		}
	}
	
	return holdings, nil
}

// CalculateTotalValue calculates the total value of a portfolio
func (r *portfolioRepository) CalculateTotalValue(ctx context.Context, portfolioID uint) (float64, error) {
	var totalValue float64
	
	if err := r.db.WithContext(ctx).Model(&models.PortfolioHolding{}).
		Where("portfolio_id = ?", portfolioID).
		Select("COALESCE(SUM(value), 0)").
		Scan(&totalValue).Error; err != nil {
		return 0, fmt.Errorf("failed to calculate total value: %w", err)
	}
	
	return totalValue, nil
}

// GetPortfolioSummary retrieves portfolio summary with analytics
func (r *portfolioRepository) GetPortfolioSummary(ctx context.Context, portfolioID uint) (*entities.PortfolioSummary, error) {
	// This is a simplified implementation
	// In a real implementation, you would calculate various metrics
	holdings, err := r.GetHoldings(ctx, portfolioID)
	if err != nil {
		return nil, fmt.Errorf("failed to get holdings for summary: %w", err)
	}
	
	var totalValue, totalPnL float64
	allocations := make([]entities.AssetAllocation, len(holdings))
	
	for i, holding := range holdings {
		totalValue += holding.Value
		totalPnL += holding.PnL
		
		allocations[i] = entities.AssetAllocation{
			Symbol:     holding.Symbol,
			Name:       holding.Symbol, // In real implementation, fetch full name
			Value:      holding.Value,
			Percentage: 0, // Will be calculated after total is known
		}
	}
	
	// Calculate percentages
	for i := range allocations {
		if totalValue > 0 {
			allocations[i].Percentage = (allocations[i].Value / totalValue) * 100
		}
	}
	
	var totalPnLPercent float64
	if totalValue > 0 {
		totalPnLPercent = (totalPnL / (totalValue - totalPnL)) * 100
	}
	
	return &entities.PortfolioSummary{
		TotalValue:        totalValue,
		TotalPnL:          totalPnL,
		TotalPnLPercent:   totalPnLPercent,
		AllocationByAsset: allocations,
	}, nil
}

// mapToEntity converts a database model to domain entity
func (r *portfolioRepository) mapToEntity(dbPortfolio *models.Portfolio) *entities.Portfolio {
	holdings := make([]entities.PortfolioHolding, len(dbPortfolio.Holdings))
	for i, dbHolding := range dbPortfolio.Holdings {
		holdings[i] = entities.PortfolioHolding{
			ID:           dbHolding.ID,
			PortfolioID:  dbHolding.PortfolioID,
			Symbol:       dbHolding.Symbol,
			Amount:       dbHolding.Amount,
			AveragePrice: dbHolding.AveragePrice,
			CurrentPrice: dbHolding.CurrentPrice,
			Value:        dbHolding.Value,
			PnL:          dbHolding.PnL,
			PnLPercent:   dbHolding.PnLPercent,
			CreatedAt:    dbHolding.CreatedAt,
			UpdatedAt:    dbHolding.UpdatedAt,
		}
	}
	
	return &entities.Portfolio{
		ID:          dbPortfolio.ID,
		UserID:      dbPortfolio.UserID,
		Name:        dbPortfolio.Name,
		Holdings:    holdings,
		TotalValue:  dbPortfolio.TotalValue,
		RiskLevel:   dbPortfolio.RiskLevel,
		LastUpdated: dbPortfolio.LastUpdated,
		CreatedAt:   dbPortfolio.CreatedAt,
		UpdatedAt:   dbPortfolio.UpdatedAt,
	}
}

// mapToModel converts a domain entity to database model
func (r *portfolioRepository) mapToModel(portfolio *entities.Portfolio) *models.Portfolio {
	holdings := make([]models.PortfolioHolding, len(portfolio.Holdings))
	for i, holding := range portfolio.Holdings {
		holdings[i] = models.PortfolioHolding{
			ID:           holding.ID,
			PortfolioID:  holding.PortfolioID,
			Symbol:       holding.Symbol,
			Amount:       holding.Amount,
			AveragePrice: holding.AveragePrice,
			CurrentPrice: holding.CurrentPrice,
			Value:        holding.Value,
			PnL:          holding.PnL,
			PnLPercent:   holding.PnLPercent,
			CreatedAt:    holding.CreatedAt,
			UpdatedAt:    holding.UpdatedAt,
		}
	}
	
	return &models.Portfolio{
		ID:          portfolio.ID,
		UserID:      portfolio.UserID,
		Name:        portfolio.Name,
		Holdings:    holdings,
		TotalValue:  portfolio.TotalValue,
		RiskLevel:   portfolio.RiskLevel,
		LastUpdated: portfolio.LastUpdated,
		CreatedAt:   portfolio.CreatedAt,
		UpdatedAt:   portfolio.UpdatedAt,
	}
}