# Testing Guide

This document describes how to run and understand the tests for the RavenDB Go library.

## Prerequisites

1. **RavenDB Server**: You need a running RavenDB instance
   - Download from: https://ravendb.net/downloads
   - Default URL: `http://localhost:8080`
   - Or use Docker: `docker run -p 8080:8080 ravendb/ravendb`

2. **Go Environment**: Go 1.25 or later
3. **Test Configuration**: `test_config.toml` file (see configuration section below)

## Test Configuration

Test configurations are stored in the `config/` directory. Several predefined configurations are available:

- **`config/test_config.toml`** - Standard test configuration
- **`config/local_config.toml`** - Local development environment
- **`config/docker_config.toml`** - Docker containerized environment  
- **`config/ci_config.toml`** - Continuous Integration environment

### Example Configuration (`config/test_config.toml`):
```toml
[database]
urls = ["http://localhost:8080"]
database = "RavenDBLibTestDB"

[test]
timeout = 30
clean_before_tests = true
clean_after_tests = true
```

### Configuration Options

- **`database.urls`**: Array of RavenDB server URLs
- **`database.database`**: Test database name (will be created if it doesn't exist)
- **`test.timeout`**: Test timeout in seconds
- **`test.clean_before_tests`**: Whether to clean the database before running tests
- **`test.clean_after_tests`**: Whether to clean the database after running tests

## Running Tests

### Quick Start (PowerShell)

```powershell
# List available test environments
.\scripts\test-env-manager.ps1 list

# Setup test environment
.\scripts\test-env-manager.ps1 setup test

# Run tests with default configuration
.\scripts\run-tests.ps1

# Run tests for specific environment
.\scripts\test-env-manager.ps1 test docker

# Run tests with coverage
.\scripts\run-tests.ps1 -Coverage
```

### Manual Test Execution

```powershell
# Run all tests
go test -v ./...

# Run specific test categories
go test -v -run "TestDatabase.*"    # Database connection tests
go test -v -run "TestCollection.*"  # Collection service tests

# Run with coverage
go test -v -cover -coverprofile=coverage.out ./...

# Generate coverage report
go tool cover -html=coverage.out -o coverage.html

# PowerShell script options
.\scripts\run-tests.ps1 -Pattern "TestDatabase.*"  # Specific tests
.\scripts\run-tests.ps1 -ConfigFile "config\local_config.toml"  # Custom config
.\scripts\run-tests.ps1 -Timeout 120  # Custom timeout
```

## Test Categories

### 1. Database Connection Tests (`TestDatabaseConnection`)
- **Purpose**: Verify RavenDB server connectivity
- **Tests**:
  - Database service creation
  - Database initialization
  - Connection status verification

### 2. Database Availability Tests (`TestDatabaseAvailability`)
- **Purpose**: Ensure database is accessible and operational
- **Tests**:
  - Database name retrieval
  - Document store access
  - Basic database operations

### 3. Basic CRUD Tests (`TestBasicCRUDOperations`)
- **Purpose**: Test fundamental document operations
- **Tests**:
  - **Store and Load**: Create and retrieve documents
  - **Document Exists**: Check document existence
  - **Update Document**: Modify existing documents
  - **Store Multiple**: Bulk document creation
  - **Delete Document**: Single document removal
  - **Delete Multiple**: Bulk document removal

### 4. Collection Service Tests (`TestCollectionService`)
- **Purpose**: Test type-safe collection operations
- **Tests**:
  - **Typed CRUD**: Store/load with strong typing
  - **Multiple Operations**: Bulk typed operations
  - **Update Operations**: Modify typed documents
  - **Query All**: Retrieve all documents in collection
  - **Query by Field**: Filter by specific field values
  - **Query by Range**: Filter by value ranges
  - **Search**: Full-text search across fields
  - **Query with Options**: Custom query parameters
  - **Count**: Document counting

### 5. Generic Query Tests (`TestGenericQueryOperations`)
- **Purpose**: Test generic query functions across document types
- **Tests**:
  - **Generic Query All**: Query any document type
  - **Generic Query by Field**: Type-safe field filtering
  - **Generic Query by Range**: Type-safe range queries
  - **Generic Search**: Type-safe text search

## Test Data

The tests use two main document types:

### TestUser
```go
type TestUser struct {
    ID       string    `json:"id"`
    Name     string    `json:"name"`
    Email    string    `json:"email"`
    Age      int       `json:"age"`
    IsActive bool      `json:"isActive"`
    Created  time.Time `json:"created"`
}
```

### TestProduct
```go
type TestProduct struct {
    ID          string  `json:"id"`
    Name        string  `json:"name"`
    Description string  `json:"description"`
    Price       float64 `json:"price"`
    Category    string  `json:"category"`
    InStock     bool    `json:"inStock"`
}
```

## Test Scenarios

### Connection and Database Tests
1. **Connection Test**: Verifies the library can connect to RavenDB
2. **Database Creation**: Tests automatic database creation if it doesn't exist
3. **Status Check**: Validates database status information

### Document Operations Tests
1. **Insert Documents**: Store various document types
2. **Access Documents**: Load documents by ID
3. **Query Documents**: Use different query patterns:
   - Query all documents in a collection
   - Query by specific field values
   - Query by value ranges
   - Full-text search across multiple fields
   - Custom queries with parameters
4. **Update Documents**: Modify existing documents
5. **Delete Documents**: Remove single and multiple documents

### Advanced Query Tests
1. **Pagination**: Test skip/take parameters
2. **Sorting**: Test ordering by fields (ascending/descending)
3. **Filtering**: Test complex where clauses
4. **Parameters**: Test parameterized queries
5. **Type Safety**: Ensure generic queries maintain type safety

## Expected Results

### Successful Test Run
- All connections established
- Database created/accessed successfully
- All CRUD operations complete without errors
- Query operations return expected results
- Type safety maintained throughout
- No memory leaks or resource issues

### Common Issues

1. **Connection Failed**: Check if RavenDB server is running
2. **Database Access Denied**: Verify server configuration and permissions
3. **Test Timeout**: Increase timeout in `test_config.toml`
4. **Document Not Found**: May indicate timing issues or improper cleanup

## Coverage Goals

- **Minimum Coverage**: 80%
- **Target Coverage**: 90%+

The test suite covers:
- All public interface methods
- Error handling paths
- Edge cases (empty results, invalid parameters)
- Type safety scenarios
- Concurrent operations (where applicable)

## Continuous Integration

For CI/CD pipelines:

```bash
# Install dependencies
go mod download

# Run tests with JUnit output
go test -v ./... -coverprofile=coverage.out -json > test-results.json

# Generate coverage report
go tool cover -html=coverage.out -o coverage.html
```

## Troubleshooting

### RavenDB Connection Issues
1. Verify RavenDB server is running: `curl http://localhost:8080`
2. Check server logs for errors
3. Ensure no firewall blocking connections
4. Verify correct URL in `test_config.toml`

### Test Failures
1. Check test output for specific error messages
2. Verify test data isn't conflicting
3. Ensure database permissions are correct
4. Check for resource cleanup issues

### Performance Issues
1. Monitor RavenDB server resource usage
2. Check for connection leaks
3. Verify proper session cleanup
4. Consider reducing test data size