package database

import (
	"fmt"
	"time"

	"gorm.io/gorm"
	"crypto-indicator-dashboard/pkg/logger"
)

// TimescaleManager handles TimescaleDB hypertable setup and management
type TimescaleManager struct {
	db     *gorm.DB
	logger logger.Logger
}

// NewTimescaleManager creates a new TimescaleDB manager
func NewTimescaleManager(db *gorm.DB, logger logger.Logger) *TimescaleManager {
	return &TimescaleManager{
		db:     db,
		logger: logger,
	}
}

// SetupHypertables creates and configures TimescaleDB hypertables for time-series data
func (tm *TimescaleManager) SetupHypertables() error {
	tm.logger.Info("Setting up TimescaleDB hypertables...")

	// Enable TimescaleDB extension
	if err := tm.enableTimescaleExtension(); err != nil {
		return fmt.Errorf("failed to enable TimescaleDB extension: %w", err)
	}

	// Create time-series tables
	tables := []HypertableConfig{
		{
			TableName:    "price_data",
			TimeColumn:   "timestamp",
			ChunkInterval: "1 day",
			Schema: `
				CREATE TABLE IF NOT EXISTS price_data (
					id SERIAL PRIMARY KEY,
					timestamp TIMESTAMPTZ NOT NULL,
					asset_symbol VARCHAR(10) NOT NULL,
					price_usd DECIMAL(20,8) NOT NULL,
					market_cap DECIMAL(30,2),
					volume_24h DECIMAL(30,2),
					data_source VARCHAR(50) NOT NULL,
					reliability_score DECIMAL(5,2),
					created_at TIMESTAMPTZ DEFAULT NOW()
				);
			`,
		},
		{
			TableName:    "indicator_data",
			TimeColumn:   "timestamp",
			ChunkInterval: "1 day",
			Schema: `
				CREATE TABLE IF NOT EXISTS indicator_data (
					id SERIAL PRIMARY KEY,
					timestamp TIMESTAMPTZ NOT NULL,
					indicator_type VARCHAR(50) NOT NULL,
					value DECIMAL(20,8) NOT NULL,
					metadata JSONB,
					confidence_level DECIMAL(5,2),
					data_source VARCHAR(50) NOT NULL,
					created_at TIMESTAMPTZ DEFAULT NOW()
				);
			`,
		},
		{
			TableName:    "market_metrics",
			TimeColumn:   "timestamp",
			ChunkInterval: "1 hour",
			Schema: `
				CREATE TABLE IF NOT EXISTS market_metrics (
					id SERIAL PRIMARY KEY,
					timestamp TIMESTAMPTZ NOT NULL,
					metric_name VARCHAR(100) NOT NULL,
					metric_value DECIMAL(20,8) NOT NULL,
					additional_data JSONB,
					asset_symbol VARCHAR(10),
					data_source VARCHAR(50) NOT NULL,
					created_at TIMESTAMPTZ DEFAULT NOW()
				);
			`,
		},
		{
			TableName:    "rainbow_chart_data",
			TimeColumn:   "timestamp",
			ChunkInterval: "1 day",
			Schema: `
				CREATE TABLE IF NOT EXISTS rainbow_chart_data (
					id SERIAL PRIMARY KEY,
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
			`,
		},
		{
			TableName:    "network_metrics",
			TimeColumn:   "timestamp",
			ChunkInterval: "1 hour",
			Schema: `
				CREATE TABLE IF NOT EXISTS network_metrics (
					id SERIAL PRIMARY KEY,
					timestamp TIMESTAMPTZ NOT NULL,
					network VARCHAR(20) NOT NULL,
					hash_rate DECIMAL(30,2),
					difficulty DECIMAL(30,2),
					block_height BIGINT,
					total_supply DECIMAL(20,8),
					transaction_count BIGINT,
					fees_total DECIMAL(20,8),
					mempool_size INTEGER,
					data_source VARCHAR(50) NOT NULL,
					created_at TIMESTAMPTZ DEFAULT NOW()
				);
			`,
		},
	}

	// Create tables and hypertables
	for _, tableConfig := range tables {
		if err := tm.createHypertable(tableConfig); err != nil {
			return fmt.Errorf("failed to create hypertable %s: %w", tableConfig.TableName, err)
		}
	}

	// Create indexes for better query performance
	if err := tm.createIndexes(); err != nil {
		return fmt.Errorf("failed to create indexes: %w", err)
	}

	tm.logger.Info("TimescaleDB hypertables setup completed successfully")
	return nil
}

// HypertableConfig defines configuration for a TimescaleDB hypertable
type HypertableConfig struct {
	TableName     string
	TimeColumn    string
	ChunkInterval string
	Schema        string
}

// enableTimescaleExtension enables the TimescaleDB extension
func (tm *TimescaleManager) enableTimescaleExtension() error {
	tm.logger.Info("Enabling TimescaleDB extension...")
	
	query := "CREATE EXTENSION IF NOT EXISTS timescaledb CASCADE;"
	if err := tm.db.Exec(query).Error; err != nil {
		return fmt.Errorf("failed to enable TimescaleDB extension: %w", err)
	}
	
	tm.logger.Info("TimescaleDB extension enabled successfully")
	return nil
}

// createHypertable creates a table and converts it to a hypertable
func (tm *TimescaleManager) createHypertable(config HypertableConfig) error {
	tm.logger.Info("Creating hypertable", "table", config.TableName)

	// Create the table
	if err := tm.db.Exec(config.Schema).Error; err != nil {
		return fmt.Errorf("failed to create table %s: %w", config.TableName, err)
	}

	// Check if table is already a hypertable
	var isHypertable bool
	checkQuery := `
		SELECT EXISTS (
			SELECT 1 FROM timescaledb_information.hypertables 
			WHERE hypertable_name = ?
		);
	`
	if err := tm.db.Raw(checkQuery, config.TableName).Scan(&isHypertable).Error; err != nil {
		tm.logger.Warn("Could not check if table is already a hypertable", "table", config.TableName, "error", err)
	}

	// Convert to hypertable if not already one
	if !isHypertable {
		hypertableQuery := fmt.Sprintf(
			"SELECT create_hypertable('%s', '%s', chunk_time_interval => interval '%s');",
			config.TableName,
			config.TimeColumn,
			config.ChunkInterval,
		)
		
		if err := tm.db.Exec(hypertableQuery).Error; err != nil {
			return fmt.Errorf("failed to create hypertable %s: %w", config.TableName, err)
		}
		
		tm.logger.Info("Hypertable created successfully", "table", config.TableName)
	} else {
		tm.logger.Info("Table is already a hypertable", "table", config.TableName)
	}

	return nil
}

// createIndexes creates optimized indexes for time-series queries
func (tm *TimescaleManager) createIndexes() error {
	tm.logger.Info("Creating TimescaleDB indexes...")

	indexes := []string{
		// Price data indexes
		"CREATE INDEX IF NOT EXISTS idx_price_data_symbol_time ON price_data (asset_symbol, timestamp DESC);",
		"CREATE INDEX IF NOT EXISTS idx_price_data_source ON price_data (data_source);",
		
		// Indicator data indexes
		"CREATE INDEX IF NOT EXISTS idx_indicator_type_time ON indicator_data (indicator_type, timestamp DESC);",
		"CREATE INDEX IF NOT EXISTS idx_indicator_confidence ON indicator_data (confidence_level) WHERE confidence_level > 70;",
		
		// Market metrics indexes
		"CREATE INDEX IF NOT EXISTS idx_market_metrics_name_time ON market_metrics (metric_name, timestamp DESC);",
		"CREATE INDEX IF NOT EXISTS idx_market_metrics_asset ON market_metrics (asset_symbol, timestamp DESC);",
		
		// Rainbow chart indexes
		"CREATE INDEX IF NOT EXISTS idx_rainbow_chart_time ON rainbow_chart_data (timestamp DESC);",
		"CREATE INDEX IF NOT EXISTS idx_rainbow_chart_band ON rainbow_chart_data (current_band);",
		
		// Network metrics indexes
		"CREATE INDEX IF NOT EXISTS idx_network_metrics_network_time ON network_metrics (network, timestamp DESC);",
		"CREATE INDEX IF NOT EXISTS idx_network_metrics_block_height ON network_metrics (block_height DESC);",
	}

	for _, indexQuery := range indexes {
		if err := tm.db.Exec(indexQuery).Error; err != nil {
			tm.logger.Warn("Failed to create index", "query", indexQuery, "error", err)
		}
	}

	tm.logger.Info("TimescaleDB indexes created successfully")
	return nil
}

// SetupDataRetentionPolicies configures automatic data retention
func (tm *TimescaleManager) SetupDataRetentionPolicies() error {
	tm.logger.Info("Setting up data retention policies...")

	policies := []RetentionPolicy{
		{
			TableName: "price_data",
			Interval:  "3 years", // Keep price data for 3 years
		},
		{
			TableName: "indicator_data",
			Interval:  "2 years", // Keep indicator data for 2 years
		},
		{
			TableName: "market_metrics",
			Interval:  "1 year", // Keep market metrics for 1 year
		},
		{
			TableName: "rainbow_chart_data",
			Interval:  "5 years", // Keep rainbow chart data for 5 years (historical analysis)
		},
		{
			TableName: "network_metrics",
			Interval:  "1 year", // Keep network metrics for 1 year
		},
	}

	for _, policy := range policies {
		if err := tm.addRetentionPolicy(policy); err != nil {
			tm.logger.Warn("Failed to add retention policy", "table", policy.TableName, "error", err)
		}
	}

	tm.logger.Info("Data retention policies setup completed")
	return nil
}

// RetentionPolicy defines data retention configuration
type RetentionPolicy struct {
	TableName string
	Interval  string
}

// addRetentionPolicy adds a retention policy to a hypertable
func (tm *TimescaleManager) addRetentionPolicy(policy RetentionPolicy) error {
	// Remove existing policy if any
	removeQuery := fmt.Sprintf("SELECT remove_retention_policy('%s', if_exists => true);", policy.TableName)
	tm.db.Exec(removeQuery) // Ignore errors

	// Add new retention policy
	addQuery := fmt.Sprintf(
		"SELECT add_retention_policy('%s', INTERVAL '%s');",
		policy.TableName,
		policy.Interval,
	)

	if err := tm.db.Exec(addQuery).Error; err != nil {
		return fmt.Errorf("failed to add retention policy for %s: %w", policy.TableName, err)
	}

	tm.logger.Info("Retention policy added", "table", policy.TableName, "interval", policy.Interval)
	return nil
}

// OptimizeHypertables runs maintenance tasks on hypertables
func (tm *TimescaleManager) OptimizeHypertables() error {
	tm.logger.Info("Running hypertable optimization...")

	// Recompute chunk statistics
	tables := []string{"price_data", "indicator_data", "market_metrics", "rainbow_chart_data", "network_metrics"}
	
	for _, table := range tables {
		// Recompute chunk statistics for better query planning
		statsQuery := fmt.Sprintf("SELECT recompute_chunk_stats('%s');", table)
		if err := tm.db.Exec(statsQuery).Error; err != nil {
			tm.logger.Warn("Failed to recompute chunk stats", "table", table, "error", err)
		}
	}

	tm.logger.Info("Hypertable optimization completed")
	return nil
}

// GetTableStats returns statistics about TimescaleDB tables
func (tm *TimescaleManager) GetTableStats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Get hypertable information
	var hypertables []map[string]interface{}
	hypertableQuery := `
		SELECT 
			hypertable_name,
			num_chunks,
			table_size,
			index_size,
			total_size
		FROM timescaledb_information.hypertables
		JOIN timescaledb_information.hypertable_detailed_size 
		ON hypertable_name = hypertable_schema||'.'||hypertable_name;
	`
	
	if err := tm.db.Raw(hypertableQuery).Scan(&hypertables).Error; err != nil {
		return nil, fmt.Errorf("failed to get hypertable stats: %w", err)
	}

	stats["hypertables"] = hypertables
	stats["last_updated"] = time.Now()
	stats["total_hypertables"] = len(hypertables)

	return stats, nil
}