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

// BlockchainClient handles Blockchain.com API interactions
type BlockchainClient struct {
	baseURL    string
	httpClient *http.Client
	logger     logger.Logger
}

// NewBlockchainClient creates a new Blockchain.com API client
func NewBlockchainClient(logger logger.Logger) *BlockchainClient {
	return &BlockchainClient{
		baseURL: "https://blockchain.info",
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		logger: logger,
	}
}

// BitcoinStats represents Bitcoin network statistics
type BitcoinStats struct {
	MarketPriceUSD         float64 `json:"market_price_usd"`
	HashRate               float64 `json:"hash_rate"`
	TotalFeesBTC           float64 `json:"total_fees_btc"`
	NTransactions          int64   `json:"n_transactions"`
	TransactionRate        float64 `json:"transaction_rate"`
	OutputVolume           float64 `json:"output_volume"`
	EstimatedBTCValue      float64 `json:"estimated_btc_sent"`
	EstimatedTxValueUSD    float64 `json:"estimated_transaction_volume_usd"`
	TotalBTC               float64 `json:"total_btc"`
	MarketCap              float64 `json:"market_cap"`
	TradeVolumeUSD         float64 `json:"trade_volume_usd"`
	Blocks                 int64   `json:"blocks_size"`
	NextRetarget           int64   `json:"nextretarget"`
	Difficulty             float64 `json:"difficulty"`
	EstimatedTxValue       float64 `json:"estimated_transaction_volume"`
	BlocksCount            int64   `json:"n_blocks_total"`
	MinutesBetweenBlocks   float64 `json:"minutes_between_blocks"`
	Timestamp              float64 `json:"timestamp"`
}

// SingleStatValue represents a single statistic value from Blockchain.com
type SingleStatValue struct {
	Name        string  `json:"name"`
	Unit        string  `json:"unit"`
	Period      string  `json:"period"`
	Description string  `json:"description"`
	Status      string  `json:"status"`
	Values      []struct {
		X float64 `json:"x"` // Timestamp
		Y float64 `json:"y"` // Value
	} `json:"values"`
}

// ChartData represents historical chart data
type ChartData struct {
	Status      string `json:"status"`
	Name        string `json:"name"`
	Unit        string `json:"unit"`
	Period      string `json:"period"`
	Description string `json:"description"`
	Values      []struct {
		X float64 `json:"x"` // Timestamp (Unix)
		Y float64 `json:"y"` // Value
	} `json:"values"`
}

// PoolsData represents mining pool distribution
type PoolsData struct {
	Pools []struct {
		PoolName string `json:"pool_name"`
		Blocks   int    `json:"blocks"`
	} `json:"pools"`
}

// GetBitcoinStats retrieves comprehensive Bitcoin network statistics
func (bc *BlockchainClient) GetBitcoinStats() (*BitcoinStats, error) {
	endpoint := "/stats?format=json"
	
	data, err := bc.makeRequest(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Bitcoin stats: %w", err)
	}

	var stats BitcoinStats
	if err := json.Unmarshal(data, &stats); err != nil {
		return nil, fmt.Errorf("failed to unmarshal Bitcoin stats: %w", err)
	}

	bc.logger.Info("Successfully fetched Bitcoin stats", 
		"price_usd", stats.MarketPriceUSD,
		"hash_rate", stats.HashRate,
		"difficulty", stats.Difficulty)

	return &stats, nil
}

// GetBitcoinPrice retrieves current Bitcoin price from Blockchain.com
func (bc *BlockchainClient) GetBitcoinPrice() (float64, error) {
	stats, err := bc.GetBitcoinStats()
	if err != nil {
		return 0, fmt.Errorf("failed to get Bitcoin price: %w", err)
	}
	return stats.MarketPriceUSD, nil
}

// GetHashRate retrieves current network hash rate
func (bc *BlockchainClient) GetHashRate() (float64, error) {
	stats, err := bc.GetBitcoinStats()
	if err != nil {
		return 0, fmt.Errorf("failed to get hash rate: %w", err)
	}
	return stats.HashRate, nil
}

// GetDifficulty retrieves current mining difficulty
func (bc *BlockchainClient) GetDifficulty() (float64, error) {
	stats, err := bc.GetBitcoinStats()
	if err != nil {
		return 0, fmt.Errorf("failed to get difficulty: %w", err)
	}
	return stats.Difficulty, nil
}

// GetSingleStat retrieves a specific statistic
func (bc *BlockchainClient) GetSingleStat(statName string) (*SingleStatValue, error) {
	endpoint := fmt.Sprintf("/single/%s?format=json", statName)
	
	data, err := bc.makeRequest(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch single stat %s: %w", statName, err)
	}

	var stat SingleStatValue
	if err := json.Unmarshal(data, &stat); err != nil {
		return nil, fmt.Errorf("failed to unmarshal single stat: %w", err)
	}

	bc.logger.Info("Successfully fetched single stat", "stat", statName, "values_count", len(stat.Values))
	return &stat, nil
}

// GetChartData retrieves historical chart data for specific metrics
func (bc *BlockchainClient) GetChartData(chartType string, timespan *string) (*ChartData, error) {
	endpoint := fmt.Sprintf("/charts/%s?format=json", chartType)
	if timespan != nil {
		endpoint += fmt.Sprintf("&timespan=%s", *timespan)
	}
	
	data, err := bc.makeRequest(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch chart data for %s: %w", chartType, err)
	}

	var chartData ChartData
	if err := json.Unmarshal(data, &chartData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal chart data: %w", err)
	}

	bc.logger.Info("Successfully fetched chart data", 
		"chart_type", chartType, 
		"values_count", len(chartData.Values))

	return &chartData, nil
}

// GetHashRateHistory retrieves historical hash rate data
func (bc *BlockchainClient) GetHashRateHistory(timespan string) (*ChartData, error) {
	return bc.GetChartData("hash-rate", &timespan)
}

// GetDifficultyHistory retrieves historical difficulty data
func (bc *BlockchainClient) GetDifficultyHistory(timespan string) (*ChartData, error) {
	return bc.GetChartData("difficulty", &timespan)
}

// GetTransactionCountHistory retrieves historical transaction count
func (bc *BlockchainClient) GetTransactionCountHistory(timespan string) (*ChartData, error) {
	return bc.GetChartData("n-transactions", &timespan)
}

// GetBlockSizeHistory retrieves historical average block size
func (bc *BlockchainClient) GetBlockSizeHistory(timespan string) (*ChartData, error) {
	return bc.GetChartData("avg-block-size", &timespan)
}

// GetMempoolSize retrieves current mempool transaction count
func (bc *BlockchainClient) GetMempoolSize() (int64, error) {
	endpoint := "/q/unconfirmedcount"
	
	data, err := bc.makeRequest(endpoint)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch mempool size: %w", err)
	}

	var count int64
	if err := json.Unmarshal(data, &count); err != nil {
		return 0, fmt.Errorf("failed to unmarshal mempool size: %w", err)
	}

	return count, nil
}

// GetLatestBlockHeight retrieves the latest block height
func (bc *BlockchainClient) GetLatestBlockHeight() (int64, error) {
	endpoint := "/q/getblockcount"
	
	data, err := bc.makeRequest(endpoint)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch block height: %w", err)
	}

	var height int64
	if err := json.Unmarshal(data, &height); err != nil {
		return 0, fmt.Errorf("failed to unmarshal block height: %w", err)
	}

	return height, nil
}

// GetTotalBitcoinsInCirculation retrieves total bitcoins in circulation
func (bc *BlockchainClient) GetTotalBitcoinsInCirculation() (float64, error) {
	endpoint := "/q/totalbc"
	
	data, err := bc.makeRequest(endpoint)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch total bitcoins: %w", err)
	}

	var total float64
	if err := json.Unmarshal(data, &total); err != nil {
		return 0, fmt.Errorf("failed to unmarshal total bitcoins: %w", err)
	}

	// Convert from satoshis to bitcoins
	return total / 100000000, nil
}

// GetMiningPoolDistribution retrieves mining pool distribution
func (bc *BlockchainClient) GetMiningPoolDistribution() (*PoolsData, error) {
	endpoint := "/pools?format=json"
	
	data, err := bc.makeRequest(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch mining pools: %w", err)
	}

	var pools PoolsData
	if err := json.Unmarshal(data, &pools); err != nil {
		return nil, fmt.Errorf("failed to unmarshal mining pools: %w", err)
	}

	bc.logger.Info("Successfully fetched mining pool distribution", "pools_count", len(pools.Pools))
	return &pools, nil
}

// GetNetworkSummary provides a comprehensive network summary
func (bc *BlockchainClient) GetNetworkSummary() (map[string]interface{}, error) {
	stats, err := bc.GetBitcoinStats()
	if err != nil {
		return nil, fmt.Errorf("failed to get network summary: %w", err)
	}

	blockHeight, _ := bc.GetLatestBlockHeight()
	mempoolSize, _ := bc.GetMempoolSize()
	totalBTC, _ := bc.GetTotalBitcoinsInCirculation()

	summary := map[string]interface{}{
		"price_usd":             stats.MarketPriceUSD,
		"market_cap":            stats.MarketCap,
		"hash_rate":             stats.HashRate,
		"difficulty":            stats.Difficulty,
		"block_height":          blockHeight,
		"mempool_size":          mempoolSize,
		"total_btc":             totalBTC,
		"transaction_rate":      stats.TransactionRate,
		"minutes_between_blocks": stats.MinutesBetweenBlocks,
		"total_fees_btc":        stats.TotalFeesBTC,
		"trade_volume_usd":      stats.TradeVolumeUSD,
		"last_updated":          time.Now().Unix(),
	}

	return summary, nil
}

// makeRequest makes an HTTP request to the Blockchain.com API
func (bc *BlockchainClient) makeRequest(endpoint string) ([]byte, error) {
	reqURL := bc.baseURL + endpoint

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("User-Agent", "CryptoIndicatorDashboard/1.0")

	bc.logger.Debug("Making Blockchain.com API request", 
		"url", reqURL,
		"endpoint", endpoint)

	resp, err := bc.httpClient.Do(req)
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
		bc.logger.Error("Blockchain.com API request failed", 
			"status_code", resp.StatusCode,
			"response", string(body))
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	return body, nil
}

// HealthCheck performs a health check on the Blockchain.com service
func (bc *BlockchainClient) HealthCheck() error {
	// Try to fetch Bitcoin price as a simple health check
	_, err := bc.GetBitcoinPrice()
	if err != nil {
		return fmt.Errorf("Blockchain.com health check failed: %w", err)
	}
	return nil
}