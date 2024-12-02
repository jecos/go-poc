package database

import (
	"fmt"
	"go-poc/internal/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func InitDB() (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?interpolateParams=true",
		config.Conf.DBUsername, config.Conf.DBPassword, config.Conf.DBHost, config.Conf.DBPort, config.Conf.DBName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}
