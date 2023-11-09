// Package gormet provides functionality for working with GORM ORM.
package gormet

// Config represents the configuration options for the gormet package.
type Config struct {
	Paginate bool
}

// DefaultConfig returns a pointer to a Config with default values.
func DefaultConfig() *Config {
	return &Config{
		Paginate: true, // Entity paginated search is
	}
}
