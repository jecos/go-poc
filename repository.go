package main

import (
	"fmt"
	"github.com/Goldziher/go-utils/sliceutils"
	"go-poc/models"
	"gorm.io/gorm"
	"log"
)

type Repository interface {
	CheckDatabaseConnection() string
	GetOccurrences(seqId int, userFilter *Query) ([]Occurrence, error)
	CountOccurrences(seqId int) (int64, error)
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

func (r *MySQLRepository) GetOccurrences(seqId int, userQuery *Query) ([]Occurrence, error) {
	var columns = sliceutils.Map(userQuery.SelectedFields, func(field Field, index int, slice []Field) string {
		return fmt.Sprintf("%s.%s as %s", field.Table.Alias, field.Name, field.Name)
	})
	if columns == nil {
		columns = []string{"o.locus_id"}
	}

	tx := r.db.Table("occurrences o").Select(columns).Where("o.seq_id = ?", seqId)
	if hasFieldFromTable(userQuery.FilteredFields, models.VariantTable) || hasFieldFromTable(userQuery.SelectedFields, models.VariantTable) {
		tx = tx.Joins("JOIN variants v ON v.locus_id=o.locus_id")
	}

	if userQuery.Filters != nil {
		filters, params := userQuery.Filters.ToSQL()
		//args := append([]interface{}{seqId}, params...)
		tx.Where(filters, params...)

	}
	var occurrences []Occurrence
	err := tx.Find(&occurrences).Error
	if err != nil {
		log.Fatal("error fetching users:", err)
	}

	return occurrences, err
}

func hasFieldFromTable(fields []Field, table Table) bool {
	return sliceutils.Some(fields, func(field Field, index int, slice []Field) bool {
		return field.Table == table
	})
}

func (r *MySQLRepository) CountOccurrences(seqId int) (int64, error) {
	var count int64
	err := r.db.Table("occurrences o").Where("o.seq_id = ?", seqId).Count(&count).Error
	if err != nil {
		log.Fatal("error fetching users:", err)
	}
	return count, err
}
