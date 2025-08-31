package services

import (
	"fmt"
	"strings"

	"github.com/ravendb/ravendb-go-client"
	"github.com/ternarybob/ravendb/interfaces"
)

// Query is a generic method that queries documents of a specific type T.
func Query[T any](service interfaces.IRavenDBService, collection string, options *interfaces.QueryOptions) (*interfaces.GenericQueryResult[T], error) {
	store := service.GetStore().(*ravendb.DocumentStore)
	session, err := store.OpenSession(service.GetDatabase())
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
	rqlQuery.WriteString(fmt.Sprintf("from @all_docs where @metadata.'@collection' = '%s'", collection))

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

// QueryAll is a generic method that queries all documents of a specific type
func QueryAll[T any](service interfaces.IRavenDBService, collection string) (*interfaces.GenericQueryResult[T], error) {
	options := &interfaces.QueryOptions{
		Skip: 0,
		Take: 1024, // Default page size
	}

	return Query[T](service, collection, options)
}

// QueryByField is a generic method that queries documents by a specific field value
func QueryByField[T any](service interfaces.IRavenDBService, collection, fieldName string, fieldValue interface{}, options *interfaces.QueryOptions) (*interfaces.GenericQueryResult[T], error) {
	if options == nil {
		options = &interfaces.QueryOptions{}
	}

	// Build where clause using the actual field name
	options.WhereClause = fmt.Sprintf("%s = $value", fieldName)
	if options.Parameters == nil {
		options.Parameters = make(map[string]interface{})
	}
	options.Parameters["value"] = fieldValue

	return Query[T](service, collection, options)
}

// QueryByRange is a generic method that queries documents within a range of values
func QueryByRange[T any](service interfaces.IRavenDBService, collection, fieldName string, minValue, maxValue interface{}, options *interfaces.QueryOptions) (*interfaces.GenericQueryResult[T], error) {
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

	return Query[T](service, collection, options)
}

// Search is a generic method that performs a full-text search across documents
func Search[T any](service interfaces.IRavenDBService, collection, searchTerm string, searchFields []string, options *interfaces.QueryOptions) (*interfaces.GenericQueryResult[T], error) {
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

	return Query[T](service, collection, options)
}
