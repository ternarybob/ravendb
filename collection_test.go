package ravendb

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/ternarybob/ravendb/interfaces"
)

func TestCollectionService(t *testing.T) {
	testConfig, err := LoadTestConfig(testConfigPath)
	require.NoError(t, err, "Failed to load test configuration")

	db, cleanup := SetupTestDatabase(t, testConfig)
	defer cleanup()

	// Create collection service
	userCollection := NewCollection[TestUser](db, "Users")

	t.Run("StoreAndLoadTypedDocument", func(t *testing.T) {
		user := TestUser{
			ID:       "users/typed-1",
			Name:     "Alice Cooper",
			Email:    "alice@example.com",
			Age:      28,
			IsActive: true,
			Created:  time.Now(),
		}

		// Store document
		err := userCollection.Store("users/typed-1", user)
		assert.NoError(t, err, "Failed to store typed document")

		// Load document
		loadedUser, err := userCollection.LoadByID("users/typed-1")
		assert.NoError(t, err, "Failed to load typed document")
		assert.NotNil(t, loadedUser)

		// Verify document contents
		assert.Equal(t, user.ID, loadedUser.ID)
		assert.Equal(t, user.Name, loadedUser.Name)
		assert.Equal(t, user.Email, loadedUser.Email)
		assert.Equal(t, user.Age, loadedUser.Age)
		assert.Equal(t, user.IsActive, loadedUser.IsActive)
	})

	t.Run("StoreMultipleTypedDocuments", func(t *testing.T) {
		users := map[string]TestUser{
			"users/typed-2": {
				ID:       "users/typed-2",
				Name:     "Bob Wilson",
				Email:    "bob.wilson@example.com",
				Age:      32,
				IsActive: true,
				Created:  time.Now(),
			},
			"users/typed-3": {
				ID:       "users/typed-3",
				Name:     "Carol Davis",
				Email:    "carol@example.com",
				Age:      29,
				IsActive: false,
				Created:  time.Now(),
			},
			"users/typed-4": {
				ID:       "users/typed-4",
				Name:     "David Brown",
				Email:    "david@example.com",
				Age:      45,
				IsActive: true,
				Created:  time.Now(),
			},
		}

		err := userCollection.StoreMultiple(users)
		assert.NoError(t, err, "Failed to store multiple typed documents")

		// Verify documents exist
		for id := range users {
			exists, err := userCollection.Exists(id)
			assert.NoError(t, err)
			assert.True(t, exists, "Document %s should exist", id)
		}
	})

	t.Run("LoadMultipleByIDs", func(t *testing.T) {
		ids := []string{"users/typed-2", "users/typed-3", "users/typed-4"}
		users, err := userCollection.LoadMultipleByIDs(ids)
		assert.NoError(t, err, "Failed to load multiple documents by IDs")
		assert.Len(t, users, 3, "Should load 3 documents")

		// Verify loaded documents
		userMap := make(map[string]TestUser)
		for _, user := range users {
			userMap[user.ID] = user
		}

		assert.Contains(t, userMap, "users/typed-2")
		assert.Contains(t, userMap, "users/typed-3")
		assert.Contains(t, userMap, "users/typed-4")
	})

	t.Run("UpdateTypedDocument", func(t *testing.T) {
		// Load existing document
		user, err := userCollection.LoadByID("users/typed-2")
		require.NoError(t, err)
		require.NotNil(t, user)

		// Update fields
		user.Age = 33
		user.IsActive = false

		// Update document
		err = userCollection.Update("users/typed-2", *user)
		assert.NoError(t, err, "Failed to update typed document")

		// Load and verify
		updatedUser, err := userCollection.LoadByID("users/typed-2")
		assert.NoError(t, err)
		assert.Equal(t, 33, updatedUser.Age)
		assert.False(t, updatedUser.IsActive)
	})

	t.Run("QueryAllDocuments", func(t *testing.T) {
		results, err := userCollection.QueryAll()
		assert.NoError(t, err, "Failed to query all documents")
		assert.GreaterOrEqual(t, len(results.Results), 4, "Should have at least 4 users")
		assert.False(t, results.HasMore) // Since we're within default page size
	})

	t.Run("QueryByField", func(t *testing.T) {
		// Query active users
		activeUsers, err := userCollection.QueryByField("isActive", true, nil)
		assert.NoError(t, err, "Failed to query by field")
		
		// Verify all returned users are active
		for _, user := range activeUsers.Results {
			assert.True(t, user.IsActive, "All returned users should be active")
		}

		// Query inactive users
		inactiveUsers, err := userCollection.QueryByField("isActive", false, nil)
		assert.NoError(t, err, "Failed to query inactive users")
		
		// Verify all returned users are inactive
		for _, user := range inactiveUsers.Results {
			assert.False(t, user.IsActive, "All returned users should be inactive")
		}
	})

	t.Run("QueryByRange", func(t *testing.T) {
		// Query users by age range
		ageRangeUsers, err := userCollection.QueryByRange("age", 25, 35, nil)
		assert.NoError(t, err, "Failed to query by age range")

		// Verify all returned users are in age range
		for _, user := range ageRangeUsers.Results {
			assert.GreaterOrEqual(t, user.Age, 25, "User age should be >= 25")
			assert.LessOrEqual(t, user.Age, 35, "User age should be <= 35")
		}
	})

	t.Run("SearchDocuments", func(t *testing.T) {
		// Search for users with "Bob" in name or email
		searchResults, err := userCollection.Search("Bob", []string{"name", "email"}, nil)
		assert.NoError(t, err, "Failed to search documents")

		// Verify search results contain users with "Bob"
		found := false
		for _, user := range searchResults.Results {
			if user.Name == "Bob Wilson" {
				found = true
				break
			}
		}
		assert.True(t, found, "Should find Bob Wilson in search results")
	})

	t.Run("QueryWithOptions", func(t *testing.T) {
		// Query with custom options
		options := &interfaces.QueryOptions{
			Skip:      0,
			Take:      2,
			OrderBy:   "age",
			OrderDesc: true, // Order by age descending
		}

		results, err := userCollection.Query(options)
		assert.NoError(t, err, "Failed to query with options")
		assert.LessOrEqual(t, len(results.Results), 2, "Should return at most 2 results")
		assert.Equal(t, 0, results.Skip)
		assert.Equal(t, 2, results.Take)

		// Verify ordering (should be descending by age)
		if len(results.Results) >= 2 {
			assert.GreaterOrEqual(t, results.Results[0].Age, results.Results[1].Age, 
				"Results should be ordered by age descending")
		}
	})

	t.Run("CountDocuments", func(t *testing.T) {
		count, err := userCollection.Count()
		assert.NoError(t, err, "Failed to count documents")
		assert.GreaterOrEqual(t, count, 4, "Should have at least 4 users")
	})

	t.Run("DeleteTypedDocuments", func(t *testing.T) {
		// Delete single document
		err := userCollection.Delete("users/typed-1")
		assert.NoError(t, err, "Failed to delete typed document")

		// Verify deletion
		exists, err := userCollection.Exists("users/typed-1")
		assert.NoError(t, err)
		assert.False(t, exists, "Document should be deleted")

		// Delete multiple documents
		ids := []string{"users/typed-2", "users/typed-3", "users/typed-4"}
		err = userCollection.DeleteMultiple(ids)
		assert.NoError(t, err, "Failed to delete multiple typed documents")

		// Verify deletions
		for _, id := range ids {
			exists, err := userCollection.Exists(id)
			assert.NoError(t, err)
			assert.False(t, exists, "Document %s should be deleted", id)
		}
	})
}

func TestGenericQueryOperations(t *testing.T) {
	testConfig, err := LoadTestConfig(testConfigPath)
	require.NoError(t, err, "Failed to load test configuration")

	db, cleanup := SetupTestDatabase(t, testConfig)
	defer cleanup()

	// Set up test data
	productCollection := NewCollection[TestProduct](db, "Products")
	
	products := map[string]TestProduct{
		"products/laptop": {
			ID:          "products/laptop",
			Name:        "Gaming Laptop",
			Description: "High-performance gaming laptop",
			Price:       1299.99,
			Category:    "Electronics",
			InStock:     true,
		},
		"products/mouse": {
			ID:          "products/mouse",
			Name:        "Wireless Mouse",
			Description: "Ergonomic wireless mouse",
			Price:       49.99,
			Category:    "Electronics",
			InStock:     true,
		},
		"products/keyboard": {
			ID:          "products/keyboard",
			Name:        "Mechanical Keyboard",
			Description: "RGB mechanical keyboard",
			Price:       129.99,
			Category:    "Electronics",
			InStock:     false,
		},
	}

	err = productCollection.StoreMultiple(products)
	require.NoError(t, err, "Failed to store test products")

	t.Run("GenericQueryAll", func(t *testing.T) {
		results, err := QueryAll[TestProduct](db, "Products")
		assert.NoError(t, err, "Failed to query all products generically")
		assert.GreaterOrEqual(t, len(results.Results), 3, "Should have at least 3 products")
	})

	t.Run("GenericQueryByField", func(t *testing.T) {
		// Query products in stock
		inStockProducts, err := QueryByField[TestProduct](db, "Products", "inStock", true, nil)
		assert.NoError(t, err, "Failed to query products by field generically")
		
		// Verify all returned products are in stock
		for _, product := range inStockProducts.Results {
			assert.True(t, product.InStock, "All returned products should be in stock")
		}
	})

	t.Run("GenericQueryByRange", func(t *testing.T) {
		// Query products by price range
		options := &interfaces.QueryOptions{
			Take: 10,
		}

		priceRangeProducts, err := QueryByRange[TestProduct](db, "Products", "price", 50.0, 200.0, options)
		assert.NoError(t, err, "Failed to query products by price range generically")

		// Verify all returned products are in price range
		for _, product := range priceRangeProducts.Results {
			assert.GreaterOrEqual(t, product.Price, 50.0, "Product price should be >= 50")
			assert.LessOrEqual(t, product.Price, 200.0, "Product price should be <= 200")
		}
	})

	t.Run("GenericSearch", func(t *testing.T) {
		// Search for products with "laptop" in name or description
		searchResults, err := Search[TestProduct](db, "Products", "laptop", []string{"name", "description"}, nil)
		assert.NoError(t, err, "Failed to search products generically")

		// Verify search results
		found := false
		for _, product := range searchResults.Results {
			if product.Name == "Gaming Laptop" {
				found = true
				break
			}
		}
		assert.True(t, found, "Should find Gaming Laptop in search results")
	})

	// Clean up test products
	t.Cleanup(func() {
		productIDs := []string{"products/laptop", "products/mouse", "products/keyboard"}
		productCollection.DeleteMultiple(productIDs)
	})
}