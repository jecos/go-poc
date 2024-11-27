package main

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/georgysavva/scany/sqlscan"
)

type Repository interface {
	CheckDatabaseConnection() string
	GetOccurrences(seqId int, selectedCols []string, userFilter *Filter, joinedTables []string) ([]Occurrence, error)
	CountOccurrences(seqId int) (int, error)
}

type MySQLRepository struct {
	db *sql.DB
}

func NewMySQLRepository(db *sql.DB) *MySQLRepository {
	return &MySQLRepository{db: db}
}

func (r *MySQLRepository) CheckDatabaseConnection() string {
	if err := r.db.Ping(); err != nil {
		return "down"
	}
	return "up"
}

type Filter struct {
	userFilters string
	userParams  []interface{}
}

func (r *MySQLRepository) GetOccurrences(seqId int, selectedCols []string, userFilter *Filter, joinedTables []string) ([]Occurrence, error) {
	if len(selectedCols) == 0 || (len(selectedCols) == 1 && selectedCols[0] == "") {
		selectedCols = []string{"locus_id"}
	}
	// Define allowed selectedCols
	allowedColumns := map[string]bool{
		"seq_id":                 true,
		"locus_id":               true,
		"quality":                true,
		"filter":                 true,
		"zygosity":               true,
		"pf":                     true,
		"af":                     true,
		"gnomad_v3_af":           true,
		"hgvsg":                  true,
		"omim_inheritance_code":  true,
		"ad_ratio":               true,
		"variant_class":          true,
		"vep_impact":             true,
		"symbol":                 true,
		"clinvar_interpretation": true,
		"mane_select":            true,
		"canonical":              true,
	}
	// Validate requested selectedCols
	var validColumns []string
	for _, col := range selectedCols {
		if allowedColumns[col] {
			validColumns = append(validColumns, col)
		} else {
			return nil, fmt.Errorf("invalid column: %s", col)
		}
	}
	var (
		rows *sql.Rows
		err  error
	)

	if userFilter != nil {
		query := fmt.Sprintf("SELECT %s FROM occurrences where seq_id = ? and %s", strings.Join(validColumns, ", "), userFilter.userFilters)
		args := append([]interface{}{seqId}, userFilter.userParams...)
		rows, err = r.db.Query(query, args...)
	} else {
		query := fmt.Sprintf("SELECT %s FROM occurrences where seq_id = ?", strings.Join(validColumns, ", "))
		rows, err = r.db.Query(query, seqId)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var occurrences []Occurrence
	err = sqlscan.ScanAll(&occurrences, rows)
	if err != nil {
		return nil, err
	}

	return occurrences, nil
}

func (r *MySQLRepository) CountOccurrences(seqId int) (int, error) {
	query := `SELECT COUNT(1) FROM occurrences where seq_id = ?`
	var count int
	err := r.db.QueryRow(query, seqId).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
