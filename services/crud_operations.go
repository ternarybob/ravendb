package services

import (
	"fmt"
	"reflect"

	"github.com/ravendb/ravendb-go-client"
)

// Store stores a document with the specified ID
func (ds *DatabaseService) Store(id string, document interface{}) error {
	store := ds.GetStore().(*ravendb.DocumentStore)
	session, err := store.OpenSession(ds.GetDatabase())
	if err != nil {
		return fmt.Errorf("failed to open session: %w", err)
	}
	defer session.Close()

	// Use the proper Store API from the documentation
	if id != "" {
		err = session.StoreWithID(document, id)
	} else {
		err = session.Store(document)
	}

	if err != nil {
		return fmt.Errorf("failed to store document: %w", err)
	}

	return session.SaveChanges()
}

// StoreMultiple stores multiple documents in a single transaction
func (ds *DatabaseService) StoreMultiple(documents map[string]interface{}) error {
	store := ds.GetStore().(*ravendb.DocumentStore)
	session, err := store.OpenSession(ds.GetDatabase())
	if err != nil {
		return fmt.Errorf("failed to open session: %w", err)
	}
	defer session.Close()

	for id, document := range documents {
		if id != "" {
			err = session.StoreWithID(document, id)
		} else {
			err = session.Store(document)
		}
		if err != nil {
			return fmt.Errorf("failed to store document with ID %s: %w", id, err)
		}
	}

	return session.SaveChanges()
}

// LoadByID loads a document by ID into the result interface
func (ds *DatabaseService) LoadByID(id string, result interface{}) error {
	store := ds.GetStore().(*ravendb.DocumentStore)
	session, err := store.OpenSession(ds.GetDatabase())
	if err != nil {
		return fmt.Errorf("failed to open session: %w", err)
	}
	defer session.Close()

	err = session.Load(result, id)
	if err != nil {
		return fmt.Errorf("failed to load document: %w", err)
	}

	return nil
}

// LoadMultipleByIDs loads multiple documents by their IDs
func (ds *DatabaseService) LoadMultipleByIDs(ids []string, results interface{}) error {
	store := ds.GetStore().(*ravendb.DocumentStore)
	session, err := store.OpenSession(ds.GetDatabase())
	if err != nil {
		return fmt.Errorf("failed to open session: %w", err)
	}
	defer session.Close()

	// Load documents one by one and collect them
	resultsSlice := reflect.ValueOf(results).Elem()
	for _, id := range ids {
		var doc interface{}
		err = session.Load(&doc, id)
		if err != nil {
			return fmt.Errorf("failed to load document %s: %w", id, err)
		}
		if doc != nil {
			resultsSlice = reflect.Append(resultsSlice, reflect.ValueOf(doc))
		}
	}
	reflect.ValueOf(results).Elem().Set(resultsSlice)

	return nil
}

// Update updates an existing document
func (ds *DatabaseService) Update(id string, updates map[string]interface{}) error {
	store := ds.GetStore().(*ravendb.DocumentStore)
	session, err := store.OpenSession(ds.GetDatabase())
	if err != nil {
		return fmt.Errorf("failed to open session: %w", err)
	}
	defer session.Close()

	// Load the document first
	var document map[string]interface{}
	err = session.Load(&document, id)
	if err != nil {
		return fmt.Errorf("failed to load document for update: %w", err)
	}

	if document == nil {
		return fmt.Errorf("document with ID %s not found", id)
	}

	// Apply updates
	for key, value := range updates {
		document[key] = value
	}

	// Store the updated document
	err = session.StoreWithID(document, id)
	if err != nil {
		return fmt.Errorf("failed to store updated document: %w", err)
	}

	return session.SaveChanges()
}

// Delete removes a document by ID
func (ds *DatabaseService) Delete(id string) error {
	store := ds.GetStore().(*ravendb.DocumentStore)
	session, err := store.OpenSession(ds.GetDatabase())
	if err != nil {
		return fmt.Errorf("failed to open session: %w", err)
	}
	defer session.Close()

	session.Delete(id)
	return session.SaveChanges()
}

// DeleteMultiple removes multiple documents by their IDs
func (ds *DatabaseService) DeleteMultiple(ids []string) error {
	store := ds.GetStore().(*ravendb.DocumentStore)
	session, err := store.OpenSession(ds.GetDatabase())
	if err != nil {
		return fmt.Errorf("failed to open session: %w", err)
	}
	defer session.Close()

	for _, id := range ids {
		session.Delete(id)
	}

	return session.SaveChanges()
}

// Utility Methods

// Exists checks if a document with the given ID exists
func (ds *DatabaseService) Exists(id string) (bool, error) {
	store := ds.GetStore().(*ravendb.DocumentStore)
	session, err := store.OpenSession(ds.GetDatabase())
	if err != nil {
		return false, fmt.Errorf("failed to open session: %w", err)
	}
	defer session.Close()

	var document map[string]interface{}
	err = session.Load(&document, id)
	if err != nil {
		return false, fmt.Errorf("failed to check document existence: %w", err)
	}

	return document != nil, nil
}

// CountDocuments returns the total number of documents in a collection
func (ds *DatabaseService) CountDocuments(collection string) (int, error) {
	store := ds.GetStore().(*ravendb.DocumentStore)
	session, err := store.OpenSession(ds.GetDatabase())
	if err != nil {
		return 0, fmt.Errorf("failed to open session: %w", err)
	}
	defer session.Close()

	// Simplified implementation for now
	// TODO: Implement proper document counting when RavenDB query API is clarified
	return 0, nil
}
