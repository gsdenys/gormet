package gormet

// Pagination contains the response with actual content and additional data for paginated searches.
type Pagination[T any] struct {
	Response   Response[T]    `json:"response"`
	criteria   interface{}    `json:"-"`
	repository *Repository[T] `json:"-"`
}

// func (p Pagination[T]) MarshalJSON() ([]byte, error) {
// 	if p.repository.PageSize == 0 {
// 		// If PageSize is not specified, marshal only the Entities and TotalCount from Response[T].
// 		return json.Marshal(struct {
// 			Entities   []T  `json:"entities"`
// 			TotalCount uint `json:"totalCount"`
// 		}{
// 			Entities:   p.Response.Entities,
// 			TotalCount: p.Response.TotalCount,
// 		})
// 	}

// 	return json.Marshal(p.Response)
// }

// PaginatedResponse represents the paginated search results including entities, total count, and pagination details.
type Response[T any] struct {
	Entities    []T   `json:"entities"`
	TotalCount  int64 `json:"totalCount"`
	Page        int   `json:"page"`
	PageSize    int   `json:"pageSize"`
	TotalPages  int64 `json:"totalPages"`
	HasNextPage bool  `json:"hasNextPage"`
	HasPrevPage bool  `json:"hasPrevPage"`
}

// func (r Response[T]) MarshalJSON() ([]byte, error) {
// 	if r.PageSize == 0 {
// 		return json.Marshal(struct {
// 			Entities   []T  `json:"entities"`
// 			TotalCount uint `json:"totalCount"`
// 		}{
// 			Entities:   r.Entities,
// 			TotalCount: r.TotalCount,
// 		})
// 	}

// 	return json.Marshal(r)
// }

// Search performs a paginated search for entities in the database based on given criteria.
//
// This method takes a GORM query condition, performs a paginated search using GORM's Find method,
// and returns the paginated results. The paginated results include the entities found and additional
// information such as total count and pagination details.
//
// Usage:
// query := "your_column = ?"
// args := []interface{}{"your_value"}
// pagination, err := repo.Search(page, query, args...)
//
//	if err != nil {
//	    // Handle error
//	}
//
// // Use pagination for further processing
//
// Parameters:
// - page: The page number for pagination (starting from 1).
// - query: GORM query condition.
// - args: Arguments for the query condition.
//
// Returns:
// - A structure containing the paginated search results, including entities, total count, and pagination details.
// - An error if the search operation encounters any issues.
func (r *Repository[T]) Search(page int, query interface{}, args ...interface{}) (Pagination[T], error) {
	var offset int = int(r.PageSize) * (page - 1)
	var limit int = int(r.PageSize)

	// Perform the paginated search using GORM's Find method.
	// Debug mode is enabled for more detailed logs during development.
	entities := make([]T, 0)
	searchResult := r.db.Debug().Where(query, args...).Offset(offset).Limit(limit).Find(&entities)

	// Check if the Find operation encountered an error.
	if searchResult.Error != nil {
		// Return the encountered error.
		return Pagination[T]{}, searchResult.Error
	}

	// Get the total count for the entire search without pagination.
	var totalCount int64
	r.db.Model(new(T)).Where(query, args...).Count(&totalCount)

	// Create a PaginatedResponse for the initial response.
	response := Response[T]{
		Entities:    entities,
		TotalCount:  totalCount,
		Page:        page,
		PageSize:    int(r.PageSize),
		TotalPages:  (totalCount + int64(limit-1)) / int64(r.PageSize),
		HasNextPage: int64(page*limit) < totalCount,
		HasPrevPage: page > 1,
	}

	// Create a Pagination with the initial response, criteria, and repository.
	pagination := Pagination[T]{
		Response:   response,
		criteria:   query,
		repository: r,
	}

	return pagination, nil
}

// SearchAll performs a paginated search for all entities in the database based on given criteria.
//
// This method takes a GORM query condition, performs a paginated search using GORM's Find method,
// and returns the paginated results. The paginated results include all entities found and additional
// information such as total count and pagination details.
//
// Usage:
// query := "your_column = ?"
// args := []interface{}{"your_value"}
// pagination, err := repo.SearchAll(query, args...)
//
//	if err != nil {
//	    // Handle error
//	}
//
// // Use pagination for further processing
//
// Parameters:
// - query: GORM query condition.
// - args: Arguments for the query condition.
//
// Returns:
// - A structure containing the paginated search results, including entities, total count, and pagination details.
// - An error if the search operation encounters any issues.
func (r *Repository[T]) SearchAll(query interface{}, args ...interface{}) (Pagination[T], error) {
	return r.Search(0, query, args...)
}
