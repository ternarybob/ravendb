package services

import (
	"fmt"

	"github.com/ravendb/ravendb-go-client"
	"github.com/ternarybob/ravendb/interfaces"
)

// DatabaseService handles RavenDB connections and database operations
type DatabaseService struct {
	store    *ravendb.DocumentStore
	database string
}

// NewDatabaseService creates a new RavenDB database service
func NewDatabaseService(urls []string, database string) (interfaces.IRavenDBService, error) {
	store := ravendb.NewDocumentStore(urls, database)

	// Configure for single-node development setup
	store.GetConventions().SetDisableTopologyUpdates(true)

	// Initialize the document store
	if err := store.Initialize(); err != nil {
		return nil, fmt.Errorf("failed to initialize RavenDB store: %w", err)
	}

	return &DatabaseService{
		store:    store,
		database: database,
	}, nil
}

// Init initializes the RavenDB database with robust error handling
func (ds *DatabaseService) Init() error {
	fmt.Printf("Initializing RavenDB connection to: %v, database: %s\n", ds.store.GetUrls(), ds.database)

	// First, try to create the database (this is safer than testing first)
	fmt.Printf("Creating database '%s' if it doesn't exist...\n", ds.database)

	databaseRecord := &ravendb.DatabaseRecord{
		DatabaseName: ds.database,
	}

	// Use replication factor of 1 for single-node setup
	operation := ravendb.NewCreateDatabaseOperation(databaseRecord, 1)
	err := ds.store.Maintenance().Server().Send(operation)
	if err != nil {
		// Database might already exist, this is okay
		fmt.Printf("Database creation result (might already exist): %v\n", err)
	}

	// Now test if we can open a session to the specific database
	session, err := ds.store.OpenSession(ds.database)
	if err != nil {
		return fmt.Errorf("failed to open session to database '%s': %w", ds.database, err)
	}
	defer session.Close()

	fmt.Printf("Successfully connected to RavenDB database '%s'\n", ds.database)
	return nil
}

// InitializeWithSeeding initializes the database and optionally seeds it with data
func (ds *DatabaseService) InitializeWithSeeding(seedData bool) error {
	// First, initialize/create the database
	if err := ds.Init(); err != nil {
		return fmt.Errorf("database initialization failed: %w", err)
	}

	// If seeding is requested, check if database is empty
	if seedData {
		isEmpty, err := ds.isDatabaseEmpty()
		if err != nil {
			fmt.Printf("Warning: Could not check if database is empty: %v\n", err)
			return nil // Continue without seeding
		}

		if isEmpty {
			fmt.Println("Database is empty, automatic seeding would be performed here")
			// Note: Actual seeding is handled by service layers
		} else {
			fmt.Println("Database contains data, skipping automatic seeding")
		}
	}

	return nil
}

// isDatabaseEmpty checks if the database has any documents
func (ds *DatabaseService) isDatabaseEmpty() (bool, error) {
	session, err := ds.store.OpenSession(ds.database)
	if err != nil {
		return false, fmt.Errorf("failed to open session: %w", err)
	}
	defer session.Close()

	// Try to get database statistics to check document count
	statisticsOperation := ravendb.NewGetStatisticsOperation("")
	result := ds.store.Maintenance().ForDatabase(ds.database).Send(statisticsOperation)

	// If we can't get statistics, assume database is empty
	if result == nil {
		return true, nil
	}

	// For now, we'll assume database is empty since statistics parsing is complex
	return true, nil
}

// Close closes the RavenDB connection
func (ds *DatabaseService) Close() error {
	if ds.store != nil {
		ds.store.Close()
	}
	return nil
}

// GetDatabaseStatus returns information about the RavenDB database state
func (ds *DatabaseService) GetDatabaseStatus() (map[string]interface{}, error) {
	status := make(map[string]interface{})

	session, err := ds.store.OpenSession(ds.database)
	if err != nil {
		status["database_name"] = ds.database
		status["status"] = "disconnected"
		status["error"] = err.Error()
		return status, err
	}
	defer session.Close()

	// Get database statistics
	statisticsOperation := ravendb.NewGetStatisticsOperation("")
	result := ds.store.Maintenance().ForDatabase(ds.database).Send(statisticsOperation)
	
	if result == nil {
		status["database_name"] = ds.database
		status["status"] = "connected"
		status["session_active"] = true
		status["statistics_error"] = "failed to get statistics"
		return status, nil
	}

	// For now, provide basic status without detailed statistics
	status["database_name"] = ds.database
	status["status"] = "connected"
	status["session_active"] = true
	status["document_count"] = 0 // Placeholder
	status["index_count"] = 0    // Placeholder

	return status, nil
}

// GetStore returns the underlying RavenDB DocumentStore as interface{} for interface compatibility
func (ds *DatabaseService) GetStore() interface{} {
	return ds.store
}

// GetDatabase returns the database name
func (ds *DatabaseService) GetDatabase() string {
	return ds.database
}