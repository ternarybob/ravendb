package examples

import (
	"fmt"
	"log"

	ravendb "github.com/ternarybob/ravendb"
	"github.com/ternarybob/ravendb/interfaces"
)

// User represents a sample document type
type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Age      int    `json:"age"`
	IsActive bool   `json:"isActive"`
}

// Product represents another sample document type
type Product struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Category    string  `json:"category"`
	InStock     bool    `json:"inStock"`
}

// BasicUsageExample demonstrates basic usage of the RavenDB library
func BasicUsageExample() {
	// Create a local development configuration
	config := ravendb.NewLocalConfig("ExampleDB")

	// Initialize database service
	db, err := ravendb.NewDatabase(config)
	if err != nil {
		log.Fatalf("Failed to create database service: %v", err)
	}
	defer db.Close()

	// Initialize the database
	if err := db.Init(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	fmt.Println("=== Basic CRUD Operations ===")
	basicCRUDExample(db)

	fmt.Println("\n=== Collection Service Example ===")
	collectionServiceExample(db)

	fmt.Println("\n=== Generic Query Example ===")
	genericQueryExample(db)
}

func basicCRUDExample(db interfaces.IRavenDBService) {
	// Store a document
	user := User{
		ID:       "users/1",
		Name:     "John Doe",
		Email:    "john@example.com",
		Age:      30,
		IsActive: true,
	}

	if err := db.Store("users/1", user); err != nil {
		log.Printf("Failed to store user: %v", err)
		return
	}
	fmt.Println("✓ Stored user document")

	// Load the document
	var loadedUser User
	if err := db.LoadByID("users/1", &loadedUser); err != nil {
		log.Printf("Failed to load user: %v", err)
		return
	}
	fmt.Printf("✓ Loaded user: %s (%s)\n", loadedUser.Name, loadedUser.Email)

	// Check if document exists
	exists, err := db.Exists("users/1")
	if err != nil {
		log.Printf("Failed to check existence: %v", err)
		return
	}
	fmt.Printf("✓ User exists: %v\n", exists)

	// Update the document
	updates := map[string]interface{}{
		"age":      31,
		"isActive": false,
	}
	if err := db.Update("users/1", updates); err != nil {
		log.Printf("Failed to update user: %v", err)
		return
	}
	fmt.Println("✓ Updated user document")
}

func collectionServiceExample(db interfaces.IRavenDBService) {
	// Create a typed collection service for User documents
	userCollection := ravendb.NewCollection[User](db, "Users")

	// Store multiple users
	users := map[string]User{
		"users/2": {
			ID:       "users/2",
			Name:     "Jane Smith",
			Email:    "jane@example.com",
			Age:      25,
			IsActive: true,
		},
		"users/3": {
			ID:       "users/3",
			Name:     "Bob Johnson",
			Email:    "bob@example.com",
			Age:      35,
			IsActive: true,
		},
	}

	if err := userCollection.StoreMultiple(users); err != nil {
		log.Printf("Failed to store multiple users: %v", err)
		return
	}
	fmt.Println("✓ Stored multiple users")

	// Query all users
	allUsers, err := userCollection.QueryAll()
	if err != nil {
		log.Printf("Failed to query all users: %v", err)
		return
	}
	fmt.Printf("✓ Found %d users total\n", len(allUsers.Results))

	// Query users by field
	activeUsers, err := userCollection.QueryByField("isActive", true, nil)
	if err != nil {
		log.Printf("Failed to query active users: %v", err)
		return
	}
	fmt.Printf("✓ Found %d active users\n", len(activeUsers.Results))

	// Search users by name
	searchResults, err := userCollection.Search("John", []string{"name", "email"}, nil)
	if err != nil {
		log.Printf("Failed to search users: %v", err)
		return
	}
	fmt.Printf("✓ Search found %d users matching 'John'\n", len(searchResults.Results))
}

func genericQueryExample(db interfaces.IRavenDBService) {
	// Create product documents for demonstration
	productCollection := ravendb.NewCollection[Product](db, "Products")

	products := map[string]Product{
		"products/1": {
			ID:          "products/1",
			Name:        "Laptop",
			Description: "High-performance laptop",
			Price:       999.99,
			Category:    "Electronics",
			InStock:     true,
		},
		"products/2": {
			ID:          "products/2",
			Name:        "Mouse",
			Description: "Wireless optical mouse",
			Price:       29.99,
			Category:    "Electronics",
			InStock:     false,
		},
	}

	if err := productCollection.StoreMultiple(products); err != nil {
		log.Printf("Failed to store products: %v", err)
		return
	}

	// Query products in price range using the generic query function
	options := &interfaces.QueryOptions{
		Take: 10,
	}

	priceRangeResults, err := ravendb.QueryByRange[Product](
		db, "Products", "price", 20.0, 1000.0, options,
	)
	if err != nil {
		log.Printf("Failed to query products by price range: %v", err)
		return
	}
	fmt.Printf("✓ Found %d products in price range $20-$1000\n", len(priceRangeResults.Results))

	// Query with custom options
	customOptions := &interfaces.QueryOptions{
		OrderBy:   "price",
		OrderDesc: true,
		Take:      5,
	}

	expensiveProducts, err := ravendb.QueryAll[Product](db, "Products")
	if err != nil {
		log.Printf("Failed to query expensive products: %v", err)
		return
	}
	fmt.Printf("✓ Query with custom options returned %d products\n", len(expensiveProducts.Results))
}
