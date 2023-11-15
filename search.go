package gormet

// Pagination contains the response with actual content and additional data for paginated searches.
type Pagination[T any] struct {
	Response   Response[T]    `json:"response"`
	criteria   interface{}    `json:"-"`
	repository *Repository[T] `json:"-"`
}

// Response represents the paginated search results including entities, total count, and pagination details.
type Response[T any] struct {
	Entities    []T   `json:"entities"`
	TotalCount  int64 `json:"totalCount"`
	Page        uint  `json:"page"`
	PageSize    uint  `json:"pageSize"`
	TotalPages  int64 `json:"totalPages"`
	HasNextPage bool  `json:"hasNextPage"`
	HasPrevPage bool  `json:"hasPrevPage"`
}

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
func (r *Repository[T]) Search(page uint, query interface{}, args ...interface{}) (Pagination[T], error) {
	var offset int = getOffset(page, r.PageSize)
	var limit int = getLimit(r.PageSize)

	var entities []T
	var count int64
	var err error

	if entities, err = r.executeSearch(offset, limit, query, args...); err != nil {
		return Pagination[T]{}, err
	}

	if count, err = r.countRows(query, args...); err != nil {
		return Pagination[T]{}, err
	}

	// Create a Pagination with the initial response, criteria, and repository.
	pagination := Pagination[T]{
		Response: Response[T]{
			Entities:    entities,
			TotalCount:  count,
			Page:        page,
			PageSize:    r.PageSize,
			TotalPages:  countTotalPages(count, limit, r.PageSize),
			HasNextPage: getHasNextPage(page, limit, count),
			HasPrevPage: getHasPreviousPage(page),
		},
		criteria:   query,
		repository: r,
	}

	return pagination, nil
}

// executeSearch performs the paginated search using GORM's Find method.
func (r *Repository[T]) executeSearch(offset int, limit int, query interface{}, args ...interface{}) ([]T, error) {
	entities := make([]T, 0)
	searchResult := r.db.Debug().Where(query, args...).Offset(offset).Limit(limit).Find(&entities)

	return entities, searchResult.Error
}

// countRows gets the total count for the entire search without pagination.
func (r *Repository[T]) countRows(query interface{}, args ...interface{}) (int64, error) {
	var totalCount int64
	result := r.db.Model(new(T)).Where(query, args...).Count(&totalCount)

	return totalCount, result.Error
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
func (r *Repository[T]) SearchAll(query interface{}, args ...interface{}) ([]T, error) {
	var entities []T
	var err error

	if entities, err = r.executeSearch(-1, -1, query, args...); err != nil {
		return []T{}, err
	}

	return entities, nil
}

// getOffset calculates the offset based on the page number and page size.
func getOffset(page uint, pageSize uint) int {
	if page == 0 {
		return -1
	}

	return int(pageSize * (page - 1))
}

// getLimit calculates the limit based on the page size.
func getLimit(pageSize uint) int {
	if pageSize == 0 {
		return -1
	}

	return int(pageSize)
}

// countTotalPages calculates the total number of pages based on total count, limit, and page size.
func countTotalPages(totalCount int64, limit int, pageSize uint) int64 {
	return (totalCount + int64(limit-1)) / int64(pageSize)
}

// getHasNextPage checks if there is a next page based on the current page, limit, and total count.
func getHasNextPage(page uint, limit int, totalCount int64) bool {
	return int64(int(page)*limit) < totalCount
}

// getHasPreviousPage checks if there is a previous page based on the current page.
func getHasPreviousPage(page uint) bool {
	return page > 1
}
