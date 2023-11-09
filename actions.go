package gormet

import (
	"errors"
	"fmt"
)

// Create inserts a new entity into the database using the repository.
func (r *Repository[T]) DeleteById(id interface{}) error {
	if id == nil {
		return errors.New("the id should not be nil")
	}

	return r.db.Debug().Delete(new(T), fmt.Sprintf("%s = ?", r.pkName), id).Error
}

// Create inserts a new entity into the database using the repository.
func (r *Repository[T]) Remove(entity *T) error {

	return nil
}

func (r *Repository[T]) Update(entity *T) error {
	if entity == nil {
		return errors.New("the entity should not be nil")
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
	result := r.db.Debug().First(entity, fmt.Sprintf("%s = ?", r.pkName), id)

	// Check for errors during the database query
	if result.Error != nil {
		return nil, result.Error
	}

	// Return the retrieved entity and no errors
	return entity, nil
}
