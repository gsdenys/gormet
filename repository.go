// Package gormet provides generic repository functionality for GORM ORM with validation support.
package gormet

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

// Repository represents a generic Repository for GORM ORM.
type Repository[T any] struct {
	db     *gorm.DB // db represents the GORM database connection.
	config *Config  // config holds the configuration settings for the repository.
}

// CreateRepository creates a new repository with the given database connection and configuration.
func CreateRepository[T any](db *gorm.DB, conf *Config) *Repository[T] {
	if conf == nil {
		conf = DefaultConfig() // Set default configuration if not provided.
	}

	return &Repository[T]{
		db:     db,
		config: conf,
	}
}

// New creates and returns a new repository with the given database connection.
func New[T any](db *gorm.DB) *Repository[T] {
	return (*Repository[T])(CreateRepository[T](db, nil))
}

// Create inserts a new entity into the database using the repository.
func (r *Repository[T]) Create(entity *T) error {
	if entity == nil {
		return errors.New("The entity cannot be nil")
	}

	// Validate the entity if validation is enabled in the configuration.
	if err := isValid[T](entity, r.config.Validate); err != nil {
		return err
	}

	// Create the entity in the database.
	if result := r.db.Create(entity); result.Error != nil {
		return result.Error
	}

	return nil
}

// Create inserts a new entity into the database using the repository.
func (r *Repository[T]) DeleteById(id interface{}) error {
	if id == nil {
		return errors.New("the id cannot be nil")
	}

	return r.db.Delete(new(T), id).Error
}

// Create inserts a new entity into the database using the repository.
func (r *Repository[T]) Remove(entity *T) error {

	return nil
}

func (r *Repository[T]) Update(entity *T) error {
	if entity == nil {
		return errors.New("the entity cannot be nil")
	}

	// Validate the entity if validation is enabled in the configuration.
	if err := isValid[T](entity, r.config.Validate); err != nil {
		return err
	}

	// Update the entity in the database.
	if result := r.db.Save(entity); result.Error != nil {
		return result.Error
	}

	return nil
}

// Get retrieves a single entity from the database based on the provided filter criteria.
// It takes a pointer to the repository, an entity object as a filter, and returns a pointer to the retrieved entity and an error, if any.
func (r *Repository[T]) Get(entity T) (*T, error) {
	// Create a new instance of the entity to store the retrieved data
	resp := new(T)

	// Query the database to find the first record that matches the provided filter (entity)
	result := r.db.First(resp, entity)

	// Check for errors during the database query
	if result.Error != nil {
		return nil, result.Error
	}

	// Return the retrieved entity and no errors
	return resp, nil
}

// GetById retrieves a single entity from the database based on its unique identifier (id).
// It takes a pointer to the repository and the id of the entity, and returns a pointer to the retrieved entity and an error, if any.
func (r *Repository[T]) GetById(id interface{}) (*T, error) {
	// Check if the provided id is nil
	if id == nil {
		return nil, errors.New("the id cannot be nil")
	}

	// Create a new instance of the entity to store the retrieved data
	entity := new(T)

	// Query the database to find the first record that matches the provided id
	result := r.db.First(entity, id)

	// Check for errors during the database query
	if result.Error != nil {
		return nil, result.Error
	}

	// Return the retrieved entity and no errors
	return entity, nil
}

// isValid checks the validity of the entity based on the validation configuration.
func isValid[T any](entity *T, validable bool) error {
	if validable {
		validate := validator.New()
		return validate.Struct(entity)
	}

	return nil
}
