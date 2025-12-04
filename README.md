# Anonymous Query Scan

[![Go](https://img.shields.io/badge/go-1.22+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Go Reference](https://pkg.go.dev/badge/github.com/naughtyGitCat/anonymous-query-scan.svg)](https://pkg.go.dev/github.com/naughtyGitCat/anonymous-query-scan)

A Go library that solves MySQL anonymous query scanning issues with proper type conversion and marshalling.

## 🚀 Overview

The `go-sql-driver/mysql` driver has limitations when it comes to scanning query results without predefined structs. This library provides enhanced anonymous query scanning capabilities with intelligent type conversion, particularly for MySQL data types.

## ✨ Features

- 🔍 **Anonymous Query Scanning**: Scan SQL query results without predefined structs
- 🕒 **Smart Time Handling**: Proper conversion of MySQL time-related types (DATETIME, DATE, TIMESTAMP)
- 🗺️ **Multiple Output Formats**: Support for both slice and map-based result formats
- 🌍 **Timezone Aware**: Configurable timezone handling (UTC vs Local time)
- 📦 **JSON Support**: Enhanced JSON field handling
- 🔧 **Type Safe**: Intelligent type conversion based on MySQL column types

## 📦 Installation

```bash
go get github.com/naughtyGitCat/anonymous-query-scan/mysql
```

## 🛠️ Quick Start

```go
package main

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "log"

    _ "github.com/go-sql-driver/mysql"
    mysql "github.com/naughtyGitCat/anonymous-query-scan/mysql"
)

func main() {
    // Connect to database
    db, err := sql.Open("mysql", "user:password@tcp(localhost:3306)/database")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // Execute query
    rows, err := db.Query("SELECT id, name, created_at FROM users LIMIT 5")
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()

    // Scan results as map (recommended)
    mappedRows, err := mysql.ScanAnonymousMappedRows(rows)
    if err != nil {
        log.Fatal(err)
    }

    // Convert to JSON
    jsonBytes, err := json.MarshalIndent(mappedRows, "", "  ")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(string(jsonBytes))
}
```

## 📖 API Reference

### Core Functions

#### `ScanAnonymousMappedRows(rows *sql.Rows) ([]map[string]any, error)`

Scans query results into a slice of maps with intelligent MySQL type conversion.

**Features:**
- Time-related columns converted to local time
- TIMESTAMP columns converted to numbers (Unix timestamp)
- JSON columns properly parsed
- Column names as map keys

**Return format:**
```go
[
  {"id": 1, "name": "John", "created_at": "2024-01-01T10:00:00+08:00"},
  {"id": 2, "name": "Jane", "created_at": "2024-01-02T11:30:00+08:00"}
]
```

#### `ScanAnonymousRows(rows *sql.Rows) ([][]any, error)`

Scans query results into a slice of slices with type conversion.

**Return format:**
```go
[
  [1, "John", "2024-01-01T10:00:00Z"],
  [2, "Jane", "2024-01-02T11:30:00Z"]
]
```

### Deprecated Functions

#### `DeprecatedScanAnonymousMappedRows(rows *sql.Rows) ([]map[string]any, error)`

Legacy function with UTC time conversion and string JSON handling.

#### `DeprecatedScanAnonymousRows(rows *sql.Rows) ([][]any, error)`

Legacy function with UTC time conversion and basic type handling.

## 🔄 Type Conversion

| MySQL Type | Go Type | Notes |
|------------|---------|--------|
| TINYINT, SMALLINT, INT | `int64` | |
| BIGINT | `int64` | |
| FLOAT, DOUBLE | `float64` | |
| DECIMAL | `string` | Preserves precision |
| VARCHAR, TEXT | `string` | |
| DATE, DATETIME | `time.Time` | Local timezone |
| TIMESTAMP | `int64` | Unix timestamp (new) / `time.Time` (deprecated) |
| JSON | `interface{}` | Parsed JSON object |
| BLOB, BINARY | `[]byte` | |

## 🔧 Advanced Usage

### Custom Error Handling

```go
mappedRows, err := mysql.ScanAnonymousMappedRows(rows)
if err != nil {
    switch {
    case strings.Contains(err.Error(), "duplicate column name"):
        // Handle duplicate column names
        log.Printf("Query contains duplicate column names: %v", err)
    case strings.Contains(err.Error(), "convert value failed"):
        // Handle type conversion errors
        log.Printf("Type conversion failed: %v", err)
    default:
        log.Printf("Scanning failed: %v", err)
    }
    return
}
```

### Working with Different Databases

While optimized for MySQL, the library also supports basic functionality with other SQL databases:

```go
// SQLite example
import _ "github.com/mattn/go-sqlite3"

db, err := sql.Open("sqlite3", "./test.db")
// ... use the same scanning functions
```

## ⚡ Performance Considerations

- Use `ScanAnonymousMappedRows` for most use cases as it provides better data access
- Use `ScanAnonymousRows` when memory usage is critical and you don't need column names
- The library scans all rows into memory - consider pagination for large result sets
- Type conversion is performed on each value - cache results when possible

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [go-sql-driver/mysql](https://github.com/go-sql-driver/mysql) - The underlying MySQL driver
- [Grafana Plugin SDK](https://github.com/grafana/grafana-plugin-sdk-go) - Inspiration for type conversion patterns

## 📚 References

- [MySQL 8.4 Data Types](https://dev.mysql.com/doc/refman/8.4/en/data-types.html)
- [Go SQL Driver Examples](https://github.com/go-sql-driver/mysql/wiki/Examples)
- [MySQL Driver Issue #407](https://github.com/go-sql-driver/mysql/issues/407) - Context for this library's creation
- [Go Reflection Deep Dive](https://stackoverflow.com/questions/55995917/reverse-of-reflect-typeof)
- [Converting SQL Rows to Typed JSON in Go](https://stackoverflow.com/questions/42774467/how-to-convert-sql-rows-to-typed-json-in-golang)
- [Marshal Byte Array as JSON Array in Go](https://stackoverflow.com/questions/14177862/how-to-marshal-a-byte-uint8-array-as-json-array-in-go/78662958#78662958)
- [Marshal Byte to JSON Issues](https://stackoverflow.com/questions/34089750/marshal-byte-to-json-giving-a-strange-string)
