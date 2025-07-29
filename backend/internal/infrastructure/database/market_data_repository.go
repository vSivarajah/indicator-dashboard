package database

import (
	"context"
	"time"
	"crypto-indicator-dashboard/internal/domain/entities"
	"crypto-indicator-dashboard/internal/domain/repositories"
	"crypto-indicator-dashboard/pkg/errors"
	"crypto-indicator-dashboard/pkg/logger"
	"gorm.io/gorm"
)

// marketDataRepository implements the MarketDataRepository interface
type marketDataRepository struct {
	db     *gorm.DB
	logger logger.Logger
}

// NewMarketDataRepository creates a new instance of market data repository
func NewMarketDataRepository(db *gorm.DB, logger logger.Logger) repositories.MarketDataRepository {
	return &marketDataRepository{
		db:     db,
		logger: logger,
	}
}

// StorePriceData saves crypto price data to the database
func (r *marketDataRepository) StorePriceData(ctx context.Context, priceData *entities.CryptoPrice) error {
	r.logger.Debug("Saving price data", "symbol", priceData.Symbol, "price", priceData.Price)

	if err := r.db.WithContext(ctx).Create(priceData).Error; err != nil {
		r.logger.Error("Failed to save price data", "error", err, "symbol", priceData.Symbol)
		return errors.Wrap(err, errors.ErrorTypeInternal, "failed to save price data")
	}

	return nil
}

// GetPriceHistory retrieves historical crypto price data for a symbol
func (r *marketDataRepository) GetPriceHistory(ctx context.Context, symbol string, from, to time.Time) ([]entities.CryptoPrice, error) {
	r.logger.Debug("Retrieving price history", "symbol", symbol, "from", from, "to", to)

	var priceData []entities.CryptoPrice
	if err := r.db.WithContext(ctx).
		Where("symbol = ? AND created_at BETWEEN ? AND ?", symbol, from, to).
		Order("created_at ASC").
		Find(&priceData).Error; err != nil {
		r.logger.Error("Failed to retrieve price history", "error", err, "symbol", symbol)
		return nil, errors.Wrap(err, errors.ErrorTypeInternal, "failed to retrieve price history")
	}

	return priceData, nil
}

// GetLatestPrice retrieves the latest price for a symbol
func (r *marketDataRepository) GetLatestPrice(ctx context.Context, symbol string) (*entities.CryptoPrice, error) {
	r.logger.Debug("Retrieving latest price", "symbol", symbol)

	var priceData entities.CryptoPrice
	if err := r.db.WithContext(ctx).
		Where("symbol = ?", symbol).
		Order("created_at DESC").
		First(&priceData).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NotFound("price_data")
		}
		r.logger.Error("Failed to retrieve latest price", "error", err, "symbol", symbol)
		return nil, errors.Wrap(err, errors.ErrorTypeInternal, "failed to retrieve latest price")
	}

	return &priceData, nil
}

// StoreDominanceData saves Bitcoin dominance data to the database
func (r *marketDataRepository) StoreDominanceData(ctx context.Context, dominanceData *entities.BitcoinDominance) error {
	r.logger.Debug("Saving dominance data", "dominance", dominanceData.CurrentDominance, "source", dominanceData.DataSource)

	if err := r.db.WithContext(ctx).Create(dominanceData).Error; err != nil {
		r.logger.Error("Failed to save dominance data", "error", err, "dominance", dominanceData.CurrentDominance)
		return errors.Wrap(err, errors.ErrorTypeInternal, "failed to save dominance data")
	}

	return nil
}

// GetDominanceHistory retrieves historical Bitcoin dominance data
func (r *marketDataRepository) GetDominanceHistory(ctx context.Context, from, to time.Time) ([]entities.BitcoinDominance, error) {
	r.logger.Debug("Retrieving dominance history", "from", from, "to", to)

	var dominanceData []entities.BitcoinDominance
	if err := r.db.WithContext(ctx).
		Where("created_at BETWEEN ? AND ?", from, to).
		Order("created_at ASC").
		Find(&dominanceData).Error; err != nil {
		r.logger.Error("Failed to retrieve dominance history", "error", err)
		return nil, errors.Wrap(err, errors.ErrorTypeInternal, "failed to retrieve dominance history")
	}

	return dominanceData, nil
}

// GetLatestDominance retrieves the latest Bitcoin dominance data
func (r *marketDataRepository) GetLatestDominance(ctx context.Context) (*entities.BitcoinDominance, error) {
	r.logger.Debug("Retrieving latest dominance data")

	var dominanceData entities.BitcoinDominance
	if err := r.db.WithContext(ctx).
		Order("created_at DESC").
		First(&dominanceData).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NotFound("dominance_data")
		}
		r.logger.Error("Failed to retrieve latest dominance data", "error", err)
		return nil, errors.Wrap(err, errors.ErrorTypeInternal, "failed to retrieve latest dominance data")
	}

	return &dominanceData, nil
}

// SaveMarketMetrics saves market metrics to the database
func (r *marketDataRepository) SaveMarketMetrics(ctx context.Context, metrics *entities.MarketMetrics) error {
	r.logger.Debug("Saving market metrics")

	if err := r.db.WithContext(ctx).Create(metrics).Error; err != nil {
		r.logger.Error("Failed to save market metrics", "error", err)
		return errors.Wrap(err, errors.ErrorTypeInternal, "failed to save market metrics")
	}

	return nil
}

// GetMarketMetricsHistory retrieves historical market metrics
func (r *marketDataRepository) GetMarketMetricsHistory(ctx context.Context, from, to time.Time) ([]entities.MarketMetrics, error) {
	r.logger.Debug("Retrieving market metrics history", "from", from, "to", to)

	var metrics []entities.MarketMetrics
	if err := r.db.WithContext(ctx).
		Where("created_at BETWEEN ? AND ?", from, to).
		Order("created_at ASC").
		Find(&metrics).Error; err != nil {
		r.logger.Error("Failed to retrieve market metrics history", "error", err)
		return nil, errors.Wrap(err, errors.ErrorTypeInternal, "failed to retrieve market metrics history")
	}

	return metrics, nil
}

// GetLatestMarketMetrics retrieves the latest market metrics
func (r *marketDataRepository) GetLatestMarketMetrics(ctx context.Context) (*entities.MarketMetrics, error) {
	r.logger.Debug("Retrieving latest market metrics")

	var metrics entities.MarketMetrics
	if err := r.db.WithContext(ctx).
		Order("created_at DESC").
		First(&metrics).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NotFound("market_metrics")
		}
		r.logger.Error("Failed to retrieve latest market metrics", "error", err)
		return nil, errors.Wrap(err, errors.ErrorTypeInternal, "failed to retrieve latest market metrics")
	}

	return &metrics, nil
}