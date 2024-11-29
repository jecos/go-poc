package main

import (
	"github.com/gin-gonic/gin"
	"go-poc/models"
	"net/http"
	"strconv"
)

func statusHandler(repo Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		status := repo.CheckDatabaseConnection()
		c.JSON(http.StatusOK, gin.H{
			"status": status,
		})
	}
}

func occurrencesListHandler(repo Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			body  ListBody
			query Query
		)

		// Bind JSON to the struct
		if err := c.ShouldBindJSON(&body); err != nil {
			// Return a 400 Bad Request if validation fails
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		query, err := BuildQuery(body.SelectedFields, body.SQON, &models.OccurrencesFields)
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

func occurrencesCountHandler(repo Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			body  CountBody
			query Query
		)

		// Bind JSON to the struct
		if err := c.ShouldBindJSON(&body); err != nil {
			// Return a 400 Bad Request if validation fails
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		query, err := BuildQuery(nil, body.SQON, &models.OccurrencesFields)
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

func occurrencesAggregateHandler(repo Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			body  AggregationBody
			query Query
		)

		// Bind JSON to the struct
		if err := c.ShouldBindJSON(&body); err != nil {
			// Return a 400 Bad Request if validation fails
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		selected := []string{body.Field}
		query, err := BuildQuery(selected, body.SQON, &models.OccurrencesFields)
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
