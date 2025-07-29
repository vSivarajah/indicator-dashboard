package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"crypto-indicator-dashboard/internal/domain/services"
	"crypto-indicator-dashboard/pkg/logger"
)

// cacheServiceImpl implements the CacheService interface
type cacheServiceImpl struct {
	redisCache   services.CacheService
	fallbackCache map[string]fallbackCacheItem
	logger       logger.Logger
}

// fallbackCacheItem represents an item in the fallback cache
type fallbackCacheItem struct {
	Data       []byte
	ExpiresAt  time.Time
}

// NewCacheService creates a new cache service with Redis primary and in-memory fallback
func NewCacheService(redisCache services.CacheService, logger logger.Logger) services.CacheService {
	return &cacheServiceImpl{
		redisCache:    redisCache,
		fallbackCache: make(map[string]fallbackCacheItem),
		logger:        logger,
	}
}

// GetOrSet gets a value from cache or sets it using the provided function
func (c *cacheServiceImpl) GetOrSet(ctx context.Context, key string, dest interface{}, expiration interface{}, setFunc func() (interface{}, error)) error {
	// Try to get from Redis first
	if c.redisCache != nil {
		err := c.redisCache.Get(ctx, key, dest)
		if err == nil {
			c.logger.Debug("Cache hit from Redis", "key", key)
			return nil
		}
		c.logger.Debug("Cache miss from Redis", "key", key, "error", err)
	}
	
	// Try fallback cache
	if item, exists := c.fallbackCache[key]; exists {
		if time.Now().Before(item.ExpiresAt) {
			if err := json.Unmarshal(item.Data, dest); err == nil {
				c.logger.Debug("Cache hit from fallback", "key", key)
				return nil
			}
		} else {
			// Remove expired item
			delete(c.fallbackCache, key)
		}
	}
	
	c.logger.Debug("Cache miss, executing set function", "key", key)
	
	// Execute the set function to get fresh data
	value, err := setFunc()
	if err != nil {
		return fmt.Errorf("failed to execute set function: %w", err)
	}
	
	// Set in cache
	if err := c.Set(ctx, key, value, expiration); err != nil {
		c.logger.Warn("Failed to set cache", "key", key, "error", err)
	}
	
	// Marshal to dest
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}
	
	return json.Unmarshal(data, dest)
}

// Get retrieves a value from cache
func (c *cacheServiceImpl) Get(ctx context.Context, key string, dest interface{}) error {
	// Try Redis first
	if c.redisCache != nil {
		err := c.redisCache.Get(ctx, key, dest)
		if err == nil {
			return nil
		}
	}
	
	// Try fallback cache
	if item, exists := c.fallbackCache[key]; exists {
		if time.Now().Before(item.ExpiresAt) {
			return json.Unmarshal(item.Data, dest)
		} else {
			delete(c.fallbackCache, key)
		}
	}
	
	return fmt.Errorf("key not found in cache: %s", key)
}

// Set stores a value in cache
func (c *cacheServiceImpl) Set(ctx context.Context, key string, value interface{}, expiration interface{}) error {
	var exp time.Duration
	
	switch v := expiration.(type) {
	case time.Duration:
		exp = v
	case int:
		exp = time.Duration(v) * time.Second
	case int64:
		exp = time.Duration(v) * time.Second
	default:
		exp = 5 * time.Minute // default expiration
	}
	
	// Try to set in Redis
	if c.redisCache != nil {
		if err := c.redisCache.Set(ctx, key, value, exp); err == nil {
			c.logger.Debug("Set cache in Redis", "key", key, "expiration", exp)
			return nil
		} else {
			c.logger.Warn("Failed to set Redis cache", "key", key, "error", err)
		}
	}
	
	// Set in fallback cache
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value for fallback cache: %w", err)
	}
	
	c.fallbackCache[key] = fallbackCacheItem{
		Data:      data,
		ExpiresAt: time.Now().Add(exp),
	}
	
	c.logger.Debug("Set cache in fallback", "key", key, "expiration", exp)
	return nil
}

// Exists checks if a key exists in cache
func (c *cacheServiceImpl) Exists(ctx context.Context, key string) bool {
	// Check Redis first (note: interface is different for Redis cache)
	// For now, we'll skip Redis exists check and use Get for existence checking
	
	// Check fallback cache
	if item, exists := c.fallbackCache[key]; exists {
		if time.Now().Before(item.ExpiresAt) {
			return true
		} else {
			// Remove expired item
			delete(c.fallbackCache, key)
		}
	}
	
	return false
}

// Delete removes a value from cache
func (c *cacheServiceImpl) Delete(ctx context.Context, key string) error {
	// Delete from Redis
	if c.redisCache != nil {
		if err := c.redisCache.Delete(ctx, key); err != nil {
			c.logger.Warn("Failed to delete from Redis cache", "key", key, "error", err)
		}
	}
	
	// Delete from fallback cache
	delete(c.fallbackCache, key)
	
	c.logger.Debug("Deleted from cache", "key", key)
	return nil
}

// Clear clears all cache entries
func (c *cacheServiceImpl) Clear(ctx context.Context) error {
	// Clear Redis
	if c.redisCache != nil {
		if err := c.redisCache.Clear(ctx); err != nil {
			c.logger.Warn("Failed to clear Redis cache", "error", err)
		}
	}
	
	// Clear fallback cache
	c.fallbackCache = make(map[string]fallbackCacheItem)
	
	c.logger.Info("Cleared all cache")
	return nil
}

// HealthCheck performs a health check on the cache service
func (c *cacheServiceImpl) HealthCheck(ctx context.Context) error {
	testKey := "health_check_test"
	testValue := "test_value"
	
	// Test set and get
	if err := c.Set(ctx, testKey, testValue, 10*time.Second); err != nil {
		return fmt.Errorf("cache health check failed on set: %w", err)
	}
	
	var result string
	if err := c.Get(ctx, testKey, &result); err != nil {
		return fmt.Errorf("cache health check failed on get: %w", err)
	}
	
	if result != testValue {
		return fmt.Errorf("cache health check failed: expected %s, got %s", testValue, result)
	}
	
	// Clean up
	c.Delete(ctx, testKey)
	
	return nil
}

// cleanupExpired removes expired items from fallback cache (should be called periodically)
func (c *cacheServiceImpl) cleanupExpired() {
	now := time.Now()
	for key, item := range c.fallbackCache {
		if now.After(item.ExpiresAt) {
			delete(c.fallbackCache, key)
		}
	}
}

// StartCleanupRoutine starts a background routine to clean up expired cache items
func (c *cacheServiceImpl) StartCleanupRoutine() {
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		
		for range ticker.C {
			c.cleanupExpired()
		}
	}()
}