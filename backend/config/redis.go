package config

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client
var ctx = context.Background()

// ConnectRedis initializes Redis connection
func ConnectRedis() {
	addr := "localhost:6379"
	password := ""
	db := 0 // default DB

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	// Test connection
	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Printf("Redis connection failed: %v", err)
		log.Println("Running without Redis cache")
		RedisClient = nil
		return
	}

	log.Println("Redis connected successfully")
}

// GetRedisContext returns the context for Redis operations
func GetRedisContext() context.Context {
	return ctx
}