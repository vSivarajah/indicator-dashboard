package database

import (
	"context"
	"crypto-indicator-dashboard/internal/domain/entities"
	"crypto-indicator-dashboard/internal/domain/repositories"
	"crypto-indicator-dashboard/pkg/errors"
	"crypto-indicator-dashboard/pkg/logger"
	"time"

	"gorm.io/gorm"
)

// dcaRepository implements the DCARepository interface
type dcaRepository struct {
	db     *gorm.DB
	logger logger.Logger
}

// NewDCARepository creates a new instance of DCA repository
func NewDCARepository(db *gorm.DB, logger logger.Logger) repositories.DCARepository {
	return &dcaRepository{
		db:     db,
		logger: logger,
	}
}

// CreateStrategy saves a new DCA strategy to the database
func (r *dcaRepository) CreateStrategy(ctx context.Context, strategy *entities.DCAStrategy) error {
	r.logger.Info("Creating new DCA strategy", 
		"user_id", strategy.UserID, 
		"name", strategy.Name,
		"symbol", strategy.Symbol)

	if err := r.db.WithContext(ctx).Create(strategy).Error; err != nil {
		r.logger.Error("Failed to create DCA strategy", 
			"error", err, 
			"user_id", strategy.UserID,
			"name", strategy.Name)
		return errors.Wrap(err, errors.ErrorTypeInternal, "failed to create DCA strategy")
	}

	r.logger.Info("Successfully created DCA strategy", 
		"id", strategy.ID, 
		"user_id", strategy.UserID,
		"name", strategy.Name)
	return nil
}

// GetStrategyByID retrieves a DCA strategy by its ID
func (r *dcaRepository) GetStrategyByID(ctx context.Context, id uint) (*entities.DCAStrategy, error) {
	r.logger.Debug("Retrieving DCA strategy by ID", "id", id)

	var strategy entities.DCAStrategy
	if err := r.db.WithContext(ctx).First(&strategy, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			r.logger.Debug("DCA strategy not found", "id", id)
			return nil, errors.NotFound("dca_strategy")
		}
		r.logger.Error("Failed to retrieve DCA strategy", "error", err, "id", id)
		return nil, errors.Wrap(err, errors.ErrorTypeInternal, "failed to retrieve DCA strategy")
	}

	return &strategy, nil
}

// GetStrategiesByUserID retrieves all DCA strategies for a user
func (r *dcaRepository) GetStrategiesByUserID(ctx context.Context, userID string) ([]entities.DCAStrategy, error) {
	r.logger.Debug("Retrieving DCA strategies for user", "user_id", userID)

	var strategies []entities.DCAStrategy
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&strategies).Error; err != nil {
		r.logger.Error("Failed to retrieve user DCA strategies", "error", err, "user_id", userID)
		return nil, errors.Wrap(err, errors.ErrorTypeInternal, "failed to retrieve user DCA strategies")
	}

	r.logger.Debug("Retrieved DCA strategies", "count", len(strategies), "user_id", userID)
	return strategies, nil
}

// UpdateStrategy modifies an existing DCA strategy
func (r *dcaRepository) UpdateStrategy(ctx context.Context, strategy *entities.DCAStrategy) error {
	r.logger.Info("Updating DCA strategy", 
		"id", strategy.ID, 
		"user_id", strategy.UserID,
		"name", strategy.Name)

	strategy.UpdatedAt = time.Now()
	
	if err := r.db.WithContext(ctx).Save(strategy).Error; err != nil {
		r.logger.Error("Failed to update DCA strategy", 
			"error", err, 
			"id", strategy.ID)
		return errors.Wrap(err, errors.ErrorTypeInternal, "failed to update DCA strategy")
	}

	r.logger.Info("Successfully updated DCA strategy", "id", strategy.ID)
	return nil
}

// DeleteStrategy removes a DCA strategy from the database
func (r *dcaRepository) DeleteStrategy(ctx context.Context, id uint) error {
	r.logger.Info("Deleting DCA strategy", "id", id)

	result := r.db.WithContext(ctx).Delete(&entities.DCAStrategy{}, id)
	if err := result.Error; err != nil {
		r.logger.Error("Failed to delete DCA strategy", "error", err, "id", id)
		return errors.Wrap(err, errors.ErrorTypeInternal, "failed to delete DCA strategy")
	}

	if result.RowsAffected == 0 {
		r.logger.Debug("DCA strategy not found for deletion", "id", id)
		return errors.NotFound("dca_strategy")
	}

	r.logger.Info("Successfully deleted DCA strategy", "id", id)
	return nil
}

// CreatePurchase saves a new DCA purchase to the database
func (r *dcaRepository) CreatePurchase(ctx context.Context, purchase *entities.DCAPurchase) error {
	r.logger.Debug("Creating DCA purchase", 
		"strategy_id", purchase.StrategyID,
		"amount", purchase.Amount,
		"price", purchase.Price)

	if err := r.db.WithContext(ctx).Create(purchase).Error; err != nil {
		r.logger.Error("Failed to create DCA purchase", 
			"error", err, 
			"strategy_id", purchase.StrategyID)
		return errors.Wrap(err, errors.ErrorTypeInternal, "failed to create DCA purchase")
	}

	r.logger.Debug("Successfully created DCA purchase", "id", purchase.ID)
	return nil
}

// GetPurchasesByStrategy retrieves all purchases for a DCA strategy
func (r *dcaRepository) GetPurchasesByStrategy(ctx context.Context, strategyID uint) ([]entities.DCAPurchase, error) {
	r.logger.Debug("Retrieving purchases for strategy", "strategy_id", strategyID)

	var purchases []entities.DCAPurchase
	if err := r.db.WithContext(ctx).
		Where("strategy_id = ?", strategyID).
		Order("created_at DESC").
		Find(&purchases).Error; err != nil {
		r.logger.Error("Failed to retrieve strategy purchases", "error", err, "strategy_id", strategyID)
		return nil, errors.Wrap(err, errors.ErrorTypeInternal, "failed to retrieve strategy purchases")
	}

	r.logger.Debug("Retrieved purchases", "count", len(purchases), "strategy_id", strategyID)
	return purchases, nil
}

// SaveSimulation saves a DCA simulation result
func (r *dcaRepository) SaveSimulation(ctx context.Context, simulation *entities.DCASimulation) error {
	r.logger.Debug("Saving DCA simulation", 
		"user_id", simulation.UserID,
		"symbol", simulation.Symbol,
		"amount", simulation.Amount)

	if err := r.db.WithContext(ctx).Create(simulation).Error; err != nil {
		r.logger.Error("Failed to save DCA simulation", 
			"error", err, 
			"user_id", simulation.UserID)
		return errors.Wrap(err, errors.ErrorTypeInternal, "failed to save DCA simulation")
	}

	return nil
}

// GetSimulationByID retrieves a DCA simulation by its ID
func (r *dcaRepository) GetSimulationByID(ctx context.Context, id uint) (*entities.DCASimulation, error) {
	r.logger.Debug("Retrieving DCA simulation by ID", "id", id)

	var simulation entities.DCASimulation
	if err := r.db.WithContext(ctx).First(&simulation, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			r.logger.Debug("DCA simulation not found", "id", id)
			return nil, errors.NotFound("dca_simulation")
		}
		r.logger.Error("Failed to retrieve DCA simulation", "error", err, "id", id)
		return nil, errors.Wrap(err, errors.ErrorTypeInternal, "failed to retrieve DCA simulation")
	}

	return &simulation, nil
}

// GetSimulationsByUser retrieves all DCA simulations for a user
func (r *dcaRepository) GetSimulationsByUser(ctx context.Context, userID string) ([]entities.DCASimulation, error) {
	r.logger.Debug("Retrieving DCA simulations for user", "user_id", userID)

	var simulations []entities.DCASimulation
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&simulations).Error; err != nil {
		r.logger.Error("Failed to retrieve user DCA simulations", "error", err, "user_id", userID)
		return nil, errors.Wrap(err, errors.ErrorTypeInternal, "failed to retrieve user DCA simulations")
	}

	r.logger.Debug("Retrieved DCA simulations", "count", len(simulations), "user_id", userID)
	return simulations, nil
}