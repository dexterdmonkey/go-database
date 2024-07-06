# Database Package

[![Version](https://img.shields.io/badge/version-0.0.1-blue.svg)]()
[![License](https://img.shields.io/badge/license-Apache%202.0-green.svg)]()

This package provides functionality for creating and managing database connections using GORM (Go Object Relational Mapper).

## Features

- **PostgreSQL Implementation**: Implements an interface for PostgreSQL database connections using GORM.
- **Connection Pool Management**: Methods to set maximum and minimum connection pools.
- **Custom Logging**: Customizable logging setup, including debug mode for detailed query and transaction logs.

## Installation

```bash
go get github.com/dexterdmonkey/go-database
```

## Usage

### Initialize PostgreSQL Database

```go
package main

import (
	"fmt"
	"time"

	"github.com/dexterdmonkey/go-database"
)

func main() {
	cfg := &database.Config{
		Host:              "localhost",
		Port:              5432,
		User:              "username",
		Pass:              "password",
		Name:              "dbname",
		MaxConnectionPool: 20,
		MinConnectionPool: 5,
		Timezone:          "Asia/Jakarta",
	}

	db, err := database.CreatePostgreSQL(cfg)
	if err != nil {
		fmt.Println("Error creating PostgreSQL connection:", err)
		return
	}

	defer db.Close()

	// Use the database connection...
}
```

### Set Connection Pool Limits

```go
// Set maximum connection pool limit
err := db.SetMaxConnectionPool(20)
if err != nil {
	fmt.Println("Error setting max connection pool:", err)
}

// Set minimum connection pool
err = db.SetMinConnectionPool(5)
if err != nil {
	fmt.Println("Error setting min connection pool:", err)
}
```

### Enable Debug Mode

```go
// Enable debug mode for detailed logging
db.DebugMode()

// Execute queries or transactions...
```

## API Reference

### `CreatePostgreSQL(cfg *Config) (*PostgreSQL, error)`

Creates a new PostgreSQL database connection.

- `cfg`: Configuration parameters including database credentials and connection settings.

### `SetMaxConnectionPool(n int) error`

Sets the maximum number of open connections to the database.

- `n`: Maximum number of open connections. Set to 0 or a negative value for unlimited connections.

### `SetMinConnectionPool(n int) error`

Sets the minimum number of idle connections to the database.

- `n`: Minimum number of idle connections. Set to 0 or a negative value to disable idle connections.

### `SetLogger(writer logger.Writer)`

Sets a custom logger for the database.

- `writer`: Custom logger writer implementing the `gorm.io/gorm/logger.Writer` interface.

### `DebugMode()`

Enables debug mode for detailed logging of SQL queries and transactions.

- Enables detailed logging including SQL statements, execution time, and affected rows.

## License

This package is licensed under the Apache License 2.0. See the [LICENSE](LICENSE) file for details.

## Author

Written by dexterdmonkey.
```
Feel free to customize and expand the README further based on additional features, examples, or specific usage scenarios of your package.