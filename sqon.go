package main

import (
	"encoding/json"
	"fmt"
	"github.com/Goldziher/go-utils/sliceutils"
	"strings"
)

func findMetaByName(metas []FieldMetadata, name string) *FieldMetadata {
	return sliceutils.Find(metas, func(meta FieldMetadata, index int, slice []FieldMetadata) bool {
		return meta.FieldName == name
	})

}

type SQON struct {
	Field   string      `json:"field,omitempty"`   // Field to filter on (for leaf nodes)
	Value   interface{} `json:"value,omitempty"`   // Value(s) for the filter
	Content []SQON      `json:"content,omitempty"` // Nested SQON (for "not" or nested filters)
	Op      string      `json:"op,omitempty"`      // Operation at this node
}

// Parse - parses a JSON string into an SQON structure and collects metadata for visited fields.
func Parse(jsonData string, fieldMetadata []FieldMetadata) (*SQON, []FieldMetadata, error) {
	var sqon SQON
	var visitedFields []FieldMetadata

	err := json.Unmarshal([]byte(jsonData), &sqon)
	if err != nil {
		return nil, nil, err
	}

	visitedFields, err = validateAndCollectMetadata(sqon, fieldMetadata, visitedFields)
	if err != nil {
		return nil, nil, err
	}

	return &sqon, visitedFields, err
}

var validOps = map[string]bool{
	"and": true, "or": true, "not": true,
	"in": true, "not-in": true, "<=": true, ">=": true, "<": true, ">": true,
	"between": true, "all": true,
}

// validateAndCollectMetadata validates the SQON structure and collects visited field metadata.
func validateAndCollectMetadata(sqon SQON, fieldMetadata []FieldMetadata, visitedFields []FieldMetadata) ([]FieldMetadata, error) {

	// Check if the operation is valid
	if !validOps[sqon.Op] {
		return nil, fmt.Errorf("invalid operation: %s", sqon.Op)
	}

	// Special validation for "not" operation
	if sqon.Op == "not" {
		if len(sqon.Content) != 1 {
			return nil, fmt.Errorf("'not' operation must have exactly one child")
		}
	}
	if sqon.Field != "" {
		if sqon.Content != nil {
			return nil, fmt.Errorf("a sqon cannot have both content and field defined: %s", sqon.Field)
		}
		if sqon.Value == nil {
			return nil, fmt.Errorf("value must be defined: %s", sqon.Field)
		}
		// Check if the field is allowed
		meta := findMetaByName(fieldMetadata, sqon.Field)

		if meta == nil || !meta.IsAllowed {
			return nil, fmt.Errorf("unauthorized or unknown field: %s", sqon.Field)
		}
		// Add to visited fields
		return sliceutils.EnsureUniqueAndAppend(visitedFields, *meta), nil

	} else {
		var newVisitedFields []FieldMetadata
		for _, content := range sqon.Content {
			newFields, err := validateAndCollectMetadata(content, fieldMetadata, visitedFields)
			if err != nil {
				return nil, err
			}
			newVisitedFields = sliceutils.Unique(append(newVisitedFields, newFields...))
		}
		return newVisitedFields, nil
	}

}

// ToSQL generates a SQL query and its associated parameters from an SQON object.
func ToSQL(sqon *SQON, fieldMetadata []FieldMetadata) (string, []interface{}, error) {
	var params []interface{}
	sql, err := generateSQL(*sqon, fieldMetadata, &params)
	return sql, params, err
}

// generateSQL is a recursive helper function to build SQL and parameters from an SQON.
func generateSQL(sqon SQON, fieldMetadata []FieldMetadata, params *[]interface{}) (string, error) {
	switch sqon.Op {
	case "and", "or":
		var parts []string
		for _, content := range sqon.Content {
			part, err := generateSQLPart(content, fieldMetadata, params)
			if err != nil {
				return "", err
			}
			parts = append(parts, part)

		}
		return fmt.Sprintf("(%s)", strings.Join(parts, fmt.Sprintf(" %s ", strings.ToUpper(sqon.Op)))), nil
	case "not":
		content := sqon.Content[0]
		part, err := generateSQLPart(content, fieldMetadata, params)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("NOT (%s)", part), nil
	default:
		return "", fmt.Errorf("unsupported operation: %s", sqon.Op)
	}
}

// generateSQLPart generates a part of the SQL query for a given SQONContent.
func generateSQLPart(content SQON, fieldMetadata []FieldMetadata, params *[]interface{}) (string, error) {
	var part string
	var err error
	if content.Content != nil {
		part, err = generateSQL(content, fieldMetadata, params)
	} else {
		part, err = generateCondition(content, fieldMetadata, params)
	}
	if err != nil {
		return "", err
	}
	return part, nil
}

// generateCondition generates a SQL condition for a leaf node.
func generateCondition(sqon SQON, fieldMetadata []FieldMetadata, params *[]interface{}) (string, error) {
	meta := findMetaByName(fieldMetadata, sqon.Field)
	field := meta.FieldName

	switch sqon.Op {
	case "in", "not-in":
		placeholder := generatePlaceholders(len(sqon.Value.([]interface{})), params, sqon.Value.([]interface{}))
		operator := "IN"
		if sqon.Op == "not-in" {
			operator = "NOT IN"
		}
		return fmt.Sprintf("%s %s (%s)", field, operator, placeholder), nil
	case "<=", ">=", "<", ">":
		*params = append(*params, sqon.Value)
		return fmt.Sprintf("%s %s ?", field, sqon.Op), nil
	case "between":
		values := sqon.Value.([]interface{})
		*params = append(*params, values[0], values[1])
		return fmt.Sprintf("%s BETWEEN ? AND ?", field), nil
	default:
		return "", fmt.Errorf("unsupported operation for field '%s': %s", field, sqon.Op)
	}
}

// generatePlaceholders generates SQL placeholders for a list of values.
func generatePlaceholders(count int, params *[]interface{}, values []interface{}) string {
	for _, value := range values {
		*params = append(*params, value)
	}
	return strings.TrimSuffix(strings.Repeat("?, ", count), ", ")
}
