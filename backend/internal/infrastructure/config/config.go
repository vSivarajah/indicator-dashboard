package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds all configuration settings
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	External ExternalConfig
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port            string
	Host            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration
	Environment     string
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
	MaxConns int
	MinConns int
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

// ExternalConfig holds external API configuration
type ExternalConfig struct {
	CoinGeckoAPIKey     string
	CoinMarketCapAPIKey string
	AlternativeAPI      string
	RateLimitDelay      time.Duration
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	config := &Config{
		Server: ServerConfig{
			Port:            getEnv("PORT", "8080"),
			Host:            getEnv("HOST", "localhost"),
			ReadTimeout:     getDurationEnv("READ_TIMEOUT", 15*time.Second),
			WriteTimeout:    getDurationEnv("WRITE_TIMEOUT", 15*time.Second),
			IdleTimeout:     getDurationEnv("IDLE_TIMEOUT", 60*time.Second),
			ShutdownTimeout: getDurationEnv("SHUTDOWN_TIMEOUT", 10*time.Second),
			Environment:     getEnv("ENVIRONMENT", "development"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "password"),
			DBName:   getEnv("DB_NAME", "crypto_dashboard"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
			MaxConns: getIntEnv("DB_MAX_CONNS", 25),
			MinConns: getIntEnv("DB_MIN_CONNS", 5),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getIntEnv("REDIS_DB", 0),
		},
		External: ExternalConfig{
			CoinGeckoAPIKey:     getEnv("COINGECKO_API_KEY", ""),
			CoinMarketCapAPIKey: getEnv("COINMARKETCAP_API_KEY", "f3ea5727-a012-4b0e-8e81-4d6b515c35e4"),
			AlternativeAPI:      getEnv("ALTERNATIVE_API_URL", "https://api.alternative.me"),
			RateLimitDelay:      getDurationEnv("RATE_LIMIT_DELAY", 100*time.Millisecond),
		},
	}

	return config, nil
}

// GetDSN returns the database connection string
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode)
}

// GetRedisAddr returns the Redis connection address
func (c *RedisConfig) GetRedisAddr() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}

// IsDevelopment returns true if running in development mode
func (c *ServerConfig) IsDevelopment() bool {
	return c.Environment == "development"
}

// IsProduction returns true if running in production mode
func (c *ServerConfig) IsProduction() bool {
	return c.Environment == "production"
}

// Helper functions for environment variable parsing

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getIntEnv(key string, fallback int) int {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil {
			return parsed
		}
	}
	return fallback
}

func getDurationEnv(key string, fallback time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if parsed, err := time.ParseDuration(value); err == nil {
			return parsed
		}
	}
	return fallback
}