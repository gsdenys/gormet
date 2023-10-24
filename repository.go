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

// isValid checks the validity of the entity based on the validation configuration.
func isValid[T any](entity *T, validable bool) error {
	if validable {
		validate := validator.New()
		return validate.Struct(entity)
	}

	return nil
}
