package entities

import (
	"time"
)

// DCAStrategy represents a dollar cost averaging strategy
type DCAStrategy struct {
	ID               uint       `json:"id"`
	UserID           string     `json:"user_id"`
	Name             string     `json:"name"`
	Symbol           string     `json:"symbol"` // BTC, ETH, etc.
	Amount           float64    `json:"amount"` // Amount per purchase
	Frequency        string     `json:"frequency"` // daily, weekly, monthly
	StartDate        time.Time  `json:"start_date"`
	EndDate          *time.Time `json:"end_date"` // Optional end date
	IsActive         bool       `json:"is_active"`
	TotalInvested    float64    `json:"total_invested"`
	TotalQuantity    float64    `json:"total_quantity"`
	AveragePrice     float64    `json:"average_price"`
	CurrentValue     float64    `json:"current_value"`
	TotalReturn      float64    `json:"total_return"`
	TotalReturnPct   float64    `json:"total_return_pct"`
	PurchaseCount    int        `json:"purchase_count"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

// DCAPurchase represents individual DCA purchases
type DCAPurchase struct {
	ID           uint        `json:"id"`
	StrategyID   uint        `json:"strategy_id"`
	Strategy     DCAStrategy `json:"strategy"`
	Date         time.Time   `json:"date"`
	Amount       float64     `json:"amount"` // USD amount invested
	Price        float64     `json:"price"`  // Price per coin at time of purchase
	Quantity     float64     `json:"quantity"` // Quantity purchased
	MarketCap    float64     `json:"market_cap"` // Market cap at time of purchase
	MVRVZScore   float64     `json:"mvrv_zscore"` // MVRV Z-Score at time of purchase
	FearGreed    int         `json:"fear_greed"` // Fear & Greed index at purchase
	IsSimulated  bool        `json:"is_simulated"` // True for backtesting
	CreatedAt    time.Time   `json:"created_at"`
}

// DCASimulation represents backtesting results
type DCASimulation struct {
	ID                     uint      `json:"id"`
	UserID                 string    `json:"user_id"`
	Symbol                 string    `json:"symbol"`
	Amount                 float64   `json:"amount"`
	Frequency              string    `json:"frequency"`
	StartDate              time.Time `json:"start_date"`
	EndDate                time.Time `json:"end_date"`
	TotalInvested          float64   `json:"total_invested"`
	TotalQuantity          float64   `json:"total_quantity"`
	FinalValue             float64   `json:"final_value"`
	TotalReturn            float64   `json:"total_return"`
	TotalReturnPct         float64   `json:"total_return_pct"`
	AnnualizedReturn       float64   `json:"annualized_return"`
	MaxDrawdown            float64   `json:"max_drawdown"`
	MaxDrawdownPct         float64   `json:"max_drawdown_pct"`
	SharpeRatio            float64   `json:"sharpe_ratio"`
	PurchaseCount          int       `json:"purchase_count"`
	BestPurchaseDate       time.Time `json:"best_purchase_date"`
	WorstPurchaseDate      time.Time `json:"worst_purchase_date"`
	AvgMVRVAtPurchase      float64   `json:"avg_mvrv_at_purchase"`
	AvgFearGreedAtPurchase int       `json:"avg_fear_greed_at_purchase"`
	CreatedAt              time.Time `json:"created_at"`
}

// DCARequest represents a DCA simulation request
type DCARequest struct {
	UserID     string    `json:"user_id"`
	Symbol     string    `json:"symbol" binding:"required"`
	Amount     float64   `json:"amount" binding:"required,gt=0"`
	Frequency  string    `json:"frequency" binding:"required,oneof=daily weekly monthly"`
	StartDate  time.Time `json:"start_date" binding:"required"`
	EndDate    time.Time `json:"end_date" binding:"required"`
	IsBacktest bool      `json:"is_backtest"`
}