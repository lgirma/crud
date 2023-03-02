package crud

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizeFilterWithEmpty(t *testing.T) {
	res := NormalizeFilter(nil, 20)
	assert.NotNil(t, res)
	assert.Equal(t, 20, res.Limit)
}

func TestNormalizeFilterWithDefaults(t *testing.T) {
	res := NormalizeFilter(Paged(2, 9), 20)
	assert.NotNil(t, res)
	assert.Equal(t, 9, res.Limit)
	assert.Equal(t, 2, res.Page)
}

func TestNormalizeFilterWithSortInfos(t *testing.T) {
	filter := Paged(2, 9)
	filter.Sort = "name:asc,age:desc,salary:desc"
	res := NormalizeFilter(filter, 20)
	assert.NotNil(t, res)
	assert.Equal(t, 9, res.Limit)
	assert.Equal(t, 2, res.Page)

	assert.Len(t, filter.SortBy, 3)
	assert.Equal(t, "name", filter.SortBy[0].Column)
	assert.False(t, filter.SortBy[0].Desc)
	assert.Equal(t, "age", filter.SortBy[1].Column)
	assert.True(t, filter.SortBy[1].Desc)
	assert.Equal(t, "salary", filter.SortBy[2].Column)
	assert.True(t, filter.SortBy[2].Desc)
}

func TestGetOrderByQuery(t *testing.T) {
	q := GetOrderByQuery(PagedAndSorted(1, 1, []SortInfo{
		{Column: "col1", Desc: false},
		{Column: "col2", Desc: true},
		{Column: "col3", Desc: false},
		{Column: "col4", Desc: true},
	}))

	assert.Equal(t, "col1,col2 desc,col3,col4 desc", q)
}
