package main

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	data2 "go-poc/internal/repository"
	"go-poc/internal/server"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"
)

func testList(t *testing.T, data string, body string, expected string) {
	ParallelTestWithDb(t, data, func(t *testing.T, db *gorm.DB) {
		repo := data2.NewMySQLRepository(db)
		router := gin.Default()
		router.POST("/occurrences/:seq_id/list", server.OccurrencesListHandler(repo))

		req, _ := http.NewRequest("POST", "/occurrences/1/list", bytes.NewBufferString(body))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
func testCount(t *testing.T, data string, body string, expected int) {
	ParallelTestWithDb(t, data, func(t *testing.T, db *gorm.DB) {
		repo := data2.NewMySQLRepository(db)
		router := gin.Default()
		router.POST("/occurrences/:seq_id/count", server.OccurrencesCountHandler(repo))

		req, _ := http.NewRequest("POST", "/occurrences/1/count", bytes.NewBufferString(body))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, fmt.Sprintf(`{"count":%d}`, expected), w.Body.String())
	})
}
func testAggregation(t *testing.T, data string, body string, expected string) {
	ParallelTestWithDb(t, data, func(t *testing.T, db *gorm.DB) {
		repo := data2.NewMySQLRepository(db)
		router := gin.Default()
		router.POST("/occurrences/:seq_id/aggregate", server.OccurrencesAggregateHandler(repo))

		req, _ := http.NewRequest("POST", "/occurrences/1/aggregate", bytes.NewBufferString(body))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}

func TestIntegrationOccurrencesList(t *testing.T) {
	testList(t, "simple", "{}", `[{"locus_id":1000}]`)
}

func TestIntegrationOccurrencesListSqon(t *testing.T) {
	body := `{
			"selected_fields":[
				"seq_id","locus_id","filter","zygosity","pf","af","hgvsg","ad_ratio","variant_class"
			],
			"sqon":{
				"op":"in",
				"field": "filter",
				"value": "PASS"
		}
		}`
	expected := `[{"ad_ratio":1, "af":0.01, "filter":"PASS", "hgvsg":"hgvsg1", "locus_id":1000, "pf":0.99, "seq_id":1, "variant_class":"class1", "zygosity":"HET"}]`
	testList(t, "multiple", body, expected)

}

func TestIntegrationOccurrencesCount(t *testing.T) {
	testCount(t, "simple", "{}", 1)

}
func TestIntegrationCountSqon(t *testing.T) {
	body := `{
			"selected_fields":[
				"seq_id","locus_id","filter","zygosity","pf","af","hgvsg","ad_ratio","variant_class"
			],
			"sqon":{
				"op":"in",
				"field": "filter",
				"value": "PASS"
		    }
		}`
	testCount(t, "multiple", body, 1)
}

func TestIntegrationAggregation(t *testing.T) {
	body := `{
			"field": "zygosity",
			"sqon": {
				"op": "and",
				"content": [
					{
						"op":"in",
						"field": "filter",
						"value": "PASS"
					},
					{
						"op": "in",
						"field": "zygosity",
						"value": "HOM"
					}            
		
				]
			},
			"size": 10
		}`
	expected := `[{"key": "HET", "count": 2}, {"key": "HOM", "count": 1}]`
	testAggregation(t, "aggregation", body, expected)
}
