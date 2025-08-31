package interfaces

// QueryOptions provides flexible query configuration
type QueryOptions struct {
	Skip         int                    `json:"skip,omitempty"`
	Take         int                    `json:"take,omitempty"`
	OrderBy      string                 `json:"orderBy,omitempty"`
	OrderDesc    bool                   `json:"orderDesc,omitempty"`
	WhereClause  string                 `json:"whereClause,omitempty"`
	Parameters   map[string]interface{} `json:"parameters,omitempty"`
	IncludeTotal bool                   `json:"includeTotal,omitempty"`
}

// QueryResult contains paginated query results
type QueryResult struct {
	Results    []interface{} `json:"results"`
	TotalCount int           `json:"totalCount,omitempty"`
	Skip       int           `json:"skip"`
	Take       int           `json:"take"`
	HasMore    bool          `json:"hasMore"`
}

// GenericQueryResult contains typed paginated query results using generics
type GenericQueryResult[T any] struct {
	Results    []T  `json:"results"`
	TotalCount int  `json:"totalCount,omitempty"`
	Skip       int  `json:"skip"`
	Take       int  `json:"take"`
	HasMore    bool `json:"hasMore"`
}

// IRavenDBService defines the comprehensive interface for RavenDB operations
type IRavenDBService interface {
	// Database lifecycle
	Init() error
	InitializeWithSeeding(seedData bool) error
	Close() error
	GetDatabaseStatus() (map[string]interface{}, error)

	// Basic CRUD operations
	Store(id string, document interface{}) error
	LoadByID(id string, result interface{}) error
	Delete(id string) error

	// Enhanced CRUD operations
	StoreMultiple(documents map[string]interface{}) error
	LoadMultipleByIDs(ids []string, results interface{}) error
	Update(id string, updates map[string]interface{}) error
	DeleteMultiple(ids []string) error

	// Utility methods
	Exists(id string) (bool, error)
	CountDocuments(collection string) (int, error)

	// Additional database service specific methods
	GetStore() interface{} // Returns the underlying DocumentStore as interface{}
	GetDatabase() string
}

// IRavenCollectionService defines the interface for generic collection operations
type IRavenCollectionService[T any] interface {
	// CRUD Operations
	Store(id string, document T) error
	StoreMultiple(documents map[string]T) error
	LoadByID(id string) (*T, error)
	LoadMultipleByIDs(ids []string) ([]T, error)
	Update(id string, document T) error
	Delete(id string) error
	DeleteMultiple(ids []string) error

	// Query Operations
	Query(options *QueryOptions) (*GenericQueryResult[T], error)
	QueryAll() (*GenericQueryResult[T], error)
	QueryByField(fieldName string, fieldValue interface{}, options *QueryOptions) (*GenericQueryResult[T], error)
	QueryByRange(fieldName string, minValue, maxValue interface{}, options *QueryOptions) (*GenericQueryResult[T], error)
	Search(searchTerm string, searchFields []string, options *QueryOptions) (*GenericQueryResult[T], error)

	// Utility Operations
	Exists(id string) (bool, error)
	Count() (int, error)
}