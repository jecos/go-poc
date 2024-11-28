package main

import (
	"github.com/gin-gonic/gin"
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
		//columnsParam := c.Query("columns")
		//columns := strings.Split(columnsParam, ",")
		seqID, err := strconv.Atoi(c.Param("seq_id"))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		occurrences, err := repo.GetOccurrences(seqID, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, occurrences)
	}
}

func occurrencesCountHandler(repo Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		seqID, err := strconv.Atoi(c.Param("seq_id"))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		count, err := repo.CountOccurrences(seqID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"count": count})
	}
}
