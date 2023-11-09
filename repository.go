// Package gormet provides generic repository functionality for GORM ORM with validation support.
package gormet

import (
	"fmt"

	"gorm.io/gorm"
)

// Repository is a generic repository type that provides
// CRUD operations for a given model that is represented by a GORM model.
type Repository[T any] struct {
	db     *gorm.DB // The database connection handle.
	config *Config  // Configuration settings for the repository.
	pkName string   // The name of the primary key field in the database table.
}

// New creates and returns a new instance of Repository for a specific model type T,
// with the provided database connection and optional configuration settings.
// It automatically determines the primary key field for the model type T.
//
// Usage:
// repo, err := New[YourModelType](db, nil)
//
//	if err != nil {
//	    // Handle error
//	}
//
// Parameters:
//   - db: A *gorm.DB instance representing the database connection.
//
// Returns:
// - A pointer to a newly created Repository for type T if successful.
// - An error if there is a failure in determining the primary key or other initializations.
func New[T any](db *gorm.DB) (*Repository[T], error) {
	// Initialize a variable to hold the name of the primary key field.
	var pkName string
	var err error

	// Retrieve the primary key field name using the getPrimaryKeyFieldName function.
	// The new(T) creates a new instance of the model type T, which is required for reflection.
	if pkName, err = getPrimaryKeyFieldName(db, new(T)); err != nil {
		// If there's an error in retrieving the primary key, return nil and the error.
		return nil, fmt.Errorf("impossible to retrieve primary key: %v", err)
	}

	// Create a new Repository instance for the model type T with the database connection,
	// configuration, and primary key name.
	repo := &Repository[T]{
		db:     db,
		config: DefaultConfig(),
		pkName: pkName,
	}

	// Return the newly created repository and nil error (indicating success).
	return repo, nil
}

func (r *Repository[T]) Paginate(paginate bool) *Repository[T] {
	r.config.Paginate = paginate
	return r
}

func (r *Repository[T]) Validate(validate bool) *Repository[T] {
	r.config.Validate = validate
	return r
}
