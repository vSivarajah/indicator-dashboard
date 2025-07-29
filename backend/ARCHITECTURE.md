# Crypto Indicator Dashboard - Backend Architecture

## Overview

This document describes the clean architecture implementation of the Crypto Indicator Dashboard backend, built using Go with the Gin web framework. The architecture follows Domain-Driven Design (DDD) principles and implements a layered approach for maintainability, testability, and scalability.

## Architecture Layers

```
┌─────────────────────────────────────────────────────────────┐
│                    Presentation Layer                       │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐  │
│  │   Handlers  │  │ Middleware  │  │    HTTP Routes      │  │
│  └─────────────┘  └─────────────┘  └─────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                   Application Layer                         │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐  │
│  │  Use Cases  │  │     DTOs    │  │   Service Impls     │  │
│  └─────────────┘  └─────────────┘  └─────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                     Domain Layer                            │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐  │
│  │  Entities   │  │ Interfaces  │  │  Domain Services    │  │
│  └─────────────┘  └─────────────┘  └─────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                 Infrastructure Layer                        │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐  │
│  │  Database   │  │    Cache    │  │  External APIs      │  │
│  └─────────────┘  └─────────────┘  └─────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

## Directory Structure

```
backend/
├── cmd/
│   └── server/
│       └── main.go                    # Application entry point
├── internal/
│   ├── application/                   # Application Layer
│   │   ├── dto/                       # Data Transfer Objects
│   │   ├── services/                  # Service Implementations
│   │   └── usecases/                  # Use Case Orchestration
│   ├── domain/                        # Domain Layer
│   │   ├── entities/                  # Core Business Entities
│   │   ├── repositories/              # Repository Interfaces
│   │   └── services/                  # Domain Service Interfaces
│   ├── infrastructure/                # Infrastructure Layer
│   │   ├── cache/                     # Cache Implementations
│   │   ├── config/                    # Configuration & DI
│   │   └── database/                  # Database Implementations
│   └── presentation/                  # Presentation Layer
│       ├── handlers/                  # HTTP Request Handlers
│       └── middleware/                # HTTP Middleware
├── pkg/                              # Shared Packages
│   ├── errors/                       # Error Handling
│   └── logger/                       # Logging Utilities
├── models/                           # Legacy Database Models
└── services/                         # Service Layer (Market Analysis & Time-Series)
```

## Layer Responsibilities

### 1. Presentation Layer (`internal/presentation/`)

**Purpose**: Handles HTTP requests and responses, authentication, validation, and API documentation.

**Components**:
- **Handlers**: Process HTTP requests and delegate to use cases
- **Middleware**: Cross-cutting concerns (logging, CORS, rate limiting, authentication)
- **Routes**: URL routing configuration

**Key Files**:
- `handlers/indicator_handler.go` - Market indicator endpoints
- `handlers/portfolio_handler.go` - Portfolio management endpoints
- `middleware/cors.go` - CORS configuration
- `middleware/logging.go` - Request/response logging
- `middleware/rate_limiting.go` - Rate limiting implementation

**Dependencies**: Can only depend on Application layer interfaces

### 2. Application Layer (`internal/application/`)

**Purpose**: Orchestrates business logic, implements use cases, and coordinates between different domain services.

**Components**:
- **Use Cases**: High-level business workflows
- **DTOs**: Data transfer objects for API communication
- **Service Implementations**: Concrete implementations of domain service interfaces

**Key Files**:
- `usecases/indicator_usecase.go` - Indicator analysis workflows
- `usecases/portfolio_usecase.go` - Portfolio management workflows
- `services/mvrv_service_impl.go` - MVRV calculation implementation
- `dto/indicator_dto.go` - Indicator data transfer objects

**Dependencies**: Can depend on Domain layer interfaces and Infrastructure layer for implementation

### 3. Domain Layer (`internal/domain/`)

**Purpose**: Contains core business logic, entities, and defines contracts for external dependencies.

**Components**:
- **Entities**: Core business objects with behavior
- **Repository Interfaces**: Data access contracts
- **Service Interfaces**: Business logic contracts
- **Value Objects**: Immutable data objects

**Key Files**:
- `entities/indicator.go` - Market indicator entities
- `entities/portfolio.go` - Portfolio and holdings entities
- `entities/dca.go` - Dollar Cost Averaging entities
- `repositories/indicator_repository.go` - Indicator data access interface
- `services/indicator_service.go` - Indicator calculation interfaces

**Dependencies**: No dependencies on other layers (pure business logic)

### 4. Infrastructure Layer (`internal/infrastructure/`)

**Purpose**: Implements external concerns like database access, caching, external APIs, and configuration.

**Components**:
- **Database**: Repository implementations using GORM with PostgreSQL/TimescaleDB
- **Time-Series Storage**: Optimized time-series data management and retention
- **Cache**: Redis and in-memory cache implementations
- **Configuration**: Dependency injection and configuration management
- **External Services**: Multi-source API clients with consensus algorithms

**Key Files**:
- `database/indicator_repository.go` - Database implementation for indicators
- `database/dca_repository.go` - Database implementation for DCA strategies
- `database/timescale_setup.go` - TimescaleDB hypertables and optimization
- `database/postgres_timeseries_setup.go` - PostgreSQL time-series tables and indexes
- `cache/redis_cache.go` - Redis and mock cache implementations
- `config/dependencies.go` - Dependency injection container
- `external/coincap_client.go` - CoinCap API integration with authentication
- `external/blockchain_client.go` - Blockchain.com network metrics client

**Dependencies**: Implements Domain layer interfaces

## Design Patterns Used

### 1. Dependency Injection

All dependencies are injected through constructor functions and interfaces:

```go
// Service constructor with dependencies
func NewMVRVService(
    indicatorRepo repositories.IndicatorRepository,
    marketDataRepo repositories.MarketDataRepository,
    cache cache.CacheService,
    logger logger.Logger,
) services.IndicatorService {
    return &mvrvServiceImpl{
        indicatorRepo:  indicatorRepo,
        marketDataRepo: marketDataRepo,
        cache:          cache,
        logger:         logger,
    }
}
```

### 2. Repository Pattern

Data access is abstracted through repository interfaces:

```go
type IndicatorRepository interface {
    Create(ctx context.Context, indicator *entities.Indicator) error
    GetByID(ctx context.Context, id uint) (*entities.Indicator, error)
    GetHistoricalData(ctx context.Context, name string, from, to time.Time) ([]entities.Indicator, error)
    // ... other methods
}
```

### 3. Service Layer Pattern

Business logic is encapsulated in service interfaces:

```go
type IndicatorService interface {
    Calculate(ctx context.Context, params map[string]interface{}) (*entities.Indicator, error)
    GetHistoricalData(ctx context.Context, period string) ([]entities.Indicator, error)
    GetLatest(ctx context.Context) (*entities.Indicator, error)
}
```

### 4. Factory Pattern

Services and repositories are created through factory functions:

```go
func NewDependencies(config *Config) (*Dependencies, error) {
    deps := &Dependencies{Config: config}
    deps.initDatabase()
    deps.initCache()
    deps.initRepositories()
    deps.initDomainServices()
    return deps, nil
}
```

## Data Flow

### Request Processing Flow

```
1. HTTP Request → Handler
2. Handler → Use Case
3. Use Case → Domain Service
4. Domain Service → Repository/External API
5. Repository → Database/Cache
6. Response flows back through the same layers
```

### Example: MVRV Indicator Calculation

```
┌─────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Client    │───▶│ IndicatorHandler │───▶│ MVRVServiceImpl │
└─────────────┘    └─────────────────┘    └─────────────────┘
                            │                       │
                            ▼                       ▼
┌─────────────┐    ┌─────────────────┐    ┌─────────────────┐
│  Database   │◀───│IndicatorRepo    │◀───│  Cache Service  │
└─────────────┘    └─────────────────┘    └─────────────────┘
```

## Key Architectural Benefits

### 1. **Testability**
- Each layer can be tested in isolation
- Dependencies are injected as interfaces
- Mock implementations for testing

### 2. **Maintainability**
- Clear separation of concerns
- Single responsibility principle
- Loose coupling between layers

### 3. **Scalability**
- Easy to add new features
- Horizontal scaling through stateless design
- Caching layer for performance

### 4. **Flexibility**
- Easy to swap implementations (database, cache, external APIs)
- Support for multiple data sources
- Graceful degradation when services are unavailable

## Configuration Management

### Environment Configuration

```go
type Config struct {
    Server   ServerConfig
    Database DatabaseConfig
    Redis    RedisConfig
    External ExternalConfig
}
```

### Dependency Injection Container

```go
type Dependencies struct {
    // Infrastructure
    DB     *gorm.DB
    Redis  *redis.Client
    Cache  cache.CacheService
    Logger logger.Logger
    
    // Repositories
    PortfolioRepo  repositories.PortfolioRepository
    IndicatorRepo  repositories.IndicatorRepository
    
    // Domain Services
    PortfolioService  services.PortfolioService
    IndicatorService  services.IndicatorService
}
```

## Error Handling Strategy

### Structured Error Types

```go
type AppError struct {
    Type       ErrorType `json:"type"`
    Message    string    `json:"message"`
    Details    string    `json:"details,omitempty"`
    StatusCode int       `json:"-"`
    Cause      error     `json:"-"`
}
```

### Error Flow

1. Domain layer returns specific error types
2. Application layer wraps with context
3. Presentation layer converts to HTTP responses
4. All errors are logged with context

## Logging Strategy

### Structured Logging

```go
logger.Info("Processing MVRV calculation",
    "user_id", userID,
    "symbol", "BTC",
    "cache_hit", cacheHit,
)
```

### Log Levels
- **DEBUG**: Detailed tracing information
- **INFO**: General operational information
- **WARN**: Warning conditions
- **ERROR**: Error conditions requiring attention

## Caching Strategy

### Multi-Layer Caching

1. **Application Cache**: In-memory cache for frequently accessed data
2. **Redis Cache**: Distributed cache for session data and computed results
3. **Database Cache**: Query result caching at the database level

### Cache Keys

```
Pattern: {service}:{operation}:{identifier}
Example: mvrv:calculation:btc_latest
```

## Database Design

### Transaction Management

```go
func (tm *TransactionManager) WithTransaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
    return tm.db.WithContext(ctx).Transaction(fn)
}
```

### Connection Pooling

- Maximum open connections: Configurable
- Connection lifetime: 1 hour
- Idle timeout: 10 minutes

## Security Considerations

### Input Validation
- Request validation middleware
- Parameter sanitization
- SQL injection prevention through ORM

### Rate Limiting
- Per-IP rate limiting
- Configurable limits per endpoint
- Redis-backed rate limiting

### CORS Configuration
- Configurable allowed origins
- Preflight request handling
- Secure headers

## Performance Optimizations

### Database
- Connection pooling
- Query optimization
- Proper indexing
- Batch operations for bulk inserts

### Caching
- Redis for distributed caching
- In-memory cache for hot data
- Cache invalidation strategies

### HTTP
- Response compression
- Keep-alive connections
- Timeouts and circuit breakers

## Monitoring and Observability

### Metrics
- Request/response times
- Error rates
- Cache hit rates
- Database connection metrics

### Health Checks
- Database connectivity
- Redis connectivity
- External API availability

### Logging
- Structured JSON logs
- Request tracing
- Error tracking with stack traces

## Time-Series Data Architecture

### TimescaleDB Integration

**Hypertable Management**:
- Automatic table partitioning by time intervals
- Optimized chunk management for query performance
- Data retention policies with automated cleanup
- Real-time compression and optimization

**Table Structure**:
```sql
-- Price data with multi-source validation
price_data_points (
    timestamp, asset_symbol, price_usd, market_cap, 
    volume_24h, data_source, reliability_score
)

-- Indicator calculations with metadata
indicator_data_points (
    timestamp, indicator_type, value, metadata, 
    confidence_level, data_source
)

-- Network metrics for blockchain analysis
network_metric_points (
    timestamp, network, hash_rate, difficulty, 
    block_height, total_supply, transaction_count
)

-- Rainbow chart analysis data
rainbow_chart_data_points (
    timestamp, bitcoin_price, log_regression_price,
    current_band, cycle_position, risk_level
)
```

**Performance Optimizations**:
- Time-based indexing strategies
- Batch insertion for high-throughput writes
- Query optimization for time-range queries
- Automatic table statistics updates

### Multi-Source Data Pipeline

**Consensus Algorithm**:
```go
// Price consensus calculation with confidence scoring
func calculateConsensus(sources map[string]float64) (float64, float64) {
    mean := calculateMean(sources)
    deviation := calculateStandardDeviation(sources, mean)
    confidence := calculateConfidenceLevel(deviation, len(sources))
    return mean, confidence
}
```

**API Integration Strategy**:
- CoinCap API: Primary source with authentication
- CoinGecko API: Secondary source for validation
- Blockchain.com: Network metrics without authentication
- Health monitoring with automatic failover
- Rate limiting and error handling per source

**Data Quality Assurance**:
- Multi-source validation and outlier detection
- Confidence scoring based on source agreement
- Historical data validation and backfilling
- Real-time monitoring of data quality metrics

## Current Architecture Status

### Implemented Components ✅
1. **Time-Series Database**: PostgreSQL with TimescaleDB support
2. **Multi-Source APIs**: CoinCap, CoinGecko, Blockchain.com integration
3. **Consensus Pricing**: Real-time price validation and confidence scoring
4. **Rainbow Chart Analysis**: Complete mathematical implementation
5. **Network Metrics**: Bitcoin blockchain statistics tracking
6. **Data Retention**: Automated cleanup with configurable policies
7. **Performance Optimization**: Indexes, caching, and query optimization

### In Progress 🚧
1. **Background Job System**: Scheduled data collection and processing
2. **API Rate Management**: Smart rate limiting and request optimization
3. **ML Data Pipeline**: Feature engineering and model training preparation

### Planned Features 📋
1. **WebSocket**: Real-time indicator updates
2. **Advanced Analytics**: Machine learning model integration
3. **Microservices**: Split into indicator, portfolio, and analytics services
4. **Event Sourcing**: Track state changes for audit trails
5. **GraphQL**: Alternative API for flexible data fetching
6. **Authentication**: JWT-based user authentication

### Technical Improvements
1. **Circuit Breakers**: For external API resilience
2. **Distributed Tracing**: Request tracing across services
3. **Metrics Collection**: Prometheus/Grafana integration
4. **Container Orchestration**: Kubernetes deployment
5. **Load Testing**: Performance benchmarking

## Development Guidelines

### Adding New Features

1. **Start with Domain**: Define entities and interfaces
2. **Implement Application Layer**: Create use cases and DTOs
3. **Add Infrastructure**: Implement repositories and external services
4. **Create Handlers**: Add HTTP endpoints
5. **Add Tests**: Unit and integration tests
6. **Update Documentation**: Keep architecture docs current

### Code Standards

- Follow Go naming conventions
- Use dependency injection
- Write comprehensive tests
- Document public interfaces
- Handle errors gracefully
- Use structured logging

### Testing Strategy

- **Unit Tests**: Test individual components in isolation
- **Integration Tests**: Test layer interactions
- **End-to-End Tests**: Test complete workflows
- **Performance Tests**: Test under load
- **Contract Tests**: Test interface compliance

This architecture provides a solid foundation for building scalable, maintainable, and testable cryptocurrency market analysis applications.