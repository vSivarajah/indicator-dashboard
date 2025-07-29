package errors

import (
	"fmt"
	"net/http"
)

// ErrorType represents the type of error
type ErrorType string

const (
	ErrorTypeValidation    ErrorType = "VALIDATION_ERROR"
	ErrorTypeNotFound      ErrorType = "NOT_FOUND"
	ErrorTypeUnauthorized  ErrorType = "UNAUTHORIZED"
	ErrorTypeForbidden     ErrorType = "FORBIDDEN"
	ErrorTypeConflict      ErrorType = "CONFLICT"
	ErrorTypeInternal      ErrorType = "INTERNAL_ERROR"
	ErrorTypeExternal      ErrorType = "EXTERNAL_SERVICE_ERROR"
	ErrorTypeRateLimit     ErrorType = "RATE_LIMIT_ERROR"
	ErrorTypeTimeout       ErrorType = "TIMEOUT_ERROR"
)

// AppError represents an application error
type AppError struct {
	Type       ErrorType `json:"type"`
	Message    string    `json:"message"`
	Details    string    `json:"details,omitempty"`
	StatusCode int       `json:"-"`
	Cause      error     `json:"-"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (caused by: %v)", e.Type, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

// Unwrap returns the cause of the error
func (e *AppError) Unwrap() error {
	return e.Cause
}

// New creates a new AppError
func New(errorType ErrorType, message string) *AppError {
	return &AppError{
		Type:       errorType,
		Message:    message,
		StatusCode: getStatusCode(errorType),
	}
}

// Wrap wraps an existing error with additional context
func Wrap(err error, errorType ErrorType, message string) *AppError {
	return &AppError{
		Type:       errorType,
		Message:    message,
		StatusCode: getStatusCode(errorType),
		Cause:      err,
	}
}

// Validation creates a validation error
func Validation(message string, details ...string) *AppError {
	err := &AppError{
		Type:       ErrorTypeValidation,
		Message:    message,
		StatusCode: http.StatusBadRequest,
	}
	if len(details) > 0 {
		err.Details = details[0]
	}
	return err
}

// NotFound creates a not found error
func NotFound(resource string) *AppError {
	return &AppError{
		Type:       ErrorTypeNotFound,
		Message:    fmt.Sprintf("%s not found", resource),
		StatusCode: http.StatusNotFound,
	}
}

// Unauthorized creates an unauthorized error
func Unauthorized(message string) *AppError {
	return &AppError{
		Type:       ErrorTypeUnauthorized,
		Message:    message,
		StatusCode: http.StatusUnauthorized,
	}
}

// Forbidden creates a forbidden error
func Forbidden(message string) *AppError {
	return &AppError{
		Type:       ErrorTypeForbidden,
		Message:    message,
		StatusCode: http.StatusForbidden,
	}
}

// Conflict creates a conflict error
func Conflict(message string) *AppError {
	return &AppError{
		Type:       ErrorTypeConflict,
		Message:    message,
		StatusCode: http.StatusConflict,
	}
}

// Internal creates an internal server error
func Internal(message string, cause error) *AppError {
	return &AppError{
		Type:       ErrorTypeInternal,
		Message:    message,
		StatusCode: http.StatusInternalServerError,
		Cause:      cause,
	}
}

// External creates an external service error
func External(service, message string, cause error) *AppError {
	return &AppError{
		Type:       ErrorTypeExternal,
		Message:    fmt.Sprintf("%s service error: %s", service, message),
		StatusCode: http.StatusBadGateway,
		Cause:      cause,
	}
}

// RateLimit creates a rate limit error
func RateLimit(message string) *AppError {
	return &AppError{
		Type:       ErrorTypeRateLimit,
		Message:    message,
		StatusCode: http.StatusTooManyRequests,
	}
}

// Timeout creates a timeout error
func Timeout(operation string) *AppError {
	return &AppError{
		Type:       ErrorTypeTimeout,
		Message:    fmt.Sprintf("%s operation timed out", operation),
		StatusCode: http.StatusRequestTimeout,
	}
}

// getStatusCode returns the HTTP status code for an error type
func getStatusCode(errorType ErrorType) int {
	switch errorType {
	case ErrorTypeValidation:
		return http.StatusBadRequest
	case ErrorTypeNotFound:
		return http.StatusNotFound
	case ErrorTypeUnauthorized:
		return http.StatusUnauthorized
	case ErrorTypeForbidden:
		return http.StatusForbidden
	case ErrorTypeConflict:
		return http.StatusConflict
	case ErrorTypeExternal:
		return http.StatusBadGateway
	case ErrorTypeRateLimit:
		return http.StatusTooManyRequests
	case ErrorTypeTimeout:
		return http.StatusRequestTimeout
	case ErrorTypeInternal:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

// IsType checks if an error is of a specific type
func IsType(err error, errorType ErrorType) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Type == errorType
	}
	return false
}

// NewNotFoundError creates a new not found error with entity and identifier
func NewNotFoundError(entity, identifier string) *AppError {
	return &AppError{
		Type:       ErrorTypeNotFound,
		Message:    fmt.Sprintf("%s not found: %s", entity, identifier),
		StatusCode: http.StatusNotFound,
		Details:    fmt.Sprintf("entity: %s, identifier: %s", entity, identifier),
	}
}

// GetStatusCode extracts the HTTP status code from an error
func GetStatusCode(err error) int {
	if appErr, ok := err.(*AppError); ok {
		return appErr.StatusCode
	}
	return http.StatusInternalServerError
}