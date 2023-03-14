package crud

import (
	"math"
	"strings"
)

type SortInfo struct {
	Column string
	Desc   bool
}

type DataFilter struct {
	Page   int    `form:"page"`
	Limit  int    `form:"limit"`
	Offset int    `form:"offset"`
	Sort   string `form:"sort"`
	SortBy []SortInfo
	FindBy map[string]any
}

func Paged(page int, limit int) *DataFilter {
	return &DataFilter{Page: page, Limit: limit}
}

func PagedAndSorted(page int, limit int, sortBy []SortInfo) *DataFilter {
	return &DataFilter{Page: page, Limit: limit, SortBy: sortBy}
}

type PagedList[T any] struct {
	List        []T  `json:"list"`
	TotalCount  int  `json:"totalCount"`
	Page        int  `json:"page"`
	Limit       int  `json:"limit"`
	HasNext     bool `json:"hasNext"`
	HasPrevious bool `json:"hasPrevious"`
	TotalPages  int  `json:"totalPages"`
	Skip        int  `json:"skip"`
}

func NewPagedList[T any](list []T, totalCount int, filter *DataFilter) *PagedList[T] {
	totalPages := int(math.Ceil(float64(totalCount) / float64(filter.Limit)))
	return &PagedList[T]{
		List:        list,
		TotalCount:  totalCount,
		Page:        filter.Page,
		Limit:       filter.Limit,
		HasNext:     filter.Page < totalPages-1,
		HasPrevious: filter.Page != 0,
		TotalPages:  totalPages,
		Skip:        filter.Page * filter.Limit,
	}
}

func NormalizeFilter(filter *DataFilter, defaultPageSize int) *DataFilter {
	if filter == nil {
		filter = &DataFilter{Page: 0, Limit: defaultPageSize}
	} else {
		if filter.Limit < 1 {
			filter.Limit = defaultPageSize
		}
	}
	if len(filter.Sort) > 0 {
		sortInfos := strings.Split(filter.Sort, ",")
		for _, s := range sortInfos {
			sortSpec := strings.Split(s, ":")
			sortInfo := SortInfo{Column: sortSpec[0]}
			if len(sortSpec) > 1 {
				sortInfo.Desc = strings.ToLower(sortSpec[1]) == "desc"
			}

			filter.SortBy = append(filter.SortBy, sortInfo)
		}
	}
	if filter.Offset > 0 {
		filter.Page = filter.Offset / filter.Limit
	}
	return filter
}

func GetOrderByQuery(filter *DataFilter) string {
	result := make([]string, 0)
	for _, sort := range filter.SortBy {
		item := sort.Column
		if sort.Desc {
			item += " desc"
		}
		result = append(result, item)
	}
	return strings.Join(result, ",")
}

type CrudMetadata struct {
	List   any
	Detail any
	Create any
	Update any
	Find   any
	SortBy []string
}
