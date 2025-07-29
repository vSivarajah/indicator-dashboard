package usecases

import (
	"context"
	"fmt"
	"crypto-indicator-dashboard/internal/domain/entities"
	"crypto-indicator-dashboard/internal/domain/repositories"
	"crypto-indicator-dashboard/internal/domain/services"
	"crypto-indicator-dashboard/internal/application/dto"
)

// PortfolioUseCase handles portfolio-related business logic
type PortfolioUseCase struct {
	portfolioRepo   repositories.PortfolioRepository
	portfolioSvc    services.PortfolioService
	riskAnalysisSvc services.RiskAnalysisService
}

// NewPortfolioUseCase creates a new portfolio use case
func NewPortfolioUseCase(
	portfolioRepo repositories.PortfolioRepository,
	portfolioSvc services.PortfolioService,
	riskAnalysisSvc services.RiskAnalysisService,
) *PortfolioUseCase {
	return &PortfolioUseCase{
		portfolioRepo:   portfolioRepo,
		portfolioSvc:    portfolioSvc,
		riskAnalysisSvc: riskAnalysisSvc,
	}
}

// CreatePortfolio creates a new portfolio for a user
func (uc *PortfolioUseCase) CreatePortfolio(ctx context.Context, req *dto.CreatePortfolioRequest) (*dto.PortfolioResponse, error) {
	// Validate request
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}
	
	// Create portfolio entity
	portfolio := &entities.Portfolio{
		UserID: req.UserID,
		Name:   req.Name,
	}
	
	// Save portfolio
	if err := uc.portfolioRepo.Create(ctx, portfolio); err != nil {
		return nil, fmt.Errorf("failed to create portfolio: %w", err)
	}
	
	return dto.NewPortfolioResponse(portfolio), nil
}

// GetPortfolio retrieves a portfolio by ID
func (uc *PortfolioUseCase) GetPortfolio(ctx context.Context, portfolioID uint) (*dto.PortfolioResponse, error) {
	portfolio, err := uc.portfolioRepo.GetByID(ctx, portfolioID)
	if err != nil {
		return nil, fmt.Errorf("failed to get portfolio: %w", err)
	}
	
	return dto.NewPortfolioResponse(portfolio), nil
}

// GetUserPortfolios retrieves all portfolios for a user
func (uc *PortfolioUseCase) GetUserPortfolios(ctx context.Context, userID string) (*dto.PortfolioListResponse, error) {
	portfolios, err := uc.portfolioRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user portfolios: %w", err)
	}
	
	return dto.NewPortfolioListResponse(portfolios), nil
}

// AddHolding adds a new holding to a portfolio
func (uc *PortfolioUseCase) AddHolding(ctx context.Context, req *dto.AddHoldingRequest) (*dto.HoldingResponse, error) {
	// Validate request
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}
	
	// Verify portfolio exists
	_, err := uc.portfolioRepo.GetByID(ctx, req.PortfolioID)
	if err != nil {
		return nil, fmt.Errorf("portfolio not found: %w", err)
	}
	
	// Create holding
	holding := &entities.PortfolioHolding{
		PortfolioID:  req.PortfolioID,
		Symbol:       req.Symbol,
		Amount:       req.Amount,
		AveragePrice: req.AveragePrice,
	}
	
	if err := uc.portfolioRepo.AddHolding(ctx, req.PortfolioID, holding); err != nil {
		return nil, fmt.Errorf("failed to add holding: %w", err)
	}
	
	return dto.NewHoldingResponse(holding), nil
}

// GetPortfolioSummary retrieves portfolio summary with analytics
func (uc *PortfolioUseCase) GetPortfolioSummary(ctx context.Context, portfolioID uint) (*dto.PortfolioSummaryResponse, error) {
	// Get portfolio
	portfolio, err := uc.portfolioRepo.GetByID(ctx, portfolioID)
	if err != nil {
		return nil, fmt.Errorf("failed to get portfolio: %w", err)
	}
	
	// Get portfolio summary
	summary, err := uc.portfolioRepo.GetPortfolioSummary(ctx, portfolioID)
	if err != nil {
		return nil, fmt.Errorf("failed to get portfolio summary: %w", err)
	}
	
	// Calculate risk metrics
	riskMetrics, err := uc.riskAnalysisSvc.AnalyzePortfolioRisk(ctx, portfolio)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate risk metrics: %w", err)
	}
	
	summary.RiskMetrics = *riskMetrics
	
	return dto.NewPortfolioSummaryResponse(summary), nil
}

// UpdateHolding updates an existing holding
func (uc *PortfolioUseCase) UpdateHolding(ctx context.Context, req *dto.UpdateHoldingRequest) error {
	// Validate request
	if err := req.Validate(); err != nil {
		return fmt.Errorf("invalid request: %w", err)
	}
	
	// Update holding
	holding := &entities.PortfolioHolding{
		ID:           req.HoldingID,
		Amount:       req.Amount,
		AveragePrice: req.AveragePrice,
	}
	
	if err := uc.portfolioRepo.UpdateHolding(ctx, holding); err != nil {
		return fmt.Errorf("failed to update holding: %w", err)
	}
	
	return nil
}

// RemoveHolding removes a holding from a portfolio
func (uc *PortfolioUseCase) RemoveHolding(ctx context.Context, holdingID uint) error {
	if err := uc.portfolioRepo.RemoveHolding(ctx, holdingID); err != nil {
		return fmt.Errorf("failed to remove holding: %w", err)
	}
	
	return nil
}