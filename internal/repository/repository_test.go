package repository

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"go-poc/internal/types"
	"go-poc/test/testutils"
	"gorm.io/gorm"
	"os"
	"testing"
)

func TestCheckDatabaseConnection(t *testing.T) {
	testutils.ParallelTestWithDb(t, "simple", func(t *testing.T, db *gorm.DB) {
		repo := New(db)
		status := repo.CheckDatabaseConnection()
		assert.Equal(t, "up", status)

	})
}

func TestGetOccurrences(t *testing.T) {
	testutils.ParallelTestWithDb(t, "simple", func(t *testing.T, db *gorm.DB) {

		repo := New(db)
		query := types.Query{
			SelectedFields: types.OccurrencesFields,
		}
		occurrences, err := repo.GetOccurrences(1, &query)
		assert.NoError(t, err)
		if assert.Len(t, occurrences, 1) {
			assert.Equal(t, 1, occurrences[0].SeqId)
			assert.EqualValues(t, 1000, occurrences[0].LocusId)
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

func TestGetOccurrencesWithPartialColumns(t *testing.T) {
	testutils.ParallelTestWithDb(t, "simple", func(t *testing.T, db *gorm.DB) {
		repo := New(db)
		query := types.Query{

			SelectedFields: []types.Field{types.SeqIdField, types.LocusIdField, types.AdRatioField, types.FilterField},
		}
		occurrences, err := repo.GetOccurrences(1, &query)
		assert.NoError(t, err)
		if assert.Len(t, occurrences, 1) {
			assert.Equal(t, 1, occurrences[0].SeqId)
			assert.EqualValues(t, 1000, occurrences[0].LocusId)
			assert.Equal(t, "PASS", occurrences[0].Filter)
			assert.Empty(t, occurrences[0].VepImpact)
		}
	})
}

func TestGetOccurrencesWithNoColumns(t *testing.T) {
	testutils.ParallelTestWithDb(t, "simple", func(t *testing.T, db *gorm.DB) {

		repo := New(db)
		query := types.Query{}
		occurrences, err := repo.GetOccurrences(1, &query)
		assert.NoError(t, err)
		assert.Len(t, occurrences, 1)

		if assert.Len(t, occurrences, 1) {
			assert.EqualValues(t, 1000, occurrences[0].LocusId)
			assert.Empty(t, occurrences[0].Filter)
		}
	})
}

func TestCountOccurrences(t *testing.T) {
	testutils.ParallelTestWithDb(t, "simple", func(t *testing.T, db *gorm.DB) {
		repo := New(db)
		count, err := repo.CountOccurrences(1, nil)
		assert.NoError(t, err)
		assert.EqualValues(t, 1, count)
	})
}

func TestCountOccurrencesFilter(t *testing.T) {
	testutils.ParallelTestWithDb(t, "multiple", func(t *testing.T, db *gorm.DB) {

		repo := New(db)

		query := types.Query{
			Filters: &types.ComparisonNode{
				Operator: "in",
				Value:    "PASS",
				Field:    types.FilterField,
			},
			SelectedFields: types.OccurrencesFields,
		}
		c, err := repo.CountOccurrences(1, &query)

		if assert.NoError(t, err) {
			assert.EqualValues(t, 1, c)
		}
	})
}

func TestGetOccurrencesFilter(t *testing.T) {
	testutils.ParallelTestWithDb(t, "multiple", func(t *testing.T, db *gorm.DB) {

		repo := New(db)

		query := types.Query{
			Filters: &types.ComparisonNode{
				Operator: "in",
				Value:    "PASS",
				Field:    types.FilterField,
			},
			SelectedFields: types.OccurrencesFields,
		}
		occurrences, err := repo.GetOccurrences(1, &query)
		assert.NoError(t, err)
		if assert.Len(t, occurrences, 1) {
			assert.Equal(t, 1, occurrences[0].SeqId)
			assert.EqualValues(t, 1000, occurrences[0].LocusId)
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

func TestGetOccurrencesLimit(t *testing.T) {
	testutils.ParallelTestWithDb(t, "pagination", func(t *testing.T, db *gorm.DB) {

		repo := New(db)

		query := types.Query{
			SelectedFields: types.OccurrencesFields,
			Pagination: &types.Pagination{
				Limit:  5,
				Offset: 0,
			},
		}
		occurrences, err := repo.GetOccurrences(1, &query)
		assert.NoError(t, err)
		assert.Len(t, occurrences, 5)
	})
}
func TestGetOccurrencesLimitAndOffset(t *testing.T) {
	testutils.ParallelTestWithDb(t, "pagination", func(t *testing.T, db *gorm.DB) {

		repo := New(db)

		query := types.Query{
			SelectedFields: types.OccurrencesFields,
			Pagination: &types.Pagination{
				Limit:  12,
				Offset: 5,
			},
		}
		occurrences, err := repo.GetOccurrences(1, &query)
		assert.NoError(t, err)
		if assert.Len(t, occurrences, 12) {
			assert.Equal(t, occurrences[0].LocusId, 1006)
			assert.Equal(t, occurrences[:1][0].LocusId, 1018)
		}
	})
}

func TestMain(m *testing.M) {
	testutils.SetupContainer()
	code := m.Run()
	testutils.StopContainer()
	os.Exit(code)
}
