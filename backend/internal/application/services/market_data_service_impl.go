package services

import (
	"context"
	"fmt"
	"time"
	"crypto-indicator-dashboard/internal/domain/entities"
	"crypto-indicator-dashboard/internal/domain/repositories"
	"crypto-indicator-dashboard/internal/domain/services"
	"crypto-indicator-dashboard/internal/infrastructure/external"
	"crypto-indicator-dashboard/pkg/logger"
)

// marketDataServiceImpl implements the MarketDataService interface
type marketDataServiceImpl struct {
	repo              repositories.MarketDataRepository
	coinMarketCapClient *external.CoinMarketCapClient
	tradingViewScraper  *external.TradingViewScraper
	cacheService      services.CacheService
	logger            logger.Logger
}

// NewMarketDataService creates a new market data service implementation
func NewMarketDataService(
	repo repositories.MarketDataRepository,
	coinMarketCapClient *external.CoinMarketCapClient,
	tradingViewScraper *external.TradingViewScraper,
	cacheService services.CacheService,
	logger logger.Logger,
) services.MarketDataService {
	return &marketDataServiceImpl{
		repo:                repo,
		coinMarketCapClient: coinMarketCapClient,
		tradingViewScraper:  tradingViewScraper,
		cacheService:        cacheService,
		logger:              logger,
	}
}

// GetCryptoPrices retrieves current cryptocurrency prices from CoinMarketCap
func (s *marketDataServiceImpl) GetCryptoPrices(ctx context.Context, symbols []string) (map[string]*entities.CryptoPrice, error) {
	cacheKey := fmt.Sprintf("crypto_prices_%v", symbols)
	
	// Try to get from cache first
	var cachedPrices map[string]*entities.CryptoPrice
	if err := s.cacheService.GetOrSet(ctx, cacheKey, &cachedPrices, 2*time.Minute, func() (interface{}, error) {
		return s.fetchCryptoPricesFromAPI(ctx, symbols)
	}); err != nil {
		s.logger.Error("Failed to get crypto prices from cache", "error", err, "symbols", symbols)
		// Fallback to direct API call
		return s.fetchCryptoPricesFromAPI(ctx, symbols)
	}
	
	return cachedPrices, nil
}

// fetchCryptoPricesFromAPI fetches prices directly from CoinMarketCap API
func (s *marketDataServiceImpl) fetchCryptoPricesFromAPI(ctx context.Context, symbols []string) (map[string]*entities.CryptoPrice, error) {
	s.logger.Info("Fetching crypto prices from CoinMarketCap API", "symbols", symbols)
	
	response, err := s.coinMarketCapClient.GetLatestQuotes(symbols, "USD")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch quotes from CoinMarketCap: %w", err)
	}
	
	prices := make(map[string]*entities.CryptoPrice)
	for symbol, data := range response.Data {
		if usdQuote, exists := data.Quote["USD"]; exists {
			price := &entities.CryptoPrice{
				Symbol:           symbol,
				Name:             data.Name,
				Price:            usdQuote.Price,
				Volume24h:        usdQuote.Volume24h,
				MarketCap:        usdQuote.MarketCap,
				PercentChange1h:  usdQuote.PercentChange1h,
				PercentChange24h: usdQuote.PercentChange24h,
				PercentChange7d:  usdQuote.PercentChange7d,
				PercentChange30d: usdQuote.PercentChange30d,
				LastUpdated:      usdQuote.LastUpdated,
				DataSource:       "CoinMarketCap",
			}
			prices[symbol] = price
			
			// Store in database for historical tracking
			if err := s.repo.StorePriceData(ctx, price); err != nil {
				s.logger.Warn("Failed to store price data", "error", err, "symbol", symbol)
			}
		}
	}
	
	s.logger.Info("Successfully fetched crypto prices", "count", len(prices), "symbols", symbols)
	return prices, nil
}

// GetBitcoinDominance retrieves Bitcoin dominance from multiple sources
func (s *marketDataServiceImpl) GetBitcoinDominance(ctx context.Context) (*entities.BitcoinDominance, error) {
	cacheKey := "bitcoin_dominance"
	
	// Try to get from cache first
	var cachedDominance *entities.BitcoinDominance
	if err := s.cacheService.GetOrSet(ctx, cacheKey, &cachedDominance, 5*time.Minute, func() (interface{}, error) {
		return s.fetchBitcoinDominanceFromSources(ctx)
	}); err != nil {
		s.logger.Error("Failed to get Bitcoin dominance from cache", "error", err)
		// Fallback to direct fetch
		return s.fetchBitcoinDominanceFromSources(ctx)
	}
	
	return cachedDominance, nil
}

// fetchBitcoinDominanceFromSources fetches Bitcoin dominance from multiple sources
func (s *marketDataServiceImpl) fetchBitcoinDominanceFromSources(ctx context.Context) (*entities.BitcoinDominance, error) {
	s.logger.Info("Fetching Bitcoin dominance from multiple sources")
	
	var primaryDominance, secondaryDominance float64
	var primarySource, secondarySource string
	var primaryErr, secondaryErr error
	
	// Try CoinMarketCap first
	primaryDominance, primaryErr = s.coinMarketCapClient.GetBitcoinDominance()
	if primaryErr == nil {
		primarySource = "CoinMarketCap"
		s.logger.Info("Got Bitcoin dominance from CoinMarketCap", "dominance", primaryDominance)
	}
	
	// Try TradingView as secondary source
	tvData, secondaryErr := s.tradingViewScraper.GetBitcoinDominanceWithFallback()
	if secondaryErr == nil {
		secondaryDominance = tvData.CurrentDominance
		secondarySource = "TradingView"
		s.logger.Info("Got Bitcoin dominance from TradingView", "dominance", secondaryDominance)
	}
	
	// Determine which source to use
	var finalDominance float64
	var finalSource string
	var confidence float64 = 1.0
	
	if primaryErr == nil && secondaryErr == nil {
		// Both sources available - compare and use average if close
		diff := abs(primaryDominance - secondaryDominance)
		if diff < 2.0 { // If difference is less than 2%, average them
			finalDominance = (primaryDominance + secondaryDominance) / 2
			finalSource = "CoinMarketCap + TradingView (averaged)"
			confidence = 0.95
			s.logger.Info("Using averaged Bitcoin dominance", 
				"cmc_dominance", primaryDominance,
				"tv_dominance", secondaryDominance,
				"final_dominance", finalDominance)
		} else {
			// Large difference, prefer CoinMarketCap
			finalDominance = primaryDominance
			finalSource = primarySource
			confidence = 0.8
			s.logger.Warn("Large difference between dominance sources", 
				"cmc_dominance", primaryDominance,
				"tv_dominance", secondaryDominance,
				"using", finalSource)
		}
	} else if primaryErr == nil {
		finalDominance = primaryDominance
		finalSource = primarySource
		confidence = 0.9
	} else if secondaryErr == nil {
		finalDominance = secondaryDominance
		finalSource = secondarySource
		confidence = 0.85
	} else {
		return nil, fmt.Errorf("failed to fetch Bitcoin dominance from any source: cmc_error=%v, tv_error=%v", primaryErr, secondaryErr)
	}
	
	// Create dominance entity
	dominance := &entities.BitcoinDominance{
		CurrentDominance:    finalDominance,
		PreviousDominance:   0, // Would need historical data
		Change24h:          0,  // Would need historical data
		ChangePercent24h:   0,  // Would need historical data
		LastUpdated:        time.Now(),
		DataSource:         finalSource,
		Confidence:         confidence,
	}
	
	// If we have TradingView data with change information, use it
	if secondaryErr == nil && tvData.ChangePercent24h != 0 {
		dominance.ChangePercent24h = tvData.ChangePercent24h
		dominance.Change24h = tvData.Change24h
		dominance.PreviousDominance = tvData.PreviousDominance
	}
	
	// Store in database for historical tracking
	if err := s.repo.StoreDominanceData(ctx, dominance); err != nil {
		s.logger.Warn("Failed to store dominance data", "error", err)
	}
	
	s.logger.Info("Successfully determined Bitcoin dominance", 
		"dominance", finalDominance,
		"source", finalSource,
		"confidence", confidence)
	
	return dominance, nil
}

// GetMultipleCryptoPrices is a convenience method for getting common crypto prices
func (s *marketDataServiceImpl) GetMultipleCryptoPrices(ctx context.Context) (map[string]*entities.CryptoPrice, error) {
	commonSymbols := []string{"BTC", "ETH", "BNB", "SOL", "ADA", "XRP", "DOT", "AVAX", "MATIC", "LINK"}
	return s.GetCryptoPrices(ctx, commonSymbols)
}

// GetTopCryptoPrices gets prices for top N cryptocurrencies by market cap
func (s *marketDataServiceImpl) GetTopCryptoPrices(ctx context.Context, count int) (map[string]*entities.CryptoPrice, error) {
	// This would require a different CoinMarketCap endpoint for top coins by market cap
	// For now, return common major cryptocurrencies
	symbols := []string{"BTC", "ETH", "BNB", "SOL", "ADA", "XRP", "DOT", "AVAX", "MATIC", "LINK"}
	if count < len(symbols) {
		symbols = symbols[:count]
	}
	return s.GetCryptoPrices(ctx, symbols)
}

// RefreshAllMarketData refreshes all market data from external sources
func (s *marketDataServiceImpl) RefreshAllMarketData(ctx context.Context) error {
	s.logger.Info("Refreshing all market data")
	
	// Refresh crypto prices
	_, err := s.GetMultipleCryptoPrices(ctx)
	if err != nil {
		s.logger.Error("Failed to refresh crypto prices", "error", err)
		return fmt.Errorf("failed to refresh crypto prices: %w", err)
	}
	
	// Refresh Bitcoin dominance
	_, err = s.GetBitcoinDominance(ctx)
	if err != nil {
		s.logger.Error("Failed to refresh Bitcoin dominance", "error", err)
		return fmt.Errorf("failed to refresh Bitcoin dominance: %w", err)
	}
	
	s.logger.Info("Successfully refreshed all market data")
	return nil
}

// HealthCheck performs health checks on all external data sources
func (s *marketDataServiceImpl) HealthCheck(ctx context.Context) map[string]error {
	results := make(map[string]error)
	
	// Check CoinMarketCap
	if err := s.coinMarketCapClient.HealthCheck(); err != nil {
		results["coinmarketcap"] = err
	} else {
		results["coinmarketcap"] = nil
	}
	
	// Check TradingView scraper
	if err := s.tradingViewScraper.HealthCheck(); err != nil {
		results["tradingview"] = err
	} else {
		results["tradingview"] = nil
	}
	
	return results
}

// Helper function to calculate absolute value
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}