package usecases

import (
	"context"
	"fmt"
	"crypto-indicator-dashboard/internal/domain/services"
	"crypto-indicator-dashboard/internal/application/dto"
)

// IndicatorUseCase handles indicator-related business logic
type IndicatorUseCase struct {
	mvrvSvc        services.MVRVService
	dominanceSvc   services.DominanceService
	fearGreedSvc   services.FearGreedService
	bubbleRiskSvc  services.BubbleRiskService
	macroSvc       services.MacroService
	marketCycleSvc services.MarketCycleService
}

// NewIndicatorUseCase creates a new indicator use case
func NewIndicatorUseCase(
	mvrvSvc services.MVRVService,
	dominanceSvc services.DominanceService,
	fearGreedSvc services.FearGreedService,
	bubbleRiskSvc services.BubbleRiskService,
	macroSvc services.MacroService,
	marketCycleSvc services.MarketCycleService,
) *IndicatorUseCase {
	return &IndicatorUseCase{
		mvrvSvc:        mvrvSvc,
		dominanceSvc:   dominanceSvc,
		fearGreedSvc:   fearGreedSvc,
		bubbleRiskSvc:  bubbleRiskSvc,
		macroSvc:       macroSvc,
		marketCycleSvc: marketCycleSvc,
	}
}

// GetMVRVIndicator retrieves MVRV Z-Score analysis
func (uc *IndicatorUseCase) GetMVRVIndicator(ctx context.Context) (*dto.MVRVResponse, error) {
	result, err := uc.mvrvSvc.GetMVRVZScore(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get MVRV data: %w", err)
	}
	
	return dto.NewMVRVResponse(result), nil
}

// GetDominanceIndicator retrieves Bitcoin dominance analysis
func (uc *IndicatorUseCase) GetDominanceIndicator(ctx context.Context) (*dto.DominanceResponse, error) {
	result, err := uc.dominanceSvc.GetDominanceAnalysis(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get dominance data: %w", err)
	}
	
	return dto.NewDominanceResponse(result), nil
}

// GetFearGreedIndicator retrieves Fear & Greed index analysis
func (uc *IndicatorUseCase) GetFearGreedIndicator(ctx context.Context) (*dto.FearGreedResponse, error) {
	result, err := uc.fearGreedSvc.GetFearGreedAnalysis(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get Fear & Greed data: %w", err)
	}
	
	return dto.NewFearGreedResponse(result), nil
}

// GetBubbleRiskIndicator retrieves bubble risk analysis
func (uc *IndicatorUseCase) GetBubbleRiskIndicator(ctx context.Context) (*dto.BubbleRiskResponse, error) {
	result, err := uc.bubbleRiskSvc.GetBubbleRiskAnalysis(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get bubble risk data: %w", err)
	}
	
	return dto.NewBubbleRiskResponse(result), nil
}

// GetInflationIndicator retrieves inflation analysis
func (uc *IndicatorUseCase) GetInflationIndicator(ctx context.Context) (*dto.InflationResponse, error) {
	result, err := uc.macroSvc.GetInflationAnalysis(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get inflation data: %w", err)
	}
	
	return dto.NewInflationResponse(result), nil
}

// GetInterestRateIndicator retrieves interest rate analysis
func (uc *IndicatorUseCase) GetInterestRateIndicator(ctx context.Context) (*dto.InterestRateResponse, error) {
	result, err := uc.macroSvc.GetInterestRateAnalysis(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get interest rate data: %w", err)
	}
	
	return dto.NewInterestRateResponse(result), nil
}

// GetMarketCycle retrieves current market cycle analysis
func (uc *IndicatorUseCase) GetMarketCycle(ctx context.Context) (*dto.MarketCycleResponse, error) {
	result, err := uc.marketCycleSvc.GetCurrentCycle(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get market cycle data: %w", err)
	}
	
	return dto.NewMarketCycleResponse(result), nil
}

// GetChartData retrieves chart data for a specific indicator
func (uc *IndicatorUseCase) GetChartData(ctx context.Context, indicator string) (*dto.ChartDataResponse, error) {
	var chartData map[string]interface{}
	var err error
	
	switch indicator {
	case "mvrv":
		chartData, err = uc.mvrvSvc.GetHistoricalMVRVChart(ctx)
	case "dominance":
		chartData, err = uc.dominanceSvc.GetDominanceChart(ctx)
	case "fear-greed":
		chartData, err = uc.fearGreedSvc.GetFearGreedChart(ctx)
	case "bubble-risk":
		chartData, err = uc.bubbleRiskSvc.GetBubbleRiskChart(ctx)
	default:
		return nil, fmt.Errorf("unsupported indicator: %s", indicator)
	}
	
	if err != nil {
		return nil, fmt.Errorf("failed to get chart data for %s: %w", indicator, err)
	}
	
	return dto.NewChartDataResponse(indicator, chartData), nil
}