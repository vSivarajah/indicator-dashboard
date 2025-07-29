package external

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
	"crypto-indicator-dashboard/pkg/logger"
)

// CoinMarketCapClient handles CoinMarketCap API interactions
type CoinMarketCapClient struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
	logger     logger.Logger
}

// NewCoinMarketCapClient creates a new CoinMarketCap API client
func NewCoinMarketCapClient(apiKey string, logger logger.Logger) *CoinMarketCapClient {
	return &CoinMarketCapClient{
		apiKey:  apiKey,
		baseURL: "https://pro-api.coinmarketcap.com/v1",
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		logger: logger,
	}
}

// CryptoCurrency represents a cryptocurrency from CoinMarketCap
type CryptoCurrency struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Symbol string `json:"symbol"`
	Slug   string `json:"slug"`
}

// Quote represents price quote data
type Quote struct {
	Price            float64   `json:"price"`
	Volume24h        float64   `json:"volume_24h"`
	VolumeChange24h  float64   `json:"volume_change_24h"`
	PercentChange1h  float64   `json:"percent_change_1h"`
	PercentChange24h float64   `json:"percent_change_24h"`
	PercentChange7d  float64   `json:"percent_change_7d"`
	PercentChange30d float64   `json:"percent_change_30d"`
	MarketCap        float64   `json:"market_cap"`
	MarketCapDominance float64 `json:"market_cap_dominance"`
	FullyDilutedMarketCap float64 `json:"fully_diluted_market_cap"`
	LastUpdated      time.Time `json:"last_updated"`
}

// CryptoPriceData represents complete price data for a cryptocurrency
type CryptoPriceData struct {
	ID                int                    `json:"id"`
	Name              string                 `json:"name"`
	Symbol            string                 `json:"symbol"`
	Slug              string                 `json:"slug"`
	NumMarketPairs    int                    `json:"num_market_pairs"`
	DateAdded         time.Time              `json:"date_added"`
	Tags              []string               `json:"tags"`
	MaxSupply         *float64               `json:"max_supply"`
	CirculatingSupply float64                `json:"circulating_supply"`
	TotalSupply       float64                `json:"total_supply"`
	Quote             map[string]Quote       `json:"quote"`
	LastUpdated       time.Time              `json:"last_updated"`
}

// LatestQuotesResponse represents the response from latest quotes endpoint
type LatestQuotesResponse struct {
	Status struct {
		Timestamp    time.Time `json:"timestamp"`
		ErrorCode    int       `json:"error_code"`
		ErrorMessage *string   `json:"error_message"`
		Elapsed      int       `json:"elapsed"`
		CreditCount  int       `json:"credit_count"`
		Notice       *string   `json:"notice"`
	} `json:"status"`
	Data map[string]CryptoPriceData `json:"data"`
}

// GlobalMetricsData represents global cryptocurrency market data
type GlobalMetricsData struct {
	ActiveCryptocurrencies int `json:"active_cryptocurrencies"`
	TotalCryptocurrencies int `json:"total_cryptocurrencies"`
	ActiveMarketPairs     int `json:"active_market_pairs"`
	ActiveExchanges       int `json:"active_exchanges"`
	TotalExchanges        int `json:"total_exchanges"`
	EthDominance          float64 `json:"eth_dominance"`
	BtcDominance          float64 `json:"btc_dominance"`
	EthDominanceYesterday float64 `json:"eth_dominance_yesterday"`
	BtcDominanceYesterday float64 `json:"btc_dominance_yesterday"`
	EthDominance24hPercentageChange float64 `json:"eth_dominance_24h_percentage_change"`
	BtcDominance24hPercentageChange float64 `json:"btc_dominance_24h_percentage_change"`
	DefiVolumeYesterday   float64 `json:"defi_volume_yesterday"`
	DefiVolume24h         float64 `json:"defi_volume_24h"`
	DefiVolume24hReported float64 `json:"defi_volume_24h_reported"`
	DefiMarketCap         float64 `json:"defi_market_cap"`
	DefiVolume24hPercentageChange float64 `json:"defi_volume_24h_percentage_change"`
	StablecoinVolume24h   float64 `json:"stablecoin_volume_24h"`
	StablecoinVolume24hReported float64 `json:"stablecoin_volume_24h_reported"`
	StablecoinVolume24hPercentageChange float64 `json:"stablecoin_volume_24h_percentage_change"`
	StablecoinMarketCap   float64 `json:"stablecoin_market_cap"`
	DerivativesVolume24h  float64 `json:"derivatives_volume_24h"`
	DerivativesVolume24hReported float64 `json:"derivatives_volume_24h_reported"`
	DerivativesVolume24hPercentageChange float64 `json:"derivatives_volume_24h_percentage_change"`
	Quote                 map[string]Quote `json:"quote"`
	LastUpdated           time.Time `json:"last_updated"`
}

// GlobalMetricsResponse represents the response from global metrics endpoint
type GlobalMetricsResponse struct {
	Status struct {
		Timestamp    time.Time `json:"timestamp"`
		ErrorCode    int       `json:"error_code"`
		ErrorMessage *string   `json:"error_message"`
		Elapsed      int       `json:"elapsed"`
		CreditCount  int       `json:"credit_count"`
		Notice       *string   `json:"notice"`
	} `json:"status"`
	Data GlobalMetricsData `json:"data"`
}

// GetLatestQuotes retrieves latest price quotes for specified cryptocurrencies
func (c *CoinMarketCapClient) GetLatestQuotes(symbols []string, convert string) (*LatestQuotesResponse, error) {
	if convert == "" {
		convert = "USD"
	}

	params := url.Values{}
	// Join symbols with comma for CoinMarketCap API
	symbolsStr := ""
	for i, symbol := range symbols {
		if i > 0 {
			symbolsStr += ","
		}
		symbolsStr += symbol
	}
	params.Set("symbol", symbolsStr)
	params.Set("convert", convert)

	endpoint := "/cryptocurrency/quotes/latest"
	data, err := c.makeRequest(endpoint, params)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch latest quotes: %w", err)
	}

	var response LatestQuotesResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal latest quotes response: %w", err)
	}

	if response.Status.ErrorCode != 0 {
		errorMsg := "unknown error"
		if response.Status.ErrorMessage != nil {
			errorMsg = *response.Status.ErrorMessage
		}
		return nil, fmt.Errorf("CoinMarketCap API error: %s (code: %d)", errorMsg, response.Status.ErrorCode)
	}

	c.logger.Info("Successfully fetched latest quotes", 
		"symbols", symbols, 
		"convert", convert,
		"credit_count", response.Status.CreditCount)

	return &response, nil
}

// GetGlobalMetrics retrieves global cryptocurrency market metrics
func (c *CoinMarketCapClient) GetGlobalMetrics(convert string) (*GlobalMetricsResponse, error) {
	if convert == "" {
		convert = "USD"
	}

	params := url.Values{}
	params.Set("convert", convert)

	endpoint := "/global-metrics/quotes/latest"
	data, err := c.makeRequest(endpoint, params)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch global metrics: %w", err)
	}

	var response GlobalMetricsResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal global metrics response: %w", err)
	}

	if response.Status.ErrorCode != 0 {
		errorMsg := "unknown error"
		if response.Status.ErrorMessage != nil {
			errorMsg = *response.Status.ErrorMessage
		}
		return nil, fmt.Errorf("CoinMarketCap API error: %s (code: %d)", errorMsg, response.Status.ErrorCode)
	}

	c.logger.Info("Successfully fetched global metrics", 
		"convert", convert,
		"btc_dominance", response.Data.BtcDominance,
		"credit_count", response.Status.CreditCount)

	return &response, nil
}

// GetPriceBySymbol is a convenience method to get price for a single symbol
func (c *CoinMarketCapClient) GetPriceBySymbol(symbol, convert string) (float64, error) {
	response, err := c.GetLatestQuotes([]string{symbol}, convert)
	if err != nil {
		return 0, err
	}

	if data, exists := response.Data[symbol]; exists {
		if quote, exists := data.Quote[convert]; exists {
			return quote.Price, nil
		}
		return 0, fmt.Errorf("convert currency %s not found in response", convert)
	}
	
	return 0, fmt.Errorf("symbol %s not found in response", symbol)
}

// GetBitcoinDominance retrieves Bitcoin dominance from global metrics
func (c *CoinMarketCapClient) GetBitcoinDominance() (float64, error) {
	response, err := c.GetGlobalMetrics("USD")
	if err != nil {
		return 0, fmt.Errorf("failed to get Bitcoin dominance: %w", err)
	}

	return response.Data.BtcDominance, nil
}

// makeRequest makes an HTTP request to the CoinMarketCap API
func (c *CoinMarketCapClient) makeRequest(endpoint string, params url.Values) ([]byte, error) {
	reqURL := c.baseURL + endpoint
	if len(params) > 0 {
		reqURL += "?" + params.Encode()
	}

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add required headers
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Encoding", "deflate, gzip")
	req.Header.Set("X-CMC_PRO_API_KEY", c.apiKey)

	c.logger.Debug("Making CoinMarketCap API request", 
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
		c.logger.Error("CoinMarketCap API request failed", 
			"status_code", resp.StatusCode,
			"response", string(body))
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	return body, nil
}

// Health check for the CoinMarketCap service
func (c *CoinMarketCapClient) HealthCheck() error {
	// Try to fetch Bitcoin price as a simple health check
	_, err := c.GetPriceBySymbol("BTC", "USD")
	if err != nil {
		return fmt.Errorf("CoinMarketCap health check failed: %w", err)
	}
	return nil
}