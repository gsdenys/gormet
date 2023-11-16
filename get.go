package gormet

import (
	"errors"
	"fmt"
)

// Get retrieves a single entity from the database based on the provided filter criteria.
// It takes a pointer to the repository, an entity object as a filter, and returns a pointer to the retrieved entity and an error, if any.
//
// Example:
//
//	userRepo, err := NewUserRepository(db)
//	if err != nil {
//		// Handle error
//	}
//
//	// Create a sample user entity for filtering
//	filterUser := User{ID: 1}
//
//	// Retrieve the user entity from the database based on the filter
//	retrievedUser, err := userRepo.Get(filterUser)
//	if err != nil {
//		// Handle error
//	}
//
// Parameters:
// - entity: An object of the entity type with the filter criteria.
//
// Returns:
// - A pointer to the retrieved entity.
// - An error if the retrieval operation encounters any issues.
func (r *Repository[T]) Get(entity T) (*T, error) {
	// Create a new instance of the entity to store the retrieved data
	retrievedEntity := new(T)

	// Query the database to find the first record that matches the provided filter (entity)
	result := r.db.First(retrievedEntity, entity)

	// Check for errors during the database query
	if result.Error != nil {
		return nil, result.Error
	}

	// Return the retrieved entity and no errors
	return retrievedEntity, nil
}

// GetById retrieves a single entity from the database based on its unique identifier (id).
// It takes a pointer to the repository and the id of the entity, and returns a pointer to the retrieved entity and an error, if any.
//
// Example:
//
//	userRepo, err := NewUserRepository(db)
//	if err != nil {
//		// Handle error
//	}
//
//	// Provide the unique identifier for the user to be retrieved
//	userId := 1
//
//	// Retrieve the user entity from the database based on the unique identifier
//	retrievedUser, err := userRepo.GetById(userId)
//	if err != nil {
//		// Handle error
//	}
//
// Parameters:
// - id: The unique identifier of the entity.
//
// Returns:
// - A pointer to the retrieved entity.
// - An error if the retrieval operation encounters any issues, including if the provided id is nil.
func (r *Repository[T]) GetById(id interface{}) (*T, error) {
	// Check if the provided id is nil
	if id == nil {
		return nil, errors.New("the id should not be nil")
	}

	// Create a new instance of the entity to store the retrieved data
	retrievedEntity := new(T)

	// Query the database to find the first record that matches the provided id
	result := r.db.First(retrievedEntity, fmt.Sprintf("%s = ?", r.pkName), id)

	// Check for errors during the database query
	if result.Error != nil {
		return nil, result.Error
	}

	// Return the retrieved entity and no errors
	return retrievedEntity, nil
}
