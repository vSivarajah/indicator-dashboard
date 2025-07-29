package models

import (
	"time"
	"gorm.io/gorm"
	"crypto-indicator-dashboard/internal/domain/entities"
)

// Indicator represents a market indicator
type Indicator struct {
	ID          uint      `json:"id" gorm:"primarykey"`
	Name        string    `json:"name" gorm:"not null;index"`
	Type        string    `json:"type" gorm:"not null"` // crypto, macro, on-chain
	Value       string    `json:"value" gorm:"not null"`
	NumericValue float64  `json:"numeric_value"`
	Change      string    `json:"change"`
	RiskLevel   string    `json:"risk_level"` // low, medium, high
	Status      string    `json:"status"`
	Description string    `json:"description"`
	Source      string    `json:"source"`
	Timestamp   time.Time `json:"timestamp" gorm:"not null;index"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// PriceData represents historical price data
type PriceData struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	Symbol    string    `json:"symbol" gorm:"not null;index"`
	Price     float64   `json:"price" gorm:"not null"`
	Volume    float64   `json:"volume"`
	MarketCap float64   `json:"market_cap"`
	Timestamp time.Time `json:"timestamp" gorm:"not null;index"`
	CreatedAt time.Time `json:"created_at"`
}

// OnChainData represents blockchain-specific metrics
type OnChainData struct {
	ID               uint      `json:"id" gorm:"primarykey"`
	Symbol           string    `json:"symbol" gorm:"not null;index"`
	MarketValue      float64   `json:"market_value"`
	RealizedValue    float64   `json:"realized_value"`
	MVRVRatio        float64   `json:"mvrv_ratio"`
	MVRVZScore       float64   `json:"mvrv_zscore"`
	ActiveAddresses  uint64    `json:"active_addresses"`
	TransactionCount uint64    `json:"transaction_count"`
	NetworkHashRate  float64   `json:"network_hash_rate"`
	Timestamp        time.Time `json:"timestamp" gorm:"not null;index"`
	CreatedAt        time.Time `json:"created_at"`
}

// MacroData represents macroeconomic indicators
type MacroData struct {
	ID           uint      `json:"id" gorm:"primarykey"`
	Indicator    string    `json:"indicator" gorm:"not null;index"` // inflation, interest_rate, etc.
	Value        float64   `json:"value" gorm:"not null"`
	Change       float64   `json:"change"`
	Country      string    `json:"country" gorm:"default:'US'"`
	Source       string    `json:"source"`
	ReleaseDate  time.Time `json:"release_date"`
	Timestamp    time.Time `json:"timestamp" gorm:"not null;index"`
	CreatedAt    time.Time `json:"created_at"`
}

// Portfolio represents a user's portfolio
type Portfolio struct {
	ID          uint              `json:"id" gorm:"primarykey"`
	UserID      string            `json:"user_id" gorm:"not null;index"`
	Name        string            `json:"name" gorm:"not null"`
	Holdings    []PortfolioHolding `json:"holdings" gorm:"foreignKey:PortfolioID"`
	TotalValue  float64           `json:"total_value"`
	RiskLevel   string            `json:"risk_level"`
	LastUpdated time.Time         `json:"last_updated"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

// PortfolioHolding represents individual holdings in a portfolio
type PortfolioHolding struct {
	ID           uint    `json:"id" gorm:"primarykey"`
	PortfolioID  uint    `json:"portfolio_id" gorm:"not null;index"`
	Symbol       string  `json:"symbol" gorm:"not null"`
	Amount       float64 `json:"amount" gorm:"not null"`
	AveragePrice float64 `json:"average_price"`
	CurrentPrice float64 `json:"current_price"`
	Value        float64 `json:"value"`
	PnL          float64 `json:"pnl"`
	PnLPercent   float64 `json:"pnl_percent"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// MarketCycle represents market cycle analysis
type MarketCycle struct {
	ID               uint      `json:"id" gorm:"primarykey"`
	Stage            string    `json:"stage" gorm:"not null"` // bear, early_bull, mid_bull, late_bull
	Confidence       float64   `json:"confidence"`
	DominanceLevel   float64   `json:"dominance_level"`
	FearGreedIndex   int       `json:"fear_greed_index"`
	MVRVZScore       float64   `json:"mvrv_zscore"`
	BubbleRisk       string    `json:"bubble_risk"`
	EstimatedDuration int      `json:"estimated_duration"` // months
	Timestamp        time.Time `json:"timestamp" gorm:"not null;index"`
	CreatedAt        time.Time `json:"created_at"`
}

// DCAStrategy represents a dollar cost averaging strategy
type DCAStrategy struct {
	ID               uint      `json:"id" gorm:"primarykey"`
	UserID           string    `json:"user_id" gorm:"not null;index"`
	Name             string    `json:"name" gorm:"not null"`
	Symbol           string    `json:"symbol" gorm:"not null"` // BTC, ETH, etc.
	Amount           float64   `json:"amount" gorm:"not null"` // Amount per purchase
	Frequency        string    `json:"frequency" gorm:"not null"` // daily, weekly, monthly
	StartDate        time.Time `json:"start_date" gorm:"not null"`
	EndDate          *time.Time `json:"end_date"` // Optional end date
	IsActive         bool      `json:"is_active" gorm:"default:true"`
	TotalInvested    float64   `json:"total_invested" gorm:"default:0"`
	TotalQuantity    float64   `json:"total_quantity" gorm:"default:0"`
	AveragePrice     float64   `json:"average_price" gorm:"default:0"`
	CurrentValue     float64   `json:"current_value" gorm:"default:0"`
	TotalReturn      float64   `json:"total_return" gorm:"default:0"`
	TotalReturnPct   float64   `json:"total_return_pct" gorm:"default:0"`
	PurchaseCount    int       `json:"purchase_count" gorm:"default:0"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// DCAPurchase represents individual DCA purchases
type DCAPurchase struct {
	ID           uint      `json:"id" gorm:"primarykey"`
	StrategyID   uint      `json:"strategy_id" gorm:"not null;index"`
	Strategy     DCAStrategy `json:"strategy" gorm:"foreignKey:StrategyID"`
	Date         time.Time `json:"date" gorm:"not null;index"`
	Amount       float64   `json:"amount" gorm:"not null"` // USD amount invested
	Price        float64   `json:"price" gorm:"not null"` // Price per coin at time of purchase
	Quantity     float64   `json:"quantity" gorm:"not null"` // Quantity purchased
	MarketCap    float64   `json:"market_cap"` // Market cap at time of purchase
	MVRVZScore   float64   `json:"mvrv_zscore"` // MVRV Z-Score at time of purchase
	FearGreed    int       `json:"fear_greed"` // Fear & Greed index at purchase
	IsSimulated  bool      `json:"is_simulated" gorm:"default:false"` // True for backtesting
	CreatedAt    time.Time `json:"created_at"`
}

// DCASimulation represents backtesting results
type DCASimulation struct {
	ID                uint      `json:"id" gorm:"primarykey"`
	UserID            string    `json:"user_id" gorm:"not null;index"`
	Symbol            string    `json:"symbol" gorm:"not null"`
	Amount            float64   `json:"amount" gorm:"not null"`
	Frequency         string    `json:"frequency" gorm:"not null"`
	StartDate         time.Time `json:"start_date" gorm:"not null"`
	EndDate           time.Time `json:"end_date" gorm:"not null"`
	TotalInvested     float64   `json:"total_invested"`
	TotalQuantity     float64   `json:"total_quantity"`
	FinalValue        float64   `json:"final_value"`
	TotalReturn       float64   `json:"total_return"`
	TotalReturnPct    float64   `json:"total_return_pct"`
	AnnualizedReturn  float64   `json:"annualized_return"`
	MaxDrawdown       float64   `json:"max_drawdown"`
	MaxDrawdownPct    float64   `json:"max_drawdown_pct"`
	SharpeRatio       float64   `json:"sharpe_ratio"`
	PurchaseCount     int       `json:"purchase_count"`
	BestPurchaseDate  time.Time `json:"best_purchase_date"`
	WorstPurchaseDate time.Time `json:"worst_purchase_date"`
	AvgMVRVAtPurchase float64   `json:"avg_mvrv_at_purchase"`
	AvgFearGreedAtPurchase int  `json:"avg_fear_greed_at_purchase"`
	CreatedAt         time.Time `json:"created_at"`
}

// AutoMigrate runs database migrations
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		// Legacy models
		&Indicator{},
		&PriceData{},
		&OnChainData{},
		&MacroData{},
		&Portfolio{},
		&PortfolioHolding{},
		&MarketCycle{},
		&DCAStrategy{},
		&DCAPurchase{},
		&DCASimulation{},
		// New architecture entities
		&entities.CryptoPrice{},
		&entities.BitcoinDominance{},
		&entities.MarketMetrics{},
		&entities.PriceAlert{},
		&entities.TradingPair{},
		&entities.MarketData{},
	)
}