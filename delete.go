package gormet

import (
	"errors"
	"fmt"
)

// DeleteById removes an entity from the database using its ID.
//
// This method takes an ID as an argument, ensures it is not nil, and then uses GORM's Delete method
// to delete the corresponding record from the database. The repository's primary key name is used
// in the query condition. If the operation is successful, it returns nil. If the operation fails,
// an error is returned, which could be due to database connectivity issues or other constraints.
//
// Usage:
// err := repo.DeleteById(id)
//
//	if err != nil {
//	    // Handle error
//	}
//
// Parameters:
//   - id: An interface{} representing the ID of the entity to be deleted from the database.
//     It should not be nil.
//
// Returns:
// - nil if the entity is successfully deleted from the database.
// - An error if the ID is nil or if GORM encounters any issues while deleting the record.
func (r *Repository[T]) DeleteById(id interface{}) error {
	if id == nil {
		return errors.New("the ID should not be nil")
	}

	condition := fmt.Sprintf("%s = ?", r.pkName)
	deleteResult := r.db.Debug().Delete(new(T), condition, id)

	if deleteResult.Error != nil {
		return deleteResult.Error
	}

	if deleteResult.RowsAffected == 0 {
		return fmt.Errorf("no register found")
	}

	return nil
}
