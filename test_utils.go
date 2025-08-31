package ravendb

import (
	"os"
	"testing"

	"github.com/pelletier/go-toml/v2"
	"github.com/ternarybob/ravendb/services"
)

// TestConfig holds the test configuration
type TestConfig struct {
	Database DatabaseConfig `toml:"database"`
	Test     TestSettings   `toml:"test"`
}

// DatabaseConfig holds database connection settings
type DatabaseConfig struct {
	URLs     []string `toml:"urls"`
	Database string   `toml:"database"`
}

// TestSettings holds test-specific settings
type TestSettings struct {
	Timeout          int  `toml:"timeout"`
	CleanBeforeTests bool `toml:"clean_before_tests"`
	CleanAfterTests  bool `toml:"clean_after_tests"`
}

// LoadTestConfig loads test configuration from TOML file
func LoadTestConfig(filepath string) (*TestConfig, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var config TestConfig
	err = toml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

// SetupTestDatabase creates a test database service with the given config
func SetupTestDatabase(t *testing.T, testConfig *TestConfig) (*services.DatabaseService, func()) {
	config := &Config{
		URLs:     testConfig.Database.URLs,
		Database: testConfig.Database.Database,
	}

	db, err := NewDatabase(config)
	if err != nil {
		t.Fatalf("Failed to create database service: %v", err)
	}

	// Initialize database
	err = db.Init()
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}

	// Return cleanup function
	cleanup := func() {
		if testConfig.Test.CleanAfterTests {
			// Clean up test database if needed
			// This would involve dropping the database or clearing collections
		}
		db.Close()
	}

	return db.(*services.DatabaseService), cleanup
}
