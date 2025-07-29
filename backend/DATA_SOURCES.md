# Data Sources Documentation

This document provides comprehensive information about all external data sources used in the Cryptocurrency Indicator Dashboard.

## Overview

The system is designed to use **free data sources only**, avoiding expensive API subscriptions while maintaining high data quality through multi-source validation and consensus algorithms.

## Primary Data Sources

### 1. CoinCap API (rest.coincap.io/v3)

**Status**: ✅ Implemented and Active  
**Authentication**: API Key Required (Bearer Token)  
**Rate Limits**: Enhanced limits with API key  
**Cost**: Free with registration  

**Endpoints Used**:
```
GET /assets                    # List all assets
GET /assets/{id}              # Single asset data  
GET /assets/{id}/history      # Historical price data
GET /markets                  # Market data
```

**Data Provided**:
- Real-time cryptocurrency prices
- Market capitalization data
- 24-hour trading volume
- Asset rankings and metadata
- Historical price data (daily/hourly)
- Global market statistics

**Implementation**: `internal/infrastructure/external/coincap_client.go`

**Configuration**:
```go
// API Key setup
apiKey := "12529d2a4759537c59fa8bf185ad2783e674b612f2b509af6ce7278610bb299d"
client := external.NewCoinCapClient(apiKey, logger)
```

**Example Response**:
```json
{
  "data": {
    "id": "bitcoin",
    "rank": "1",
    "symbol": "BTC",
    "name": "Bitcoin",
    "supply": "19580000.0000000000000000",
    "maxSupply": "21000000.0000000000000000",
    "marketCapUsd": "2340000000000.0000000000000000",
    "volumeUsd24Hr": "15400000000.0000000000000000",
    "priceUsd": "117500.0000000000000000"
  }
}
```

### 2. CoinGecko API (api.coingecko.com)

**Status**: ✅ Implemented and Active  
**Authentication**: None Required (Free Tier)  
**Rate Limits**: 10-30 requests per minute  
**Cost**: Free  

**Endpoints Used**:
```
GET /simple/price             # Current prices
GET /global                   # Global market data
GET /coins/{id}/market_chart  # Historical data
```

**Data Provided**:
- Real-time cryptocurrency prices
- Market cap and volume data
- Global market statistics
- Historical price charts
- Market dominance data

**Implementation**: `services/data_fetcher.go`

**Example Usage**:
```go
fetcher := NewDataFetcher()
priceData, err := fetcher.FetchPriceData()
// Returns Bitcoin, Ethereum prices in USD
```

### 3. Blockchain.com API (blockstream.info)

**Status**: ✅ Implemented and Active  
**Authentication**: None Required  
**Rate Limits**: No strict limits  
**Cost**: Free  

**Endpoints Used**:
```
GET /stats?format=json        # Bitcoin network statistics
```

**Data Provided**:
- Bitcoin network hash rate
- Mining difficulty
- Total Bitcoin supply
- Transaction count and fees
- Mempool statistics
- Block height and timing

**Implementation**: `internal/infrastructure/external/blockchain_client.go`

**Example Response**:
```json
{
  "hash_rate": 450000000000000000000,
  "difficulty": 72032680272743.31,
  "timestamp": 1706447600,
  "n_btc_mined": 1968750000000,
  "n_tx": 948542,
  "totalbc": 1968750000000,
  "total_fees_btc": 45679230000
}
```

### 4. Alternative.me Fear & Greed Index

**Status**: ✅ Implemented and Active  
**Authentication**: None Required  
**Rate Limits**: Reasonable usage  
**Cost**: Free  

**Endpoints Used**:
```
GET /fng/                     # Fear & Greed Index
```

**Data Provided**:
- Current fear & greed score (0-100)
- Historical fear & greed data
- Sentiment classification
- Component breakdown

**Implementation**: `services/fear_greed_service.go`

## Consensus Pricing Algorithm

### Multi-Source Validation

The system combines data from multiple sources to create consensus pricing with confidence scoring:

```go
type AggregatedPrice struct {
    USD             float64           `json:"usd"`
    Sources         map[string]float64 `json:"sources"`
    Consensus       float64           `json:"consensus"`
    Deviation       float64           `json:"deviation"`
    ConfidenceLevel float64           `json:"confidence_level"`
}
```

### Confidence Calculation

```go
func calculateConfidenceLevel(deviation float64, sourceCount int) float64 {
    // Base confidence on number of sources
    baseConfidence := math.Min(float64(sourceCount)*25, 75)
    
    // Adjust based on price deviation
    deviationPenalty := math.Min(deviation/100*25, 25)
    
    confidence := baseConfidence - deviationPenalty
    return math.Max(0, math.Min(100, confidence))
}
```

### Data Quality Metrics

- **Source Agreement**: Standard deviation between price sources
- **Reliability Score**: Historical accuracy of each data source
- **Confidence Level**: Overall confidence in consensus price (0-100%)
- **Data Freshness**: Timestamp validation and staleness detection

## API Health Monitoring

### Health Check System

```go
func (mda *MarketDataAggregator) HealthCheck() map[string]error {
    results := make(map[string]error)
    
    // Check CoinCap
    results["coincap"] = mda.coinCapClient.HealthCheck()
    
    // Check CoinGecko  
    _, err := mda.coinGeckoFetcher.FetchPriceData()
    results["coingecko"] = err
    
    return results
}
```

### Fallback Strategies

1. **Primary Source Failure**: Automatically fallback to secondary sources
2. **Multiple Source Failure**: Use cached data with staleness warnings
3. **Complete Failure**: Return error with last known good data
4. **Graceful Degradation**: Continue operation with reduced confidence

## Rate Limiting & Error Handling

### Rate Limiting Strategy

```go
// Per-source rate limiting
type RateLimiter struct {
    requests map[string]int
    windows  map[string]time.Time
    limits   map[string]int
}
```

### Error Handling

- **Network Errors**: Retry with exponential backoff
- **Rate Limits**: Queue requests and implement delays
- **Invalid Data**: Validate and filter outliers
- **API Changes**: Version detection and adaptation

## Data Storage Pipeline

### Time-Series Storage

All external data is stored in PostgreSQL time-series tables:

```sql
-- Price data with source tracking
CREATE TABLE price_data_points (
    timestamp TIMESTAMPTZ NOT NULL,
    asset_symbol VARCHAR(10) NOT NULL,
    price_usd DECIMAL(20,8) NOT NULL,
    data_source VARCHAR(50) NOT NULL,
    reliability_score DECIMAL(5,2)
);

-- Network metrics
CREATE TABLE network_metric_points (
    timestamp TIMESTAMPTZ NOT NULL,
    network VARCHAR(20) NOT NULL,
    hash_rate DECIMAL(30,2),
    difficulty DECIMAL(30,2),
    data_source VARCHAR(50) NOT NULL
);
```

### Data Retention Policies

- **Price Data**: 3 years retention
- **Indicator Data**: 2 years retention  
- **Network Metrics**: 1 year retention
- **Rainbow Chart Data**: 5 years retention (historical analysis)

## Usage Examples

### Basic Price Fetching

```go
// Initialize aggregator
aggregator := services.NewMarketDataAggregator(logger, coinCapAPIKey)

// Get aggregated Bitcoin data
btcData, err := aggregator.GetAggregatedAssetData("BTC")
if err != nil {
    log.Printf("Error: %v", err)
    return
}

fmt.Printf("Consensus Price: $%.2f\n", btcData.Price.Consensus)
fmt.Printf("Confidence: %.1f%%\n", btcData.Price.ConfidenceLevel)
fmt.Printf("Sources: %v\n", btcData.DataSources)
```

### Network Metrics Retrieval

```go
// Initialize Blockchain.com client
client := external.NewBlockchainClient(logger)

// Get Bitcoin network stats
stats, err := client.GetBitcoinStats()
if err != nil {
    log.Printf("Error: %v", err)
    return
}

fmt.Printf("Hash Rate: %.2e H/s\n", stats.HashRate)
fmt.Printf("Difficulty: %.2f\n", stats.Difficulty)
fmt.Printf("Total Supply: %.8f BTC\n", stats.TotalSupply)
```

### Historical Data Access

```go
// Get historical price data
history, err := aggregator.GetHistoricalPriceData("BTC", 30) // Last 30 days
if err != nil {
    log.Printf("Error: %v", err)
    return
}

for _, point := range history {
    fmt.Printf("Date: %s, Price: $%.2f\n", 
        point["date"], point["price"])
}
```

## Performance Considerations

### Caching Strategy

- **Redis Cache**: 5-10 minute cache for aggregated data
- **Application Cache**: In-memory cache for frequently accessed data
- **Database Cache**: Query result caching at PostgreSQL level

### Optimization Techniques

- **Batch Requests**: Combine multiple API calls where possible
- **Parallel Processing**: Fetch from multiple sources simultaneously
- **Connection Pooling**: Reuse HTTP connections for efficiency
- **Data Compression**: Compress stored historical data

## Monitoring & Alerting

### Key Metrics

- **API Response Times**: Track latency for each data source
- **Error Rates**: Monitor failed requests and error types
- **Data Quality**: Track consensus confidence and source agreement
- **Cache Hit Rates**: Monitor caching effectiveness

### Alerting Conditions

- API downtime or high error rates
- Data quality degradation (low confidence scores)
- Price discrepancies exceeding thresholds
- Network connectivity issues

## Future Data Source Additions

### Planned Integrations

1. **FRED API**: US economic data (inflation, interest rates)
2. **Messari API**: Advanced on-chain metrics
3. **Glassnode**: Additional network and market metrics
4. **DeFiPulse**: DeFi protocol data and TVL metrics

### Integration Guidelines

1. **Free Tier First**: Always start with free tier evaluation
2. **Multi-Source**: Never rely on single source for critical data
3. **Consensus Integration**: Add to existing consensus algorithm
4. **Health Monitoring**: Implement comprehensive health checks
5. **Error Handling**: Graceful degradation and fallback strategies
6. **Documentation**: Complete API documentation and examples

## Troubleshooting

### Common Issues

1. **API Key Issues**: Verify CoinCap API key is correctly configured
2. **Rate Limiting**: Implement delays and respect API limits
3. **Network Timeouts**: Increase timeout values for slow connections
4. **Data Inconsistencies**: Check source reliability scores
5. **Cache Issues**: Verify Redis connection and cache expiration

### Debug Commands

```bash
# Test individual APIs
go run test_coincap.go
go run test_blockchain.go

# Test aggregated data
go run test_aggregator.go

# Monitor database
go run test_postgres_timeseries.go
```

This comprehensive data source documentation ensures reliable, free access to high-quality cryptocurrency market data while maintaining system resilience and data accuracy.