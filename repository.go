package main

import (
	"database/sql"
	"fmt"
	"github.com/Goldziher/go-utils/sliceutils"
	"strings"

	"github.com/georgysavva/scany/sqlscan"
)

type Repository interface {
	CheckDatabaseConnection() string
	GetOccurrences(seqId int, userFilter *Query) ([]Occurrence, error)
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

func (r *MySQLRepository) GetOccurrences(seqId int, userQuery *Query) ([]Occurrence, error) {
	var columns = sliceutils.Map(userQuery.SelectedFields, func(field Field, index int, slice []Field) string {
		return fmt.Sprintf("%s.%s as %s", field.Table.Alias, field.Name, field.Name)
	})
	if columns == nil {
		columns = []string{"o.locus_id"}
	}
	columnsPart := strings.Join(columns, ", ")
	var (
		rows *sql.Rows
		err  error
	)

	if userQuery.Filters != nil {
		filters, params := userQuery.Filters.ToSQL()
		query := fmt.Sprintf("SELECT %s FROM occurrences o JOIN variants v ON v.locus_id=o.locus_id WHERE o.seq_id = ? AND %s", columnsPart, filters)
		args := append([]interface{}{seqId}, params...)
		rows, err = r.db.Query(query, args...)
	} else {
		query := fmt.Sprintf("SELECT %s FROM occurrences o JOIN variants v ON v.locus_id=o.locus_id WHERE o.seq_id = ?", columnsPart)
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
