package crud

import "math"

type DataFilter struct {
	CurrentPage  int
	ItemsPerPage int
	SortBy       []string
}

func Paged(currentPage int, itemsPerPage int) *DataFilter {
	return &DataFilter{CurrentPage: currentPage, ItemsPerPage: itemsPerPage}
}

func PagedAndSorted(currentPage int, itemsPerPage int, sortBy []string) *DataFilter {
	return &DataFilter{CurrentPage: currentPage, ItemsPerPage: itemsPerPage, SortBy: sortBy}
}

type PagedList[T any] struct {
	List         []T
	TotalCount   int
	CurrentPage  int
	ItemsPerPage int
	HasNext      bool
	HasPrevious  bool
	TotalPages   int
	Skip         int
}

func NewPagedList[T any](list []T, totalCount int, filter *DataFilter) *PagedList[T] {
	totalPages := int(math.Ceil(float64(totalCount) / float64(filter.ItemsPerPage)))
	return &PagedList[T]{
		List:         list,
		TotalCount:   totalCount,
		CurrentPage:  filter.CurrentPage,
		ItemsPerPage: filter.ItemsPerPage,
		HasNext:      filter.CurrentPage < totalPages-1,
		HasPrevious:  filter.CurrentPage != 0,
		TotalPages:   totalPages,
		Skip:         filter.CurrentPage * filter.ItemsPerPage,
	}
}

func normalizeFilter(filter *DataFilter, defaultPageSize int) *DataFilter {
	if filter == nil {
		filter = &DataFilter{CurrentPage: 0, ItemsPerPage: defaultPageSize}
	} else {
		if filter.ItemsPerPage < 1 {
			filter.ItemsPerPage = defaultPageSize
		}
	}
	return filter
}
