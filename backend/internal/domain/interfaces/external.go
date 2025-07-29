package interfaces

import (
	"context"
	"crypto-indicator-dashboard/internal/infrastructure/external"
	"time"
)

// CoinCapClient defines the interface for CoinCap API interactions
type CoinCapClient interface {
	GetAssets(limit int) (*external.AssetsResponse, error)
	GetAsset(assetID string) (*external.AssetResponse, error)
	GetAssetHistory(assetID, interval string, start, end *time.Time) (*external.HistoryResponse, error)
	GetMarkets(assetID string, limit int) (*external.MarketsResponse, error)
	GetBitcoinPrice() (float64, error)
	GetTop10Assets() (*external.AssetsResponse, error)
	GetBitcoinHistoricalData(interval string, days int) (*external.HistoryResponse, error)
	GetGlobalMarketData() (map[string]interface{}, error)
	HealthCheck() error
}

// CoinMarketCapClient defines the interface for CoinMarketCap API interactions
type CoinMarketCapClient interface {
	GetLatestQuotes(symbols []string) (map[string]interface{}, error)
	GetGlobalMetrics() (map[string]interface{}, error)
	GetHistoricalData(symbol string, start, end time.Time) ([]map[string]interface{}, error)
	HealthCheck() error
}

// BlockchainClient defines the interface for blockchain data interactions
type BlockchainClient interface {
	GetNetworkStats() (map[string]interface{}, error)
	GetBlockHeight() (int64, error)
	GetHashRate() (float64, error)
	GetDifficulty() (float64, error)
	GetMempoolSize() (int64, error)
	HealthCheck() error
}

// HTTPClient defines the interface for HTTP interactions
type HTTPClient interface {
	Get(ctx context.Context, url string, headers map[string]string) ([]byte, error)
	Post(ctx context.Context, url string, body []byte, headers map[string]string) ([]byte, error)
	Put(ctx context.Context, url string, body []byte, headers map[string]string) ([]byte, error)
	Delete(ctx context.Context, url string, headers map[string]string) ([]byte, error)
}

// CacheClient defines the interface for cache operations
type CacheClient interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value []byte, expiration time.Duration) error
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
	TTL(ctx context.Context, key string) (time.Duration, error)
	Ping(ctx context.Context) error
}

// MessageQueue defines the interface for message queue operations
type MessageQueue interface {
	Publish(ctx context.Context, topic string, message []byte) error
	Subscribe(ctx context.Context, topic string, handler func([]byte) error) error
	Close() error
}

// MetricsCollector defines the interface for metrics collection
type MetricsCollector interface {
	IncrementCounter(name string, labels map[string]string)
	RecordHistogram(name string, value float64, labels map[string]string)
	RecordGauge(name string, value float64, labels map[string]string)
	StartTimer(name string, labels map[string]string) TimerStop
}

// TimerStop represents a function to stop a timer
type TimerStop func()

// NotificationService defines the interface for sending notifications
type NotificationService interface {
	SendEmail(ctx context.Context, to, subject, body string) error
	SendSMS(ctx context.Context, to, message string) error
	SendWebhook(ctx context.Context, url string, payload interface{}) error
	SendPushNotification(ctx context.Context, deviceToken, title, message string) error
}

// HealthChecker defines the interface for health check operations
type HealthChecker interface {
	HealthCheck() error
	Name() string
}

// APIRateLimiter defines the interface for API rate limiting
type APIRateLimiter interface {
	Allow(key string) bool
	Reset(key string) error
	GetRemaining(key string) int
	GetResetTime(key string) time.Time
}