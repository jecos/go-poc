package models

type FieldMetadata struct {
	FieldName  string // Name of the field
	IsAllowed  bool   // Whether the field is authorized for filtering
	CustomOp   string // Custom operation, e.g., "array_contains"
	DefaultOp  string // Default operation to use if no custom one exists
	TableName  string // Name of the table where the field is located
	TableAlias string // Alias of the table where the field is located
}
