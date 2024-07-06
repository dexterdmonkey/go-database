/*
Package database provides types and functions for managing database configurations.

Version: 0.0.1
License: Apache License 2.0

Author: dexterdmonkey

This package defines the Config struct for configuring database connections,
providing fields for essential database connection parameters such as host, port,
user credentials, database name, connection pool sizes, and timezone.

Example usage:

    cfg := database.Config{
        Host:              "localhost",
        Port:              5432,
        User:              "user",
        Pass:              "password",
        Name:              "mydatabase",
        MaxConnectionPool: 10,
        MinConnectionPool: 2,
        Timezone:          "Asia/Jakarta",
    }

    fmt.Println(cfg.DSN()) // Output: "user=user password=password dbname=mydatabase port=5432 host=localhost sslmode=disable TimeZone=Asia/Jakarta"

    fmt.Println(cfg.String()) // Output: "user=user password=password dbname=mydatabase port=5432 host=localhost min-pool=2 max-pool=10"

*/

package database

import "fmt"

// Config holds configuration parameters for connecting to a database.
type Config struct {
	Host              string // Database host address.
	Port              int    // Database port number.
	User              string // Database user name.
	Pass              string // Database password.
	Name              string // Database name.
	MaxConnectionPool int    // Maximum size of the connection pool. Set to <= 0 for unlimited connections. Default is 0.
	MinConnectionPool int    // Minimum size of the connection pool. Set to <= 0 for no connection pooling. Default is 0.
	Timezone          string // Timezone of the database server. Default is "Asia/Jakarta".
}

// String returns a formatted string representation of the Config, including connection details and pool settings.
func (cfg Config) String() string {
	return fmt.Sprintf(
		"user=%s password=%s dbname=%s port=%d host=%s min-pool=%d max-pool=%d",
		cfg.User, cfg.Pass, cfg.Name, cfg.Port, cfg.Host, cfg.MinConnectionPool, cfg.MaxConnectionPool,
	)
}

// DSN returns the Data Source Name (DSN) string used for connecting to the database.
func (cfg Config) DSN() string {
	return fmt.Sprintf(
		"user=%s password=%s dbname=%s port=%d host=%s sslmode=disable TimeZone=%s",
		cfg.User, cfg.Pass, cfg.Name, cfg.Port, cfg.Host, cfg.Timezone,
	)
}
