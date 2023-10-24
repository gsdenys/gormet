// Package gormet provides generic repository functionality for GORM ORM with validation support.
package gormet

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

// Repository represents a generic repository for GORM ORM.
type Repository[T any] struct {
	DB     *gorm.DB // DB represents the GORM database connection.
	Config *Config  // Config holds the configuration settings for the repository.
}

// CreateRepository creates a new repository with the given database connection and configuration.
func CreateRepository[T any](db *gorm.DB, conf *Config) *Repository[T] {
	if conf == nil {
		conf = DefaultConfig() // Set default configuration if not provided.
	}

	return &Repository[T]{
		DB:     db,
		Config: conf,
	}
}

// New creates and returns a new repository with the given database connection.
func New[T any](db *gorm.DB) *Repository[T] {
	return (*Repository[T])(CreateRepository[T](db, nil))
}

// Create inserts a new entity into the database using the repository.
func (r *Repository[T]) Create(entity *T) error {
	if entity == nil {
		return errors.New("the entity cannot be nil")
	}

	// Validate the entity if validation is enabled in the configuration.
	if err := isValid[T](entity, r.Config.Validate); err != nil {
		return err
	}

	// Create the entity in the database.
	if result := r.DB.Create(entity); result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *Repository[T]) Update(entity *T) error {
	if entity == nil {
		return errors.New("the entity cannot be nil")
	}

	// Validate the entity if validation is enabled in the configuration.
	if err := isValid[T](entity, r.Config.Validate); err != nil {
		return err
	}

	// Update the entity in the database.
	if result := r.DB.Save(entity); result.Error != nil {
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
	result := r.DB.Debug().Where(&entity).First(resp)

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
	result := r.DB.First(entity, id)

	// Check for errors during the database query
	if result.Error != nil {
		return nil, result.Error
	}

	// If no records are found, return an error indicating that the entity was not found
	if result.RowsAffected == 0 {
		return nil, errors.New("entity not found")
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
