package entities

import (
	"time"
)

// Indicator represents a market indicator
type Indicator struct {
	ID           uint                   `json:"id" gorm:"primaryKey"`
	Name         string                 `json:"name" gorm:"not null"`
	Type         string                 `json:"type" gorm:"not null"` // crypto, macro, on-chain
	Value        float64                `json:"value"`
	StringValue  string                 `json:"string_value,omitempty"`
	Change       string                 `json:"change"`
	RiskLevel    string                 `json:"risk_level"` // low, medium, high
	Status       string                 `json:"status"`
	Description  string                 `json:"description"`
	Source       string                 `json:"source"`
	Confidence   float64                `json:"confidence"` // 0.0 to 1.0
	Metadata     map[string]interface{} `json:"metadata" gorm:"serializer:json"`
	Timestamp    time.Time              `json:"timestamp"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
}

// TableName returns the table name for Indicator
func (Indicator) TableName() string {
	return "indicators"
}

// MVRVData represents MVRV calculation data
type MVRVData struct {
	Date          time.Time `json:"date"`
	Price         float64   `json:"price"`
	MarketCap     float64   `json:"market_cap"`
	RealizedCap   float64   `json:"realized_cap"`
	MVRVRatio     float64   `json:"mvrv_ratio"`
	MVRVZScore    float64   `json:"mvrv_zscore"`
	CircSupply    float64   `json:"circulating_supply"`
}

// MVRVResult represents the final MVRV analysis
type MVRVResult struct {
	CurrentZScore    float64            `json:"current_zscore"`
	MVRVRatio        float64            `json:"mvrv_ratio"`
	MarketCap        float64            `json:"market_cap"`
	RealizedCap      float64            `json:"realized_cap"`
	Price            float64            `json:"price"`
	RiskLevel        string             `json:"risk_level"`
	Status           string             `json:"status"`
	HistoricalData   []MVRVData         `json:"historical_data"`
	LastUpdated      time.Time          `json:"last_updated"`
	ZScoreThresholds map[string]float64 `json:"zscore_thresholds"`
}

// DominanceResult represents Bitcoin dominance analysis
type DominanceResult struct {
	CurrentDominance  float64     `json:"current_dominance"`
	Change24h         float64     `json:"change_24h"`
	Change7d          float64     `json:"change_7d"`
	Change30d         float64     `json:"change_30d"`
	Trend             string      `json:"trend"`
	TrendStrength     string      `json:"trend_strength"`
	RiskLevel         string      `json:"risk_level"`
	Status            string      `json:"status"`
	MarketCycleStage  string      `json:"market_cycle_stage"`
	AltSeasonSignal   bool        `json:"alt_season_signal"`
	CriticalLevels    map[string]float64 `json:"critical_levels"`
	LastUpdated       time.Time   `json:"last_updated"`
}

// FearGreedResult represents Fear & Greed index analysis
type FearGreedResult struct {
	CurrentValue          int              `json:"current_value"`
	Change24h             int              `json:"change_24h"`
	Change7d              int              `json:"change_7d"`
	Classification        string           `json:"classification"`
	RiskLevel             string           `json:"risk_level"`
	Status                string           `json:"status"`
	Components            map[string]int   `json:"components"`
	TradingRecommendation string           `json:"trading_recommendation"`
	DataSource            string           `json:"data_source"`
	NextUpdate            time.Time        `json:"next_update"`
	LastUpdated           time.Time        `json:"last_updated"`
}

// BubbleRiskResult represents bubble risk analysis
type BubbleRiskResult struct {
	CurrentRiskScore      float64            `json:"current_risk_score"`
	RiskCategory          string             `json:"risk_category"`
	ConfidenceLevel       float64            `json:"confidence_level"`
	RiskLevel             string             `json:"risk_level"`
	Status                string             `json:"status"`
	Components            map[string]float64 `json:"components"`
	TradingRecommendation string             `json:"trading_recommendation"`
	DataSource            string             `json:"data_source"`
	CriticalLevels        map[string]float64 `json:"critical_levels"`
	LastUpdated           time.Time          `json:"last_updated"`
}

// MarketCycle represents market cycle analysis
type MarketCycle struct {
	ID                uint      `json:"id" gorm:"primaryKey"`
	Stage             string    `json:"stage"` // bear, early_bull, mid_bull, late_bull
	Confidence        float64   `json:"confidence"`
	DominanceLevel    float64   `json:"dominance_level"`
	FearGreedIndex    int       `json:"fear_greed_index"`
	MVRVZScore        float64   `json:"mvrv_zscore"`
	BubbleRisk        string    `json:"bubble_risk"`
	EstimatedDuration int       `json:"estimated_duration"` // months
	Timestamp         time.Time `json:"timestamp"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// TableName returns the table name for MarketCycle
func (MarketCycle) TableName() string {
	return "market_cycles"
}