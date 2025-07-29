package logger

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"
	"gorm.io/gorm/logger"
)

// Logger defines the logging interface
type Logger interface {
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	With(args ...interface{}) Logger
	WithContext(ctx context.Context) Logger
}

// slogLogger implements Logger using slog
type slogLogger struct {
	logger *slog.Logger
}

// New creates a new logger instance
func New(environment string) Logger {
	var handler slog.Handler
	
	if environment == "production" {
		// JSON handler for production
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
			AddSource: true,
		})
	} else {
		// Text handler for development
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
			AddSource: true,
		})
	}
	
	return &slogLogger{
		logger: slog.New(handler),
	}
}

// Debug logs a debug message
func (l *slogLogger) Debug(msg string, args ...interface{}) {
	l.logger.Debug(msg, args...)
}

// Info logs an info message
func (l *slogLogger) Info(msg string, args ...interface{}) {
	l.logger.Info(msg, args...)
}

// Warn logs a warning message
func (l *slogLogger) Warn(msg string, args ...interface{}) {
	l.logger.Warn(msg, args...)
}

// Error logs an error message
func (l *slogLogger) Error(msg string, args ...interface{}) {
	l.logger.Error(msg, args...)
}

// With adds structured context to the logger
func (l *slogLogger) With(args ...interface{}) Logger {
	return &slogLogger{
		logger: l.logger.With(args...),
	}
}

// WithContext adds context to the logger
func (l *slogLogger) WithContext(ctx context.Context) Logger {
	// Extract trace ID or other context values if needed
	return l
}

// GormLogger implements gorm.logger.Interface
type GormLogger struct {
	logger Logger
}

// NewGormLogger creates a new GORM logger
func NewGormLogger(logger Logger) logger.Interface {
	return &GormLogger{
		logger: logger,
	}
}

// LogMode sets the log level (for gorm.logger.Interface)
func (l *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	return l
}

// Info logs info messages (for gorm.logger.Interface)
func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	l.logger.Info(fmt.Sprintf(msg, data...))
}

// Warn logs warning messages (for gorm.logger.Interface)
func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	l.logger.Warn(fmt.Sprintf(msg, data...))
}

// Error logs error messages (for gorm.logger.Interface)
func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	l.logger.Error(fmt.Sprintf(msg, data...))
}

// Trace logs SQL queries (for gorm.logger.Interface)
func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()
	
	if err != nil {
		l.logger.Error("SQL Error",
			"error", err,
			"elapsed", elapsed,
			"rows", rows,
			"sql", sql,
		)
	} else {
		l.logger.Debug("SQL Query",
			"elapsed", elapsed,
			"rows", rows,
			"sql", sql,
		)
	}
}