package config

import (
	"context"
	"crypto-indicator-dashboard/internal/application/services"
	"crypto-indicator-dashboard/internal/application/usecases"
	"crypto-indicator-dashboard/internal/domain/repositories"
	domainServices "crypto-indicator-dashboard/internal/domain/services"
	"crypto-indicator-dashboard/internal/infrastructure/cache"
	"crypto-indicator-dashboard/internal/infrastructure/database"
	"crypto-indicator-dashboard/internal/infrastructure/external"
	"crypto-indicator-dashboard/pkg/logger"

	"github.com/go-redis/redis/v8"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Dependencies holds all application dependencies
type Dependencies struct {
	// Configuration
	Config *Config

	// Infrastructure
	DB     *gorm.DB
	Redis  *redis.Client
	Logger logger.Logger
	Cache  domainServices.CacheService

	// Repositories
	PortfolioRepo  repositories.PortfolioRepository
	IndicatorRepo  repositories.IndicatorRepository
	MarketDataRepo repositories.MarketDataRepository
	DCARepo        repositories.DCARepository

	// Domain Services
	PortfolioService  domainServices.PortfolioService
	IndicatorService  domainServices.IndicatorService
	DCAService        domainServices.DCAService
	MarketDataService domainServices.MarketDataService

	// External API Clients
	CoinMarketCapClient *external.CoinMarketCapClient
	TradingViewScraper  *external.TradingViewScraper

	// Use Cases
	PortfolioUseCase *usecases.PortfolioUseCase
	IndicatorUseCase *usecases.IndicatorUseCase
}

// NewDependencies creates and wires up all application dependencies
func NewDependencies(config *Config) (*Dependencies, error) {
	deps := &Dependencies{
		Config: config,
	}

	// Initialize logger
	deps.Logger = logger.New(config.Server.Environment)

	// Initialize database
	if err := deps.initDatabase(); err != nil {
		deps.Logger.Error("Failed to initialize database", "error", err)
		// Continue without database for graceful degradation
	}

	// Initialize Redis
	if err := deps.initRedis(); err != nil {
		deps.Logger.Error("Failed to initialize Redis", "error", err)
		// Continue without Redis for graceful degradation
	}

	// Initialize external clients
	deps.initExternalClients()

	// Initialize cache
	deps.initCache()

	// Initialize repositories
	deps.initRepositories()

	// Initialize domain services
	deps.initDomainServices()

	// Initialize use cases
	deps.initUseCases()

	return deps, nil
}

// initDatabase initializes the database connection
func (d *Dependencies) initDatabase() error {
	db, err := gorm.Open(postgres.Open(d.Config.Database.GetDSN()), &gorm.Config{
		Logger: logger.NewGormLogger(d.Logger),
	})
	if err != nil {
		return err
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	sqlDB.SetMaxOpenConns(d.Config.Database.MaxConns)
	sqlDB.SetMaxIdleConns(d.Config.Database.MinConns)

	d.DB = db
	return nil
}

// initRedis initializes the Redis connection
func (d *Dependencies) initRedis() error {
	rdb := redis.NewClient(&redis.Options{
		Addr:     d.Config.Redis.GetRedisAddr(),
		Password: d.Config.Redis.Password,
		DB:       d.Config.Redis.DB,
	})

	// Test connection
	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		return err
	}

	d.Redis = rdb
	return nil
}

// initExternalClients initializes external API clients
func (d *Dependencies) initExternalClients() {
	// Initialize CoinMarketCap client
	if d.Config.External.CoinMarketCapAPIKey != "" {
		d.CoinMarketCapClient = external.NewCoinMarketCapClient(
			d.Config.External.CoinMarketCapAPIKey,
			d.Logger,
		)
	}

	// Initialize TradingView scraper
	d.TradingViewScraper = external.NewTradingViewScraper(d.Logger)
}

// initCache initializes the cache service
func (d *Dependencies) initCache() {
	// Create a Redis cache service if available, otherwise use nil (will use fallback)
	var redisCache domainServices.CacheService
	if d.Redis != nil {
		// For now, we'll use nil for Redis and rely on fallback
		redisCache = nil
	}

	// Use our cache service implementation with fallback
	d.Cache = cache.NewCacheService(redisCache, d.Logger)
}

// initRepositories initializes all repositories
func (d *Dependencies) initRepositories() {
	if d.DB != nil {
		d.PortfolioRepo = database.NewPortfolioRepository(d.DB)
		d.IndicatorRepo = database.NewIndicatorRepository(d.DB, d.Logger)
		d.MarketDataRepo = database.NewMarketDataRepository(d.DB, d.Logger)
		d.DCARepo = database.NewDCARepository(d.DB, d.Logger)
	}
}

// initDomainServices initializes domain services
func (d *Dependencies) initDomainServices() {
	// Initialize market data service
	if d.MarketDataRepo != nil && d.CoinMarketCapClient != nil && d.TradingViewScraper != nil {
		d.MarketDataService = services.NewMarketDataService(
			d.MarketDataRepo,
			d.CoinMarketCapClient,
			d.TradingViewScraper,
			d.Cache,
			d.Logger,
		)
	}
}

// initUseCases initializes use cases
func (d *Dependencies) initUseCases() {
	// Note: These will be properly initialized once domain services are migrated
}

// Cleanup gracefully closes all connections
func (d *Dependencies) Cleanup() error {
	if d.Redis != nil {
		if err := d.Redis.Close(); err != nil {
			d.Logger.Error("Failed to close Redis connection", "error", err)
		}
	}

	if d.DB != nil {
		sqlDB, err := d.DB.DB()
		if err == nil {
			if err := sqlDB.Close(); err != nil {
				d.Logger.Error("Failed to close database connection", "error", err)
			}
		}
	}

	return nil
}

// TransactionManager provides database transaction management
type TransactionManager struct {
	db *gorm.DB
}

// NewTransactionManager creates a new transaction manager
func NewTransactionManager(db *gorm.DB) *TransactionManager {
	return &TransactionManager{db: db}
}

// WithTransaction executes a function within a database transaction
func (tm *TransactionManager) WithTransaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return tm.db.WithContext(ctx).Transaction(fn)
}
