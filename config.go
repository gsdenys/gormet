// Package gormet provides functionality for working with GORM ORM.
package gormet

// Config represents the configuration options for the gormet package.
type Config struct {
	// Validate indicates whether entity validation is enabled or not.
	Validate bool
}

// DefaultConfig returns a pointer to a Config with default values.
func DefaultConfig() *Config {
	return &Config{
		Validate: true, // Entity validation is enabled by default.
	}
}
