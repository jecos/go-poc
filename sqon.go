package main

import (
	"encoding/json"
)

type SQON struct {
	Field   string      `json:"field,omitempty"`   // Field to filter on (for leaf nodes)
	Value   interface{} `json:"value,omitempty"`   // Value(s) for the filter
	Content []SQON      `json:"content,omitempty"` // Nested SQON (for "not" or nested filters)
	Op      string      `json:"op,omitempty"`      // Operation at this node
}

type ListQuery struct {
	SelectedFields []string `json:"selected_fields"`
	SQON           *SQON    `json:"sqon"`
	Limit          int64    `json:"limit"`
	Offset         int64    `json:"offset"`
}

type CountQuery struct {
	SQON *SQON `json:"sqon"`
}

// Parse - parses a JSON string into an SQON structure and collects metadata for visited fields.
func Parse(jsonData string) (*SQON, error) {
	var sqon SQON
	return &sqon, json.Unmarshal([]byte(jsonData), &sqon)
}
