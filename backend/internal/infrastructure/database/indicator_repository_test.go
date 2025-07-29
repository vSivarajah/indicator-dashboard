package database

import (
	"context"
	"crypto-indicator-dashboard/internal/domain/entities"
	"crypto-indicator-dashboard/internal/testutil"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// IndicatorRepositoryTestSuite provides integration tests for IndicatorRepository
type IndicatorRepositoryTestSuite struct {
	suite.Suite
	testDB *testutil.TestDB
	repo   *indicatorRepository
	ctx    context.Context
}

func (suite *IndicatorRepositoryTestSuite) SetupSuite() {
	suite.testDB = testutil.NewTestDB(suite.T())
	suite.ctx = context.Background()

	// Manually create table to avoid GORM auto-migration conflicts
	err := suite.testDB.DB.Exec(`
		CREATE TABLE IF NOT EXISTS indicators (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			type TEXT NOT NULL,
			value REAL,
			string_value TEXT,
			change TEXT,
			risk_level TEXT,
			status TEXT,
			description TEXT,
			source TEXT,
			confidence REAL,
			metadata TEXT,
			timestamp DATETIME,
			created_at DATETIME,
			updated_at DATETIME
		)
	`).Error
	require.NoError(suite.T(), err, "Failed to create indicators table")

	// Initialize repository
	suite.repo = NewIndicatorRepository(suite.testDB.DB, suite.testDB.Logger).(*indicatorRepository)
}

func (suite *IndicatorRepositoryTestSuite) TearDownSuite() {
	suite.testDB.Cleanup()
}

func (suite *IndicatorRepositoryTestSuite) SetupTest() {
	// Clean up before each test by deleting all records
	suite.testDB.DB.Exec("DELETE FROM indicators")
}

func (suite *IndicatorRepositoryTestSuite) TestCreate_Success() {
	indicator := &entities.Indicator{
		Name:        "mvrv",
		Type:        "market",
		Value:       2.45,
		StringValue: "2.45",
		Change:      "+0.12",
		RiskLevel:   "medium",
		Status:      "MEDIUM: Testing resistance",
		Description: "MVRV Z-Score indicator",
		Source:      "coingecko",
		Confidence:  0.85,
		Metadata: map[string]interface{}{
			"mvrv_ratio":   1.8,
			"market_cap":   850000000000.0,
			"realized_cap": 472222222222.0,
		},
		Timestamp: time.Now(),
	}

	err := suite.repo.Create(suite.ctx, indicator)

	require.NoError(suite.T(), err, "Create should not return error")
	assert.NotZero(suite.T(), indicator.ID, "ID should be set after creation")
	assert.NotZero(suite.T(), indicator.CreatedAt, "CreatedAt should be set")
	assert.NotZero(suite.T(), indicator.UpdatedAt, "UpdatedAt should be set")
}

func (suite *IndicatorRepositoryTestSuite) TestCreate_DuplicateAllowed() {
	indicator1 := &entities.Indicator{
		Name:       "mvrv",
		Type:       "market",
		Value:      2.45,
		RiskLevel:  "medium",
		Status:     "MEDIUM",
		Confidence: 0.85,
		Timestamp:  time.Now(),
	}

	indicator2 := &entities.Indicator{
		Name:       "mvrv",
		Type:       "market",
		Value:      2.50,
		RiskLevel:  "medium",
		Status:     "MEDIUM",
		Confidence: 0.88,
		Timestamp:  time.Now().Add(time.Minute),
	}

	err1 := suite.repo.Create(suite.ctx, indicator1)
	err2 := suite.repo.Create(suite.ctx, indicator2)

	require.NoError(suite.T(), err1, "First create should succeed")
	require.NoError(suite.T(), err2, "Second create should succeed (duplicates allowed)")
	assert.NotEqual(suite.T(), indicator1.ID, indicator2.ID, "IDs should be different")
}

func (suite *IndicatorRepositoryTestSuite) TestGetByID_Success() {
	// Create test indicator
	original := &entities.Indicator{
		Name:       "fear_greed",
		Type:       "sentiment",
		Value:      72.0,
		RiskLevel:  "high",
		Status:     "GREED",
		Confidence: 0.9,
		Metadata: map[string]interface{}{
			"components": map[string]int{
				"volatility": 75,
				"momentum":   80,
			},
		},
		Timestamp: time.Now(),
	}

	err := suite.repo.Create(suite.ctx, original)
	require.NoError(suite.T(), err)

	// Retrieve by ID
	retrieved, err := suite.repo.GetByID(suite.ctx, original.ID)

	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), retrieved)
	testutil.AssertIndicatorEqual(suite.T(), original, retrieved)
	
	// Verify metadata is preserved
	assert.Equal(suite.T(), original.Metadata["components"], retrieved.Metadata["components"])
}

func (suite *IndicatorRepositoryTestSuite) TestGetByID_NotFound() {
	nonExistentID := uint(99999)
	
	result, err := suite.repo.GetByID(suite.ctx, nonExistentID)

	assert.Error(suite.T(), err, "Should return error for non-existent ID")
	assert.Nil(suite.T(), result, "Result should be nil for non-existent ID")
}

func (suite *IndicatorRepositoryTestSuite) TestGetLatest_Success() {
	now := time.Now()
	
	// Create multiple indicators with different timestamps
	indicators := []*entities.Indicator{
		{
			Name:      "dominance",
			Type:      "market",
			Value:     54.2,
			Timestamp: now.Add(-2 * time.Hour),
		},
		{
			Name:      "dominance",
			Type:      "market",
			Value:     54.5,
			Timestamp: now.Add(-1 * time.Hour), // This should be the latest
		},
		{
			Name:      "dominance",
			Type:      "market",
			Value:     54.0,
			Timestamp: now.Add(-3 * time.Hour),
		},
	}

	// Insert all indicators
	for _, indicator := range indicators {
		err := suite.repo.Create(suite.ctx, indicator)
		require.NoError(suite.T(), err)
	}

	// Get latest
	latest, err := suite.repo.GetLatest(suite.ctx, "dominance")

	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), latest)
	assert.Equal(suite.T(), 54.5, latest.Value, "Should return the most recent indicator")
	assert.Equal(suite.T(), "dominance", latest.Name)
}

func (suite *IndicatorRepositoryTestSuite) TestGetLatest_NotFound() {
	result, err := suite.repo.GetLatest(suite.ctx, "non_existent")

	assert.Error(suite.T(), err, "Should return error for non-existent indicator")
	assert.Nil(suite.T(), result, "Result should be nil")
}

func (suite *IndicatorRepositoryTestSuite) TestGetHistoricalData_Success() {
	now := time.Now()
	from := now.Add(-7 * 24 * time.Hour)
	to := now

	// Create historical data
	testData := []*entities.Indicator{
		{Name: "mvrv", Type: "market", Value: 1.5, Timestamp: now.Add(-8 * 24 * time.Hour)}, // Outside range
		{Name: "mvrv", Type: "market", Value: 2.0, Timestamp: now.Add(-6 * 24 * time.Hour)}, // In range
		{Name: "mvrv", Type: "market", Value: 2.2, Timestamp: now.Add(-4 * 24 * time.Hour)}, // In range
		{Name: "mvrv", Type: "market", Value: 2.5, Timestamp: now.Add(-2 * 24 * time.Hour)}, // In range
		{Name: "dominance", Type: "market", Value: 55.0, Timestamp: now.Add(-3 * 24 * time.Hour)}, // Different indicator
		{Name: "mvrv", Type: "market", Value: 3.0, Timestamp: now.Add(1 * time.Hour)},       // Future (outside range)
	}

	for _, indicator := range testData {
		err := suite.repo.Create(suite.ctx, indicator)
		require.NoError(suite.T(), err)
	}

	// Get historical data
	results, err := suite.repo.GetHistoricalData(suite.ctx, "mvrv", from, to)

	require.NoError(suite.T(), err)
	assert.Len(suite.T(), results, 3, "Should return 3 indicators within date range")
	
	// Verify all results are MVRV indicators within date range
	for _, result := range results {
		assert.Equal(suite.T(), "mvrv", result.Name)
		assert.True(suite.T(), result.Timestamp.After(from) || result.Timestamp.Equal(from))
		assert.True(suite.T(), result.Timestamp.Before(to) || result.Timestamp.Equal(to))
	}

	// Verify chronological order (oldest first)
	for i := 1; i < len(results); i++ {
		assert.True(suite.T(), results[i].Timestamp.After(results[i-1].Timestamp), 
			"Results should be ordered chronologically")
	}
}

func (suite *IndicatorRepositoryTestSuite) TestGetHistoricalData_EmptyResult() {
	from := time.Now().Add(-7 * 24 * time.Hour)
	to := time.Now()

	results, err := suite.repo.GetHistoricalData(suite.ctx, "non_existent", from, to)

	require.NoError(suite.T(), err)
	assert.Empty(suite.T(), results, "Should return empty slice for non-existent indicator")
}

func (suite *IndicatorRepositoryTestSuite) TestUpdate_Success() {
	// Create original indicator
	original := &entities.Indicator{
		Name:       "bubble_risk",
		Type:       "market",
		Value:      45.0,
		RiskLevel:  "medium",
		Status:     "MEDIUM",
		Confidence: 0.75,
		Timestamp:  time.Now(),
	}

	err := suite.repo.Create(suite.ctx, original)
	require.NoError(suite.T(), err)

	// Update indicator
	original.Value = 55.0
	original.RiskLevel = "high"
	original.Status = "HIGH"
	original.Confidence = 0.80

	err = suite.repo.Update(suite.ctx, original)
	require.NoError(suite.T(), err)

	// Verify update
	updated, err := suite.repo.GetByID(suite.ctx, original.ID)
	require.NoError(suite.T(), err)

	assert.Equal(suite.T(), 55.0, updated.Value)
	assert.Equal(suite.T(), "high", updated.RiskLevel) 
	assert.Equal(suite.T(), "HIGH", updated.Status)
	assert.Equal(suite.T(), 0.80, updated.Confidence)
	assert.True(suite.T(), updated.UpdatedAt.After(updated.CreatedAt), "UpdatedAt should be newer than CreatedAt")
}

func (suite *IndicatorRepositoryTestSuite) TestUpdate_NotFound() {
	nonExistent := &entities.Indicator{
		ID:         99999,
		Name:       "test",
		Type:       "market",
		Value:      1.0,
		Timestamp:  time.Now(),
	}

	err := suite.repo.Update(suite.ctx, nonExistent)
	assert.Error(suite.T(), err, "Should return error when updating non-existent indicator")
}

func (suite *IndicatorRepositoryTestSuite) TestDelete_Success() {
	// Create indicator to delete
	indicator := &entities.Indicator{
		Name:      "test_delete",
		Type:      "market",
		Value:     1.0,
		Timestamp: time.Now(),
	}

	err := suite.repo.Create(suite.ctx, indicator)
	require.NoError(suite.T(), err)

	// Delete indicator
	err = suite.repo.Delete(suite.ctx, indicator.ID)
	require.NoError(suite.T(), err)

	// Verify deletion
	deleted, err := suite.repo.GetByID(suite.ctx, indicator.ID)
	assert.Error(suite.T(), err, "Should return error for deleted indicator")
	assert.Nil(suite.T(), deleted, "Deleted indicator should not be found")
}

func (suite *IndicatorRepositoryTestSuite) TestDelete_NotFound() {
	err := suite.repo.Delete(suite.ctx, 99999)
	assert.Error(suite.T(), err, "Should return error when deleting non-existent indicator")
}

func (suite *IndicatorRepositoryTestSuite) TestConcurrentAccess() {
	// Test concurrent creates
	const numGoroutines = 10
	const indicatorsPerGoroutine = 5

	results := make(chan error, numGoroutines)
	
	for i := 0; i < numGoroutines; i++ {
		go func(goroutineID int) {
			for j := 0; j < indicatorsPerGoroutine; j++ {
				indicator := &entities.Indicator{
					Name:      "concurrent_test",
					Type:      "market",
					Value:     float64(goroutineID*100 + j),
					Timestamp: time.Now(),
				}
				
				err := suite.repo.Create(suite.ctx, indicator)
				if err != nil {
					results <- err
					return
				}
			}
			results <- nil
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < numGoroutines; i++ {
		err := <-results
		assert.NoError(suite.T(), err, "Concurrent create should not fail")
	}

	// Verify all indicators were created
	historical, err := suite.repo.GetHistoricalData(suite.ctx, "concurrent_test", 
		time.Now().Add(-time.Hour), time.Now().Add(time.Hour))
	require.NoError(suite.T(), err)
	assert.Len(suite.T(), historical, numGoroutines*indicatorsPerGoroutine, 
		"All concurrent indicators should be created")
}

func (suite *IndicatorRepositoryTestSuite) TestLargeMetadata() {
	// Test with large metadata object
	largeMetadata := make(map[string]interface{})
	
	// Create nested structure with many fields
	for i := 0; i < 100; i++ {
		largeMetadata[fmt.Sprintf("field_%d", i)] = map[string]interface{}{
			"value":       float64(i),
			"description": fmt.Sprintf("This is field number %d with some description", i),
			"nested": map[string]string{
				"key1": "value1",
				"key2": "value2",
				"key3": "value3",
			},
		}
	}

	indicator := &entities.Indicator{
		Name:      "large_metadata_test",
		Type:      "market",
		Value:     1.0,
		Metadata:  largeMetadata,
		Timestamp: time.Now(),
	}

	// Should handle large metadata without issues
	err := suite.repo.Create(suite.ctx, indicator)
	require.NoError(suite.T(), err, "Should handle large metadata")

	// Retrieve and verify metadata is intact
	retrieved, err := suite.repo.GetByID(suite.ctx, indicator.ID)
	require.NoError(suite.T(), err)

	assert.Len(suite.T(), retrieved.Metadata, 100, "All metadata fields should be preserved")
	
	// Spot check some values
	field0, exists := retrieved.Metadata["field_0"]
	assert.True(suite.T(), exists, "field_0 should exist")
	
	field0Map, ok := field0.(map[string]interface{})
	assert.True(suite.T(), ok, "field_0 should be a map")
	assert.Equal(suite.T(), float64(0), field0Map["value"], "Nested value should be preserved")
}

// Run the test suite
func TestIndicatorRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(IndicatorRepositoryTestSuite))
}

// Table-driven tests for various scenarios
func TestIndicatorRepository_EdgeCases(t *testing.T) {
	testDB := testutil.NewTestDB(t)
	defer testDB.Cleanup()

	// Manually create table to avoid GORM auto-migration conflicts
	err := testDB.DB.Exec(`
		CREATE TABLE IF NOT EXISTS indicators (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			type TEXT NOT NULL,
			value REAL,
			string_value TEXT,
			change TEXT,
			risk_level TEXT,
			status TEXT,
			description TEXT,
			source TEXT,
			confidence REAL,
			metadata TEXT,
			timestamp DATETIME,
			created_at DATETIME,
			updated_at DATETIME
		)
	`).Error
	require.NoError(t, err)

	repo := NewIndicatorRepository(testDB.DB, testDB.Logger).(*indicatorRepository)

	ctx := context.Background()

	tests := []struct {
		name          string
		indicator     *entities.Indicator
		expectError   bool
		errorContains string
	}{
		{
			name: "Empty name",
			indicator: &entities.Indicator{
				Name:      "",
				Type:      "market",
				Value:     1.0,
				Timestamp: time.Now(),
			},
			expectError:   true,
			errorContains: "name",
		},
		{
			name: "Nil metadata",
			indicator: &entities.Indicator{
				Name:      "test_nil_metadata",
				Type:      "market",
				Value:     1.0,
				Metadata:  nil,
				Timestamp: time.Now(),
			},
			expectError: false,
		},
		{
			name: "Zero timestamp",
			indicator: &entities.Indicator{
				Name:      "test_zero_time",
				Type:      "market",
				Value:     1.0,
				Timestamp: time.Time{},
			},
			expectError: false, // Should be allowed
		},
		{
			name: "Very long strings",
			indicator: &entities.Indicator{
				Name:        string(make([]byte, 1000)), // Very long name
				Type:        "market",
				StringValue: string(make([]byte, 2000)), // Very long string value
				Value:       1.0,
				Timestamp:   time.Now(),
			},
			expectError: false, // Database should handle or truncate
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Create(ctx, tt.indicator)
			
			if tt.expectError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}