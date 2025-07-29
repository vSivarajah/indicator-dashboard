package cache

import (
	"context"
	"crypto-indicator-dashboard/pkg/errors"
	"crypto-indicator-dashboard/pkg/logger"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
)

// CacheService defines the interface for cache operations
type CacheService interface {
	Get(ctx context.Context, key string, dest interface{}) error
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
	FlushAll(ctx context.Context) error
	GetOrSet(ctx context.Context, key string, dest interface{}, fetcher func() (interface{}, error), expiration time.Duration) error
}

// redisCache implements CacheService using Redis
type redisCache struct {
	client *redis.Client
	logger logger.Logger
}

// NewRedisCache creates a new Redis cache service
func NewRedisCache(client *redis.Client, logger logger.Logger) CacheService {
	return &redisCache{
		client: client,
		logger: logger,
	}
}

// Get retrieves a value from cache and unmarshals it into dest
func (c *redisCache) Get(ctx context.Context, key string, dest interface{}) error {
	c.logger.Debug("Getting value from cache", "key", key)

	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			c.logger.Debug("Cache miss", "key", key)
			return errors.NotFound("cache_key")
		}
		c.logger.Error("Failed to get value from cache", "error", err, "key", key)
		return errors.Wrap(err, errors.ErrorTypeExternal, "failed to get value from cache")
	}

	if err := json.Unmarshal([]byte(val), dest); err != nil {
		c.logger.Error("Failed to unmarshal cached value", "error", err, "key", key)
		return errors.Wrap(err, errors.ErrorTypeInternal, "failed to unmarshal cached value")
	}

	c.logger.Debug("Cache hit", "key", key)
	return nil
}

// Set stores a value in cache with expiration
func (c *redisCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	c.logger.Debug("Setting value in cache", "key", key, "expiration", expiration)

	data, err := json.Marshal(value)
	if err != nil {
		c.logger.Error("Failed to marshal value for cache", "error", err, "key", key)
		return errors.Wrap(err, errors.ErrorTypeInternal, "failed to marshal value for cache")
	}

	if err := c.client.Set(ctx, key, data, expiration).Err(); err != nil {
		c.logger.Error("Failed to set value in cache", "error", err, "key", key)
		return errors.Wrap(err, errors.ErrorTypeExternal, "failed to set value in cache")
	}

	c.logger.Debug("Successfully set value in cache", "key", key)
	return nil
}

// Delete removes a value from cache
func (c *redisCache) Delete(ctx context.Context, key string) error {
	c.logger.Debug("Deleting value from cache", "key", key)

	result, err := c.client.Del(ctx, key).Result()
	if err != nil {
		c.logger.Error("Failed to delete value from cache", "error", err, "key", key)
		return errors.Wrap(err, errors.ErrorTypeExternal, "failed to delete value from cache")
	}

	if result == 0 {
		c.logger.Debug("Key not found in cache", "key", key)
		return errors.NotFound("cache_key")
	}

	c.logger.Debug("Successfully deleted value from cache", "key", key)
	return nil
}

// Exists checks if a key exists in cache
func (c *redisCache) Exists(ctx context.Context, key string) (bool, error) {
	c.logger.Debug("Checking if key exists in cache", "key", key)

	result, err := c.client.Exists(ctx, key).Result()
	if err != nil {
		c.logger.Error("Failed to check key existence in cache", "error", err, "key", key)
		return false, errors.Wrap(err, errors.ErrorTypeExternal, "failed to check key existence in cache")
	}

	exists := result > 0
	c.logger.Debug("Key existence check result", "key", key, "exists", exists)
	return exists, nil
}

// FlushAll removes all keys from cache
func (c *redisCache) FlushAll(ctx context.Context) error {
	c.logger.Info("Flushing all cache data")

	if err := c.client.FlushAll(ctx).Err(); err != nil {
		c.logger.Error("Failed to flush cache", "error", err)
		return errors.Wrap(err, errors.ErrorTypeExternal, "failed to flush cache")
	}

	c.logger.Info("Successfully flushed all cache data")
	return nil
}

// GetOrSet retrieves a value from cache or sets it if not found
func (c *redisCache) GetOrSet(ctx context.Context, key string, dest interface{}, fetcher func() (interface{}, error), expiration time.Duration) error {
	c.logger.Debug("GetOrSet operation", "key", key, "expiration", expiration)

	// Try to get from cache first
	err := c.Get(ctx, key, dest)
	if err == nil {
		c.logger.Debug("Found value in cache", "key", key)
		return nil
	}

	// If not found or error other than not found, fetch new value
	if !errors.IsType(err, errors.ErrorTypeNotFound) {
		c.logger.Warn("Cache get operation failed, fetching fresh data", "error", err, "key", key)
	}

	c.logger.Debug("Cache miss, fetching fresh data", "key", key)
	
	// Fetch fresh data
	value, err := fetcher()
	if err != nil {
		c.logger.Error("Failed to fetch fresh data", "error", err, "key", key)
		return errors.Wrap(err, errors.ErrorTypeExternal, "failed to fetch fresh data")
	}

	// Set in cache for future use
	if setErr := c.Set(ctx, key, value, expiration); setErr != nil {
		c.logger.Warn("Failed to cache fresh data", "error", setErr, "key", key)
		// Don't return error here, as we still have the data
	}

	// Marshal and unmarshal to populate dest with the correct type
	data, err := json.Marshal(value)
	if err != nil {
		c.logger.Error("Failed to marshal fetched value", "error", err, "key", key)
		return errors.Wrap(err, errors.ErrorTypeInternal, "failed to marshal fetched value")
	}

	if err := json.Unmarshal(data, dest); err != nil {
		c.logger.Error("Failed to unmarshal fetched value", "error", err, "key", key)
		return errors.Wrap(err, errors.ErrorTypeInternal, "failed to unmarshal fetched value")
	}

	c.logger.Debug("Successfully fetched and cached fresh data", "key", key)
	return nil
}

// mockCache implements CacheService for testing or when Redis is not available
type mockCache struct {
	data   map[string]cacheItem
	logger logger.Logger
}

type cacheItem struct {
	value      []byte
	expiration time.Time
}

// NewMockCache creates a new mock cache service
func NewMockCache(logger logger.Logger) CacheService {
	return &mockCache{
		data:   make(map[string]cacheItem),
		logger: logger,
	}
}

// Get retrieves a value from mock cache
func (c *mockCache) Get(ctx context.Context, key string, dest interface{}) error {
	c.logger.Debug("Getting value from mock cache", "key", key)

	item, exists := c.data[key]
	if !exists || time.Now().After(item.expiration) {
		if exists && time.Now().After(item.expiration) {
			delete(c.data, key)
		}
		c.logger.Debug("Mock cache miss", "key", key)
		return errors.NotFound("cache_key")
	}

	if err := json.Unmarshal(item.value, dest); err != nil {
		c.logger.Error("Failed to unmarshal cached value", "error", err, "key", key)
		return errors.Wrap(err, errors.ErrorTypeInternal, "failed to unmarshal cached value")
	}

	c.logger.Debug("Mock cache hit", "key", key)
	return nil
}

// Set stores a value in mock cache
func (c *mockCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	c.logger.Debug("Setting value in mock cache", "key", key, "expiration", expiration)

	data, err := json.Marshal(value)
	if err != nil {
		c.logger.Error("Failed to marshal value for mock cache", "error", err, "key", key)
		return errors.Wrap(err, errors.ErrorTypeInternal, "failed to marshal value for cache")
	}

	c.data[key] = cacheItem{
		value:      data,
		expiration: time.Now().Add(expiration),
	}

	c.logger.Debug("Successfully set value in mock cache", "key", key)
	return nil
}

// Delete removes a value from mock cache
func (c *mockCache) Delete(ctx context.Context, key string) error {
	c.logger.Debug("Deleting value from mock cache", "key", key)

	if _, exists := c.data[key]; !exists {
		c.logger.Debug("Key not found in mock cache", "key", key)
		return errors.NotFound("cache_key")
	}

	delete(c.data, key)
	c.logger.Debug("Successfully deleted value from mock cache", "key", key)
	return nil
}

// Exists checks if a key exists in mock cache
func (c *mockCache) Exists(ctx context.Context, key string) (bool, error) {
	c.logger.Debug("Checking if key exists in mock cache", "key", key)

	item, exists := c.data[key]
	if exists && time.Now().After(item.expiration) {
		delete(c.data, key)
		exists = false
	}

	c.logger.Debug("Key existence check result", "key", key, "exists", exists)
	return exists, nil
}

// FlushAll removes all keys from mock cache
func (c *mockCache) FlushAll(ctx context.Context) error {
	c.logger.Info("Flushing all mock cache data")
	c.data = make(map[string]cacheItem)
	c.logger.Info("Successfully flushed all mock cache data")
	return nil
}

// GetOrSet retrieves a value from mock cache or sets it if not found
func (c *mockCache) GetOrSet(ctx context.Context, key string, dest interface{}, fetcher func() (interface{}, error), expiration time.Duration) error {
	c.logger.Debug("GetOrSet operation on mock cache", "key", key, "expiration", expiration)

	// Try to get from cache first
	err := c.Get(ctx, key, dest)
	if err == nil {
		c.logger.Debug("Found value in mock cache", "key", key)
		return nil
	}

	// If not found or error other than not found, fetch new value
	if !errors.IsType(err, errors.ErrorTypeNotFound) {
		c.logger.Warn("Mock cache get operation failed, fetching fresh data", "error", err, "key", key)
	}

	c.logger.Debug("Mock cache miss, fetching fresh data", "key", key)
	
	// Fetch fresh data
	value, err := fetcher()
	if err != nil {
		c.logger.Error("Failed to fetch fresh data", "error", err, "key", key)
		return errors.Wrap(err, errors.ErrorTypeExternal, "failed to fetch fresh data")
	}

	// Set in cache for future use
	if setErr := c.Set(ctx, key, value, expiration); setErr != nil {
		c.logger.Warn("Failed to cache fresh data", "error", setErr, "key", key)
		// Don't return error here, as we still have the data
	}

	// Marshal and unmarshal to populate dest with the correct type
	data, err := json.Marshal(value)
	if err != nil {
		c.logger.Error("Failed to marshal fetched value", "error", err, "key", key)
		return errors.Wrap(err, errors.ErrorTypeInternal, "failed to marshal fetched value")
	}

	if err := json.Unmarshal(data, dest); err != nil {
		c.logger.Error("Failed to unmarshal fetched value", "error", err, "key", key)
		return errors.Wrap(err, errors.ErrorTypeInternal, "failed to unmarshal fetched value")
	}

	c.logger.Debug("Successfully fetched and cached fresh data in mock cache", "key", key)
	return nil
}