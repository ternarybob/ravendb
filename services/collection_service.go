package services

import (
	"fmt"
	"strings"

	"github.com/ravendb/ravendb-go-client"
	"github.com/ternarybob/ravendb/interfaces"
)

// CollectionService provides generic type-safe operations for a specific document type
type CollectionService[T any] struct {
	database   interfaces.IRavenDBService
	collection string
}

// NewCollectionService creates a new collection service for a specific document type
func NewCollectionService[T any](database interfaces.IRavenDBService, collection string) interfaces.IRavenCollectionService[T] {
	return &CollectionService[T]{
		database:   database,
		collection: collection,
	}
}

// CRUD Operations

// Store stores a document with the specified ID
func (cs *CollectionService[T]) Store(id string, document T) error {
	store := cs.database.GetStore().(*ravendb.DocumentStore)
	session, err := store.OpenSession(cs.database.GetDatabase())
	if err != nil {
		return fmt.Errorf("failed to open session: %w", err)
	}
	defer session.Close()

	// Use the proper Store API from the documentation - RavenDB expects a pointer
	if id != "" {
		err = session.StoreWithID(&document, id)
	} else {
		err = session.Store(&document)
	}

	if err != nil {
		return fmt.Errorf("failed to store document: %w", err)
	}

	return session.SaveChanges()
}

// StoreMultiple stores multiple documents in a single transaction
func (cs *CollectionService[T]) StoreMultiple(documents map[string]T) error {
	store := cs.database.GetStore().(*ravendb.DocumentStore)
	session, err := store.OpenSession(cs.database.GetDatabase())
	if err != nil {
		return fmt.Errorf("failed to open session: %w", err)
	}
	defer session.Close()

	for id, document := range documents {
		// Need to pass pointer to RavenDB
		doc := document // Create a copy to take address of
		if id != "" {
			err = session.StoreWithID(&doc, id)
		} else {
			err = session.Store(&doc)
		}
		if err != nil {
			return fmt.Errorf("failed to store document with ID %s: %w", id, err)
		}
	}

	return session.SaveChanges()
}

// LoadByID loads a document by ID
func (cs *CollectionService[T]) LoadByID(id string) (*T, error) {
	store := cs.database.GetStore().(*ravendb.DocumentStore)
	session, err := store.OpenSession(cs.database.GetDatabase())
	if err != nil {
		return nil, fmt.Errorf("failed to open session: %w", err)
	}
	defer session.Close()

	var result *T
	err = session.Load(&result, id)
	if err != nil {
		return nil, fmt.Errorf("failed to load document: %w", err)
	}

	// If document doesn't exist, result will be nil
	if result == nil {
		return nil, nil
	}

	return result, nil
}

// LoadMultipleByIDs loads multiple documents by their IDs
func (cs *CollectionService[T]) LoadMultipleByIDs(ids []string) ([]T, error) {
	store := cs.database.GetStore().(*ravendb.DocumentStore)
	session, err := store.OpenSession(cs.database.GetDatabase())
	if err != nil {
		return nil, fmt.Errorf("failed to open session: %w", err)
	}
	defer session.Close()

	var results []T
	for _, id := range ids {
		var doc *T
		err = session.Load(&doc, id)
		if err != nil {
			return nil, fmt.Errorf("failed to load document %s: %w", id, err)
		}
		// Check if document exists
		if doc != nil {
			results = append(results, *doc)
		}
	}

	return results, nil
}

// Update updates an existing document
func (cs *CollectionService[T]) Update(id string, document T) error {
	store := cs.database.GetStore().(*ravendb.DocumentStore)
	session, err := store.OpenSession(cs.database.GetDatabase())
	if err != nil {
		return fmt.Errorf("failed to open session: %w", err)
	}
	defer session.Close()

	// Store the updated document - RavenDB expects a pointer
	err = session.StoreWithID(&document, id)
	if err != nil {
		return fmt.Errorf("failed to store updated document: %w", err)
	}

	return session.SaveChanges()
}

// Delete removes a document by ID
func (cs *CollectionService[T]) Delete(id string) error {
	store := cs.database.GetStore().(*ravendb.DocumentStore)
	session, err := store.OpenSession(cs.database.GetDatabase())
	if err != nil {
		return fmt.Errorf("failed to open session: %w", err)
	}
	defer session.Close()

	// Load document first, then delete - same pattern as database service
	var document *T
	err = session.Load(&document, id)
	if err != nil {
		return fmt.Errorf("failed to load document for deletion: %w", err)
	}
	
	if document == nil {
		return fmt.Errorf("document with ID %s not found", id)
	}

	session.Delete(document)
	return session.SaveChanges()
}

// DeleteMultiple removes multiple documents by their IDs
func (cs *CollectionService[T]) DeleteMultiple(ids []string) error {
	store := cs.database.GetStore().(*ravendb.DocumentStore)
	session, err := store.OpenSession(cs.database.GetDatabase())
	if err != nil {
		return fmt.Errorf("failed to open session: %w", err)
	}
	defer session.Close()

	// Load each document first, then delete - same pattern as database service
	for _, id := range ids {
		var document *T
		err = session.Load(&document, id)
		if err != nil {
			return fmt.Errorf("failed to load document %s for deletion: %w", id, err)
		}
		
		if document != nil {
			session.Delete(document)
		}
		// Skip if document doesn't exist instead of failing
	}

	return session.SaveChanges()
}

// Query Operations

// Query executes a generic query with options
func (cs *CollectionService[T]) Query(options *interfaces.QueryOptions) (*interfaces.GenericQueryResult[T], error) {
	store := cs.database.GetStore().(*ravendb.DocumentStore)
	session, err := store.OpenSession(cs.database.GetDatabase())
	if err != nil {
		return nil, fmt.Errorf("failed to open session: %w", err)
	}
	defer session.Close()

	// Set default values
	if options == nil {
		options = &interfaces.QueryOptions{}
	}
	if options.Take <= 0 {
		options.Take = 25
	}
	if options.Take > 1024 {
		options.Take = 1024
	}

	// Build RQL query dynamically
	var rqlQuery strings.Builder
	// For now, use a flexible approach that works with Go struct collections
	// RavenDB Go client typically assigns collection names based on the struct type name
	rqlQuery.WriteString(fmt.Sprintf("from @all_docs where @metadata.'@collection' = '%s'", cs.collection))

	// Add WHERE clause if specified
	if options.WhereClause != "" {
		rqlQuery.WriteString(fmt.Sprintf(" AND (%s)", options.WhereClause))
	}

	// Add ORDER BY if specified
	if options.OrderBy != "" {
		if options.OrderDesc {
			rqlQuery.WriteString(fmt.Sprintf(" ORDER BY %s DESC", options.OrderBy))
		} else {
			rqlQuery.WriteString(fmt.Sprintf(" ORDER BY %s", options.OrderBy))
		}
	}

	// Add LIMIT (skip, take) for pagination
	if options.Skip > 0 || options.Take > 0 {
		skip := options.Skip
		take := options.Take
		if take <= 0 {
			take = 25
		}
		rqlQuery.WriteString(fmt.Sprintf(" LIMIT %d, %d", skip, take))
	}

	// Execute the raw query
	queryStr := rqlQuery.String()
	query := session.Advanced().RawQuery(queryStr)

	// Set parameters if provided
	if options.Parameters != nil {
		for key, value := range options.Parameters {
			query = query.AddParameter(key, value)
		}
	}

	var results []*T
	err = query.GetResults(&results)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	// Convert pointers to values
	finalResults := make([]T, len(results))
	for i, res := range results {
		if res != nil {
			finalResults[i] = *res
		}
	}

	totalCount := len(finalResults)
	hasMore := options.Take > 0 && totalCount == options.Take

	return &interfaces.GenericQueryResult[T]{
		Results:    finalResults,
		TotalCount: totalCount,
		Skip:       options.Skip,
		Take:       options.Take,
		HasMore:    hasMore,
	}, nil
}

// QueryAll queries all documents of type T
func (cs *CollectionService[T]) QueryAll() (*interfaces.GenericQueryResult[T], error) {
	options := &interfaces.QueryOptions{
		Skip: 0,
		Take: 1024, // Default page size
	}

	return cs.Query(options)
}

// QueryByField queries documents by a specific field value
func (cs *CollectionService[T]) QueryByField(fieldName string, fieldValue interface{}, options *interfaces.QueryOptions) (*interfaces.GenericQueryResult[T], error) {
	if options == nil {
		options = &interfaces.QueryOptions{}
	}

	// Build where clause using the field name
	options.WhereClause = fmt.Sprintf("%s = $value", fieldName)
	if options.Parameters == nil {
		options.Parameters = make(map[string]interface{})
	}
	options.Parameters["value"] = fieldValue

	return cs.Query(options)
}

// QueryByRange queries documents within a range of values
func (cs *CollectionService[T]) QueryByRange(fieldName string, minValue, maxValue interface{}, options *interfaces.QueryOptions) (*interfaces.GenericQueryResult[T], error) {
	if options == nil {
		options = &interfaces.QueryOptions{}
	}

	// Build where clause for range
	options.WhereClause = fmt.Sprintf("%s >= $minValue AND %s <= $maxValue", fieldName, fieldName)
	if options.Parameters == nil {
		options.Parameters = make(map[string]interface{})
	}
	options.Parameters["minValue"] = minValue
	options.Parameters["maxValue"] = maxValue

	return cs.Query(options)
}

// Search performs a full-text search across specified fields
func (cs *CollectionService[T]) Search(searchTerm string, searchFields []string, options *interfaces.QueryOptions) (*interfaces.GenericQueryResult[T], error) {
	if options == nil {
		options = &interfaces.QueryOptions{}
	}

	// Build search where clause
	var whereConditions []string
	if options.Parameters == nil {
		options.Parameters = make(map[string]interface{})
	}

	for i, field := range searchFields {
		paramName := fmt.Sprintf("searchTerm%d", i)
		whereConditions = append(whereConditions, fmt.Sprintf("search(%s, $%s)", field, paramName))
		options.Parameters[paramName] = searchTerm
	}

	if len(whereConditions) > 0 {
		options.WhereClause = fmt.Sprintf("(%s)", strings.Join(whereConditions, " OR "))
	}

	return cs.Query(options)
}

// Utility Methods

// Exists checks if a document with the given ID exists
func (cs *CollectionService[T]) Exists(id string) (bool, error) {
	store := cs.database.GetStore().(*ravendb.DocumentStore)
	session, err := store.OpenSession(cs.database.GetDatabase())
	if err != nil {
		return false, fmt.Errorf("failed to open session: %w", err)
	}
	defer session.Close()

	var document *T
	err = session.Load(&document, id)
	if err != nil {
		return false, fmt.Errorf("failed to check document existence: %w", err)
	}

	// Document exists if the result is not nil
	return document != nil, nil
}

// Count returns the total number of documents in this collection
func (cs *CollectionService[T]) Count() (int, error) {
	result, err := cs.QueryAll()
	if err != nil {
		return 0, fmt.Errorf("failed to count documents: %w", err)
	}
	return result.TotalCount, nil
}
