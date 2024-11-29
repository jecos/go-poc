package main

import (
	"github.com/gin-gonic/gin"
	"go-poc/models"
	"log"
)

type Occurrence = models.Occurrence
type Field = models.Field
type Table = models.Table

func main() {

	// Load configuration
	err := loadConfig("config.json")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database connection
	db, err := initDB(config)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Create repository
	repo := NewMySQLRepository(db)

	r := gin.Default()

	r.GET("/status", statusHandler(repo))
	r.POST("/occurrences/:seq_id/count", occurrencesCountHandler(repo))
	r.POST("/occurrences/:seq_id/list", occurrencesListHandler(repo))

	r.Run(":8080")
}
