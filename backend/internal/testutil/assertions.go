package testutil

import (
	"crypto-indicator-dashboard/internal/domain/entities"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// AssertIndicatorEqual asserts that two indicators are equal
func AssertIndicatorEqual(t *testing.T, expected, actual *entities.Indicator) {
	t.Helper()
	
	require.NotNil(t, actual, "Actual indicator should not be nil")
	assert.Equal(t, expected.Name, actual.Name, "Indicator name should match")
	assert.Equal(t, expected.Type, actual.Type, "Indicator type should match")
	assert.InDelta(t, expected.Value, actual.Value, 0.001, "Indicator value should match within delta")
	assert.Equal(t, expected.RiskLevel, actual.RiskLevel, "Risk level should match")
	assert.Equal(t, expected.Status, actual.Status, "Status should match")
	assert.InDelta(t, expected.Confidence, actual.Confidence, 0.001, "Confidence should match within delta")
}

// AssertMarketDataEqual asserts that two market data objects are equal
func AssertMarketDataEqual(t *testing.T, expected, actual *entities.MarketData) {
	t.Helper()
	
	require.NotNil(t, actual, "Actual market data should not be nil")
	assert.Equal(t, expected.Symbol, actual.Symbol, "Symbol should match")
	assert.Equal(t, expected.Name, actual.Name, "Name should match")
	assert.InDelta(t, expected.Price, actual.Price, 0.01, "Price should match within delta")
	assert.InDelta(t, expected.MarketCap, actual.MarketCap, 1000000, "Market cap should match within delta")
	assert.InDelta(t, expected.Volume24h, actual.Volume24h, 1000000, "Volume should match within delta")
	assert.InDelta(t, expected.Dominance, actual.Dominance, 0.01, "Dominance should match within delta")
}

// AssertTimeWithinRange checks if a time is within a range
func AssertTimeWithinRange(t *testing.T, actual time.Time, start, end time.Time, msg string) {
	t.Helper()
	
	assert.True(t, actual.After(start) || actual.Equal(start), 
		"%s should be after or equal to start time. actual: %v, start: %v", msg, actual, start)
	assert.True(t, actual.Before(end) || actual.Equal(end), 
		"%s should be before or equal to end time. actual: %v, end: %v", msg, actual, end)
}

// AssertRecentTime checks if a time is recent (within last minute)
func AssertRecentTime(t *testing.T, actual time.Time, msg string) {
	t.Helper()
	
	now := time.Now()
	oneMinuteAgo := now.Add(-time.Minute)
	AssertTimeWithinRange(t, actual, oneMinuteAgo, now, msg)
}

// AssertNoError is a convenience wrapper for require.NoError with context
func AssertNoError(t *testing.T, err error, msgAndArgs ...interface{}) {
	t.Helper()
	require.NoError(t, err, msgAndArgs...)
}

// AssertError is a convenience wrapper for require.Error with context
func AssertError(t *testing.T, err error, msgAndArgs ...interface{}) {
	t.Helper()
	require.Error(t, err, msgAndArgs...)
}

// AssertFloat64InRange checks if a float64 value is within a range
func AssertFloat64InRange(t *testing.T, actual, min, max float64, msg string) {
	t.Helper()
	
	assert.True(t, actual >= min, "%s should be >= %f, got %f", msg, min, actual)
	assert.True(t, actual <= max, "%s should be <= %f, got %f", msg, max, actual)
}

// AssertValidRiskLevel checks if risk level is valid
func AssertValidRiskLevel(t *testing.T, riskLevel string) {
	t.Helper()
	
	validLevels := []string{"extreme_low", "low", "medium", "high", "extreme_high"}
	assert.Contains(t, validLevels, riskLevel, "Risk level should be valid")
}

// AssertValidIndicatorType checks if indicator type is valid
func AssertValidIndicatorType(t *testing.T, indicatorType string) {
	t.Helper()
	
	validTypes := []string{"market", "on-chain", "macro", "sentiment"}
	assert.Contains(t, validTypes, indicatorType, "Indicator type should be valid")
}

// AssertPositiveFloat64 checks if a float64 value is positive
func AssertPositiveFloat64(t *testing.T, value float64, msg string) {
	t.Helper()
	assert.True(t, value > 0, "%s should be positive, got %f", msg, value)
}

// AssertConfidenceScore checks if confidence score is valid (0.0 - 1.0)
func AssertConfidenceScore(t *testing.T, confidence float64) {
	t.Helper()
	AssertFloat64InRange(t, confidence, 0.0, 1.0, "Confidence score")
}

// AssertNonEmptyString checks if string is not empty
func AssertNonEmptyString(t *testing.T, value, fieldName string) {
	t.Helper()
	assert.NotEmpty(t, value, "%s should not be empty", fieldName)
}

// AssertMapContainsKeys checks if map contains all required keys
func AssertMapContainsKeys(t *testing.T, m map[string]interface{}, keys []string, msg string) {
	t.Helper()
	
	for _, key := range keys {
		assert.Contains(t, m, key, "%s should contain key '%s'", msg, key)
	}
}

// AssertSliceNotEmpty checks if slice is not empty
func AssertSliceNotEmpty(t *testing.T, slice []interface{}, msg string) {
	t.Helper()
	assert.NotEmpty(t, slice, "%s should not be empty", msg)
}

// BenchmarkHelper provides utilities for benchmark tests
type BenchmarkHelper struct{}

// NewBenchmarkHelper creates a new benchmark helper
func NewBenchmarkHelper() *BenchmarkHelper {
	return &BenchmarkHelper{}
}

// ResetTimer resets the benchmark timer (wrapper for b.ResetTimer)
func (bh *BenchmarkHelper) ResetTimer(b *testing.B) {
	b.Helper()
	b.ResetTimer()
}

// StopTimer stops the benchmark timer (wrapper for b.StopTimer)
func (bh *BenchmarkHelper) StopTimer(b *testing.B) {
	b.Helper()
	b.StopTimer()
}

// StartTimer starts the benchmark timer (wrapper for b.StartTimer)
func (bh *BenchmarkHelper) StartTimer(b *testing.B) {
	b.Helper()
	b.StartTimer()
}