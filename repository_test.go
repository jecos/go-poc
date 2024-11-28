package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"testing"
)

var occurrenceTable = Table{
	Name:  "occurrences",
	Alias: "o",
}
var variantTable = Table{
	Name:  "variants",
	Alias: "v",
}

var filterField = Field{
	Name:          "filter",
	CanBeSelected: true,
	CanBeFiltered: true,
	Table:         occurrenceTable,
}
var seqIdField = Field{
	Name:          "seq_id",
	CanBeSelected: true,
	CanBeFiltered: true,
	Table:         occurrenceTable,
}
var locusIdField = Field{
	Name:          "locus_id",
	CanBeSelected: true,
	CanBeFiltered: true,
	Table:         occurrenceTable,
}
var zygosityField = Field{
	Name:          "zygosity",
	CanBeSelected: true,
	CanBeFiltered: true,
	Table:         occurrenceTable,
}
var adRatioField = Field{
	Name:          "ad_ratio",
	CanBeSelected: true,
	CanBeFiltered: true,
	Table:         occurrenceTable,
}
var pfField = Field{
	Name:          "pf",
	CanBeSelected: true,
	CanBeFiltered: true,
	Table:         variantTable,
}
var afField = Field{
	Name:          "af",
	CanBeSelected: true,
	CanBeFiltered: true,
	Table:         variantTable,
}
var variantClassField = Field{
	Name:          "variant_class",
	CanBeSelected: true,
	CanBeFiltered: true,
	Table:         variantTable,
}
var hgvsgField = Field{
	Name:          "hgvsg",
	CanBeSelected: true,
	CanBeFiltered: true,
	Table:         variantTable,
}

func TestMySQLRepository_CheckDatabaseConnection(t *testing.T) {
	ParallelTestWithDb(t, "simple", func(t *testing.T, db *sql.DB) {
		repo := NewMySQLRepository(db)
		status := repo.CheckDatabaseConnection()
		assert.Equal(t, "up", status)

	})
}

func TestMySQLRepository_GetOccurrences(t *testing.T) {
	ParallelTestWithDb(t, "simple", func(t *testing.T, db *sql.DB) {

		repo := NewMySQLRepository(db)
		query := Query{
			SelectedFields: []Field{seqIdField, locusIdField, filterField, zygosityField, adRatioField, pfField, afField, variantClassField, hgvsgField},
		}
		occurrences, err := repo.GetOccurrences(1, &query)
		assert.NoError(t, err)
		if assert.Len(t, occurrences, 1) {
			assert.Equal(t, 1, occurrences[0].SeqId)
			assert.Equal(t, "1000", occurrences[0].LocusId)
			assert.Equal(t, "PASS", occurrences[0].Filter)
			assert.Equal(t, "HET", occurrences[0].Zygosity)
			assert.Equal(t, 0.99, occurrences[0].Pf)
			assert.Equal(t, 0.01, occurrences[0].Af)
			assert.Equal(t, "hgvsg1", occurrences[0].Hgvsg)
			assert.Equal(t, 1.0, occurrences[0].AdRatio)
			assert.Equal(t, "class1", occurrences[0].VariantClass)
		}
	})
}

func TestMySQLRepository_GetOccurrencesWithPartialColumns(t *testing.T) {
	ParallelTestWithDb(t, "simple", func(t *testing.T, db *sql.DB) {
		repo := NewMySQLRepository(db)
		query := Query{

			SelectedFields: []Field{seqIdField, locusIdField, adRatioField, pfField, afField, variantClassField, filterField},
		}
		occurrences, err := repo.GetOccurrences(1, &query)
		assert.NoError(t, err)
		if assert.Len(t, occurrences, 1) {
			assert.Equal(t, 1, occurrences[0].SeqId)
			assert.Equal(t, "1000", occurrences[0].LocusId)
			assert.Equal(t, "PASS", occurrences[0].Filter)
			assert.Empty(t, occurrences[0].VepImpact)
		}
	})
}

func TestMySQLRepository_GetOccurrencesWithNoColumns(t *testing.T) {
	ParallelTestWithDb(t, "simple", func(t *testing.T, db *sql.DB) {

		repo := NewMySQLRepository(db)
		query := Query{}
		occurrences, err := repo.GetOccurrences(1, &query)
		assert.NoError(t, err)
		assert.Len(t, occurrences, 1)

		if assert.Len(t, occurrences, 1) {
			assert.Equal(t, "1000", occurrences[0].LocusId)
			assert.Empty(t, occurrences[0].Filter)
		}
	})
}

func TestMySQLRepository_CountOccurrences(t *testing.T) {
	ParallelTestWithDb(t, "simple", func(t *testing.T, db *sql.DB) {
		repo := NewMySQLRepository(db)
		count, err := repo.CountOccurrences(1)
		assert.NoError(t, err)
		assert.Equal(t, 1, count)
	})
}

func TestMySQLRepository_GetOccurrencesFilter(t *testing.T) {
	ParallelTestWithDb(t, "multiple", func(t *testing.T, db *sql.DB) {

		repo := NewMySQLRepository(db)

		query := Query{
			Filters: &ComparisonNode{
				Operator: "in",
				Value:    "PASS",
				Field:    filterField,
			},
			SelectedFields: []Field{seqIdField, locusIdField, zygosityField, adRatioField, pfField, afField, variantClassField, filterField, hgvsgField},
		}
		occurrences, err := repo.GetOccurrences(1, &query)
		assert.NoError(t, err)
		if assert.Len(t, occurrences, 1) {
			assert.Equal(t, 1, occurrences[0].SeqId)
			assert.Equal(t, "1000", occurrences[0].LocusId)
			assert.Equal(t, "PASS", occurrences[0].Filter)
			assert.Equal(t, "HET", occurrences[0].Zygosity)
			assert.Equal(t, 0.99, occurrences[0].Pf)
			assert.Equal(t, 0.01, occurrences[0].Af)
			assert.Equal(t, "hgvsg1", occurrences[0].Hgvsg)
			assert.Equal(t, 1.0, occurrences[0].AdRatio)
			assert.Equal(t, "class1", occurrences[0].VariantClass)
		}
	})
}
