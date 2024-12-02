package server

import (
	"bytes"
	"go-poc/internal/types"
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

func (m *MockRepository) GetOccurrences(int, *types.Query) ([]types.Occurrence, error) {
	return []types.Occurrence{
		{
			SeqId:        1,
			LocusId:      1000,
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

func (m *MockRepository) CountOccurrences(int, *types.Query) (int64, error) {
	return 15, nil
}

func (m *MockRepository) AggregateOccurrences(int, *types.Query) ([]types.Aggregation, error) {
	return []types.Aggregation{
			{Bucket: "HET", Count: 2},
			{Bucket: "HOM", Count: 1},
		},
		nil
}

func TestStatusHandler(t *testing.T) {
	repo := &MockRepository{}
	router := gin.Default()
	router.GET("/status", StatusHandler(repo))

	req, _ := http.NewRequest("GET", "/status", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"status":"up"}`, w.Body.String())
}

func TestOccurrencesListHandler(t *testing.T) {
	repo := &MockRepository{}
	router := gin.Default()
	router.POST("/occurrences/:seq_id/list", OccurrencesListHandler(repo))
	body := `{
			"selected_fields":[
				"seq_id","locus_id","filter","zygosity","pf","af","hgvsg","ad_ratio","variant_class"
			]
	}`
	req, _ := http.NewRequest("POST", "/occurrences/1/list", bytes.NewBuffer([]byte(body)))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `[{
        "seq_id": 1,
        "locus_id": 1000,
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
	router.POST("/occurrences/:seq_id/count", OccurrencesCountHandler(repo))

	req, _ := http.NewRequest("POST", "/occurrences/1/count", bytes.NewBuffer([]byte("{}")))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"count":15}`, w.Body.String())
}

func TestOccurrencesAggregateHandler(t *testing.T) {
	repo := &MockRepository{}
	router := gin.Default()
	router.POST("/occurrences/:seq_id/aggregate", OccurrencesAggregateHandler(repo))

	body := `{
			"field": "zygosity",
			"sqon":{
				"op":"in",
				"field": "filter",
				"value": "PASS"
		    },
			"size": 10
	}`
	req, _ := http.NewRequest("POST", "/occurrences/1/aggregate", bytes.NewBuffer([]byte(body)))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	expected := `[{"key": "HET", "count": 2}, {"key": "HOM", "count": 1}]`
	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, expected, w.Body.String())
}
