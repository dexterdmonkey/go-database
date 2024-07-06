/*
Package database provides a custom logger implementation for GORM, a popular ORM library for Go.
This logger allows customization of logger messages and formatting, including support for colorful output.

Version: 0.0.1
License: Apache License 2.0

Author: dexterdmonkey

This package includes a custom logger that integrates with GORM's logging interface to provide
customizable logging behavior for database operations. It supports different log levels and
formatting options, including colorful output for enhanced readability.

Example usage:

	// Initialize a new logger with custom configurations
	logger := NewLogger(os.Stdout, logger.Config{
	    Colorful: true,
	    LogLevel: logger.Info,
	})

	// Integrate the logger with a GORM database instance
	db := &PostgreSQL{
	    DB: gormDB,
	    DatabaseLogger: logger,
	}

	// Use the database instance with integrated logger for database operations
*/
package database

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm/logger"
)

// dbLogger is a custom logger implementation that integrates with GORM's logging interface.
type dbLogger struct {
	logger.Writer
	logger.Config

	// Format strings for different log levels
	infoStr, warnStr, errStr            string
	traceStr, traceWarnStr, traceErrStr string
}

// NewLogger creates a new instance of the custom database logger with the given writer and configuration.
func NewLogger(writer logger.Writer, config logger.Config) *dbLogger {
	// Customize log message format based on the configuration's Colorful setting
	if config.Colorful {
		return &dbLogger{
			Writer:       writer,
			Config:       config,
			infoStr:      "\033[0m\033[32m[info] %s\033[0m",
			warnStr:      "\033[0m\033[35m[warn] %s\033[0m",
			errStr:       "\033[0m\033[31m[error] %s\033[0m",
			traceStr:     "\033[33m[%.3fms] \033[34;1m[rows:%v]\033[0m %s",
			traceWarnStr: "\033[33m%s \033[0m\033[31;1m[%.3fms] \033[33m[rows:%v]\033[35m %s\033[0m",
			traceErrStr:  "\033[35;1m%s \033[0m\033[33m[%.3fms] \033[34;1m[rows:%v]\033[0m %s",
		}
	} else {
		return &dbLogger{
			Writer:       writer,
			Config:       config,
			infoStr:      "[info] %s",
			warnStr:      "[warn] %s",
			errStr:       "[error] %s",
			traceStr:     "[%.3fms] [rows:%v] %s",
			traceWarnStr: "%s [%.3fms] [rows:%v] %s",
			traceErrStr:  "%s [%.3fms] [rows:%v] %s",
		}
	}
}

// LogMode sets the logger's log level and returns a new logger instance with the updated settings.
func (l *dbLogger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *l
	newLogger.LogLevel = level
	return &newLogger
}

// Info logs an info level message with optional data.
func (l *dbLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Info {
		l.Printf(l.infoStr, fmt.Sprintf(msg, data...))
	}
}

// Warn logs a warning level message with optional data.
func (l *dbLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Warn {
		l.Printf(l.warnStr, fmt.Sprintf(msg, data...))
	}
}

// Error logs an error level message with optional data.
func (l *dbLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Error {
		l.Printf(l.errStr, fmt.Sprintf(msg, data...))
	}
}

// Trace logs detailed information about a database operation, including its duration and parameters.
func (l *dbLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	switch {
	case err != nil && l.LogLevel >= logger.Error:
		sql, rows := fc()
		l.Printf(l.traceErrStr, err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= logger.Warn:
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)
		l.Printf(l.traceWarnStr, slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)
	case l.LogLevel == logger.Info:
		sql, rows := fc()
		l.Printf(l.traceStr, float64(elapsed.Nanoseconds())/1e6, rows, sql)
	}
}

// ParamsFilter filters sensitive parameters from SQL statements if ParameterizedQueries is enabled.
func (l *dbLogger) ParamsFilter(ctx context.Context, sql string, params ...interface{}) (string, []interface{}) {
	if l.Config.ParameterizedQueries {
		return sql, nil
	}
	return sql, params
}
