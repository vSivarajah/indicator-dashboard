package entities

import "time"

// CryptoPrice represents cryptocurrency price data
type CryptoPrice struct {
	ID               uint      `json:"id" gorm:"primaryKey"`
	Symbol           string    `json:"symbol" gorm:"index;not null"`
	Name             string    `json:"name"`
	Price            float64   `json:"price"`
	Volume24h        float64   `json:"volume_24h"`
	MarketCap        float64   `json:"market_cap"`
	PercentChange1h  float64   `json:"percent_change_1h"`
	PercentChange24h float64   `json:"percent_change_24h"`
	PercentChange7d  float64   `json:"percent_change_7d"`
	PercentChange30d float64   `json:"percent_change_30d"`
	LastUpdated      time.Time `json:"last_updated"`
	DataSource       string    `json:"data_source"`
	CreatedAt        time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt        time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName returns the table name for CryptoPrice
func (CryptoPrice) TableName() string {
	return "crypto_prices"
}

// BitcoinDominance represents Bitcoin market dominance data
type BitcoinDominance struct {
	ID                uint      `json:"id" gorm:"primaryKey"`
	CurrentDominance  float64   `json:"current_dominance"`
	PreviousDominance float64   `json:"previous_dominance"`
	Change24h         float64   `json:"change_24h"`
	ChangePercent24h  float64   `json:"change_percent_24h"`
	LastUpdated       time.Time `json:"last_updated"`
	DataSource        string    `json:"data_source"`
	Confidence        float64   `json:"confidence"` // Confidence level (0-1)
	CreatedAt         time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt         time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName returns the table name for BitcoinDominance
func (BitcoinDominance) TableName() string {
	return "bitcoin_dominance"
}

// MarketMetrics represents overall market metrics
type MarketMetrics struct {
	ID                    uint      `json:"id" gorm:"primaryKey"`
	TotalMarketCap        float64   `json:"total_market_cap"`
	TotalVolume24h        float64   `json:"total_volume_24h"`
	BitcoinDominance      float64   `json:"bitcoin_dominance"`
	EthereumDominance     float64   `json:"ethereum_dominance"`
	ActiveCryptocurrencies int      `json:"active_cryptocurrencies"`
	ActiveExchanges       int       `json:"active_exchanges"`
	MarketCapChange24h    float64   `json:"market_cap_change_24h"`
	VolumeChange24h       float64   `json:"volume_change_24h"`
	LastUpdated           time.Time `json:"last_updated"`
	DataSource            string    `json:"data_source"`
	CreatedAt             time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt             time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName returns the table name for MarketMetrics
func (MarketMetrics) TableName() string {
	return "market_metrics"
}

// PriceAlert represents a price alert configuration
type PriceAlert struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	UserID        string    `json:"user_id" gorm:"index;not null"`
	Symbol        string    `json:"symbol" gorm:"not null"`
	AlertType     string    `json:"alert_type"` // "above", "below", "percentage_change"
	TargetPrice   float64   `json:"target_price"`
	TargetPercent float64   `json:"target_percent"`
	IsActive      bool      `json:"is_active" gorm:"default:true"`
	LastTriggered *time.Time `json:"last_triggered"`
	CreatedAt     time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName returns the table name for PriceAlert
func (PriceAlert) TableName() string {
	return "price_alerts"
}

// TradingPair represents a trading pair on an exchange
type TradingPair struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	BaseAsset  string    `json:"base_asset"`
	QuoteAsset string    `json:"quote_asset"`
	Symbol     string    `json:"symbol" gorm:"uniqueIndex"`
	Exchange   string    `json:"exchange"`
	Price      float64   `json:"price"`
	Volume24h  float64   `json:"volume_24h"`
	IsActive   bool      `json:"is_active" gorm:"default:true"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName returns the table name for TradingPair
func (TradingPair) TableName() string {
	return "trading_pairs"
}

// MarketDataSummary provides a summary of all market data
type MarketDataSummary struct {
	TotalMarketCap       float64                     `json:"total_market_cap"`
	TotalVolume24h       float64                     `json:"total_volume_24h"`
	BitcoinDominance     *BitcoinDominance           `json:"bitcoin_dominance"`
	TopCryptocurrencies  map[string]*CryptoPrice     `json:"top_cryptocurrencies"`
	MarketTrend          string                      `json:"market_trend"` // "bullish", "bearish", "sideways"
	FearGreedIndex       float64                     `json:"fear_greed_index"`
	LastUpdated          time.Time                   `json:"last_updated"`
}

// GetTrendIndicator returns a simple trend indicator based on 24h changes
func (cp *CryptoPrice) GetTrendIndicator() string {
	if cp.PercentChange24h > 5 {
		return "strong_bullish"
	} else if cp.PercentChange24h > 0 {
		return "bullish"
	} else if cp.PercentChange24h > -5 {
		return "bearish"
	} else {
		return "strong_bearish"
	}
}

// IsHighVolatility checks if the price has high volatility
func (cp *CryptoPrice) IsHighVolatility() bool {
	return abs(cp.PercentChange24h) > 10 || abs(cp.PercentChange1h) > 5
}

// GetDominanceTrend returns the dominance trend based on 24h change
func (bd *BitcoinDominance) GetDominanceTrend() string {
	if bd.ChangePercent24h > 1 {
		return "increasing"
	} else if bd.ChangePercent24h < -1 {
		return "decreasing"
	} else {
		return "stable"
	}
}

// IsAltSeasonIndicator checks if Bitcoin dominance suggests alt season
func (bd *BitcoinDominance) IsAltSeasonIndicator() bool {
	return bd.CurrentDominance < 42 && bd.ChangePercent24h < -0.5
}

// GetConfidenceLevel returns a human-readable confidence level
func (bd *BitcoinDominance) GetConfidenceLevel() string {
	if bd.Confidence >= 0.9 {
		return "high"
	} else if bd.Confidence >= 0.7 {
		return "medium"
	} else {
		return "low"
	}
}

// InflationResult represents inflation analysis results
type InflationResult struct {
	CurrentRate      float64   `json:"current_rate"`
	PreviousRate     float64   `json:"previous_rate"`
	Change           float64   `json:"change"`
	ChangePercent    float64   `json:"change_percent"`
	Trend            string    `json:"trend"` // "increasing", "decreasing", "stable"
	ImpactOnCrypto   string    `json:"impact_on_crypto"` // "positive", "negative", "neutral"
	LastUpdated      time.Time `json:"last_updated"`
	DataSource       string    `json:"data_source"`
	ConfidenceLevel  float64   `json:"confidence_level"`
}

// InterestRateResult represents interest rate analysis results  
type InterestRateResult struct {
	CurrentRate      float64   `json:"current_rate"`
	PreviousRate     float64   `json:"previous_rate"`
	Change           float64   `json:"change"`
	ChangePercent    float64   `json:"change_percent"`
	Trend            string    `json:"trend"` // "increasing", "decreasing", "stable"
	ExpectedChange   string    `json:"expected_change"` // "hike", "cut", "hold"
	ImpactOnCrypto   string    `json:"impact_on_crypto"` // "positive", "negative", "neutral"
	LastUpdated      time.Time `json:"last_updated"`
	DataSource       string    `json:"data_source"`
	ConfidenceLevel  float64   `json:"confidence_level"`
}

// MarketData represents unified market data for testing and services
type MarketData struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	Symbol        string    `json:"symbol" gorm:"index;not null"`
	Name          string    `json:"name"`
	Price         float64   `json:"price"`
	MarketCap     float64   `json:"market_cap"`
	Volume24h     float64   `json:"volume_24h"`
	Change24h     float64   `json:"change_24h"`
	Change7d      float64   `json:"change_7d"`
	Change30d     float64   `json:"change_30d"`
	Dominance     float64   `json:"dominance"`
	CircSupply    float64   `json:"circulating_supply"`
	MaxSupply     float64   `json:"max_supply"`
	Source        string    `json:"source"`
	Confidence    float64   `json:"confidence"`
	LastUpdated   time.Time `json:"last_updated"`
	CreatedAt     time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName returns the table name for MarketData
func (MarketData) TableName() string {
	return "market_data"
}

// Helper function for absolute value
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}