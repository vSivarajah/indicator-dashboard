package interfaces

import (
	"context"
	"time"
)

// CacheService defines the interface for caching operations
type CacheService interface {
	// Get retrieves a cached value by key
	Get(ctx context.Context, key string) ([]byte, error)
	
	// Set stores a value in cache with expiration
	Set(ctx context.Context, key string, value []byte, expiration time.Duration) error
	
	// Delete removes a cached value
	Delete(ctx context.Context, key string) error
	
	// GetOrSet gets a value from cache or sets it using the provided function
	GetOrSet(ctx context.Context, key string, result interface{}, fetchFn func() (interface{}, error), expiration time.Duration) error
	
	// Exists checks if a key exists in cache
	Exists(ctx context.Context, key string) (bool, error)
	
	// TTL returns the time-to-live for a key
	TTL(ctx context.Context, key string) (time.Duration, error)
	
	// Clear removes all cached values (use with caution)
	Clear(ctx context.Context) error
	
	// Keys returns all keys matching a pattern
	Keys(ctx context.Context, pattern string) ([]string, error)
	
	// Size returns the number of keys in cache
	Size(ctx context.Context) (int64, error)
	
	// HealthCheck performs a health check on the cache service
	HealthCheck(ctx context.Context) error
}

// DistributedLock defines the interface for distributed locking
type DistributedLock interface {
	// Acquire attempts to acquire a lock
	Acquire(ctx context.Context, key string, expiration time.Duration) (bool, error)
	
	// Release releases a held lock
	Release(ctx context.Context, key string) error
	
	// Extend extends the expiration time of a held lock
	Extend(ctx context.Context, key string, expiration time.Duration) error
	
	// IsLocked checks if a key is currently locked
	IsLocked(ctx context.Context, key string) (bool, error)
}

// CacheMetrics defines the interface for cache metrics
type CacheMetrics interface {
	// RecordHit records a cache hit
	RecordHit(ctx context.Context, key string)
	
	// RecordMiss records a cache miss
	RecordMiss(ctx context.Context, key string)
	
	// RecordSet records a cache set operation
	RecordSet(ctx context.Context, key string, size int64)
	
	// RecordDelete records a cache delete operation
	RecordDelete(ctx context.Context, key string)
	
	// GetHitRatio returns the cache hit ratio
	GetHitRatio(ctx context.Context) (float64, error)
	
	// GetStats returns detailed cache statistics
	GetStats(ctx context.Context) (map[string]interface{}, error)
}

// CacheConfig defines the configuration for cache services
type CacheConfig struct {
	// DefaultTTL is the default time-to-live for cache entries
	DefaultTTL time.Duration
	
	// MaxKeys is the maximum number of keys to store
	MaxKeys int64
	
	// MaxMemory is the maximum memory usage in bytes
	MaxMemory int64
	
	// KeyPrefix is the prefix for all cache keys
	KeyPrefix string
	
	// EnableMetrics enables cache metrics collection
	EnableMetrics bool
	
	// EnableCompression enables value compression
	EnableCompression bool
}

// CacheEntry represents a cache entry with metadata
type CacheEntry struct {
	Key        string
	Value      []byte
	Expiration time.Time
	Size       int64
	AccessCount int64
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// CacheOperation represents different cache operations for monitoring
type CacheOperation string

const (
	CacheOpGet    CacheOperation = "get"
	CacheOpSet    CacheOperation = "set"
	CacheOpDelete CacheOperation = "delete"
	CacheOpClear  CacheOperation = "clear"
	CacheOpKeys   CacheOperation = "keys"
	CacheOpSize   CacheOperation = "size"
)