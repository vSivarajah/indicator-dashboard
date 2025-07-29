package errors

import (
	"fmt"
	"net/http"
	"time"
)

// IndicatorError represents errors specific to indicator calculations
type IndicatorError struct {
	Code        string            `json:"code"`
	Message     string            `json:"message"`
	Details     map[string]interface{} `json:"details,omitempty"`
	Timestamp   time.Time         `json:"timestamp"`
	StatusCode  int               `json:"status_code"`
	Retryable   bool              `json:"retryable"`
	Component   string            `json:"component"`
}

func (e *IndicatorError) Error() string {
	return fmt.Sprintf("[%s] %s: %s", e.Component, e.Code, e.Message)
}

// IsRetryable returns whether the error is retryable
func (e *IndicatorError) IsRetryable() bool {
	return e.Retryable
}

// GetStatusCode returns the HTTP status code for the error
func (e *IndicatorError) GetStatusCode() int {
	return e.StatusCode
}

// Indicator error codes
const (
	ErrCodeDataFetch        = "DATA_FETCH_ERROR"
	ErrCodeCalculation      = "CALCULATION_ERROR"
	ErrCodeInvalidInput     = "INVALID_INPUT"
	ErrCodeServiceUnavail   = "SERVICE_UNAVAILABLE"
	ErrCodeThreshold        = "THRESHOLD_ERROR"
	ErrCodeCacheError       = "CACHE_ERROR"
	ErrCodeDatabaseError    = "DATABASE_ERROR"
	ErrCodeAPIError         = "API_ERROR"
	ErrCodeRateLimit        = "RATE_LIMIT_ERROR"
	ErrCodeTimeout          = "TIMEOUT_ERROR"
)

// NewIndicatorError creates a new indicator error
func NewIndicatorError(code, component, message string) *IndicatorError {
	return &IndicatorError{
		Code:        code,
		Component:   component,
		Message:     message,
		Timestamp:   time.Now(),
		StatusCode:  http.StatusInternalServerError,
		Retryable:   false,
		Details:     make(map[string]interface{}),
	}
}

// MVRV specific errors
func NewMVRVDataFetchError(source string, err error) *IndicatorError {
	return &IndicatorError{
		Code:       ErrCodeDataFetch,
		Component:  "mvrv_service",
		Message:    fmt.Sprintf("Failed to fetch MVRV data from %s", source),
		StatusCode: http.StatusServiceUnavailable,
		Retryable:  true,
		Timestamp:  time.Now(),
		Details: map[string]interface{}{
			"source":     source,
			"error":      err.Error(),
			"retry_after": 300, // 5 minutes
		},
	}
}

func NewMVRVCalculationError(reason string) *IndicatorError {
	return &IndicatorError{
		Code:       ErrCodeCalculation,
		Component:  "mvrv_service",
		Message:    fmt.Sprintf("MVRV calculation failed: %s", reason),
		StatusCode: http.StatusUnprocessableEntity,
		Retryable:  false,
		Timestamp:  time.Now(),
		Details: map[string]interface{}{
			"reason": reason,
		},
	}
}

// Dominance specific errors
func NewDominanceDataError(err error) *IndicatorError {
	return &IndicatorError{
		Code:       ErrCodeDataFetch,
		Component:  "dominance_service",
		Message:    "Failed to fetch Bitcoin dominance data",
		StatusCode: http.StatusServiceUnavailable,
		Retryable:  true,
		Timestamp:  time.Now(),
		Details: map[string]interface{}{
			"error":      err.Error(),
			"retry_after": 180, // 3 minutes
		},
	}
}

// Fear & Greed specific errors
func NewFearGreedAPIError(statusCode int, response string) *IndicatorError {
	return &IndicatorError{
		Code:       ErrCodeAPIError,
		Component:  "fear_greed_service",
		Message:    "Fear & Greed API request failed",
		StatusCode: http.StatusBadGateway,
		Retryable:  statusCode >= 500,
		Timestamp:  time.Now(),
		Details: map[string]interface{}{
			"api_status_code": statusCode,
			"api_response":    response,
			"retry_after":     120, // 2 minutes
		},
	}
}

// Cache specific errors
func NewCacheError(operation, key string, err error) *IndicatorError {
	return &IndicatorError{
		Code:       ErrCodeCacheError,
		Component:  "cache_service",
		Message:    fmt.Sprintf("Cache %s operation failed for key: %s", operation, key),
		StatusCode: http.StatusInternalServerError,
		Retryable:  true,
		Timestamp:  time.Now(),
		Details: map[string]interface{}{
			"operation": operation,
			"key":       key,
			"error":     err.Error(),
		},
	}
}

// Database specific errors
func NewDatabaseError(operation, entity string, err error) *IndicatorError {
	return &IndicatorError{
		Code:       ErrCodeDatabaseError,
		Component:  "database",
		Message:    fmt.Sprintf("Database %s operation failed for %s", operation, entity),
		StatusCode: http.StatusInternalServerError,
		Retryable:  false,
		Timestamp:  time.Now(),
		Details: map[string]interface{}{
			"operation": operation,
			"entity":    entity,
			"error":     err.Error(),
		},
	}
}

// Rate limit errors
func NewRateLimitError(service string, resetTime time.Time) *IndicatorError {
	return &IndicatorError{
		Code:       ErrCodeRateLimit,
		Component:  service,
		Message:    fmt.Sprintf("Rate limit exceeded for %s", service),
		StatusCode: http.StatusTooManyRequests,
		Retryable:  true,
		Timestamp:  time.Now(),
		Details: map[string]interface{}{
			"service":    service,
			"reset_time": resetTime.Unix(),
			"retry_after": int(time.Until(resetTime).Seconds()),
		},
	}
}

// Timeout errors
func NewTimeoutError(component, operation string, duration time.Duration) *IndicatorError {
	return &IndicatorError{
		Code:       ErrCodeTimeout,
		Component:  component,
		Message:    fmt.Sprintf("Operation '%s' timed out after %v", operation, duration),
		StatusCode: http.StatusRequestTimeout,
		Retryable:  true,
		Timestamp:  time.Now(),
		Details: map[string]interface{}{
			"operation": operation,
			"timeout":   duration.String(),
		},
	}
}

// Validation errors
func NewValidationError(component, field string, value interface{}) *IndicatorError {
	return &IndicatorError{
		Code:       ErrCodeInvalidInput,
		Component:  component,
		Message:    fmt.Sprintf("Invalid value for field '%s'", field),
		StatusCode: http.StatusBadRequest,
		Retryable:  false,
		Timestamp:  time.Now(),
		Details: map[string]interface{}{
			"field": field,
			"value": value,
		},
	}
}

// Service unavailable errors
func NewServiceUnavailableError(service, reason string) *IndicatorError {
	return &IndicatorError{
		Code:       ErrCodeServiceUnavail,
		Component:  service,
		Message:    fmt.Sprintf("Service %s is currently unavailable: %s", service, reason),
		StatusCode: http.StatusServiceUnavailable,
		Retryable:  true,
		Timestamp:  time.Now(),
		Details: map[string]interface{}{
			"service": service,
			"reason":  reason,
			"retry_after": 600, // 10 minutes
		},
	}
}

// Helper functions for error checking
func IsIndicatorError(err error) bool {
	_, ok := err.(*IndicatorError)
	return ok
}

func IsRetryableError(err error) bool {
	if indErr, ok := err.(*IndicatorError); ok {
		return indErr.IsRetryable()
	}
	return false
}

func GetErrorStatusCode(err error) int {
	if indErr, ok := err.(*IndicatorError); ok {
		return indErr.GetStatusCode()
	}
	return http.StatusInternalServerError
}

// Error details extraction
func GetErrorDetails(err error) map[string]interface{} {
	if indErr, ok := err.(*IndicatorError); ok {
		return indErr.Details
	}
	return map[string]interface{}{
		"error": err.Error(),
	}
}

// Error wrapping with context
func WrapError(err error, component, operation string) *IndicatorError {
	if indErr, ok := err.(*IndicatorError); ok {
		// Add context to existing indicator error
		indErr.Details["wrapped_from"] = component
		indErr.Details["operation"] = operation
		return indErr
	}
	
	// Create new indicator error from generic error
	return &IndicatorError{
		Code:       ErrCodeServiceUnavail,
		Component:  component,
		Message:    fmt.Sprintf("Operation '%s' failed", operation),
		StatusCode: http.StatusInternalServerError,
		Retryable:  false,
		Timestamp:  time.Now(),
		Details: map[string]interface{}{
			"operation":     operation,
			"original_error": err.Error(),
		},
	}
}