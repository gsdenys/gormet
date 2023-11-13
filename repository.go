// Package gormet provides generic repository functionality for GORM ORM with validation support.
package gormet

import (
	"fmt"

	"gorm.io/gorm"
)

// Repository is a generic repository type that provides
// CRUD operations for a given model that is represented by a GORM model.
type Repository[T any] struct {
	db       *gorm.DB // The database connection handle.
	Paginate bool     // Define if the pagination is enable or not for the repository
	pkName   string   // The name of the primary key field in the database table.
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
	//
	// >>> Why the hell is this fucking important? <<<
	//
	// This information will be used as query parameter to avoid to create request directly using object.
	// this will avoid the
	if pkName, err = getPrimaryKeyFieldName(db, new(T)); err != nil {
		// If there's an error in retrieving the primary key, return nil and the error.
		return nil, fmt.Errorf("impossible to retrieve primary key: %v", err)
	}

	// Create a new Repository instance for the model type T with the database connection,
	// configuration, and primary key name.
	repo := &Repository[T]{
		db:       db,
		Paginate: true,
		pkName:   pkName,
	}

	// Return the newly created repository and nil error (indicating success).
	return repo, nil
}

// getPrimaryKeyFieldName retrieves the name of the primary key field for a given model using the provided GORM database connection.
//
// This function takes a GORM database connection (db) and a model interface. It initializes a GORM statement (stmt) using the database connection.
// Then, it parses the model to obtain schema information. If the model cannot be parsed, an error is returned.
// The function loops through the schema fields to find the primary key. If a primary key field is found, its database name is returned.
// If no primary key is found, an error is returned indicating that no primary key field was found in the model.
//
// Usage:
// pkName, err := getPrimaryKeyFieldName(db, model)
//
//	if err != nil {
//	    // Handle error
//	}
//
// Parameters:
// - db: A *gorm.DB instance representing the database connection.
// - model: An interface{} representing the model for which the primary key field needs to be determined.
//
// Returns:
// - The name of the primary key field if found in the model.
// - An error if the model cannot be parsed or if no primary key field is found.
func getPrimaryKeyFieldName(db *gorm.DB, model interface{}) (string, error) {
	// Initialize a GORM statement using the provided database connection.
	stmt := &gorm.Statement{DB: db}

	// Parse the model to obtain schema information.
	if err := stmt.Parse(model); err != nil {
		// If parsing fails, return an empty string and the encountered error.
		return "", fmt.Errorf("data struct parse error: %s", err.Error())
	}

	// Loop through the schema fields to find the primary key.
	for _, field := range stmt.Schema.Fields {
		// If a primary key field is found, return its database name and nil error.
		if field.PrimaryKey {
			return field.DBName, nil
		}
	}

	// If no primary key field is found, return an error indicating so.
	return "", fmt.Errorf("no primary key found")
}
