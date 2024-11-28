package main

import (
	"fmt"
	"github.com/Goldziher/go-utils/sliceutils"
	"strings"
)

type FilterNode interface {
	ToSQL() (string, []interface{})
}
type FilterNodeWithChildren interface {
	ToSQL() (string, []interface{})
	GetChildren() []FilterNode
}

type AndNode struct {
	Children []FilterNode
}

func (n *AndNode) GetChildren() []FilterNode {
	return n.Children
}

type OrNode struct {
	Children []FilterNode
}

func (n *OrNode) GetChildren() []FilterNode {
	return n.Children
}

type NotNode struct {
	Child FilterNode
}

type ComparisonNode struct {
	Operator      string
	Value         interface{}
	FieldMetadata FieldMetadata
}

func (n *AndNode) ToSQL() (string, []interface{}) {
	return childrenToSQL(n, "AND")
}

func (n *OrNode) ToSQL() (string, []interface{}) {
	return childrenToSQL(n, "OR")
}

func childrenToSQL(n FilterNodeWithChildren, op string) (string, []interface{}) {
	children := n.GetChildren()
	parts := make([]string, len(children))
	var newParams []interface{}
	for i, child := range children {
		part, params := child.ToSQL()
		newParams = append(newParams, params...)
		parts[i] = part
	}
	join := strings.Join(parts, fmt.Sprintf(" %s ", op))
	return fmt.Sprintf("(%s)", join), newParams
}

func (n *NotNode) ToSQL() (string, []interface{}) {
	part, params := n.Child.ToSQL()
	return fmt.Sprintf("NOT (%s)", part), params
}

func (n *ComparisonNode) ToSQL() (string, []interface{}) {
	var (
		field  string
		params []interface{}
	)

	if n.FieldMetadata.TableAlias != "" {
		field = fmt.Sprintf("%s.%s", n.FieldMetadata.TableAlias, n.FieldMetadata.FieldName)
	} else {
		field = n.FieldMetadata.FieldName
	}

	if v, ok := n.Value.([]interface{}); ok {
		params = append(params, v...) // Flatten and append all elements
	} else {
		params = append(params, n.Value) // Append directly if not a slice
	}
	valueLength := len(params)

	switch n.Operator {
	case "in":
		placeholder := placeholders(valueLength)
		operator := "IN"
		if valueLength == 1 {
			return fmt.Sprintf("%s = %s", field, placeholder), params
		}
		return fmt.Sprintf("%s %s (%s)", field, operator, placeholder), params
	case "not-in":
		placeholder := placeholders(valueLength)
		operator := "NOT IN"
		if valueLength == 1 {
			return fmt.Sprintf("%s <> %s", field, placeholder), params
		}
		return fmt.Sprintf("%s %s (%s)", field, operator, placeholder), params
	case "<=", ">=", "<", ">":
		return fmt.Sprintf("%s %s ?", field, n.Operator), params
	case "between":
		return fmt.Sprintf("%s BETWEEN ? AND ?", field), params
	case "all":
		return "", nil //TODO: implement
	default:
		return "", nil //should not happen
	}

}
func placeholders(count int) string {
	return strings.TrimSuffix(strings.Repeat("?, ", count), ", ")
}

type Query struct {
	Filters      FilterNode
	UsedMetadata []FieldMetadata
}

func ParseSQON(sqon *SQON, metadata *[]FieldMetadata) (Query, error) {
	root, visitedMeta, err := parseSQONToAST(sqon, metadata)
	return Query{Filters: root, UsedMetadata: visitedMeta}, err
}

func parseSQONToAST(sqon *SQON, metadata *[]FieldMetadata) (FilterNode, []FieldMetadata, error) {
	if sqon.Field != "" && sqon.Content != nil {
		return nil, nil, fmt.Errorf("a sqon cannot have both content and field defined: %s", sqon.Field)
	}
	switch sqon.Op {

	case "and", "or":
		if len(sqon.Content) == 1 { // Flatten single child AND/OR nodes
			return parseSQONToAST(&sqon.Content[0], metadata)
		}
		children := make([]FilterNode, len(sqon.Content))
		var newVisitedFields []FieldMetadata
		for i, item := range sqon.Content {
			child, meta, err := parseSQONToAST(&item, metadata)
			if err != nil {
				return nil, nil, err
			}
			children[i] = child
			newVisitedFields = sliceutils.Unique(append(newVisitedFields, meta...))
		}
		if sqon.Op == "and" {
			return &AndNode{Children: children}, newVisitedFields, nil
		} else {
			return &OrNode{Children: children}, newVisitedFields, nil
		}

	case "not":
		if len(sqon.Content) != 1 {
			return nil, nil, fmt.Errorf("'not' operation must have exactly one child: %s", sqon.Field)
		}
		ast, meta, err := parseSQONToAST(&sqon.Content[0], metadata)
		if err != nil {
			return nil, nil, err
		}
		return &NotNode{Child: ast}, meta, nil

	case "in", "not-in", "<", ">", "<=", ">=", "between", "all":
		if sqon.Value == nil {
			return nil, nil, fmt.Errorf("value must be defined: %s", sqon.Field)
		}
		meta := findFieldMetaByName(metadata, sqon.Field)

		if meta == nil || !meta.IsAllowed {
			return nil, nil, fmt.Errorf("unauthorized or unknown field: %s", sqon.Field)
		}

		if sqon.Op == "between" {
			values, ok := sqon.Value.([]interface{})
			if !ok {
				return nil, nil, fmt.Errorf("value should be an array of 2 elements when operation is 'between': %s", sqon.Field)
			}
			if len(values) != 2 {
				return nil, nil, fmt.Errorf("value array should contain exactly 2 elements when operation is 'between': %s", sqon.Field)
			}
		}

		_, isMultipleValue := sqon.Value.([]interface{})
		if sqon.Op != "in" && sqon.Op != "not-in" && sqon.Op != "all" && isMultipleValue {
			return nil, nil, fmt.Errorf("operation %s must have exactly one value: %s", sqon.Op, sqon.Field)
		}

		return &ComparisonNode{
			Operator:      sqon.Op,
			Value:         sqon.Value,
			FieldMetadata: *meta,
		}, []FieldMetadata{*meta}, nil

	default:
		return nil, nil, fmt.Errorf("invalid operation: %s", sqon.Op)
	}
}

func findFieldMetaByName(metas *[]FieldMetadata, name string) *FieldMetadata {
	return sliceutils.Find(*metas, func(meta FieldMetadata, index int, slice []FieldMetadata) bool {
		return meta.FieldName == name
	})

}
