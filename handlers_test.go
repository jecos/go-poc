package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type MockRepository struct{}

func (m *MockRepository) CheckDatabaseConnection() string {
	return "up"
}

func (m *MockRepository) GetOccurrences(int, []string, *Filter, []string) ([]Occurrence, error) {
	return []Occurrence{
		{
			SeqId:        1,
			LocusId:      "locus1",
			Filter:       "PASS",
			Zygosity:     "HET",
			Pf:           0.99,
			Af:           0.01,
			Hgvsg:        "hgvsg1",
			AdRatio:      1.0,
			VariantClass: "class1",
		},
	}, nil
}

func (m *MockRepository) CountOccurrences(int) (int, error) {
	return 15, nil
}

func TestStatusHandler(t *testing.T) {
	repo := &MockRepository{}
	router := gin.Default()
	router.GET("/status", statusHandler(repo))

	req, _ := http.NewRequest("GET", "/status", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"status":"up"}`, w.Body.String())
}

func TestOccurrencesListHandler(t *testing.T) {
	repo := &MockRepository{}
	router := gin.Default()
	router.GET("/occurrences/:seq_id/list", occurrencesListHandler(repo))

	req, _ := http.NewRequest("GET", "/occurrences/1/list?columns=seq_id,locus_id,filter,zygosity,pf,af,hgvsg,ad_ratio,variant_class", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `[{
        "seq_id": 1,
        "locus_id": "locus1",
        "filter": "PASS",
        "zygosity": "HET",
        "pf": 0.99,
        "af": 0.01,
        "hgvsg": "hgvsg1",
        "ad_ratio": 1.0,
        "variant_class": "class1"
    }]`, w.Body.String())
}

func TestOccurrencesCountHandler(t *testing.T) {
	repo := &MockRepository{}
	router := gin.Default()
	router.GET("/occurrences/:seq_id/count", occurrencesCountHandler(repo))

	req, _ := http.NewRequest("GET", "/occurrences/1/count", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"count":15}`, w.Body.String())
}