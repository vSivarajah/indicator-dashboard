package services

import (
	"context"
	"crypto-indicator-dashboard/internal/domain/entities"
)

// MarketDataService defines the interface for market data operations
type MarketDataService interface {
	// GetCryptoPrices retrieves current cryptocurrency prices
	GetCryptoPrices(ctx context.Context, symbols []string) (map[string]*entities.CryptoPrice, error)
	
	// GetBitcoinDominance retrieves current Bitcoin dominance data
	GetBitcoinDominance(ctx context.Context) (*entities.BitcoinDominance, error)
	
	// GetMultipleCryptoPrices gets prices for common cryptocurrencies
	GetMultipleCryptoPrices(ctx context.Context) (map[string]*entities.CryptoPrice, error)
	
	// GetTopCryptoPrices gets prices for top N cryptocurrencies by market cap
	GetTopCryptoPrices(ctx context.Context, count int) (map[string]*entities.CryptoPrice, error)
	
	// RefreshAllMarketData refreshes all market data from external sources
	RefreshAllMarketData(ctx context.Context) error
	
	// HealthCheck performs health checks on all external data sources
	HealthCheck(ctx context.Context) map[string]error
}

// CacheService defines the interface for caching operations
type CacheService interface {
	// GetOrSet gets a value from cache or sets it using the provided function
	GetOrSet(ctx context.Context, key string, dest interface{}, expiration interface{}, setFunc func() (interface{}, error)) error
	
	// Get retrieves a value from cache
	Get(ctx context.Context, key string, dest interface{}) error
	
	// Set stores a value in cache
	Set(ctx context.Context, key string, value interface{}, expiration interface{}) error
	
	// Delete removes a value from cache
	Delete(ctx context.Context, key string) error
	
	// Exists checks if a key exists in cache
	Exists(ctx context.Context, key string) bool
	
	// Clear clears all cache entries
	Clear(ctx context.Context) error
	
	// HealthCheck performs a health check on the cache service
	HealthCheck(ctx context.Context) error
}