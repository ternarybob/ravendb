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

## Architecture

The library is structured around clean interfaces:

- **`interfaces.IRavenDBService`** - Core database operations
- **`interfaces.IRavenCollectionService[T]`** - Type-safe collection operations
- **`services.DatabaseService`** - Database connection management
- **`services.CollectionService[T]`** - Generic collection implementation

## License

MIT License