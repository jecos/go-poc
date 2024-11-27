package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var fieldMetadata = []FieldMetadata{
	{FieldName: "age", IsAllowed: true, DefaultOp: "default"},
	{FieldName: "salary", IsAllowed: true, DefaultOp: "default"},
	{FieldName: "city", IsAllowed: true, DefaultOp: "default"},
	{FieldName: "clinvar_interpretations", IsAllowed: true, CustomOp: "array_contains"},
}

func TestParseSQON(t *testing.T) {
	t.Parallel()
	jsonData := `{
		"op": "or",
		"content": [
			{ "op": "in", "field": "age", "value": [30, 40] },
			{ "op": "and", "content": [
				{ "op": "in", "field": "age", "value": [10, 20] },
				{ "op": ">=", "field": "salary", "value": 50000 }
			]},
			{ "op": "in", "field": "clinvar_interpretations", "value": ["pathogenic", "likely_pathogenic"] },
			{ "op": "not", "content": [
				{ "op": "not-in", "field": "city", "value": ["New York", "Los Angeles"] }
			]}
		]
	}`

	sqon, visitedFields, err := Parse(jsonData, fieldMetadata)
	assert.NoError(t, err)
	assert.NotNil(t, sqon)
	assert.NotEmpty(t, visitedFields)

	expectedFields := []FieldMetadata{
		{FieldName: "age", IsAllowed: true, DefaultOp: "default"},
		{FieldName: "salary", IsAllowed: true, DefaultOp: "default"},
		{FieldName: "clinvar_interpretations", IsAllowed: true, CustomOp: "array_contains"},
		{FieldName: "city", IsAllowed: true, DefaultOp: "default"},
	}
	assert.ElementsMatch(t, expectedFields, visitedFields)
}

func TestToSQL(t *testing.T) {
	t.Parallel()
	sqon := &SQON{
		Op: "or",
		Content: []SQON{
			{Op: "in", Field: "age", Value: []interface{}{30, 40}},
			{
				Op: "and",
				Content: []SQON{
					{Op: "in", Field: "age", Value: []interface{}{10, 20}},
					{Op: ">=", Field: "salary", Value: 50000},
				},
			},
			{Op: "in", Field: "clinvar_interpretations", Value: []interface{}{"pathogenic", "likely_pathogenic"}},
			{
				Op: "not",
				Content: []SQON{
					{Op: "not-in", Field: "city", Value: []interface{}{"New York", "Los Angeles"}},
				},
			},
		},
	}

	sqlQuery, params, err := ToSQL(sqon, fieldMetadata)
	assert.NoError(t, err)

	expectedSQL := `(age IN (?, ?) OR (age IN (?, ?) AND salary >= ?) OR clinvar_interpretations IN (?, ?) OR NOT (city NOT IN (?, ?)))`
	expectedParams := []interface{}{30, 40, 10, 20, 50000, "pathogenic", "likely_pathogenic", "New York", "Los Angeles"}

	assert.Equal(t, expectedSQL, sqlQuery)
	assert.Equal(t, expectedParams, params)
}

func TestInvalidSQON(t *testing.T) {
	t.Parallel()
	invalidJSON := `{
		"op": "invalid_op",
		"content": [
			{ "op": "in", "field": "age", "value": [30, 40] }
		]
	}`

	_, _, err := Parse(invalidJSON, fieldMetadata)
	assert.Error(t, err)
}

func TestUnauthorizedField(t *testing.T) {
	t.Parallel()
	invalidFieldJSON := `{
		"op": "and",
		"content": [{
			"op": "in",
			"field": "my_field",
			"value": [30, 40]
		}]
}`

	_, _, err := Parse(invalidFieldJSON, fieldMetadata)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "my_field")
	assert.ErrorContains(t, err, "unauthorized")
}

func TestOneField(t *testing.T) {
	t.Parallel()
	jsonData := `{
		"op": "in",
		"field": "age",
		"value": [30, 40]
	}`

	_, visitedFields, err := Parse(jsonData, fieldMetadata)
	assert.NoError(t, err)
	assert.Len(t, visitedFields, 1)
	assert.Equal(t, []FieldMetadata{{FieldName: "age", IsAllowed: true, DefaultOp: "default"}}, visitedFields)
}

func TestSQONWithContentAndField(t *testing.T) {
	t.Parallel()
	jsonData := `{
		"op": "and",
		"field": "age",
		"content": [{
			"op": "in",
			"field": "age",
			"value": [30, 40]
		}]
	}`

	_, _, err := Parse(jsonData, fieldMetadata)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "a sqon cannot have both content and field defined")
}

func TestSQONWithEmptyValue(t *testing.T) {
	t.Parallel()
	jsonData := `{
		"op": "in",
		"field": "age"
	}`

	_, _, err := Parse(jsonData, fieldMetadata)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "value must be defined")
}
