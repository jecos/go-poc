package server

import (
	"github.com/gin-gonic/gin"
	"go-poc/internal/data"
	"go-poc/internal/types"
	"net/http"
	"strconv"
)

func StatusHandler(repo data.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		status := repo.CheckDatabaseConnection()
		c.JSON(http.StatusOK, gin.H{
			"status": status,
		})
	}
}

func OccurrencesListHandler(repo data.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			body  types.ListBody
			query types.Query
		)

		// Bind JSON to the struct
		if err := c.ShouldBindJSON(&body); err != nil {
			// Return a 400 Bad Request if validation fails
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		query, err := types.BuildQuery(body.SelectedFields, body.SQON, &types.OccurrencesFields)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		seqID, err := strconv.Atoi(c.Param("seq_id"))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		occurrences, err := repo.GetOccurrences(seqID, &query)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}

		c.JSON(http.StatusOK, occurrences)
	}
}

func OccurrencesCountHandler(repo data.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			body  types.CountBody
			query types.Query
		)

		// Bind JSON to the struct
		if err := c.ShouldBindJSON(&body); err != nil {
			// Return a 400 Bad Request if validation fails
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		query, err := types.BuildQuery(nil, body.SQON, &types.OccurrencesFields)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		seqID, err := strconv.Atoi(c.Param("seq_id"))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		count, err := repo.CountOccurrences(seqID, &query)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"count": count})
	}
}

func OccurrencesAggregateHandler(repo data.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			body  types.AggregationBody
			query types.Query
		)

		// Bind JSON to the struct
		if err := c.ShouldBindJSON(&body); err != nil {
			// Return a 400 Bad Request if validation fails
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		selected := []string{body.Field}
		query, err := types.BuildQuery(selected, body.SQON, &types.OccurrencesFields)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		seqID, err := strconv.Atoi(c.Param("seq_id"))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		aggregation, err := repo.AggregateOccurrences(seqID, &query)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
		c.JSON(http.StatusOK, aggregation)
	}
}