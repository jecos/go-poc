package main

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"
)

func testList(t *testing.T, data string, body string, expected string) {
	ParallelTestWithDb(t, data, func(t *testing.T, db *gorm.DB) {
		repo := NewMySQLRepository(db)
		router := gin.Default()
		router.POST("/occurrences/:seq_id/list", occurrencesListHandler(repo))

		req, _ := http.NewRequest("POST", "/occurrences/1/list", bytes.NewBufferString(body))
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
	ParallelTestWithDb(t, "simple", func(t *testing.T, db *gorm.DB) {
		repo := NewMySQLRepository(db)
		router := gin.Default()
		router.GET("/occurrences/:seq_id/count", occurrencesCountHandler(repo))

		req, _ := http.NewRequest("GET", "/occurrences/1/count", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"count":1}`, w.Body.String())
	})
}
