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

// indicatorRepository implements the IndicatorRepository interface
type indicatorRepository struct {
	db *gorm.DB
	logger logger.Logger
}

// NewIndicatorRepository creates a new instance of indicator repository
func NewIndicatorRepository(db *gorm.DB, logger logger.Logger) repositories.IndicatorRepository {
	return &indicatorRepository{
		db:     db,
		logger: logger,
	}
}

// Create saves a new indicator to the database
func (r *indicatorRepository) Create(ctx context.Context, indicator *entities.Indicator) error {
	r.logger.Info("Creating new indicator", 
		"name", indicator.Name, 
		"type", indicator.Type)

	if err := r.db.WithContext(ctx).Create(indicator).Error; err != nil {
		r.logger.Error("Failed to create indicator", 
			"error", err, 
			"name", indicator.Name)
		return errors.Wrap(err, errors.ErrorTypeInternal, "failed to create indicator")
	}

	r.logger.Info("Successfully created indicator", 
		"id", indicator.ID, 
		"name", indicator.Name)
	return nil
}

// GetByID retrieves an indicator by its ID
func (r *indicatorRepository) GetByID(ctx context.Context, id uint) (*entities.Indicator, error) {
	r.logger.Debug("Retrieving indicator by ID", "id", id)

	var indicator entities.Indicator
	if err := r.db.WithContext(ctx).First(&indicator, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			r.logger.Debug("Indicator not found", "id", id)
			return nil, errors.NotFound("indicator")
		}
		r.logger.Error("Failed to retrieve indicator", "error", err, "id", id)
		return nil, errors.Wrap(err, errors.ErrorTypeInternal, "failed to retrieve indicator")
	}

	return &indicator, nil
}

// GetByName retrieves an indicator by its name
func (r *indicatorRepository) GetByName(ctx context.Context, name string) (*entities.Indicator, error) {
	r.logger.Debug("Retrieving indicator by name", "name", name)

	var indicator entities.Indicator
	if err := r.db.WithContext(ctx).Where("name = ?", name).First(&indicator).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			r.logger.Debug("Indicator not found", "name", name)
			return nil, errors.NotFound("indicator")
		}
		r.logger.Error("Failed to retrieve indicator", "error", err, "name", name)
		return nil, errors.Wrap(err, errors.ErrorTypeInternal, "failed to retrieve indicator")
	}

	return &indicator, nil
}

// GetByType retrieves all indicators of a specific type
func (r *indicatorRepository) GetByType(ctx context.Context, indicatorType string) ([]entities.Indicator, error) {
	r.logger.Debug("Retrieving indicators by type", "type", indicatorType)

	var indicators []entities.Indicator
	if err := r.db.WithContext(ctx).Where("type = ?", indicatorType).Order("created_at DESC").Find(&indicators).Error; err != nil {
		r.logger.Error("Failed to retrieve indicators", "error", err, "type", indicatorType)
		return nil, errors.Wrap(err, errors.ErrorTypeInternal, "failed to retrieve indicators")
	}

	r.logger.Debug("Retrieved indicators", "count", len(indicators), "type", indicatorType)
	return indicators, nil
}

// Update modifies an existing indicator
func (r *indicatorRepository) Update(ctx context.Context, indicator *entities.Indicator) error {
	r.logger.Info("Updating indicator", 
		"id", indicator.ID, 
		"name", indicator.Name)

	indicator.UpdatedAt = time.Now()
	
	if err := r.db.WithContext(ctx).Save(indicator).Error; err != nil {
		r.logger.Error("Failed to update indicator", 
			"error", err, 
			"id", indicator.ID)
		return errors.Wrap(err, errors.ErrorTypeInternal, "failed to update indicator")
	}

	r.logger.Info("Successfully updated indicator", "id", indicator.ID)
	return nil
}

// Delete removes an indicator from the database
func (r *indicatorRepository) Delete(ctx context.Context, id uint) error {
	r.logger.Info("Deleting indicator", "id", id)

	result := r.db.WithContext(ctx).Delete(&entities.Indicator{}, id)
	if err := result.Error; err != nil {
		r.logger.Error("Failed to delete indicator", "error", err, "id", id)
		return errors.Wrap(err, errors.ErrorTypeInternal, "failed to delete indicator")
	}

	if result.RowsAffected == 0 {
		r.logger.Debug("Indicator not found for deletion", "id", id)
		return errors.NotFound("indicator")
	}

	r.logger.Info("Successfully deleted indicator", "id", id)
	return nil
}

// GetHistoricalData retrieves historical data for an indicator within a time range
func (r *indicatorRepository) GetHistoricalData(ctx context.Context, name string, from, to time.Time) ([]entities.Indicator, error) {
	r.logger.Debug("Retrieving historical data", 
		"name", name, 
		"from", from, 
		"to", to)

	var indicators []entities.Indicator
	query := r.db.WithContext(ctx).
		Where("name = ? AND created_at BETWEEN ? AND ?", name, from, to).
		Order("created_at ASC")

	if err := query.Find(&indicators).Error; err != nil {
		r.logger.Error("Failed to retrieve historical data", 
			"error", err, 
			"name", name)
		return nil, errors.Wrap(err, errors.ErrorTypeInternal, "failed to retrieve historical data")
	}

	r.logger.Debug("Retrieved historical data", 
		"count", len(indicators), 
		"name", name)
	return indicators, nil
}

// GetLatest retrieves the most recent indicator by name
func (r *indicatorRepository) GetLatest(ctx context.Context, name string) (*entities.Indicator, error) {
	r.logger.Debug("Retrieving latest indicator", "name", name)

	var indicator entities.Indicator
	if err := r.db.WithContext(ctx).
		Where("name = ?", name).
		Order("created_at DESC").
		First(&indicator).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			r.logger.Debug("No indicator found", "name", name)
			return nil, errors.NotFound("indicator")
		}
		r.logger.Error("Failed to retrieve latest indicator", "error", err, "name", name)
		return nil, errors.Wrap(err, errors.ErrorTypeInternal, "failed to retrieve latest indicator")
	}

	return &indicator, nil
}

// GetLatestByType retrieves the most recent indicators for each name of a specific type
func (r *indicatorRepository) GetLatestByType(ctx context.Context, indicatorType string) ([]entities.Indicator, error) {
	r.logger.Debug("Retrieving latest indicators by type", "type", indicatorType)

	var indicators []entities.Indicator
	
	// Use a subquery to get the latest record for each name of the specified type
	subquery := r.db.WithContext(ctx).
		Model(&entities.Indicator{}).
		Select("name, MAX(created_at) as max_created_at").
		Where("type = ?", indicatorType).
		Group("name")

	if err := r.db.WithContext(ctx).
		Joins("JOIN (?) as latest ON indicators.name = latest.name AND indicators.created_at = latest.max_created_at", subquery).
		Where("indicators.type = ?", indicatorType).
		Find(&indicators).Error; err != nil {
		r.logger.Error("Failed to retrieve latest indicators", "error", err, "type", indicatorType)
		return nil, errors.Wrap(err, errors.ErrorTypeInternal, "failed to retrieve latest indicators")
	}

	r.logger.Debug("Retrieved latest indicators", "count", len(indicators), "type", indicatorType)
	return indicators, nil
}

// BulkCreate saves multiple indicators in a single transaction
func (r *indicatorRepository) BulkCreate(ctx context.Context, indicators []entities.Indicator) error {
	r.logger.Info("Bulk creating indicators", "count", len(indicators))

	if len(indicators) == 0 {
		return nil
	}

	if err := r.db.WithContext(ctx).CreateInBatches(indicators, 100).Error; err != nil {
		r.logger.Error("Failed to bulk create indicators", "error", err, "count", len(indicators))
		return errors.Wrap(err, errors.ErrorTypeInternal, "failed to bulk create indicators")
	}

	r.logger.Info("Successfully bulk created indicators", "count", len(indicators))
	return nil
}

// CleanupOldData removes indicators older than the specified time
func (r *indicatorRepository) CleanupOldData(ctx context.Context, olderThan time.Time) error {
	r.logger.Info("Cleaning up old indicator data", "older_than", olderThan)

	result := r.db.WithContext(ctx).
		Where("created_at < ?", olderThan).
		Delete(&entities.Indicator{})

	if err := result.Error; err != nil {
		r.logger.Error("Failed to cleanup old data", "error", err, "older_than", olderThan)
		return errors.Wrap(err, errors.ErrorTypeInternal, "failed to cleanup old data")
	}

	r.logger.Info("Successfully cleaned up old data", 
		"deleted_count", result.RowsAffected, 
		"older_than", olderThan)
	return nil
}