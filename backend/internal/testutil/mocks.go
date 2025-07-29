package testutil

import (
	"context"
	"crypto-indicator-dashboard/internal/domain/entities"
	"crypto-indicator-dashboard/internal/infrastructure/external"
	"time"

	"github.com/stretchr/testify/mock"
)

// MockHTTPClient is a mock HTTP client for testing
type MockHTTPClient struct {
	mock.Mock
}

// MockCacheService is a mock cache service for testing
type MockCacheService struct {
	mock.Mock
	data map[string]interface{}
}

// NewMockCacheService creates a new mock cache service
func NewMockCacheService() *MockCacheService {
	return &MockCacheService{
		data: make(map[string]interface{}),
	}
}

// Get retrieves a value from the mock cache
func (m *MockCacheService) Get(ctx context.Context, key string, dest interface{}) error {
	args := m.Called(ctx, key, dest)
	return args.Error(0)
}

// Set stores a value in the mock cache
func (m *MockCacheService) Set(ctx context.Context, key string, value interface{}, expiration interface{}) error {
	args := m.Called(ctx, key, value, expiration)
	m.data[key] = value
	return args.Error(0)
}

// Delete removes a value from the mock cache
func (m *MockCacheService) Delete(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	delete(m.data, key)
	return args.Error(0)
}

// GetOrSet gets a value or sets it if not found
func (m *MockCacheService) GetOrSet(ctx context.Context, key string, dest interface{}, expiration interface{}, setFunc func() (interface{}, error)) error {
	args := m.Called(ctx, key, dest, expiration, setFunc)
	
	if args.Error(0) == nil {
		// If no error, call the fetch function and store result
		if data, exists := m.data[key]; exists {
			// Simulate cache hit - copy data to dest if possible
			if ptr, ok := dest.(*interface{}); ok {
				*ptr = data
			}
		} else {
			// Simulate cache miss
			fetchedData, err := setFunc()
			if err != nil {
				return err
			}
			m.data[key] = fetchedData
			if ptr, ok := dest.(*interface{}); ok {
				*ptr = fetchedData
			}
		}
	}
	
	return args.Error(0)
}

// Exists checks if a key exists
func (m *MockCacheService) Exists(ctx context.Context, key string) bool {
	args := m.Called(ctx, key)
	return args.Bool(0)
}

// TTL returns time-to-live for a key
func (m *MockCacheService) TTL(ctx context.Context, key string) (time.Duration, error) {
	args := m.Called(ctx, key)
	return args.Get(0).(time.Duration), args.Error(1)
}

// Clear removes all cached values
func (m *MockCacheService) Clear(ctx context.Context) error {
	args := m.Called(ctx)
	m.data = make(map[string]interface{})
	return args.Error(0)
}

// Keys returns all keys matching a pattern
func (m *MockCacheService) Keys(ctx context.Context, pattern string) ([]string, error) {
	args := m.Called(ctx, pattern)
	return args.Get(0).([]string), args.Error(1)
}

// Size returns the number of keys in cache
func (m *MockCacheService) Size(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

// HealthCheck performs a health check
func (m *MockCacheService) HealthCheck(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// FlushAll removes all keys from cache
func (m *MockCacheService) FlushAll(ctx context.Context) error {
	args := m.Called(ctx)
	m.data = make(map[string]interface{})
	return args.Error(0)
}

// MockIndicatorRepository is a mock implementation of IndicatorRepository
type MockIndicatorRepository struct {
	mock.Mock
}

func (m *MockIndicatorRepository) Create(ctx context.Context, indicator *entities.Indicator) error {
	args := m.Called(ctx, indicator)
	return args.Error(0)
}

func (m *MockIndicatorRepository) GetByID(ctx context.Context, id uint) (*entities.Indicator, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Indicator), args.Error(1)
}

func (m *MockIndicatorRepository) GetByName(ctx context.Context, name string) (*entities.Indicator, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Indicator), args.Error(1)
}

func (m *MockIndicatorRepository) GetByType(ctx context.Context, indicatorType string) ([]entities.Indicator, error) {
	args := m.Called(ctx, indicatorType)
	return args.Get(0).([]entities.Indicator), args.Error(1)
}

func (m *MockIndicatorRepository) GetLatest(ctx context.Context, name string) (*entities.Indicator, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Indicator), args.Error(1)
}

func (m *MockIndicatorRepository) GetLatestByType(ctx context.Context, indicatorType string) ([]entities.Indicator, error) {
	args := m.Called(ctx, indicatorType)
	return args.Get(0).([]entities.Indicator), args.Error(1)
}

func (m *MockIndicatorRepository) GetHistoricalData(ctx context.Context, name string, from, to time.Time) ([]entities.Indicator, error) {
	args := m.Called(ctx, name, from, to)
	return args.Get(0).([]entities.Indicator), args.Error(1)
}

func (m *MockIndicatorRepository) Update(ctx context.Context, indicator *entities.Indicator) error {
	args := m.Called(ctx, indicator)
	return args.Error(0)
}

func (m *MockIndicatorRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockIndicatorRepository) BulkCreate(ctx context.Context, indicators []entities.Indicator) error {
	args := m.Called(ctx, indicators)
	return args.Error(0)
}

func (m *MockIndicatorRepository) CleanupOldData(ctx context.Context, olderThan time.Time) error {
	args := m.Called(ctx, olderThan)
	return args.Error(0)
}

// MockMarketDataRepository is a mock implementation of MarketDataRepository
type MockMarketDataRepository struct {
	mock.Mock
}

// Crypto price data operations
func (m *MockMarketDataRepository) StorePriceData(ctx context.Context, priceData *entities.CryptoPrice) error {
	args := m.Called(ctx, priceData)
	return args.Error(0)
}

func (m *MockMarketDataRepository) GetPriceHistory(ctx context.Context, symbol string, from, to time.Time) ([]entities.CryptoPrice, error) {
	args := m.Called(ctx, symbol, from, to)
	return args.Get(0).([]entities.CryptoPrice), args.Error(1)
}

func (m *MockMarketDataRepository) GetLatestPrice(ctx context.Context, symbol string) (*entities.CryptoPrice, error) {
	args := m.Called(ctx, symbol)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.CryptoPrice), args.Error(1)
}

// Bitcoin dominance operations
func (m *MockMarketDataRepository) StoreDominanceData(ctx context.Context, dominanceData *entities.BitcoinDominance) error {
	args := m.Called(ctx, dominanceData)
	return args.Error(0)
}

func (m *MockMarketDataRepository) GetDominanceHistory(ctx context.Context, from, to time.Time) ([]entities.BitcoinDominance, error) {
	args := m.Called(ctx, from, to)
	return args.Get(0).([]entities.BitcoinDominance), args.Error(1)
}

func (m *MockMarketDataRepository) GetLatestDominance(ctx context.Context) (*entities.BitcoinDominance, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.BitcoinDominance), args.Error(1)
}

// Market metrics operations
func (m *MockMarketDataRepository) SaveMarketMetrics(ctx context.Context, metrics *entities.MarketMetrics) error {
	args := m.Called(ctx, metrics)
	return args.Error(0)
}

func (m *MockMarketDataRepository) GetMarketMetricsHistory(ctx context.Context, from, to time.Time) ([]entities.MarketMetrics, error) {
	args := m.Called(ctx, from, to)
	return args.Get(0).([]entities.MarketMetrics), args.Error(1)
}

func (m *MockMarketDataRepository) GetLatestMarketMetrics(ctx context.Context) (*entities.MarketMetrics, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.MarketMetrics), args.Error(1)
}

// MockCoinCapClient is a mock implementation of CoinCap client
type MockCoinCapClient struct {
	mock.Mock
}

func (m *MockCoinCapClient) GetAssets(limit int) (*external.AssetsResponse, error) {
	args := m.Called(limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*external.AssetsResponse), args.Error(1)
}

func (m *MockCoinCapClient) GetAsset(assetID string) (*external.AssetResponse, error) {
	args := m.Called(assetID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*external.AssetResponse), args.Error(1)
}

func (m *MockCoinCapClient) GetBitcoinPrice() (float64, error) {
	args := m.Called()
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockCoinCapClient) GetGlobalMarketData() (map[string]interface{}, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (m *MockCoinCapClient) HealthCheck() error {
	args := m.Called()
	return args.Error(0)
}

// TestData provides common test data for tests
type TestData struct{}

// NewTestData creates a new test data provider
func NewTestData() *TestData {
	return &TestData{}
}

// SampleIndicator returns a sample indicator for testing
func (td *TestData) SampleIndicator() *entities.Indicator {
	return &entities.Indicator{
		ID:          1,
		Name:        "mvrv",
		Type:        "market",
		Value:       2.45,
		StringValue: "2.45",
		Change:      "+0.12",
		RiskLevel:   "medium",
		Status:      "MEDIUM: Testing resistance",
		Description: "MVRV Z-Score test indicator",
		Source:      "test",
		Confidence:  0.85,
		Metadata: map[string]interface{}{
			"mvrv_ratio":   1.8,
			"market_cap":   850000000000.0,
			"realized_cap": 472222222222.0,
			"price":        43000.0,
			"z_score":      2.45,
		},
		Timestamp: time.Now(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// SampleMarketData returns sample market data for testing
func (td *TestData) SampleMarketData() *entities.MarketData {
	return &entities.MarketData{
		ID:            1,
		Symbol:        "BTC",
		Name:          "Bitcoin",
		Price:         43000.0,
		MarketCap:     850000000000.0,
		Volume24h:     25000000000.0,
		Change24h:     2.5,
		Change7d:      5.2,
		Change30d:     12.8,
		Dominance:     54.2,
		CircSupply:    19800000.0,
		MaxSupply:     21000000.0,
		Source:        "test",
		Confidence:    0.95,
		LastUpdated:   time.Now(),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
}

// SampleCoinCapAssetResponse returns a sample CoinCap asset response
func (td *TestData) SampleCoinCapAssetResponse() *external.AssetResponse {
	return &external.AssetResponse{
		Data: external.Asset{
			ID:                "bitcoin",
			Rank:              "1",
			Symbol:            "BTC",
			Name:              "Bitcoin",
			Supply:            "19800000.0000000000000000",
			MaxSupply:         stringPtr("21000000.0000000000000000"),
			MarketCapUSD:      "850000000000.0000000000000000",
			VolumeUSD24Hr:     "25000000000.0000000000000000",
			PriceUSD:          "43000.0000000000000000",
			ChangePercent24Hr: "2.5000000000000000",
			VWAP24Hr:          stringPtr("42800.0000000000000000"),
		},
		Timestamp: time.Now().Unix(),
	}
}

// MockInfrastructureCacheService is a mock for the infrastructure cache service interface
type MockInfrastructureCacheService struct {
	mock.Mock
	data map[string]interface{}
}

// NewMockInfrastructureCacheService creates a new mock infrastructure cache service
func NewMockInfrastructureCacheService() *MockInfrastructureCacheService {
	return &MockInfrastructureCacheService{
		data: make(map[string]interface{}),
	}
}

// Get retrieves a value from the mock cache
func (m *MockInfrastructureCacheService) Get(ctx context.Context, key string, dest interface{}) error {
	args := m.Called(ctx, key, dest)
	return args.Error(0)
}

// Set stores a value in the mock cache
func (m *MockInfrastructureCacheService) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	args := m.Called(ctx, key, value, expiration)
	m.data[key] = value
	return args.Error(0)
}

// Delete removes a value from the mock cache
func (m *MockInfrastructureCacheService) Delete(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	delete(m.data, key)
	return args.Error(0)
}

// Exists checks if a key exists
func (m *MockInfrastructureCacheService) Exists(ctx context.Context, key string) (bool, error) {
	args := m.Called(ctx, key)
	return args.Bool(0), args.Error(1)
}

// FlushAll removes all cached values
func (m *MockInfrastructureCacheService) FlushAll(ctx context.Context) error {
	args := m.Called(ctx)
	m.data = make(map[string]interface{})
	return args.Error(0)
}

// GetOrSet gets a value or sets it if not found
func (m *MockInfrastructureCacheService) GetOrSet(ctx context.Context, key string, dest interface{}, fetcher func() (interface{}, error), expiration time.Duration) error {
	args := m.Called(ctx, key, dest, fetcher, expiration)
	
	if args.Error(0) == nil {
		// If no error, call the fetch function and store result
		if data, exists := m.data[key]; exists {
			// Simulate cache hit - copy data to dest if possible
			if ptr, ok := dest.(*interface{}); ok {
				*ptr = data
			}
		} else {
			// Simulate cache miss
			fetchedData, err := fetcher()
			if err != nil {
				return err
			}
			m.data[key] = fetchedData
			if ptr, ok := dest.(*interface{}); ok {
				*ptr = fetchedData
			}
		}
	}
	
	return args.Error(0)
}

// Helper function to create string pointer
func stringPtr(s string) *string {
	return &s
}