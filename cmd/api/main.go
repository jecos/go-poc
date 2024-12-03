package main

import (
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"go-poc/internal/database"
	"go-poc/internal/repository"
	"go-poc/internal/server"
	"log"
)

func main() {

	// Initialize database connection
	db, err := database.New()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Create repository
	repo := repository.New(db)

	r := gin.Default()
	r.Use(gzip.Gzip(gzip.DefaultCompression))

	r.GET("/status", server.StatusHandler(repo))
	r.POST("/occurrences/:seq_id/count", server.OccurrencesCountHandler(repo))
	r.POST("/occurrences/:seq_id/list", server.OccurrencesListHandler(repo))
	r.POST("/occurrences/:seq_id/aggregate", server.OccurrencesAggregateHandler(repo))

	r.Run(":8080")
}
