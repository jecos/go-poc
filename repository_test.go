package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"testing"
)

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
		columns := []string{"seq_id", "locus_id", "filter", "zygosity", "pf", "af", "hgvsg", "ad_ratio", "variant_class"}
		occurrences, err := repo.GetOccurrences(1, columns, nil, nil)
		assert.NoError(t, err)
		assert.Len(t, occurrences, 1)
		assert.Equal(t, 1, occurrences[0].SeqId)
		assert.Equal(t, "locus1", occurrences[0].LocusId)
		assert.Equal(t, "PASS", occurrences[0].Filter)
		assert.Equal(t, "HET", occurrences[0].Zygosity)
		assert.Equal(t, 0.99, occurrences[0].Pf)
		assert.Equal(t, 0.01, occurrences[0].Af)
		assert.Equal(t, "hgvsg1", occurrences[0].Hgvsg)
		assert.Equal(t, 1.0, occurrences[0].AdRatio)
		assert.Equal(t, "class1", occurrences[0].VariantClass)
	})
}

func TestMySQLRepository_GetOccurrencesWithPartialColumns(t *testing.T) {
	ParallelTestWithDb(t, "simple", func(t *testing.T, db *sql.DB) {
		repo := NewMySQLRepository(db)
		columns := []string{"seq_id", "locus_id", "filter"}
		occurrences, err := repo.GetOccurrences(1, columns, nil, nil)
		assert.NoError(t, err)
		assert.Len(t, occurrences, 1)
		assert.Equal(t, 1, occurrences[0].SeqId)
		assert.Equal(t, "locus1", occurrences[0].LocusId)
		assert.Equal(t, "PASS", occurrences[0].Filter)
		assert.Empty(t, occurrences[0].VepImpact)
	})
}

func TestMySQLRepository_GetOccurrencesWithNoColumns(t *testing.T) {
	ParallelTestWithDb(t, "simple", func(t *testing.T, db *sql.DB) {

		repo := NewMySQLRepository(db)
		var columns []string
		occurrences, err := repo.GetOccurrences(1, columns, nil, nil)
		assert.NoError(t, err)
		assert.Len(t, occurrences, 1)
		assert.Equal(t, "locus1", occurrences[0].LocusId)
		assert.Empty(t, occurrences[0].Filter)
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
