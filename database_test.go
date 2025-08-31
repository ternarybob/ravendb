package ravendb

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testConfigPath = "config/test_config.toml"

// Test document types for testing
type TestUser struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Age      int       `json:"age"`
	IsActive bool      `json:"isActive"`
	Created  time.Time `json:"created"`
}

type TestProduct struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Category    string  `json:"category"`
	InStock     bool    `json:"inStock"`
}

func TestDatabaseConnection(t *testing.T) {
	// Load test configuration
	testConfig, err := LoadTestConfig(testConfigPath)
	require.NoError(t, err, "Failed to load test configuration")

	// Test database service creation
	config := &Config{
		URLs:     testConfig.Database.URLs,
		Database: testConfig.Database.Database,
	}

	db, err := NewDatabase(config)
	require.NoError(t, err, "Failed to create database service")
	defer db.Close()

	// Test database initialization
	err = db.Init()
	assert.NoError(t, err, "Failed to initialize database")

	// Test database status
	status, err := db.GetDatabaseStatus()
	assert.NoError(t, err, "Failed to get database status")
	assert.Equal(t, testConfig.Database.Database, status["database_name"])
	assert.Equal(t, "connected", status["status"])
}

func TestDatabaseAvailability(t *testing.T) {
	testConfig, err := LoadTestConfig(testConfigPath)
	require.NoError(t, err, "Failed to load test configuration")

	db, cleanup := SetupTestDatabase(t, testConfig)
	defer cleanup()

	// Test that we can get database information
	dbName := db.GetDatabase()
	assert.Equal(t, testConfig.Database.Database, dbName)

	// Test that we can get the underlying store
	store := db.GetStore()
	assert.NotNil(t, store)
}

func TestBasicCRUDOperations(t *testing.T) {
	testConfig, err := LoadTestConfig(testConfigPath)
	require.NoError(t, err, "Failed to load test configuration")

	db, cleanup := SetupTestDatabase(t, testConfig)
	defer cleanup()

	t.Run("StoreAndLoadDocument", func(t *testing.T) {
		// Create test document
		user := TestUser{
			ID:       "users/test-1",
			Name:     "John Doe",
			Email:    "john@example.com",
			Age:      30,
			IsActive: true,
			Created:  time.Now(),
		}

		// Store document
		err := db.Store("users/test-1", user)
		assert.NoError(t, err, "Failed to store document")

		// Load document
		var loadedUser TestUser
		err = db.LoadByID("users/test-1", &loadedUser)
		assert.NoError(t, err, "Failed to load document")

		// Verify document contents
		assert.Equal(t, user.ID, loadedUser.ID)
		assert.Equal(t, user.Name, loadedUser.Name)
		assert.Equal(t, user.Email, loadedUser.Email)
		assert.Equal(t, user.Age, loadedUser.Age)
		assert.Equal(t, user.IsActive, loadedUser.IsActive)
	})

	t.Run("DocumentExists", func(t *testing.T) {
		// Test document existence
		exists, err := db.Exists("users/test-1")
		assert.NoError(t, err, "Failed to check document existence")
		assert.True(t, exists, "Document should exist")

		// Test non-existent document
		exists, err = db.Exists("users/non-existent")
		assert.NoError(t, err, "Failed to check non-existent document")
		assert.False(t, exists, "Document should not exist")
	})

	t.Run("UpdateDocument", func(t *testing.T) {
		// Update document
		updates := map[string]interface{}{
			"age":      31,
			"isActive": false,
		}
		err := db.Update("users/test-1", updates)
		assert.NoError(t, err, "Failed to update document")

		// Load and verify updates
		var updatedUser TestUser
		err = db.LoadByID("users/test-1", &updatedUser)
		assert.NoError(t, err, "Failed to load updated document")
		assert.Equal(t, 31, updatedUser.Age)
		assert.False(t, updatedUser.IsActive)
	})

	t.Run("StoreMultipleDocuments", func(t *testing.T) {
		users := map[string]interface{}{
			"users/test-2": TestUser{
				ID:       "users/test-2",
				Name:     "Jane Smith",
				Email:    "jane@example.com",
				Age:      25,
				IsActive: true,
				Created:  time.Now(),
			},
			"users/test-3": TestUser{
				ID:       "users/test-3",
				Name:     "Bob Johnson",
				Email:    "bob@example.com",
				Age:      35,
				IsActive: true,
				Created:  time.Now(),
			},
		}

		err := db.StoreMultiple(users)
		assert.NoError(t, err, "Failed to store multiple documents")

		// Verify documents were stored
		exists, err := db.Exists("users/test-2")
		assert.NoError(t, err)
		assert.True(t, exists)

		exists, err = db.Exists("users/test-3")
		assert.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("DeleteDocument", func(t *testing.T) {
		// Delete document
		err := db.Delete("users/test-1")
		assert.NoError(t, err, "Failed to delete document")

		// Verify document was deleted
		exists, err := db.Exists("users/test-1")
		assert.NoError(t, err, "Failed to check deleted document existence")
		assert.False(t, exists, "Document should be deleted")
	})

	t.Run("DeleteMultipleDocuments", func(t *testing.T) {
		ids := []string{"users/test-2", "users/test-3"}
		err := db.DeleteMultiple(ids)
		assert.NoError(t, err, "Failed to delete multiple documents")

		// Verify documents were deleted
		for _, id := range ids {
			exists, err := db.Exists(id)
			assert.NoError(t, err)
			assert.False(t, exists, "Document %s should be deleted", id)
		}
	})
}
