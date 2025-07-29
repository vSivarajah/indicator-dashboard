# Time-Series Database Setup Guide

This guide covers the setup, configuration, and optimization of time-series data storage for the Cryptocurrency Indicator Dashboard.

## Overview

The system supports both **PostgreSQL** (standard) and **TimescaleDB** (optimized) for time-series data storage. TimescaleDB provides better performance for large datasets, but standard PostgreSQL works perfectly for development and smaller deployments.

## Database Options

### Option 1: Standard PostgreSQL (Recommended for Development)

**Pros**:
- Easy setup, widely available
- No additional extensions required
- Works with existing PostgreSQL infrastructure
- Good performance for moderate data volumes

**Cons**:
- Limited time-series optimizations
- Manual partitioning for large datasets
- Less efficient for very large historical datasets

### Option 2: TimescaleDB (Recommended for Production)

**Pros**:
- Automatic time-based partitioning (hypertables)
- Optimized time-series queries
- Built-in compression and retention policies
- Excellent performance for large datasets

**Cons**:
- Additional extension installation required
- More complex setup and configuration
- Resource overhead for small datasets

## PostgreSQL Setup (Standard)

### 1. Database Installation

#### On macOS (Homebrew)
```bash
brew install postgresql
brew services start postgresql
```

#### On Ubuntu/Debian
```bash
sudo apt update
sudo apt install postgresql postgresql-contrib
sudo systemctl start postgresql
sudo systemctl enable postgresql
```

#### On CentOS/RHEL
```bash
sudo yum install postgresql-server postgresql-contrib
sudo postgresql-setup initdb
sudo systemctl start postgresql
sudo systemctl enable postgresql
```

### 2. Database Configuration

Create database and user:
```sql
-- Connect to PostgreSQL as superuser
sudo -u postgres psql

-- Create database
CREATE DATABASE crypto_dashboard;

-- Create user
CREATE USER dashboard_user WITH PASSWORD 'your_secure_password';

-- Grant privileges
GRANT ALL PRIVILEGES ON DATABASE crypto_dashboard TO dashboard_user;

-- Exit
\q
```

### 3. Environment Configuration

Update your environment variables:
```bash
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=dashboard_user
export DB_PASSWORD=your_secure_password
export DB_NAME=crypto_dashboard
export DB_SSLMODE=disable
```

### 4. Application Setup

The application will automatically create optimized time-series tables:

```go
// Run the PostgreSQL time-series setup
go run test_postgres_timeseries.go
```

This creates:
- Optimized table schemas for time-series data
- Performance indexes for timestamp-based queries
- Automated data retention functions
- Query optimization statistics

## TimescaleDB Setup (Advanced)

### 1. TimescaleDB Installation

#### On macOS (Homebrew)
```bash
# Add TimescaleDB repository
brew tap timescale/tap

# Install TimescaleDB
brew install timescaledb

# Setup TimescaleDB
timescaledb-tune

# Restart PostgreSQL
brew services restart postgresql
```

#### On Ubuntu/Debian
```bash
# Add TimescaleDB repository
echo "deb https://packagecloud.io/timescale/timescaledb/ubuntu/ $(lsb_release -c -s) main" | sudo tee /etc/apt/sources.list.d/timescaledb.list

# Import GPG key
wget --quiet -O - https://packagecloud.io/timescale/timescaledb/gpgkey | sudo apt-key add -

# Update and install
sudo apt update
sudo apt install timescaledb-2-postgresql-14

# Tune PostgreSQL configuration
sudo timescaledb-tune

# Restart PostgreSQL
sudo systemctl restart postgresql
```

#### On CentOS/RHEL
```bash
# Add TimescaleDB repository
sudo yum install -y https://download.postgresql.org/pub/repos/yum/reporpms/EL-7-x86_64/pgdg-redhat-repo-latest.noarch.rpm

# Install TimescaleDB
sudo yum install -y timescaledb-2-postgresql-14

# Tune configuration
sudo timescaledb-tune

# Restart PostgreSQL
sudo systemctl restart postgresql
```

### 2. Enable TimescaleDB Extension

```sql
-- Connect to your database
psql -U dashboard_user -d crypto_dashboard

-- Enable TimescaleDB extension
CREATE EXTENSION IF NOT EXISTS timescaledb CASCADE;

-- Verify installation
SELECT * FROM timescaledb_information.hypertables;
```

### 3. Application Setup with TimescaleDB

```go
// Run TimescaleDB setup (will create hypertables)
go run test_timescale.go
```

## Table Schemas

### Price Data Points
```sql
CREATE TABLE price_data_points (
    id BIGSERIAL PRIMARY KEY,
    timestamp TIMESTAMPTZ NOT NULL,
    asset_symbol VARCHAR(10) NOT NULL,
    price_usd DECIMAL(20,8) NOT NULL,
    market_cap DECIMAL(30,2),
    volume_24h DECIMAL(30,2),
    data_source VARCHAR(50) NOT NULL,
    reliability_score DECIMAL(5,2),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Create hypertable (TimescaleDB only)
SELECT create_hypertable('price_data_points', 'timestamp', chunk_time_interval => INTERVAL '1 day');
```

### Indicator Data Points
```sql
CREATE TABLE indicator_data_points (
    id BIGSERIAL PRIMARY KEY,
    timestamp TIMESTAMPTZ NOT NULL,
    indicator_type VARCHAR(50) NOT NULL,
    value DECIMAL(20,8) NOT NULL,
    metadata JSONB,
    confidence_level DECIMAL(5,2),
    data_source VARCHAR(50) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Create hypertable (TimescaleDB only)
SELECT create_hypertable('indicator_data_points', 'timestamp', chunk_time_interval => INTERVAL '1 day');
```

### Network Metric Points
```sql
CREATE TABLE network_metric_points (
    id BIGSERIAL PRIMARY KEY,
    timestamp TIMESTAMPTZ NOT NULL,
    network VARCHAR(20) NOT NULL,
    hash_rate DECIMAL(30,2),
    difficulty DECIMAL(30,2),
    block_height BIGINT,
    total_supply DECIMAL(20,8),
    transaction_count BIGINT,
    fees_total DECIMAL(20,8),
    mempool_size BIGINT,
    data_source VARCHAR(50) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Create hypertable (TimescaleDB only)
SELECT create_hypertable('network_metric_points', 'timestamp', chunk_time_interval => INTERVAL '1 hour');
```

### Rainbow Chart Data Points
```sql
CREATE TABLE rainbow_chart_data_points (
    id BIGSERIAL PRIMARY KEY,
    timestamp TIMESTAMPTZ NOT NULL,
    bitcoin_price DECIMAL(20,8) NOT NULL,
    log_regression_price DECIMAL(20,8) NOT NULL,
    current_band VARCHAR(100) NOT NULL,
    current_band_color VARCHAR(7) NOT NULL,
    cycle_position DECIMAL(5,2) NOT NULL,
    risk_level VARCHAR(20) NOT NULL,
    days_from_genesis INTEGER NOT NULL,
    band_prices JSONB NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Create hypertable (TimescaleDB only)
SELECT create_hypertable('rainbow_chart_data_points', 'timestamp', chunk_time_interval => INTERVAL '1 day');
```

## Indexing Strategy

### Time-Series Optimized Indexes

```sql
-- Price data indexes
CREATE INDEX idx_price_data_timestamp_desc ON price_data_points (timestamp DESC);
CREATE INDEX idx_price_data_symbol_timestamp ON price_data_points (asset_symbol, timestamp DESC);
CREATE INDEX idx_price_data_source_timestamp ON price_data_points (data_source, timestamp DESC);

-- Indicator data indexes
CREATE INDEX idx_indicator_data_type_timestamp ON indicator_data_points (indicator_type, timestamp DESC);
CREATE INDEX idx_indicator_data_timestamp_desc ON indicator_data_points (timestamp DESC);
CREATE INDEX idx_indicator_confidence_high ON indicator_data_points (confidence_level) WHERE confidence_level > 70;

-- Network metrics indexes
CREATE INDEX idx_network_metrics_network_timestamp ON network_metric_points (network, timestamp DESC);
CREATE INDEX idx_network_metrics_block_height_desc ON network_metric_points (block_height DESC);

-- Rainbow chart indexes
CREATE INDEX idx_rainbow_chart_timestamp_desc ON rainbow_chart_data_points (timestamp DESC);
CREATE INDEX idx_rainbow_chart_band ON rainbow_chart_data_points (current_band);
CREATE INDEX idx_rainbow_chart_risk_level ON rainbow_chart_data_points (risk_level);
```

## Data Retention Policies

### Automatic Cleanup Function

```sql
-- Create cleanup function
CREATE OR REPLACE FUNCTION cleanup_old_timeseries_data()
RETURNS void AS $$
BEGIN
    -- Clean up price data older than 3 years
    DELETE FROM price_data_points 
    WHERE timestamp < NOW() - INTERVAL '3 years';
    
    -- Clean up indicator data older than 2 years
    DELETE FROM indicator_data_points 
    WHERE timestamp < NOW() - INTERVAL '2 years';
    
    -- Clean up market metrics older than 1 year
    DELETE FROM market_metric_points 
    WHERE timestamp < NOW() - INTERVAL '1 year';
    
    -- Clean up network metrics older than 1 year
    DELETE FROM network_metric_points 
    WHERE timestamp < NOW() - INTERVAL '1 year';
    
    -- Clean up rainbow chart data older than 5 years
    DELETE FROM rainbow_chart_data_points 
    WHERE timestamp < NOW() - INTERVAL '5 years';
    
    -- Log cleanup completion
    RAISE NOTICE 'Time-series data cleanup completed at %', NOW();
END;
$$ LANGUAGE plpgsql;
```

### TimescaleDB Retention Policies

```sql
-- For TimescaleDB, use built-in retention policies
SELECT add_retention_policy('price_data_points', INTERVAL '3 years');
SELECT add_retention_policy('indicator_data_points', INTERVAL '2 years');
SELECT add_retention_policy('network_metric_points', INTERVAL '1 year');
SELECT add_retention_policy('rainbow_chart_data_points', INTERVAL '5 years');
```

### Schedule Cleanup (Cron Job)

```bash
# Add to crontab for weekly cleanup
# Run every Sunday at 2 AM
0 2 * * 0 psql -U dashboard_user -d crypto_dashboard -c "SELECT cleanup_old_timeseries_data();"
```

## Performance Optimization

### PostgreSQL Configuration

Add to `postgresql.conf`:
```ini
# Memory settings
shared_buffers = 256MB                    # 25% of RAM
effective_cache_size = 1GB               # 75% of RAM
work_mem = 4MB                           # Per-query memory
maintenance_work_mem = 64MB              # Maintenance operations

# Checkpoint settings
checkpoint_completion_target = 0.7       # Spread checkpoints
wal_buffers = 16MB                       # WAL buffer size

# Query planner
random_page_cost = 1.1                   # SSD optimization
effective_io_concurrency = 200           # Concurrent I/O

# Time-series specific
max_connections = 100                    # Connection limit
shared_preload_libraries = 'timescaledb' # TimescaleDB (if used)
```

### TimescaleDB Specific Settings

```ini
# TimescaleDB settings
timescaledb.max_background_workers = 8
timescaledb.restoring = off
```

### Connection Pooling

Use pgBouncer for connection pooling:
```ini
# pgbouncer.ini
[databases]
crypto_dashboard = host=localhost port=5432 dbname=crypto_dashboard

[pgbouncer]
pool_mode = transaction
max_client_conn = 100
default_pool_size = 25
```

## Application Integration

### Go Service Implementation

```go
package main

import (
    "crypto-indicator-dashboard/config"
    "crypto-indicator-dashboard/services"
    "crypto-indicator-dashboard/internal/infrastructure/database"
)

func main() {
    // Connect to database
    config.ConnectDatabase()
    
    // Setup time-series tables
    pgManager := database.NewPostgresTimeSeriesManager(config.DB, logger)
    err := pgManager.SetupTimeSerieTables()
    if err != nil {
        log.Fatal("Failed to setup time-series tables:", err)
    }
    
    // Create time-series service
    timeSeriesService := services.NewTimeSeriesService(config.DB, logger)
    
    // Store data
    priceData := []services.PriceDataPoint{
        {
            Timestamp:        time.Now(),
            AssetSymbol:      "BTC",
            PriceUSD:         45000.00,
            MarketCap:        &[]float64{850000000000}[0],
            DataSource:       "coincap",
            ReliabilityScore: 95.0,
        },
    }
    
    err = timeSeriesService.StorePriceData(priceData)
    if err != nil {
        log.Printf("Error storing data: %v", err)
    }
}
```

### Batch Operations

For high-throughput scenarios:
```go
// Batch insert for efficiency
func (tss *TimeSeriesService) StorePriceDataBatch(dataPoints []PriceDataPoint) error {
    // Use GORM's CreateInBatches for efficient insertion
    return tss.db.CreateInBatches(dataPoints, 1000).Error
}
```

## Monitoring & Maintenance

### Database Statistics

```sql
-- Check table sizes
SELECT 
    tablename,
    pg_size_pretty(pg_total_relation_size(tablename::regclass)) as size
FROM pg_tables 
WHERE tablename LIKE '%_points';

-- Check index usage
SELECT 
    schemaname,
    tablename,
    indexname,
    idx_scan,
    idx_tup_read,
    idx_tup_fetch
FROM pg_stat_user_indexes
ORDER BY idx_scan DESC;

-- TimescaleDB specific stats
SELECT * FROM timescaledb_information.hypertables;
SELECT * FROM timescaledb_information.chunks ORDER BY chunk_name;
```

### Performance Monitoring

```sql
-- Query performance
SELECT 
    query,
    calls,
    total_time,
    mean_time,
    rows
FROM pg_stat_statements 
ORDER BY total_time DESC 
LIMIT 10;

-- Connection monitoring
SELECT 
    count(*) as active_connections,
    state
FROM pg_stat_activity 
GROUP BY state;
```

## Backup & Recovery

### Regular Backups

```bash
#!/bin/bash
# backup_timeseries.sh

DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_DIR="/path/to/backups"

# Full database backup
pg_dump -U dashboard_user -h localhost crypto_dashboard > "$BACKUP_DIR/crypto_dashboard_$DATE.sql"

# Compress backup
gzip "$BACKUP_DIR/crypto_dashboard_$DATE.sql"

# Keep only last 7 days of backups
find "$BACKUP_DIR" -name "crypto_dashboard_*.sql.gz" -mtime +7 -delete

echo "Backup completed: crypto_dashboard_$DATE.sql.gz"
```

### Point-in-Time Recovery

Enable WAL archiving in `postgresql.conf`:
```ini
wal_level = replica
archive_mode = on
archive_command = 'cp %p /path/to/archive/%f'
```

## Troubleshooting

### Common Issues

#### 1. Connection Errors
```bash
# Check PostgreSQL status
sudo systemctl status postgresql

# Check connection
psql -U dashboard_user -h localhost -d crypto_dashboard -c "SELECT NOW();"
```

#### 2. Performance Issues
```sql
-- Check slow queries
SELECT query, mean_time, calls 
FROM pg_stat_statements 
WHERE mean_time > 1000 
ORDER BY mean_time DESC;

-- Check table bloat
SELECT 
    tablename,
    pg_size_pretty(pg_total_relation_size(tablename::regclass)) as size,
    n_tup_ins,
    n_tup_upd,
    n_tup_del
FROM pg_stat_user_tables;
```

#### 3. Disk Space Issues
```sql
-- Check database size
SELECT 
    pg_database.datname,
    pg_size_pretty(pg_database_size(pg_database.datname)) AS size
FROM pg_database;

-- Run cleanup manually
SELECT cleanup_old_timeseries_data();
```

### Performance Tuning Checklist

- [ ] Proper indexing on timestamp columns
- [ ] Regular VACUUM and ANALYZE operations
- [ ] Appropriate retention policies
- [ ] Connection pooling configured
- [ ] Query optimization with EXPLAIN
- [ ] Monitoring slow queries
- [ ] Regular backup verification

## Testing

### Functionality Tests

```bash
# Test PostgreSQL setup
go run test_postgres_timeseries.go

# Test TimescaleDB setup (if available)
go run test_timescale.go

# Test data aggregation
go run test_aggregator.go
```

### Load Testing

```go
// Load test with concurrent inserts
func loadTestInserts() {
    concurrency := 10
    recordsPerWorker := 1000
    
    var wg sync.WaitGroup
    for i := 0; i < concurrency; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            // Insert records
            for j := 0; j < recordsPerWorker; j++ {
                // Insert test data
            }
        }()
    }
    wg.Wait()
}
```

This comprehensive setup guide ensures optimal time-series database performance for cryptocurrency market analysis and provides the foundation for machine learning and advanced analytics capabilities.