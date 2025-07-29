package repositories

import (
	"context"
	"crypto-indicator-dashboard/internal/domain/entities"
	"time"
)

// IndicatorRepository defines the interface for indicator data operations
type IndicatorRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, indicator *entities.Indicator) error
	GetByID(ctx context.Context, id uint) (*entities.Indicator, error)
	GetByName(ctx context.Context, name string) (*entities.Indicator, error)
	GetByType(ctx context.Context, indicatorType string) ([]entities.Indicator, error)
	Update(ctx context.Context, indicator *entities.Indicator) error
	Delete(ctx context.Context, id uint) error
	
	// Historical data operations
	GetHistoricalData(ctx context.Context, name string, from, to time.Time) ([]entities.Indicator, error)
	GetLatest(ctx context.Context, name string) (*entities.Indicator, error)
	GetLatestByType(ctx context.Context, indicatorType string) ([]entities.Indicator, error)
	
	// Bulk operations
	BulkCreate(ctx context.Context, indicators []entities.Indicator) error
	CleanupOldData(ctx context.Context, olderThan time.Time) error
}

// MarketDataRepository defines the interface for market data operations
type MarketDataRepository interface {
	// Crypto price data operations
	StorePriceData(ctx context.Context, priceData *entities.CryptoPrice) error
	GetPriceHistory(ctx context.Context, symbol string, from, to time.Time) ([]entities.CryptoPrice, error)
	GetLatestPrice(ctx context.Context, symbol string) (*entities.CryptoPrice, error)
	
	// Bitcoin dominance operations
	StoreDominanceData(ctx context.Context, dominanceData *entities.BitcoinDominance) error
	GetDominanceHistory(ctx context.Context, from, to time.Time) ([]entities.BitcoinDominance, error)
	GetLatestDominance(ctx context.Context) (*entities.BitcoinDominance, error)
	
	// Market metrics operations
	SaveMarketMetrics(ctx context.Context, metrics *entities.MarketMetrics) error
	GetMarketMetricsHistory(ctx context.Context, from, to time.Time) ([]entities.MarketMetrics, error)
	GetLatestMarketMetrics(ctx context.Context) (*entities.MarketMetrics, error)
}