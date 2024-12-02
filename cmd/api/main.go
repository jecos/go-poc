package main

import (
	"github.com/gin-gonic/gin"
	"go-poc/internal/config"
	"go-poc/internal/data"
	"go-poc/internal/database"
	"go-poc/internal/server"
	"log"
)

func main() {

	// Load configuration
	err := config.LoadConfig("config.json")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database connection
	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Create repository
	repo := data.NewMySQLRepository(db)

	r := gin.Default()

	r.GET("/status", server.StatusHandler(repo))
	r.POST("/occurrences/:seq_id/count", server.OccurrencesCountHandler(repo))
	r.POST("/occurrences/:seq_id/list", server.OccurrencesListHandler(repo))
	r.POST("/occurrences/:seq_id/aggregate", server.OccurrencesAggregateHandler(repo))

	r.Run(":8080")
}
