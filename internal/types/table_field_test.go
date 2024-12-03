package types

import (
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestFieldGetAliasWithNonEmptyAlias(t *testing.T) {
	t.Parallel()
	f := Field{
		Name:  "name",
		Alias: "alias",
	}
	assert.Equal(t, f.GetAlias(), "alias")
}

func TestFieldGetAliasWithEmptyAlias(t *testing.T) {
	t.Parallel()
	f := Field{
		Name: "name",
	}
	assert.Equal(t, f.GetAlias(), "name")
}

func TestFindSortedFields(t *testing.T) {
	t.Parallel()
	fields := []Field{
		{Name: "field1", CanBeSorted: true},
		{Name: "field2", CanBeSorted: false},
		{Name: "field3", CanBeSorted: true},
	}
	sorted := []SortBody{
		{Field: "field1", Order: "asc"},
		{Field: "field2", Order: "desc"},
		{Field: "field3", Order: "asc"},
		{Field: "field4", Order: "asc"},
	}
	expected := []SortField{
		{Field: fields[0], Order: "asc"},
		{Field: fields[2], Order: "asc"},
	}
	result := FindSortedFields(&fields, sorted)
	assert.Equal(t, result, expected)
}

func TestFindSortedFieldsWithBadOrder(t *testing.T) {
	t.Parallel()
	fields := []Field{
		{Name: "field1", CanBeSorted: true},
		{Name: "field2", CanBeSorted: false},
		{Name: "field3", CanBeSorted: true},
	}
	sorted := []SortBody{
		{Field: "field1", Order: "bad"},
		{Field: "field2", Order: "desc"},
		{Field: "field3", Order: "asc"},
		{Field: "field4", Order: "asc"},
	}
	expected := []SortField{
		{Field: fields[2], Order: "asc"},
	}
	result := FindSortedFields(&fields, sorted)
	assert.Equal(t, result, expected)
}
