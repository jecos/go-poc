package types

import (
	"github.com/Goldziher/go-utils/sliceutils"
	"slices"
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
	CustomOp      string // Custom operation, e.g., "array_contains"
	DefaultOp     string // Default operation to use if no custom one exists
	Table         Table  // Table to which the field belongs
}

func (f *Field) GetAlias() string {
	if f.Alias != "" {
		return f.Alias
	} else {
		return f.Name
	}
}
func FindByName(fields *[]Field, name string) *Field {
	return sliceutils.Find(*fields, func(field Field, index int, slice []Field) bool {
		return field.Name == name
	})

}
func FindSelectedFields(fields *[]Field, selected []string) []Field {
	return sliceutils.Filter(*fields, func(field Field, index int, slice []Field) bool {
		return field.CanBeSelected && slices.Contains(selected, field.Name)
	})

}
