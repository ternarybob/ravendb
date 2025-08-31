package ravendb

// Config holds configuration for RavenDB connection
type Config struct {
	URLs     []string `json:"urls"`
	Database string   `json:"database"`
}

// NewConfig creates a new configuration with default values
func NewConfig(urls []string, database string) *Config {
	return &Config{
		URLs:     urls,
		Database: database,
	}
}

// NewSingleNodeConfig creates a configuration for a single-node setup
func NewSingleNodeConfig(url, database string) *Config {
	return &Config{
		URLs:     []string{url},
		Database: database,
	}
}

// NewLocalConfig creates a configuration for local development
func NewLocalConfig(database string) *Config {
	return &Config{
		URLs:     []string{"http://localhost:8080"},
		Database: database,
	}
}