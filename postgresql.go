/*
Package database provides functionality for creating and managing database connections using GORM.

Version: 0.0.1
License: Apache License 2.0

Author: dexterdmonkey
*/

package database

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// PostgreSQL implements the Interface for a PostgreSQL database using GORM.
type PostgreSQL struct {
	*gorm.DB
	*dbLogger
}

// CreatePostgreSQL initializes a new PostgreSQL database connection using the provided configuration.
func CreatePostgreSQL(cfg *Config) (*PostgreSQL, error) {
	if cfg.Timezone == "" {
		cfg.Timezone = "Asia/Jakarta"
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  cfg.DSN(),
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{})

	if err != nil {
		return nil, fmt.Errorf("failed to connect database; %s", err.Error())
	}

	db := &PostgreSQL{DB: gormDB}

	if cfg.MaxConnectionPool <= 0 {
		if err := db.SetMaxConnectionPool(cfg.MaxConnectionPool); err != nil {
			return nil, err
		}
	}

	if err := db.SetMinConnectionPool(cfg.MinConnectionPool); err != nil {
		return nil, err
	}

	return db, nil
}

// SetMaxConnectionPool sets the maximum number of open connections to the database.
// It configures the PostgreSQL database connection to allow up to 'n' concurrent open connections.
//
// Parameters:
//
//	n (int): Maximum number of open connections. Set to 0 or a negative value for unlimited connections.
//
// Returns:
//
//	error: An error if setting the maximum open connections fails.
//
// Example:
//
//	db := database.New(...)
//	err := db.SetMaxConnectionPool(20)
//	if err != nil {
//	    fmt.Println("Error setting max connection pool:", err)
//	}
func (db *PostgreSQL) SetMaxConnectionPool(n int) error {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql db; %s", err.Error())
	}

	sqlDB.SetMaxOpenConns(n)
	return nil
}

// SetMinConnectionPool sets the minimum number of idle connections to the database.
// It configures the PostgreSQL database connection to maintain at least 'n' idle connections when available.
//
// Parameters:
//
//	n (int): Minimum number of idle connections. Set to 0 or a negative value to disable idle connections.
//
// Returns:
//
//	error: An error if setting the minimum idle connections fails.
//
// Example:
//
//	db := database.New(...)
//	err := db.SetMinConnectionPool(5)
//	if err != nil {
//	    fmt.Println("Error setting min connection pool:", err)
//	}
func (db *PostgreSQL) SetMinConnectionPool(n int) error {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql db; %s", err.Error())
	}

	sqlDB.SetMaxIdleConns(n)
	return nil
}

// SetLogger sets a custom logger for the database.
func (db *PostgreSQL) SetLogger(writer logger.Writer) {
	config := logger.Config{
		SlowThreshold:             200 * time.Millisecond,
		Colorful:                  true,
		IgnoreRecordNotFoundError: false,
		LogLevel:                  logger.Warn,
	}

	db.dbLogger = NewLogger(writer, config)
	db.Logger = db.dbLogger
}

// DebugMode sets the logger to debug mode for detailed logging of SQL queries and transactions.
// When enabled, the logger will output detailed information for each SQL query or transaction executed.
// This includes logging SQL statements, execution time, and number of affected rows.
//
// Example:
//
//	db := database.New(...)
//	db.DebugMode()
//
// Notes:
//   - Debug mode should be used primarily for development and debugging purposes.
//   - Enabling debug mode may impact performance due to increased logging overhead.
func (db *PostgreSQL) DebugMode() {
	db.Logger = db.dbLogger.LogMode(logger.Info)
}
