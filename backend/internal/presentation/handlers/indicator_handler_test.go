package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"crypto-indicator-dashboard/internal/infrastructure/config"
	"crypto-indicator-dashboard/internal/testutil"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// IndicatorHandlerTestSuite provides tests for IndicatorHandler
type IndicatorHandlerTestSuite struct {
	suite.Suite
	handler *IndicatorHandler
	router  *gin.Engine
	testDB  *testutil.TestDB
}

func (suite *IndicatorHandlerTestSuite) SetupTest() {
	// Create test database
	suite.testDB = testutil.NewTestDB(suite.T())

	// Create mock dependencies
	deps := &config.Dependencies{
		Logger: suite.testDB.Logger,
		Cache:  testutil.NewMockCacheService(),
	}

	// Create handler
	suite.handler = NewIndicatorHandler(deps)

	// Set up Gin router
	gin.SetMode(gin.TestMode)
	suite.router = gin.New()

	// Register routes
	apiV1 := suite.router.Group("/api/v1")
	suite.handler.RegisterRoutes(apiV1)
}

func (suite *IndicatorHandlerTestSuite) TearDownTest() {
	suite.testDB.Cleanup()
}

func (suite *IndicatorHandlerTestSuite) TestGetMVRVIndicator_Success() {
	// Create request
	req, err := http.NewRequest("GET", "/api/v1/indicators/mvrv", nil)
	require.NoError(suite.T(), err)

	// Create response recorder
	w := httptest.NewRecorder()

	// Perform request
	suite.router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(suite.T(), err)

	assert.True(suite.T(), response["success"].(bool))
	assert.Contains(suite.T(), response, "data")

	data := response["data"].(map[string]interface{})
	assert.Contains(suite.T(), data, "value")
	assert.Contains(suite.T(), data, "risk_level")
	assert.Contains(suite.T(), data, "status")
	assert.Contains(suite.T(), data, "last_updated")
}

func (suite *IndicatorHandlerTestSuite) TestGetDominanceIndicator_Success() {
	req, err := http.NewRequest("GET", "/api/v1/indicators/dominance", nil)
	require.NoError(suite.T(), err)

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(suite.T(), err)

	assert.True(suite.T(), response["success"].(bool))
	data := response["data"].(map[string]interface{})
	assert.Contains(suite.T(), data, "value")
	assert.Contains(suite.T(), data, "change")
}

func (suite *IndicatorHandlerTestSuite) TestGetFearGreedIndicator_Success() {
	req, err := http.NewRequest("GET", "/api/v1/indicators/fear-greed", nil)
	require.NoError(suite.T(), err)

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(suite.T(), err)

	assert.True(suite.T(), response["success"].(bool))
	data := response["data"].(map[string]interface{})
	assert.Contains(suite.T(), data, "value")
	assert.Contains(suite.T(), data, "status")
	assert.Contains(suite.T(), data, "risk_level")
}

func (suite *IndicatorHandlerTestSuite) TestGetBubbleRiskIndicator_Success() {
	req, err := http.NewRequest("GET", "/api/v1/indicators/bubble-risk", nil)
	require.NoError(suite.T(), err)

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(suite.T(), err)

	assert.True(suite.T(), response["success"].(bool))
	data := response["data"].(map[string]interface{})
	assert.Contains(suite.T(), data, "value")
	assert.Equal(suite.T(), "medium", data["risk_level"])
}

func (suite *IndicatorHandlerTestSuite) TestGetChartData_MVRV() {
	req, err := http.NewRequest("GET", "/api/v1/charts/mvrv", nil)
	require.NoError(suite.T(), err)

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(suite.T(), err)

	assert.Contains(suite.T(), response, "timestamps")
	assert.Contains(suite.T(), response, "zscore_data")
	assert.Contains(suite.T(), response, "price_data")
	assert.Contains(suite.T(), response, "current_zscore")
	assert.Contains(suite.T(), response, "thresholds")
}

func (suite *IndicatorHandlerTestSuite) TestGetChartData_Dominance() {
	req, err := http.NewRequest("GET", "/api/v1/charts/dominance", nil)
	require.NoError(suite.T(), err)

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(suite.T(), err)

	assert.Contains(suite.T(), response, "timestamps")
	assert.Contains(suite.T(), response, "values")
	assert.Contains(suite.T(), response, "current")
	assert.Contains(suite.T(), response, "levels")
}

func (suite *IndicatorHandlerTestSuite) TestGetChartData_UnknownIndicator() {
	req, err := http.NewRequest("GET", "/api/v1/charts/unknown", nil)
	require.NoError(suite.T(), err)

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(suite.T(), err)

	assert.Equal(suite.T(), "unknown", response["indicator"])
	assert.Contains(suite.T(), response, "message")
	assert.Contains(suite.T(), response, "mock_data")
}

// Test suite runner
func TestIndicatorHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(IndicatorHandlerTestSuite))
}

// Table-driven tests for response validation
func TestIndicatorHandler_ResponseFormats(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	testDB := testutil.NewTestDB(t)
	defer testDB.Cleanup()
	
	deps := &config.Dependencies{
		Logger: testDB.Logger,
		Cache:  testutil.NewMockCacheService(),
	}
	
	handler := NewIndicatorHandler(deps)
	apiV1 := router.Group("/api/v1")
	handler.RegisterRoutes(apiV1)

	tests := []struct {
		name           string
		endpoint       string
		expectedFields []string
	}{
		{
			name:     "MVRV endpoint",
			endpoint: "/api/v1/indicators/mvrv",
			expectedFields: []string{"success", "data"},
		},
		{
			name:     "Dominance endpoint",
			endpoint: "/api/v1/indicators/dominance",
			expectedFields: []string{"success", "data"},
		},
		{
			name:     "Fear & Greed endpoint",
			endpoint: "/api/v1/indicators/fear-greed",
			expectedFields: []string{"success", "data"},
		},
		{
			name:     "Bubble Risk endpoint",
			endpoint: "/api/v1/indicators/bubble-risk",
			expectedFields: []string{"success", "data"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", tt.endpoint, nil)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)

			var response map[string]interface{}
			err = json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			for _, field := range tt.expectedFields {
				assert.Contains(t, response, field, "Response should contain field: %s", field)
			}
		})
	}
}

// Benchmark tests for handler performance
func BenchmarkIndicatorHandler(b *testing.B) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	testDB := testutil.NewTestDB(&testing.T{})
	defer testDB.Cleanup()
	
	deps := &config.Dependencies{
		Logger: testDB.Logger,
		Cache:  testutil.NewMockCacheService(),
	}
	
	handler := NewIndicatorHandler(deps)
	apiV1 := router.Group("/api/v1")
	handler.RegisterRoutes(apiV1)

	b.Run("GetMVRVIndicator", func(b *testing.B) {
		req, _ := http.NewRequest("GET", "/api/v1/indicators/mvrv", nil)
		
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
		}
	})

	b.Run("GetChartData", func(b *testing.B) {
		req, _ := http.NewRequest("GET", "/api/v1/charts/mvrv", nil)
		
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
		}
	})
}