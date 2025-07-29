package services

import (
	"context"
	"crypto-indicator-dashboard/internal/testutil"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
)

// BenchmarkMVRVService benchmarks MVRV service operations
func BenchmarkMVRVService(b *testing.B) {
	// Setup test environment
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

	// Mock successful operations
	mockCache.On("GetOrSet", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil)
	mockIndicatorRepo.On("Create", ctx, mock.AnythingOfType("*entities.Indicator")).Return(nil)

	b.Run("Calculate", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = service.Calculate(ctx, nil)
		}
	})

	b.Run("GetLatest", func(b *testing.B) {
		// Setup test data
		indicator := testutil.NewTestData().SampleIndicator()
		mockIndicatorRepo.On("GetLatest", ctx, "mvrv").Return(indicator, nil)

		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = service.GetLatest(ctx)
		}
	})

	b.Run("GenerateHistoricalData", func(b *testing.B) {
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

		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = service.generateHistoricalMVRVData(mockBitcoinData)
		}
	})

	b.Run("CalculateZScores", func(b *testing.B) {
		// Create test data
		data := make([]MVRVData, 100)
		for i := 0; i < 100; i++ {
			data[i] = MVRVData{
				Date:      time.Now().Add(-time.Duration(i) * 24 * time.Hour),
				MVRVRatio: 1.0 + float64(i)*0.01,
			}
		}

		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			service.calculateZScores(data)
		}
	})

	b.Run("AssessRisk", func(b *testing.B) {
		testZScores := []float64{-2.5, -1.0, 0.0, 1.5, 3.0, 8.0}

		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for _, zScore := range testZScores {
				_, _ = service.assessMVRVRisk(zScore)
			}
		}
	})
}

// BenchmarkMathOperations benchmarks mathematical operations used in indicators
func BenchmarkMathOperations(b *testing.B) {
	service := &mvrvServiceImpl{}

	b.Run("CalculateStdDev", func(b *testing.B) {
		values := make([]float64, 1000)
		for i := 0; i < 1000; i++ {
			values[i] = float64(i) * 0.01
		}
		mean := service.calculateMean(values)

		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = service.calculateStdDev(values, mean)
		}
	})

	b.Run("CalculateMean", func(b *testing.B) {
		values := make([]float64, 1000)
		for i := 0; i < 1000; i++ {
			values[i] = float64(i) * 0.01
		}

		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = service.calculateMean(values)
		}
	})
}

// BenchmarkMemoryUsage benchmarks memory usage patterns
func BenchmarkMemoryUsage(b *testing.B) {
	service := &mvrvServiceImpl{}

	b.Run("HistoricalDataGeneration", func(b *testing.B) {
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

		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			data := service.generateHistoricalMVRVData(mockBitcoinData)
			// Force garbage collection to measure actual memory usage
			_ = data[len(data)-1]
		}
	})
}

// BenchmarkConcurrentOperations benchmarks concurrent access patterns
func BenchmarkConcurrentOperations(b *testing.B) {
	service := &mvrvServiceImpl{}

	b.Run("ConcurrentRiskAssessment", func(b *testing.B) {
		testZScores := []float64{-2.5, -1.0, 0.0, 1.5, 3.0, 8.0}

		b.ReportAllocs()
		b.SetParallelism(4) // Test with 4 concurrent goroutines
		b.RunParallel(func(pb *testing.PB) {
			i := 0
			for pb.Next() {
				zScore := testZScores[i%len(testZScores)]
				_, _ = service.assessMVRVRisk(zScore)
				i++
			}
		})
	})
}