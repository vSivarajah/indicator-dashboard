package external

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
	"crypto-indicator-dashboard/pkg/logger"
)

// CoinCapClient handles CoinCap API interactions
type CoinCapClient struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
	logger     logger.Logger
}

// NewCoinCapClient creates a new CoinCap API client
func NewCoinCapClient(apiKey string, logger logger.Logger) *CoinCapClient {
	return &CoinCapClient{
		apiKey:  apiKey,
		baseURL: "https://rest.coincap.io/v3",
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		logger: logger,
	}
}

// Asset represents a cryptocurrency asset from CoinCap
type Asset struct {
	ID                string  `json:"id"`
	Rank              string  `json:"rank"`
	Symbol            string  `json:"symbol"`
	Name              string  `json:"name"`
	Supply            string  `json:"supply"`
	MaxSupply         *string `json:"maxSupply"`
	MarketCapUSD      string  `json:"marketCapUsd"`
	VolumeUSD24Hr     string  `json:"volumeUsd24Hr"`
	PriceUSD          string  `json:"priceUsd"`
	ChangePercent24Hr string  `json:"changePercent24Hr"`
	VWAP24Hr          *string `json:"vwap24Hr"`
}

// AssetsResponse represents the response from assets endpoint
type AssetsResponse struct {
	Data      []Asset `json:"data"`
	Timestamp int64   `json:"timestamp"`
}

// AssetResponse represents the response from single asset endpoint
type AssetResponse struct {
	Data      Asset `json:"data"`
	Timestamp int64 `json:"timestamp"`
}

// HistoryData represents historical price data point
type HistoryData struct {
	PriceUSD string `json:"priceUsd"`
	Time     int64  `json:"time"`
	Date     string `json:"date"`
}

// HistoryResponse represents the response from history endpoint
type HistoryResponse struct {
	Data      []HistoryData `json:"data"`
	Timestamp int64         `json:"timestamp"`
}

// Market represents exchange market data
type Market struct {
	ExchangeID      string `json:"exchangeId"`
	Rank            string `json:"rank"`
	BaseSymbol      string `json:"baseSymbol"`
	BaseID          string `json:"baseId"`
	QuoteSymbol     string `json:"quoteSymbol"`
	QuoteID         string `json:"quoteId"`
	PriceQuote      string `json:"priceQuote"`
	PriceUSD        string `json:"priceUsd"`
	VolumeUSD24Hr   string `json:"volumeUsd24Hr"`
	PercentExchange string `json:"percentExchangeVolume"`
	TradesCount24Hr string `json:"tradesCount24Hr"`
	Updated         int64  `json:"updated"`
}

// MarketsResponse represents the response from markets endpoint
type MarketsResponse struct {
	Data      []Market `json:"data"`
	Timestamp int64    `json:"timestamp"`
}

// GetAssets retrieves list of all assets
func (c *CoinCapClient) GetAssets(limit int) (*AssetsResponse, error) {
	endpoint := "/assets"
	if limit > 0 {
		endpoint += fmt.Sprintf("?limit=%d", limit)
	}
	
	data, err := c.makeRequest(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch assets: %w", err)
	}

	var response AssetsResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal assets response: %w", err)
	}

	c.logger.Info("Successfully fetched assets", "count", len(response.Data))
	return &response, nil
}

// GetAsset retrieves a specific asset by ID
func (c *CoinCapClient) GetAsset(assetID string) (*AssetResponse, error) {
	endpoint := fmt.Sprintf("/assets/%s", assetID)
	
	data, err := c.makeRequest(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch asset %s: %w", assetID, err)
	}

	var response AssetResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal asset response: %w", err)
	}

	c.logger.Info("Successfully fetched asset", "asset_id", assetID, "price", response.Data.PriceUSD)
	return &response, nil
}

// GetAssetHistory retrieves historical price data for an asset
func (c *CoinCapClient) GetAssetHistory(assetID, interval string, start, end *time.Time) (*HistoryResponse, error) {
	endpoint := fmt.Sprintf("/assets/%s/history", assetID)
	
	// Add query parameters
	params := []string{}
	if interval != "" {
		params = append(params, fmt.Sprintf("interval=%s", interval))
	}
	if start != nil {
		params = append(params, fmt.Sprintf("start=%d", start.UnixMilli()))
	}
	if end != nil {
		params = append(params, fmt.Sprintf("end=%d", end.UnixMilli()))
	}
	
	if len(params) > 0 {
		endpoint += "?" + strings.Join(params, "&")
	}
	
	data, err := c.makeRequest(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch asset history for %s: %w", assetID, err)
	}

	var response HistoryResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal history response: %w", err)
	}

	c.logger.Info("Successfully fetched asset history", 
		"asset_id", assetID, 
		"interval", interval,
		"data_points", len(response.Data))
	
	return &response, nil
}

// GetMarkets retrieves market data for an asset
func (c *CoinCapClient) GetMarkets(assetID string, limit int) (*MarketsResponse, error) {
	endpoint := "/markets"
	params := []string{}
	
	if assetID != "" {
		params = append(params, fmt.Sprintf("baseId=%s", assetID))
	}
	if limit > 0 {
		params = append(params, fmt.Sprintf("limit=%d", limit))
	}
	
	if len(params) > 0 {
		endpoint += "?" + strings.Join(params, "&")
	}
	
	data, err := c.makeRequest(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch markets: %w", err)
	}

	var response MarketsResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal markets response: %w", err)
	}

	c.logger.Info("Successfully fetched markets", "count", len(response.Data))
	return &response, nil
}

// GetBitcoinPrice retrieves current Bitcoin price
func (c *CoinCapClient) GetBitcoinPrice() (float64, error) {
	response, err := c.GetAsset("bitcoin")
	if err != nil {
		return 0, fmt.Errorf("failed to get Bitcoin price: %w", err)
	}

	var price float64
	if _, err := fmt.Sscanf(response.Data.PriceUSD, "%f", &price); err != nil {
		return 0, fmt.Errorf("failed to parse Bitcoin price: %w", err)
	}

	return price, nil
}

// GetTop10Assets retrieves top 10 assets by market cap
func (c *CoinCapClient) GetTop10Assets() (*AssetsResponse, error) {
	return c.GetAssets(10)
}

// GetBitcoinHistoricalData retrieves Bitcoin historical data for a specific period
func (c *CoinCapClient) GetBitcoinHistoricalData(interval string, days int) (*HistoryResponse, error) {
	end := time.Now()
	start := end.AddDate(0, 0, -days)
	
	return c.GetAssetHistory("bitcoin", interval, &start, &end)
}

// makeRequest makes an HTTP request to the CoinCap API
func (c *CoinCapClient) makeRequest(endpoint string) ([]byte, error) {
	reqURL := c.baseURL + endpoint

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("User-Agent", "CryptoIndicatorDashboard/1.0")
	
	// Add API key if provided
	if c.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}

	c.logger.Debug("Making CoinCap API request", 
		"url", reqURL,
		"endpoint", endpoint)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Handle gzip compression
	var reader io.Reader = resp.Body
	if strings.Contains(resp.Header.Get("Content-Encoding"), "gzip") {
		gzipReader, err := gzip.NewReader(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to create gzip reader: %w", err)
		}
		defer gzipReader.Close()
		reader = gzipReader
	}

	body, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		c.logger.Error("CoinCap API request failed", 
			"status_code", resp.StatusCode,
			"response", string(body))
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	return body, nil
}

// HealthCheck performs a health check on the CoinCap service
func (c *CoinCapClient) HealthCheck() error {
	// Try to fetch Bitcoin price as a simple health check
	_, err := c.GetBitcoinPrice()
	if err != nil {
		return fmt.Errorf("CoinCap health check failed: %w", err)
	}
	return nil
}

// GetGlobalMarketData provides global market statistics
func (c *CoinCapClient) GetGlobalMarketData() (map[string]interface{}, error) {
	// Get top 10 assets to calculate global stats
	response, err := c.GetTop10Assets()
	if err != nil {
		return nil, fmt.Errorf("failed to get global market data: %w", err)
	}

	var totalMarketCap, totalVolume float64
	var btcDominance float64

	for _, asset := range response.Data {
		if marketCap := parseFloat(asset.MarketCapUSD); marketCap > 0 {
			totalMarketCap += marketCap
			if asset.Symbol == "BTC" {
				btcDominance = marketCap
			}
		}
		if volume := parseFloat(asset.VolumeUSD24Hr); volume > 0 {
			totalVolume += volume
		}
	}

	// Calculate BTC dominance percentage
	btcDominancePercent := 0.0
	if totalMarketCap > 0 {
		btcDominancePercent = (btcDominance / totalMarketCap) * 100
	}

	return map[string]interface{}{
		"total_market_cap":    totalMarketCap,
		"total_volume_24h":    totalVolume,
		"btc_dominance":       btcDominancePercent,
		"active_cryptocurrencies": len(response.Data),
		"timestamp":           time.Now().Unix(),
	}, nil
}

// parseFloat safely parses a string to float64
func parseFloat(s string) float64 {
	var f float64
	fmt.Sscanf(s, "%f", &f)
	return f
}