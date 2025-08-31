# RavenDB Go Library

A generic, type-safe RavenDB client library for Go applications. This library provides clean interfaces, type-safe operations, and easy-to-use APIs for working with RavenDB.

## Features

- **Type-safe operations** using Go generics
- **Clean interface design** for easy testing and mocking
- **Flexible querying** with parameterized queries
- **Collection-specific services** for type-safe document operations
- **Generic query operations** that work across document types
- **Connection management** with robust error handling
- **Pagination support** for large result sets

## Installation

```bash
go get github.com/ternarybob/ravendb
```

## Quick Start

```go
package main

import (
    "log"
    ravendb "github.com/ternarybob/ravendb"
)

type User struct {
    ID    string `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

func main() {
    // Create configuration
    config := ravendb.NewLocalConfig("MyDatabase")
    
    // Initialize database
    db, err := ravendb.NewDatabase(config)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()
    
    // Initialize database
    if err := db.Init(); err != nil {
        log.Fatal(err)
    }
    
    // Create a typed collection service
    users := ravendb.NewCollection[User](db, "Users")
    
    // Store a document
    user := User{ID: "users/1", Name: "John Doe", Email: "john@example.com"}
    if err := users.Store("users/1", user); err != nil {
        log.Fatal(err)
    }
    
    // Load a document
    loadedUser, err := users.LoadByID("users/1")
    if err != nil {
        log.Fatal(err)
    }
    
    // Query documents
    results, err := users.QueryAll()
    if err != nil {
        log.Fatal(err)
    }
}
```

## Configuration

### Local Development
```go
config := ravendb.NewLocalConfig("MyDatabase")
```

### Single Node
```go
config := ravendb.NewSingleNodeConfig("http://localhost:8080", "MyDatabase")
```

### Multi-Node Cluster
```go
urls := []string{"http://node1:8080", "http://node2:8080", "http://node3:8080"}
config := ravendb.NewConfig(urls, "MyDatabase")
```

## Usage Examples

### Basic CRUD Operations

```go
// Using the database service directly
db, _ := ravendb.NewDatabase(config)

// Store document
user := User{ID: "users/1", Name: "John"}
db.Store("users/1", user)

// Load document
var loadedUser User
db.LoadByID("users/1", &loadedUser)

// Update document
updates := map[string]interface{}{"name": "John Updated"}
db.Update("users/1", updates)

// Delete document
db.Delete("users/1")
```

### Type-Safe Collection Operations

```go
// Create typed collection service
users := ravendb.NewCollection[User](db, "Users")

// Store multiple documents
userMap := map[string]User{
    "users/1": {ID: "users/1", Name: "John"},
    "users/2": {ID: "users/2", Name: "Jane"},
}
users.StoreMultiple(userMap)

// Query by field
activeUsers, _ := users.QueryByField("isActive", true, nil)

// Query by range
ageRange, _ := users.QueryByRange("age", 18, 65, nil)

// Search across fields
searchResults, _ := users.Search("john", []string{"name", "email"}, nil)
```

### Advanced Querying

```go
// Custom query options
options := &interfaces.QueryOptions{
    Skip:      0,
    Take:      10,
    OrderBy:   "name",
    OrderDesc: false,
    WhereClause: "age > $minAge",
    Parameters: map[string]interface{}{
        "minAge": 18,
    },
}

results, _ := users.Query(options)
```

### Generic Query Functions

```go
// Query any document type generically
products, _ := ravendb.QueryAll[Product](db, "Products")

// Query by field generically
expensiveProducts, _ := ravendb.QueryByField[Product](
    db, "Products", "price", 100.0, nil,
)
```

## Testing

The library includes a comprehensive test suite that covers all functionality with real RavenDB integration tests.

### Prerequisites

1. **RavenDB Server**: You need a running RavenDB instance
   ```bash
   # Download from https://ravendb.net/downloads
   # OR use Docker
   docker run -p 8080:8080 ravendb/ravendb
   ```

2. **Go 1.25+**: Ensure you have Go 1.25 or later installed

### Quick Test Setup

1. **Environment Setup**:
   ```powershell
   # Setup test environment (creates config files)
   .\scripts\test-env-manager.ps1 setup test
   ```

2. **Run Tests**:
   ```powershell
   # Run all tests
   .\scripts\run-tests.ps1
   
   # Run with coverage
   .\scripts\run-tests.ps1 -Coverage
   
   # Run specific tests
   .\scripts\run-tests.ps1 -Pattern "TestDatabase.*"
   ```

### Test Environments

The library supports multiple test environments through TOML configuration:

- **`config/test_config.toml`** - Standard test configuration
- **`config/local_config.toml`** - Local development
- **`config/docker_config.toml`** - Docker environment
- **`config/ci_config.toml`** - CI/CD environment

```powershell
# List available environments
.\scripts\test-env-manager.ps1 list

# Create new environment config
.\scripts\test-env-manager.ps1 create local

# Run tests for specific environment
.\scripts\test-env-manager.ps1 test docker -Coverage
```

### Manual Test Execution

```bash
# Run all tests manually
go test -v ./...

# Run with coverage
go test -v -cover -coverprofile=coverage.out ./...

# Generate HTML coverage report
go tool cover -html=coverage.out -o coverage.html

# Run specific test categories
go test -v -run "TestDatabase.*"    # Database tests
go test -v -run "TestCollection.*"  # Collection tests
```

### Test Categories

- **Connection Tests**: Database connectivity and initialization
- **CRUD Operations**: Store, load, update, delete operations
- **Collection Services**: Type-safe collection operations
- **Query Operations**: All query types (field, range, search, generic)
- **Integration Tests**: End-to-end scenarios with real RavenDB

### Configuration Example

```toml
# config/test_config.toml
[database]
urls = ["http://localhost:8080"]
database = "RavenDBLibTestDB"

[test]
timeout = 30
clean_before_tests = true
clean_after_tests = true
```

For detailed testing information, see [TESTING.md](TESTING.md).

## Architecture

The library is structured around clean interfaces:

- **`interfaces.IRavenDBService`** - Core database operations
- **`interfaces.IRavenCollectionService[T]`** - Type-safe collection operations
- **`services.DatabaseService`** - Database connection management
- **`services.CollectionService[T]`** - Generic collection implementation

## License

MIT License