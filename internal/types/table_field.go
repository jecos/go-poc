package types

import (
	"github.com/Goldziher/go-utils/sliceutils"
)

type Table struct {
	Name  string // Name of the table
	Alias string // Alias of the table to use in query
}
type Field struct {
	Name          string // Name of the field, correspond to column name
	Alias         string // Alias of the field to use in query
	CanBeSelected bool   // Whether the field is authorized for selection
	CanBeFiltered bool   // Whether the field is authorized for filtering
	CanBeSorted   bool   // Whether the field is authorized for sorting
	CustomOp      string // Custom operation, e.g., "array_contains"
	DefaultOp     string // Default operation to use if no custom one exists
	Table         Table  // Table to which the field belongs
}

// GetAlias returns the alias of the field if it is set, otherwise returns the name
func (f *Field) GetAlias() string {
	if f.Alias != "" {
		return f.Alias
	} else {
		return f.Name
	}
}

// FindByName returns the field with the given name from the list of fields
func FindByName(fields *[]Field, name string) *Field {
	return sliceutils.Find(*fields, func(field Field, index int, slice []Field) bool {
		return field.Name == name
	})

}

// FindSelectedFields returns the fields that can be selected from the list of string field names
func FindSelectedFields(fields *[]Field, selected []string) []Field {
	var selectedFields []Field
	for _, s := range selected {
		field := FindByName(fields, s)
		if field != nil && field.CanBeSelected {
			selectedFields = append(selectedFields, *field)
		}
	}
	return selectedFields
}

// FindSortedFields returns the fields that can be sorted from the list of SortBody
func FindSortedFields(fields *[]Field, sorted []SortBody) []SortField {
	var sortedFields []SortField
	for _, sort := range sorted {
		field := FindByName(fields, sort.Field)
		if field != nil && field.CanBeSorted && (sort.Order == "asc" || sort.Order == "desc") {
			sortedFields = append(sortedFields, SortField{Field: *field, Order: sort.Order})
		}
	}
	return sortedFields

}
