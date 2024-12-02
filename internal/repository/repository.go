package repository

import (
	"fmt"
	"github.com/Goldziher/go-utils/sliceutils"
	"go-poc/internal/types"
	"gorm.io/gorm"
	"log"
)

type Occurrence = types.Occurrence
type Aggregation = types.Aggregation
type Repository interface {
	CheckDatabaseConnection() string
	GetOccurrences(seqId int, userFilter *types.Query) ([]Occurrence, error)
	CountOccurrences(seqId int, userQuery *types.Query) (int64, error)
	AggregateOccurrences(seqId int, userQuery *types.Query) ([]Aggregation, error)
}

type MySQLRepository struct {
	db *gorm.DB
}

func NewMySQLRepository(db *gorm.DB) *MySQLRepository {
	return &MySQLRepository{db: db}
}

func (r *MySQLRepository) CheckDatabaseConnection() string {
	sqlDb, err := r.db.DB()
	if err != nil {
		log.Fatal("failed to get database object:", err)
		return "down"
	}

	if err = sqlDb.Ping(); err != nil {
		log.Fatal("failed to ping database", err)
		return "down"
	}
	return "up"
}

func (r *MySQLRepository) GetOccurrences(seqId int, userQuery *types.Query) ([]Occurrence, error) {

	tx, part, err := buildQuery(seqId, userQuery, r)
	if err != nil {
		return nil, err
	}
	var columns = sliceutils.Map(userQuery.SelectedFields, func(field types.Field, index int, slice []types.Field) string {
		return fmt.Sprintf("%s.%s as %s", field.Table.Alias, field.Name, field.GetAlias())
	})
	if columns == nil {
		columns = []string{"o.locus_id"}
	}
	var occurrences []Occurrence
	if hasFieldFromTable(userQuery.FilteredFields, types.VariantTable) || hasFieldFromTable(userQuery.SelectedFields, types.VariantTable) {
		tx = tx.Select("o.locus_id").Limit(10)
		err = r.db.Table("occurrences o, variants v").
			Select(columns).
			Where("o.seq_id = ? and part=? and v.locus_id = o.locus_id and o.locus_id in (?)", seqId, part, tx).
			Find(&occurrences).Error
	} else {
		err = tx.Select(columns).Limit(10).Find(&occurrences).Error
	}
	if err != nil {
		log.Println("error fetching occurrences:", err)
	}

	return occurrences, err

}

func buildQuery(seqId int, userQuery *types.Query, r *MySQLRepository) (*gorm.DB, int, error) {
	part, err := r.GetPart(seqId)
	if err != nil {
		return nil, 0, err
	}
	tx := r.db.Table("occurrences o").Where("o.seq_id = ? and part=? and has_alt", seqId, part)
	if userQuery != nil {
		if hasFieldFromTable(userQuery.FilteredFields, types.VariantTable) || hasFieldFromTable(userQuery.SelectedFields, types.VariantTable) {
			tx = tx.Joins("JOIN variants v ON v.locus_id=o.locus_id")
		}

		if userQuery.Filters != nil {
			filters, params := userQuery.Filters.ToSQL()
			tx.Where(filters, params...)

		}
	}
	return tx, part, nil
}

func hasFieldFromTable(fields []types.Field, table types.Table) bool {
	return sliceutils.Some(fields, func(field types.Field, index int, slice []types.Field) bool {
		return field.Table == table
	})
}

func (r *MySQLRepository) CountOccurrences(seqId int, userQuery *types.Query) (int64, error) {
	tx, _, err := buildQuery(seqId, userQuery, r)
	if err != nil {
		return 0, err
	}
	var count int64
	err = tx.Count(&count).Error
	if err != nil {
		log.Print("error fetching occurrences:", err)
	}
	return count, err

}

func (r *MySQLRepository) GetPart(seqId int) (int, error) { //TODO cache
	tx := r.db.Table("sequencing_experiment").Where("seq_id = ?", seqId).Select("part")
	var part int
	err := tx.Scan(&part).Error
	if err != nil {
		log.Print("error fetching part:", err)
	}
	return part, err
}

func (r *MySQLRepository) AggregateOccurrences(seqId int, userQuery *types.Query) ([]Aggregation, error) {
	tx, _, err := buildQuery(seqId, userQuery, r)
	if err != nil {
		return nil, err
	}
	var aggregation []Aggregation
	aggCol := userQuery.SelectedFields[0].Name
	sel := fmt.Sprintf("%s as bucket, count(1) as count", aggCol)
	err = tx.Select(sel).Group(aggCol).Find(&aggregation).Error
	if err != nil {
		log.Print("error aggregation:", err)
	}
	return aggregation, err
}
