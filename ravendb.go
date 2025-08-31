// Package ravendb provides a generic RavenDB client library for Go applications.
// It offers type-safe document operations, flexible querying, and easy-to-use interfaces.
package ravendb

import (
	"github.com/ternarybob/ravendb/interfaces"
	"github.com/ternarybob/ravendb/services"
)

// NewDatabase creates a new RavenDB database service using the provided configuration
func NewDatabase(config *Config) (interfaces.IRavenDBService, error) {
	return services.NewDatabaseService(config.URLs, config.Database)
}

// NewCollection creates a new typed collection service for the specified document type
func NewCollection[T any](database interfaces.IRavenDBService, collectionName string) interfaces.IRavenCollectionService[T] {
	return services.NewCollectionService[T](database, collectionName)
}

// Query executes a generic query on the specified collection
func Query[T any](service interfaces.IRavenDBService, collection string, options *interfaces.QueryOptions) (*interfaces.GenericQueryResult[T], error) {
	return services.Query[T](service, collection, options)
}

// QueryAll queries all documents in the specified collection
func QueryAll[T any](service interfaces.IRavenDBService, collection string) (*interfaces.GenericQueryResult[T], error) {
	return services.QueryAll[T](service, collection)
}

// QueryByField queries documents by a specific field value
func QueryByField[T any](service interfaces.IRavenDBService, collection, fieldName string, fieldValue interface{}, options *interfaces.QueryOptions) (*interfaces.GenericQueryResult[T], error) {
	return services.QueryByField[T](service, collection, fieldName, fieldValue, options)
}

// QueryByRange queries documents within a value range
func QueryByRange[T any](service interfaces.IRavenDBService, collection, fieldName string, minValue, maxValue interface{}, options *interfaces.QueryOptions) (*interfaces.GenericQueryResult[T], error) {
	return services.QueryByRange[T](service, collection, fieldName, minValue, maxValue, options)
}

// Search performs a full-text search across multiple fields
func Search[T any](service interfaces.IRavenDBService, collection, searchTerm string, searchFields []string, options *interfaces.QueryOptions) (*interfaces.GenericQueryResult[T], error) {
	return services.Search[T](service, collection, searchTerm, searchFields, options)
}
