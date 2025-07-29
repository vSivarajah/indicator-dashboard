package dto

import (
	"errors"
	"time"
	"crypto-indicator-dashboard/internal/domain/entities"
)

// CreatePortfolioRequest represents a request to create a portfolio
type CreatePortfolioRequest struct {
	UserID string `json:"user_id" binding:"required"`
	Name   string `json:"name" binding:"required,min=1,max=100"`
}

// Validate validates the create portfolio request
func (r *CreatePortfolioRequest) Validate() error {
	if r.UserID == "" {
		return errors.New("user ID is required")
	}
	if r.Name == "" {
		return errors.New("portfolio name is required")
	}
	if len(r.Name) > 100 {
		return errors.New("portfolio name must be less than 100 characters")
	}
	return nil
}

// AddHoldingRequest represents a request to add a holding
type AddHoldingRequest struct {
	PortfolioID  uint    `json:"portfolio_id" binding:"required"`
	Symbol       string  `json:"symbol" binding:"required,min=1,max=10"`
	Amount       float64 `json:"amount" binding:"required,gt=0"`
	AveragePrice float64 `json:"average_price" binding:"required,gt=0"`
}

// Validate validates the add holding request
func (r *AddHoldingRequest) Validate() error {
	if r.PortfolioID == 0 {
		return errors.New("portfolio ID is required")
	}
	if r.Symbol == "" {
		return errors.New("symbol is required")
	}
	if r.Amount <= 0 {
		return errors.New("amount must be greater than 0")
	}
	if r.AveragePrice <= 0 {
		return errors.New("average price must be greater than 0")
	}
	return nil
}

// UpdateHoldingRequest represents a request to update a holding
type UpdateHoldingRequest struct {
	HoldingID    uint    `json:"holding_id" binding:"required"`
	Amount       float64 `json:"amount" binding:"required,gt=0"`
	AveragePrice float64 `json:"average_price" binding:"required,gt=0"`
}

// Validate validates the update holding request
func (r *UpdateHoldingRequest) Validate() error {
	if r.HoldingID == 0 {
		return errors.New("holding ID is required")
	}
	if r.Amount <= 0 {
		return errors.New("amount must be greater than 0")
	}
	if r.AveragePrice <= 0 {
		return errors.New("average price must be greater than 0")
	}
	return nil
}

// PortfolioResponse represents a portfolio response
type PortfolioResponse struct {
	ID          uint                `json:"id"`
	UserID      string              `json:"user_id"`
	Name        string              `json:"name"`
	Holdings    []HoldingResponse   `json:"holdings"`
	TotalValue  float64             `json:"total_value"`
	RiskLevel   string              `json:"risk_level"`
	LastUpdated time.Time           `json:"last_updated"`
	CreatedAt   time.Time           `json:"created_at"`
}

// NewPortfolioResponse creates a new portfolio response from entity
func NewPortfolioResponse(portfolio *entities.Portfolio) *PortfolioResponse {
	holdings := make([]HoldingResponse, len(portfolio.Holdings))
	for i, holding := range portfolio.Holdings {
		holdings[i] = *NewHoldingResponse(&holding)
	}
	
	return &PortfolioResponse{
		ID:          portfolio.ID,
		UserID:      portfolio.UserID,
		Name:        portfolio.Name,
		Holdings:    holdings,
		TotalValue:  portfolio.TotalValue,
		RiskLevel:   portfolio.RiskLevel,
		LastUpdated: portfolio.LastUpdated,
		CreatedAt:   portfolio.CreatedAt,
	}
}

// HoldingResponse represents a holding response
type HoldingResponse struct {
	ID           uint      `json:"id"`
	PortfolioID  uint      `json:"portfolio_id"`
	Symbol       string    `json:"symbol"`
	Amount       float64   `json:"amount"`
	AveragePrice float64   `json:"average_price"`
	CurrentPrice float64   `json:"current_price"`
	Value        float64   `json:"value"`
	PnL          float64   `json:"pnl"`
	PnLPercent   float64   `json:"pnl_percent"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// NewHoldingResponse creates a new holding response from entity
func NewHoldingResponse(holding *entities.PortfolioHolding) *HoldingResponse {
	return &HoldingResponse{
		ID:           holding.ID,
		PortfolioID:  holding.PortfolioID,
		Symbol:       holding.Symbol,
		Amount:       holding.Amount,
		AveragePrice: holding.AveragePrice,
		CurrentPrice: holding.CurrentPrice,
		Value:        holding.Value,
		PnL:          holding.PnL,
		PnLPercent:   holding.PnLPercent,
		CreatedAt:    holding.CreatedAt,
		UpdatedAt:    holding.UpdatedAt,
	}
}

// PortfolioListResponse represents a list of portfolios
type PortfolioListResponse struct {
	Portfolios []PortfolioResponse `json:"portfolios"`
	Count      int                 `json:"count"`
}

// NewPortfolioListResponse creates a new portfolio list response
func NewPortfolioListResponse(portfolios []entities.Portfolio) *PortfolioListResponse {
	responses := make([]PortfolioResponse, len(portfolios))
	for i, portfolio := range portfolios {
		responses[i] = *NewPortfolioResponse(&portfolio)
	}
	
	return &PortfolioListResponse{
		Portfolios: responses,
		Count:      len(responses),
	}
}

// PortfolioSummaryResponse represents portfolio summary data
type PortfolioSummaryResponse struct {
	TotalValue        float64                      `json:"total_value"`
	TotalPnL          float64                      `json:"total_pnl"`
	TotalPnLPercent   float64                      `json:"total_pnl_percent"`
	DayChange         float64                      `json:"day_change"`
	DayChangePercent  float64                      `json:"day_change_percent"`
	TopPerformer      *HoldingResponse             `json:"top_performer"`
	WorstPerformer    *HoldingResponse             `json:"worst_performer"`
	AllocationByAsset []entities.AssetAllocation   `json:"allocation_by_asset"`
	RiskMetrics       entities.PortfolioRiskMetrics `json:"risk_metrics"`
}

// NewPortfolioSummaryResponse creates a new portfolio summary response
func NewPortfolioSummaryResponse(summary *entities.PortfolioSummary) *PortfolioSummaryResponse {
	var topPerformer, worstPerformer *HoldingResponse
	
	if summary.TopPerformer != nil {
		topPerformer = NewHoldingResponse(summary.TopPerformer)
	}
	if summary.WorstPerformer != nil {
		worstPerformer = NewHoldingResponse(summary.WorstPerformer)
	}
	
	return &PortfolioSummaryResponse{
		TotalValue:        summary.TotalValue,
		TotalPnL:          summary.TotalPnL,
		TotalPnLPercent:   summary.TotalPnLPercent,
		DayChange:         summary.DayChange,
		DayChangePercent:  summary.DayChangePercent,
		TopPerformer:      topPerformer,
		WorstPerformer:    worstPerformer,
		AllocationByAsset: summary.AllocationByAsset,
		RiskMetrics:       summary.RiskMetrics,
	}
}