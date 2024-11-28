package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var ageMetadata = FieldMetadata{FieldName: "age", IsAllowed: true, DefaultOp: "default"}
var salaryMetadata = FieldMetadata{FieldName: "salary", IsAllowed: true, DefaultOp: "default"}
var cityMetadata = FieldMetadata{FieldName: "city", IsAllowed: true, DefaultOp: "default"}
var hobbiesMetadata = FieldMetadata{FieldName: "hobbies", IsAllowed: true, CustomOp: "array_contains"}

var fieldMetadata = []FieldMetadata{
	ageMetadata,
	salaryMetadata,
	cityMetadata,
	hobbiesMetadata,
}

func TestParseSQONToAST(t *testing.T) {
	t.Parallel()

	sqon := SQON{
		Op: "and",
		Content: []SQON{
			{Op: "in", Field: "age", Value: []interface{}{30, 40}},
			{Op: ">", Field: "salary", Value: 50000},
		},
	}

	ast, fields, err := parseSQONToAST(&sqon, &fieldMetadata)
	if assert.NoError(t, err) {
		expectedFieldMetadata :=
			[]FieldMetadata{
				{FieldName: "age", IsAllowed: true, DefaultOp: "default"},
				{FieldName: "salary", IsAllowed: true, DefaultOp: "default"},
			}
		assert.ElementsMatch(t, expectedFieldMetadata, fields)
		andNode, ok := ast.(*AndNode)
		assert.True(t, ok)
		if assert.Len(t, andNode.Children, 2) {
			compNode1, ok := andNode.Children[0].(*ComparisonNode)
			assert.True(t, ok)
			assert.Equal(t, compNode1.FieldMetadata, ageMetadata)
			assert.Equal(t, "in", compNode1.Operator)
			assert.Equal(t, []interface{}{30, 40}, compNode1.Value)

			compNode2, ok := andNode.Children[1].(*ComparisonNode)
			assert.True(t, ok)
			assert.Equal(t, compNode2.FieldMetadata, salaryMetadata)
			assert.Equal(t, ">", compNode2.Operator)
			assert.Equal(t, 50000, compNode2.Value)
		}
	}
}

func TestInvalidSQON(t *testing.T) {
	t.Parallel()
	invalidSQON := SQON{
		Op: "invalid_op",
		Content: []SQON{
			{Op: "in", Field: "age", Value: []interface{}{30, 40}},
		},
	}

	_, _, err := parseSQONToAST(&invalidSQON, &fieldMetadata)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "invalid operation: invalid_op")
}

func TestUnauthorizedFieldSQON(t *testing.T) {
	t.Parallel()
	invalidFieldSQON := SQON{
		Op: "and",
		Content: []SQON{
			{Op: "in", Field: "my_field", Value: []interface{}{30, 40}},
		},
	}

	_, _, err := parseSQONToAST(&invalidFieldSQON, &fieldMetadata)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "my_field")
	assert.ErrorContains(t, err, "unauthorized")
}

func TestParseOneField(t *testing.T) {
	t.Parallel()
	sqon := SQON{
		Op:    "in",
		Field: "age",
		Value: []interface{}{30, 40},
	}

	_, visitedFields, err := parseSQONToAST(&sqon, &fieldMetadata)
	assert.NoError(t, err)
	assert.Len(t, visitedFields, 1)
	assert.Equal(t, []FieldMetadata{{FieldName: "age", IsAllowed: true, DefaultOp: "default"}}, visitedFields)
}

func TestParseInvalidBetweenOneParam(t *testing.T) {
	t.Parallel()
	sqon := SQON{
		Op:    "between",
		Field: "age",
		Value: 30,
	}

	_, _, err := parseSQONToAST(&sqon, &fieldMetadata)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "value should be an array of 2 elements when operation is 'between'")
}
func TestParseInvalidBetweenOneParamInArray(t *testing.T) {
	t.Parallel()
	sqon := SQON{
		Op:    "between",
		Field: "age",
		Value: []interface{}{30},
	}

	_, _, err := parseSQONToAST(&sqon, &fieldMetadata)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "value array should contain exactly 2 elements when operation is 'between'")
}
func TestParseInvalidBetweenThreeParamInArray(t *testing.T) {
	t.Parallel()
	sqon := SQON{
		Op:    "between",
		Field: "age",
		Value: []interface{}{30, 40, 50},
	}

	_, _, err := parseSQONToAST(&sqon, &fieldMetadata)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "value array should contain exactly 2 elements when operation is 'between'")
}

func TestParseInvalidSingleOperator(t *testing.T) {
	t.Parallel()
	sqon := SQON{
		Op:    ">=",
		Field: "age",
		Value: []interface{}{30, 40, 50},
	}

	_, _, err := parseSQONToAST(&sqon, &fieldMetadata)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "operation >= must have exactly one value: age")
}

func TestParseSQONWithContentAndField(t *testing.T) {
	t.Parallel()
	sqon := SQON{
		Op:    "and",
		Field: "age",
		Content: []SQON{
			{Op: "in", Field: "age", Value: []interface{}{30, 40}},
		},
	}

	_, _, err := parseSQONToAST(&sqon, &fieldMetadata)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "a sqon cannot have both content and field defined")
}

func TestParseSQONWithEmptyValue(t *testing.T) {
	t.Parallel()
	sqon := SQON{
		Op:    "in",
		Field: "age",
	}

	_, _, err := parseSQONToAST(&sqon, &fieldMetadata)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "value must be defined")
}

func TestParseSQONWithInvalidNot(t *testing.T) {
	t.Parallel()
	sqon := SQON{
		Op: "not",
		Content: []SQON{
			{Op: "in", Field: "age", Value: []interface{}{30, 40}},
			{Op: "in", Field: "salary", Value: []interface{}{30, 40}},
		},
	}

	_, _, err := parseSQONToAST(&sqon, &fieldMetadata)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "'not' operation must have exactly one child")
}

func TestParseSQONtoAST(t *testing.T) {
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
			{Op: "in", Field: "hobbies", Value: []interface{}{"soccer", "hiking"}},
			{
				Op: "not",
				Content: []SQON{
					{Op: "not-in", Field: "city", Value: []interface{}{"New York", "Los Angeles"}},
				},
			},
		},
	}

	ast, visitedFields, err := parseSQONToAST(sqon, &fieldMetadata)
	assert.NoError(t, err)
	assert.NotNil(t, ast)
	assert.NotEmpty(t, visitedFields)

	expectedFields := []FieldMetadata{
		ageMetadata, salaryMetadata, cityMetadata, hobbiesMetadata,
	}
	assert.ElementsMatch(t, expectedFields, visitedFields)
	orNode, ok := ast.(*OrNode)
	assert.True(t, ok)
	assert.Len(t, orNode.Children, 4)

	compNode1, ok := orNode.Children[0].(*ComparisonNode)
	assert.True(t, ok)
	assert.Equal(t, compNode1.FieldMetadata, ageMetadata)
	assert.Equal(t, "in", compNode1.Operator)
	assert.Equal(t, []interface{}{30, 40}, compNode1.Value)

	andNode, ok := orNode.Children[1].(*AndNode)
	assert.True(t, ok)
	assert.Len(t, andNode.Children, 2)

	compNode2, ok := andNode.Children[0].(*ComparisonNode)
	assert.True(t, ok)
	assert.Equal(t, compNode2.FieldMetadata, ageMetadata)
	assert.Equal(t, "in", compNode2.Operator)
	assert.Equal(t, []interface{}{10, 20}, compNode2.Value)

	compNode3, ok := andNode.Children[1].(*ComparisonNode)
	assert.True(t, ok)
	assert.Equal(t, compNode3.FieldMetadata, salaryMetadata)
	assert.Equal(t, ">=", compNode3.Operator)
	assert.Equal(t, 50000, compNode3.Value)

	compNode4, ok := orNode.Children[2].(*ComparisonNode)
	assert.True(t, ok)
	assert.Equal(t, compNode4.FieldMetadata, hobbiesMetadata)
	assert.Equal(t, "in", compNode4.Operator)
	assert.Equal(t, []interface{}{"soccer", "hiking"}, compNode4.Value)

	notNode, ok := orNode.Children[3].(*NotNode)
	assert.True(t, ok)
	notInNode, ok := notNode.Child.(*ComparisonNode)
	assert.True(t, ok)
	assert.Equal(t, notInNode.FieldMetadata, cityMetadata)
	assert.Equal(t, "not-in", notInNode.Operator)
	assert.Equal(t, []interface{}{"New York", "Los Angeles"}, notInNode.Value)
}

func TestParseSQONOptimizeOr(t *testing.T) {
	t.Parallel()
	sqon := SQON{
		Op: "or",
		Content: []SQON{
			{Op: "in", Field: "age", Value: []interface{}{30, 40}},
		},
	}

	ast, visitedFields, err := parseSQONToAST(&sqon, &fieldMetadata)
	assert.NoError(t, err)
	assert.ElementsMatch(t, []FieldMetadata{{FieldName: "age", IsAllowed: true, DefaultOp: "default"}}, visitedFields)
	inNode, ok := ast.(*ComparisonNode)
	if assert.True(t, ok) {
		assert.Equal(t, inNode.FieldMetadata, ageMetadata)
		assert.Equal(t, "in", inNode.Operator)
		assert.Equal(t, []interface{}{30, 40}, inNode.Value)
	}
}

func TestParseSQONOptimizeAnd(t *testing.T) {
	t.Parallel()
	sqon := SQON{
		Op: "and",
		Content: []SQON{
			{Op: "in", Field: "age", Value: []interface{}{30, 40}},
		},
	}

	ast, visitedFields, err := parseSQONToAST(&sqon, &fieldMetadata)
	assert.NoError(t, err)
	assert.ElementsMatch(t, []FieldMetadata{{FieldName: "age", IsAllowed: true, DefaultOp: "default"}}, visitedFields)
	inNode, ok := ast.(*ComparisonNode)
	if assert.True(t, ok) {
		assert.Equal(t, inNode.FieldMetadata, ageMetadata)
		assert.Equal(t, "in", inNode.Operator)
		assert.Equal(t, []interface{}{30, 40}, inNode.Value)
	}
}
func TestQueryToSQLIn(t *testing.T) {
	t.Parallel()
	node := ComparisonNode{Operator: "in", Value: []interface{}{10, 20}, FieldMetadata: ageMetadata}

	sqlQuery, params := node.ToSQL()

	expectedSQL := `age IN (?, ?)`
	expectedParams := []interface{}{10, 20}

	assert.Equal(t, expectedSQL, sqlQuery)
	assert.Equal(t, expectedParams, params)
}
func TestQueryToSQLInSingleValueInArray(t *testing.T) {
	t.Parallel()
	node := ComparisonNode{Operator: "in", Value: []interface{}{10}, FieldMetadata: ageMetadata}

	sqlQuery, params := node.ToSQL()

	expectedSQL := `age = ?`
	expectedParams := []interface{}{10}

	assert.Equal(t, expectedSQL, sqlQuery)
	assert.Equal(t, expectedParams, params)
}
func TestQueryToSQLInSingleValue(t *testing.T) {
	t.Parallel()
	node := ComparisonNode{Operator: "in", Value: 10, FieldMetadata: ageMetadata}

	sqlQuery, params := node.ToSQL()

	expectedSQL := `age = ?`
	expectedParams := []interface{}{10}

	assert.Equal(t, expectedSQL, sqlQuery)
	assert.Equal(t, expectedParams, params)
}
func TestQueryToSQLInAlias(t *testing.T) {
	t.Parallel()
	node := ComparisonNode{Operator: "in", Value: []interface{}{10, 20}, FieldMetadata: FieldMetadata{FieldName: "age", IsAllowed: true, DefaultOp: "default", TableAlias: "e"}}

	sqlQuery, params := node.ToSQL()

	expectedSQL := `e.age IN (?, ?)`
	expectedParams := []interface{}{10, 20}

	assert.Equal(t, expectedSQL, sqlQuery)
	assert.Equal(t, expectedParams, params)
}
func TestQueryToSQLBetween(t *testing.T) {
	t.Parallel()
	node := ComparisonNode{Operator: "between", Value: []interface{}{30, 40}, FieldMetadata: ageMetadata}

	sqlQuery, params := node.ToSQL()

	expectedSQL := `age BETWEEN ? AND ?`
	expectedParams := []interface{}{30, 40}

	assert.Equal(t, expectedSQL, sqlQuery)
	assert.Equal(t, expectedParams, params)
}
func TestQueryCompleteToSQL(t *testing.T) {
	t.Parallel()
	node := &OrNode{
		Children: []FilterNode{
			&ComparisonNode{Operator: "in", Value: []interface{}{30, 40}, FieldMetadata: ageMetadata},
			&AndNode{
				Children: []FilterNode{
					&ComparisonNode{Operator: "in", Value: []interface{}{10, 20}, FieldMetadata: ageMetadata},
					&ComparisonNode{Operator: ">=", Value: 50000, FieldMetadata: salaryMetadata},
				},
			},
			&ComparisonNode{Operator: "in", Value: []interface{}{"soccer", "hiking"}, FieldMetadata: hobbiesMetadata},
			&NotNode{
				Child: &ComparisonNode{Operator: "not-in", Value: []interface{}{"New York", "Los Angeles"}, FieldMetadata: cityMetadata},
			},
		},
	}

	sqlQuery, params := node.ToSQL()

	expectedSQL := `(age IN (?, ?) OR (age IN (?, ?) AND salary >= ?) OR hobbies IN (?, ?) OR NOT (city NOT IN (?, ?)))`
	expectedParams := []interface{}{30, 40, 10, 20, 50000, "soccer", "hiking", "New York", "Los Angeles"}

	assert.Equal(t, expectedSQL, sqlQuery)
	assert.Equal(t, expectedParams, params)
}
