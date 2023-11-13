package gormet

import "errors"

// Create inserts a new entity of type T into the database.
//
// This method ensures the entity is not nil before attempting to create it in the database.
// It leverages GORM's Create method, which persists the entity's data into the corresponding
// table in the database. If the operation is successful, it returns nil, indicating no error occurred.
// If the operation fails, it returns an error, which could be due to constraints like unique violations,
// missing required fields, or database connectivity issues.
//
// Usage:
// err := repo.Create(&entity)
//
//	if err != nil {
//	    // Handle error
//	}
//
// Parameters:
//   - entity: A pointer to an instance of type T that represents the entity to be created in the database.
//     The entity should be a valid non-nil pointer to a struct that GORM can map to a database table.
//
// Returns:
// - nil if the entity is successfully created in the database.
// - An error if the entity is nil or if GORM encounters any issues while creating the record.
func (r *Repository[T]) Create(entity *T) error {
	if entity == nil {
		return errors.New("the entity should not be nil")
	}

	if result := r.db.Create(entity); result.Error != nil {
		return result.Error
	}

	return nil
}
