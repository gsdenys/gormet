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
	retrievedEntity := new(T)

	result := r.db.First(retrievedEntity, entity)

	if result.Error != nil {
		return nil, result.Error
	}

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
	if id == nil {
		return nil, errors.New("the id should not be nil")
	}

	retrievedEntity := new(T)
	result := r.db.First(retrievedEntity, fmt.Sprintf("%s = ?", r.pkName), id)

	if result.Error != nil {
		return nil, result.Error
	}

	return retrievedEntity, nil
}

// GetLatest retrieves the latest entity from the database without any filter criteria.
// It takes a pointer to the repository and returns a pointer to the retrieved entity and an error, if any.
//
// Example:
//
//	userRepo, err := NewUserRepository(db)
//	if err != nil {
//		// Handle error
//	}
//
//	// Retrieve the latest user entity from the database
//	latestUser, err := userRepo.GetLatest()
//	if err != nil {
//		// Handle error
//	}
//
// Returns:
// - A pointer to the retrieved entity.
// - An error if the retrieval operation encounters any issues.
func (r *Repository[T]) GetLatest() (*T, error) {
	retrievedEntity := new(T)

	result := r.db.Order(fmt.Sprintf("%s DESC", r.pkName)).First(retrievedEntity)

	if result.Error != nil {
		return nil, result.Error
	}

	return retrievedEntity, nil
}
