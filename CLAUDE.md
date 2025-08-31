# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a generic RavenDB connection and operational library for Go applications. It provides type-safe, interface-driven operations for RavenDB document database integration.

## Architecture

### Core Components
- **`interfaces/`** - Clean interface definitions for database and collection operations
- **`services/`** - Concrete implementations of the interfaces
- **`examples/`** - Usage examples and demonstrations
- **Root package** - Main library entry points and configuration

### Key Interfaces
- `IRavenDBService` - Core database operations interface
- `IRavenCollectionService[T]` - Generic, type-safe collection operations interface

### Key Services
- `DatabaseService` - Database connection and lifecycle management
- `CollectionService[T]` - Generic collection operations implementation

## Development Commands

### Build and Test
```bash
# Build the library
go build ./...

# Run all tests
go test ./...

# Run specific test suites
go test -v -run "TestDatabase.*"    # Database tests
go test -v -run "TestCollection.*"  # Collection tests

# Run tests with coverage
go test -cover -coverprofile=coverage.out ./...

# Generate HTML coverage report
go tool cover -html=coverage.out -o coverage.html

# Format code
go fmt ./...

# Run linter (if golangci-lint is installed)
golangci-lint run
```

### Test Execution
```powershell
# Environment management
.\scripts\test-env-manager.ps1 list                    # List available environments
.\scripts\test-env-manager.ps1 setup test              # Setup test environment  
.\scripts\test-env-manager.ps1 test docker             # Run tests for specific environment

# Direct test execution
.\scripts\run-tests.ps1                                 # Run with default config
.\scripts\run-tests.ps1 -Coverage                      # Run with coverage
.\scripts\run-tests.ps1 -ConfigFile "config\local_config.toml"  # Custom config

# Manual test execution
go test -v ./... -timeout 60s
```

### Dependencies
```bash
# Download dependencies
go mod download

# Tidy dependencies
go mod tidy

# Update dependencies
go get -u ./...
```

## Design Principles

- **Interface-driven design** - All major components implement clean interfaces
- **Type safety** - Use Go generics for type-safe document operations
- **Separation of concerns** - Database management separate from document operations
- **Testability** - Interfaces enable easy mocking and testing
- **Simplicity** - Keep the API simple and intuitive

## Usage Patterns

### Database Initialization
```go
config := ravendb.NewLocalConfig("MyDatabase")
db, err := ravendb.NewDatabase(config)
if err != nil {
    log.Fatal(err)
}
defer db.Close()

if err := db.Init(); err != nil {
    log.Fatal(err)
}
```

### Type-Safe Collections
```go
// Create typed collection service
users := ravendb.NewCollection[User](db, "Users")

// Use strongly-typed operations
user, err := users.LoadByID("users/1")
results, err := users.QueryAll()
```

### Generic Queries
```go
// Query any document type
products, err := ravendb.QueryAll[Product](db, "Products")
users, err := ravendb.QueryByField[User](db, "Users", "isActive", true, nil)
```

## Testing Strategy

### Test Configuration
Tests use TOML configuration files in the `config/` directory:
- `config/test_config.toml` - Standard test configuration
- `config/local_config.toml` - Local development 
- `config/docker_config.toml` - Docker environment
- `config/ci_config.toml` - CI/CD environment

Example configuration:
```toml
[database]
urls = ["http://localhost:8080"]
database = "RavenDBLibTestDB"

[test]
timeout = 30
clean_before_tests = true
clean_after_tests = true
```

### Test Categories
- **Connection Tests**: Verify RavenDB server connectivity and database creation
- **CRUD Tests**: Basic document store, load, update, delete operations
- **Collection Tests**: Type-safe collection service operations
- **Query Tests**: All query patterns (all, by field, by range, search, generic queries)
- **Integration Tests**: End-to-end scenarios with real RavenDB instance

### Test Requirements
- Running RavenDB server (local or remote)
- Valid `test_config.toml` configuration
- Test database (created automatically if not exists)

### Coverage Goals
- Minimum: 80% code coverage
- Target: 90%+ code coverage
- All public interfaces tested
- Error handling paths covered

## Code Organization

- Keep interfaces simple and focused
- Implement concrete types in services package
- Provide convenience functions in root package
- Use examples package for documentation and testing
- Follow Go naming conventions and patterns