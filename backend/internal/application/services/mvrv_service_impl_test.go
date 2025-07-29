package services

import (
	"context"
	"crypto-indicator-dashboard/internal/domain/entities"
	"crypto-indicator-dashboard/internal/testutil"
	"crypto-indicator-dashboard/pkg/errors"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// MVRVServiceTestSuite provides test suite for MVRV service
type MVRVServiceTestSuite struct {
	suite.Suite
	service           *mvrvServiceImpl
	mockIndicatorRepo *testutil.MockIndicatorRepository
	mockMarketRepo    *testutil.MockMarketDataRepository
	mockCache         *testutil.MockInfrastructureCacheService
	testData          *testutil.TestData
	server            *httptest.Server
}

func (suite *MVRVServiceTestSuite) SetupTest() {
	suite.mockIndicatorRepo = &testutil.MockIndicatorRepository{}
	suite.mockMarketRepo = &testutil.MockMarketDataRepository{}
	suite.mockCache = testutil.NewMockInfrastructureCacheService()
	suite.testData = testutil.NewTestData()

	// Create test HTTP server
	suite.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v3/coins/bitcoin":
			suite.handleBitcoinDataRequest(w, r)
		default:
			http.NotFound(w, r)
		}
	}))

	// Create service with test dependencies using mock server
	suite.service = NewMVRVServiceWithBaseURL(
		suite.mockIndicatorRepo,
		suite.mockMarketRepo,
		suite.mockCache,
		testutil.NewTestDB(suite.T()).Logger,
		suite.server.URL, // Use mock server URL instead of real API
	).(*mvrvServiceImpl)
}

func (suite *MVRVServiceTestSuite) TearDownTest() {
	if suite.server != nil {
		suite.server.Close()
	}
}

func (suite *MVRVServiceTestSuite) handleBitcoinDataRequest(w http.ResponseWriter, r *http.Request) {
	mockData := CoinGeckoBitcoinData{
		MarketData: struct {
			CurrentPrice struct {
				USD float64 `json:"usd"`
			} `json:"current_price"`
			MarketCap struct {
				USD float64 `json:"usd"`
			} `json:"market_cap"`
			CirculatingSupply float64 `json:"circulating_supply"`
		}{
			CurrentPrice: struct {
				USD float64 `json:"usd"`
			}{USD: 43000.0},
			MarketCap: struct {
				USD float64 `json:"usd"`
			}{USD: 850000000000.0},
			CirculatingSupply: 19800000.0,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(mockData)
}

func (suite *MVRVServiceTestSuite) TestCalculate_Success() {
	ctx := context.Background()

	// Mock cache miss - return mock data directly without calling fetch function
	suite.mockCache.On("GetOrSet", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil).Run(func(args mock.Arguments) {
		// Simulate successful cache operation by setting mock Bitcoin data
		dest := args.Get(1)
		if destPtr, ok := dest.(*CoinGeckoBitcoinData); ok {
			*destPtr = CoinGeckoBitcoinData{
				MarketData: struct {
					CurrentPrice struct {
						USD float64 `json:"usd"`
					} `json:"current_price"`
					MarketCap struct {
						USD float64 `json:"usd"`
					} `json:"market_cap"`
					CirculatingSupply float64 `json:"circulating_supply"`
				}{
					CurrentPrice: struct {
						USD float64 `json:"usd"`
					}{USD: 43000.0},
					MarketCap: struct {
						USD float64 `json:"usd"`
					}{USD: 850000000000.0},
					CirculatingSupply: 19800000.0,
				},
			}
		}
	})

	// Mock successful database save
	suite.mockIndicatorRepo.On("Create", ctx, mock.AnythingOfType("*entities.Indicator")).Return(nil)

	// Execute test
	result, err := suite.service.Calculate(ctx, nil)

	// Assertions
	require.NoError(suite.T(), err, "Calculate should not return error")
	require.NotNil(suite.T(), result, "Result should not be nil")

	assert.Equal(suite.T(), "mvrv", result.Name)
	assert.Equal(suite.T(), "market", result.Type)
	assert.True(suite.T(), result.Value >= 0, "MVRV Z-Score should be calculated (can be 0)")
	assert.NotEmpty(suite.T(), result.Status, "Status should be set")
	assert.NotEmpty(suite.T(), result.RiskLevel, "Risk level should be set")
	assert.True(suite.T(), result.Confidence > 0, "Confidence should be positive")
	testutil.AssertRecentTime(suite.T(), result.Timestamp, "Timestamp should be recent")

	// Verify metadata
	assert.Contains(suite.T(), result.Metadata, "mvrv_ratio")
	assert.Contains(suite.T(), result.Metadata, "market_cap")
	assert.Contains(suite.T(), result.Metadata, "realized_cap")
	assert.Contains(suite.T(), result.Metadata, "price")
	assert.Contains(suite.T(), result.Metadata, "z_score")

	// Verify mocks were called
	suite.mockCache.AssertExpectations(suite.T())
	suite.mockIndicatorRepo.AssertExpectations(suite.T())
}

func (suite *MVRVServiceTestSuite) TestCalculate_APIFailure() {
	ctx := context.Background()

	// Mock cache miss
	suite.mockCache.On("GetOrSet", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(fmt.Errorf("API unavailable"))

	// Execute test
	result, err := suite.service.Calculate(ctx, nil)

	// Should return fallback data, not error
	require.NoError(suite.T(), err, "Calculate should return fallback data on API failure")
	require.NotNil(suite.T(), result, "Result should not be nil")

	// Verify fallback indicators
	assert.Equal(suite.T(), "mvrv", result.Name)
	assert.Equal(suite.T(), float64(0.5), result.Value) // Fallback Z-score
	assert.Equal(suite.T(), float64(0.3), result.Confidence) // Low confidence for fallback
	assert.Contains(suite.T(), result.Metadata, "fallback")
	assert.True(suite.T(), result.Metadata["fallback"].(bool))

	// No database save expected for fallback - it returns the data directly
}

func (suite *MVRVServiceTestSuite) TestGetLatest_DatabaseHit() {
	ctx := context.Background()
	expectedIndicator := suite.testData.SampleIndicator()
	expectedIndicator.Timestamp = time.Now().Add(-30 * time.Minute) // Fresh data

	suite.mockIndicatorRepo.On("GetLatest", ctx, "mvrv").Return(expectedIndicator, nil)

	result, err := suite.service.GetLatest(ctx)

	require.NoError(suite.T(), err)
	testutil.AssertIndicatorEqual(suite.T(), expectedIndicator, result)
	suite.mockIndicatorRepo.AssertExpectations(suite.T())
}

func (suite *MVRVServiceTestSuite) TestGetLatest_StaleData() {
	ctx := context.Background()
	staleIndicator := suite.testData.SampleIndicator()
	staleIndicator.Timestamp = time.Now().Add(-2 * time.Hour) // Stale data

	suite.mockIndicatorRepo.On("GetLatest", ctx, "mvrv").Return(staleIndicator, nil)

	// Mock fresh calculation
	suite.mockCache.On("GetOrSet", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil)
	suite.mockIndicatorRepo.On("Create", ctx, mock.AnythingOfType("*entities.Indicator")).Return(nil)

	result, err := suite.service.GetLatest(ctx)

	require.NoError(suite.T(), err)
	assert.True(suite.T(), result.Timestamp.After(staleIndicator.Timestamp))
	suite.mockIndicatorRepo.AssertExpectations(suite.T())
}

func (suite *MVRVServiceTestSuite) TestGetLatest_NotFound() {
	ctx := context.Background()

	suite.mockIndicatorRepo.On("GetLatest", ctx, "mvrv").
		Return((*entities.Indicator)(nil), errors.NewNotFoundError("indicator", "mvrv"))

	// Mock fresh calculation since not found
	suite.mockCache.On("GetOrSet", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil)
	suite.mockIndicatorRepo.On("Create", ctx, mock.AnythingOfType("*entities.Indicator")).Return(nil)

	result, err := suite.service.GetLatest(ctx)

	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), result)
	assert.Equal(suite.T(), "mvrv", result.Name)
	suite.mockIndicatorRepo.AssertExpectations(suite.T())
}

func (suite *MVRVServiceTestSuite) TestGetHistoricalData_Success() {
	ctx := context.Background()
	period := "30d"
	expectedData := []entities.Indicator{
		*suite.testData.SampleIndicator(),
		*suite.testData.SampleIndicator(),
	}
	
	// Set different timestamps for historical data
	expectedData[0].Timestamp = time.Now().Add(-24 * time.Hour)
	expectedData[1].Timestamp = time.Now().Add(-48 * time.Hour)

	from := time.Now().AddDate(0, 0, -30)
	to := time.Now()

	suite.mockIndicatorRepo.On("GetHistoricalData", ctx, "mvrv", mock.MatchedBy(func(t time.Time) bool {
		return t.Before(from.Add(time.Minute)) && t.After(from.Add(-time.Minute))
	}), mock.MatchedBy(func(t time.Time) bool {
		return t.Before(to.Add(time.Minute)) && t.After(to.Add(-time.Minute))
	})).Return(expectedData, nil)

	result, err := suite.service.GetHistoricalData(ctx, period)

	require.NoError(suite.T(), err)
	assert.Len(suite.T(), result, 2)
	assert.Equal(suite.T(), expectedData[0].ID, result[0].ID)
	assert.Equal(suite.T(), expectedData[1].ID, result[1].ID)
	suite.mockIndicatorRepo.AssertExpectations(suite.T())
}

func (suite *MVRVServiceTestSuite) TestAssessMVRVRisk_AllLevels() {
	testCases := []struct {
		name           string
		zScore         float64
		expectedRisk   string
		expectedStatus string
	}{
		{"Extreme High", 8.0, "extreme_high", "EXTREME: Historically top of cycle - Strong sell signal"},
		{"High", 4.0, "high", "HIGH: Approaching cycle top - Consider taking profits"},
		{"Medium", 2.0, "medium", "MEDIUM: Testing resistance - Monitor closely"},
		{"Low Positive", 1.0, "low", "LOW: Above average valuation - Neutral zone"},
		{"Low Neutral", 0.0, "low", "LOW: Fair value range - Accumulation zone"},
		{"Low Negative", -1.0, "low", "LOW: Below average - Good buying opportunity"},
		{"Extreme Low", -2.0, "extreme_low", "EXTREME: Historically bottom of cycle - Strong buy signal"},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			riskLevel, status := suite.service.assessMVRVRisk(tc.zScore)
			assert.Equal(t, tc.expectedRisk, riskLevel, "Risk level should match for Z-Score: %f", tc.zScore)
			assert.Equal(t, tc.expectedStatus, status, "Status should match for Z-Score: %f", tc.zScore)
		})
	}
}

func (suite *MVRVServiceTestSuite) TestCalculateZScores_Success() {
	// Create test data with valid MVRV ratios
	data := []MVRVData{
		{Date: time.Now().Add(-3 * 24 * time.Hour), MVRVRatio: 1.2},
		{Date: time.Now().Add(-2 * 24 * time.Hour), MVRVRatio: 1.5},
		{Date: time.Now().Add(-1 * 24 * time.Hour), MVRVRatio: 1.8},
		{Date: time.Now(), MVRVRatio: 2.0},
	}

	suite.service.calculateZScores(data)

	// Verify Z-scores were calculated
	for i, d := range data {
		assert.False(suite.T(), isNaN(d.MVRVZScore), "Z-Score should not be NaN for data point %d", i)
		assert.False(suite.T(), isInf(d.MVRVZScore), "Z-Score should not be Inf for data point %d", i)
	}

	// Verify statistical properties (mean Z-score should be close to 0)
	sum := 0.0
	for _, d := range data {
		sum += d.MVRVZScore
	}
	meanZScore := sum / float64(len(data))
	assert.InDelta(suite.T(), 0.0, meanZScore, 0.1, "Mean Z-Score should be close to 0")
}

func (suite *MVRVServiceTestSuite) TestCalculateZScores_EdgeCases() {
	testCases := []struct {
		name     string
		data     []MVRVData
		expectOK bool
	}{
		{
			name:     "Empty data",
			data:     []MVRVData{},
			expectOK: true, // Should handle gracefully
		},
		{
			name: "Single data point",
			data: []MVRVData{
				{Date: time.Now(), MVRVRatio: 1.5},
			},
			expectOK: true,
		},
		{
			name: "Invalid ratios",
			data: []MVRVData{
				{Date: time.Now(), MVRVRatio: 0.0},
				{Date: time.Now(), MVRVRatio: -1.0},
			},
			expectOK: true, // Should filter out invalid values
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			// Should not panic
			assert.NotPanics(t, func() {
				suite.service.calculateZScores(tc.data)
			})

			// Verify Z-scores are reasonable
			for i, d := range tc.data {
				assert.False(t, isNaN(d.MVRVZScore), "Z-Score should not be NaN for data point %d", i)
				assert.False(t, isInf(d.MVRVZScore), "Z-Score should not be Inf for data point %d", i)
			}
		})
	}
}

func (suite *MVRVServiceTestSuite) TestGenerateHistoricalMVRVData() {
	mockBitcoinData := &CoinGeckoBitcoinData{
		MarketData: struct {
			CurrentPrice struct {
				USD float64 `json:"usd"`
			} `json:"current_price"`
			MarketCap struct {
				USD float64 `json:"usd"`
			} `json:"market_cap"`
			CirculatingSupply float64 `json:"circulating_supply"`
		}{
			CurrentPrice:      struct{ USD float64 `json:"usd"` }{USD: 43000.0},
			MarketCap:         struct{ USD float64 `json:"usd"` }{USD: 850000000000.0},
			CirculatingSupply: 19800000.0,
		},
	}

	data := suite.service.generateHistoricalMVRVData(mockBitcoinData)

	// Verify data structure
	assert.Len(suite.T(), data, 366, "Should generate 366 data points (365 days + today)")

	// Verify data quality
	for i, d := range data {
		assert.True(suite.T(), d.Price > 0, "Price should be positive for data point %d", i)
		assert.True(suite.T(), d.MarketCap > 0, "Market cap should be positive for data point %d", i)
		assert.True(suite.T(), d.RealizedCap > 0, "Realized cap should be positive for data point %d", i)
		assert.True(suite.T(), d.MVRVRatio > 0, "MVRV ratio should be positive for data point %d", i)
		assert.True(suite.T(), d.MVRVRatio <= 10, "MVRV ratio should be reasonable for data point %d", i)
		assert.False(suite.T(), isNaN(d.MVRVZScore), "Z-Score should not be NaN for data point %d", i)
		assert.False(suite.T(), isInf(d.MVRVZScore), "Z-Score should not be Inf for data point %d", i)
	}

	// Verify chronological order
	for i := 1; i < len(data); i++ {
		assert.True(suite.T(), data[i].Date.After(data[i-1].Date), "Data should be in chronological order")
	}
}

// Benchmark tests run outside of the test suite
func BenchmarkMVRVCalculate(b *testing.B) {
	// Set up test dependencies
	mockIndicatorRepo := &testutil.MockIndicatorRepository{}
	mockMarketRepo := &testutil.MockMarketDataRepository{}
	mockCache := testutil.NewMockInfrastructureCacheService()
	testDB := testutil.NewTestDB(&testing.T{})
	defer testDB.Cleanup()

	service := NewMVRVServiceWithBaseURL(
		mockIndicatorRepo,
		mockMarketRepo,
		mockCache,
		testDB.Logger,
		"http://localhost:8999", // Use dummy URL for benchmark (won't be called)
	).(*mvrvServiceImpl)

	ctx := context.Background()
	mockCache.On("GetOrSet", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil)
	mockIndicatorRepo.On("Create", ctx, mock.AnythingOfType("*entities.Indicator")).Return(nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.Calculate(ctx, nil)
	}
}

// Test suite runner
func TestMVRVServiceTestSuite(t *testing.T) {
	suite.Run(t, new(MVRVServiceTestSuite))
}

// Table-driven tests for risk assessment
func TestMVRVRiskAssessment(t *testing.T) {
	service := &mvrvServiceImpl{}

	tests := []struct {
		name           string
		zScore         float64
		expectedRisk   string
		shouldContain  string
	}{
		{"Bubble Territory", 10.0, "extreme_high", "EXTREME"},
		{"Bull Market Peak", 5.0, "high", "HIGH"},
		{"Resistance Test", 2.5, "medium", "MEDIUM"},
		{"Fair Value", 0.8, "low", "LOW"},
		{"Accumulation Zone", -0.8, "low", "buying opportunity"},
		{"Strong Buy Signal", -2.5, "extreme_low", "Strong buy signal"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			riskLevel, status := service.assessMVRVRisk(tt.zScore)
			assert.Equal(t, tt.expectedRisk, riskLevel)
			assert.Contains(t, status, tt.shouldContain)
		})
	}
}

// Helper functions for math checks
func isNaN(f float64) bool {
	return f != f
}

func isInf(f float64) bool {
	return f > 1e308 || f < -1e308
}