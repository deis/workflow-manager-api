package handlers

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/arschles/assert"
	"github.com/deis/workflow-manager-api/data"
)

type clustersHandlerTestCase struct {
	Name  string
	DB    data.DB
	Count data.Count
}

func TestClustersHandler(t *testing.T) {
	db, err := data.NewMemDB()
	assert.NoErr(t, err)
	testCases := []clustersHandlerTestCase{
		{Name: "Valid count", DB: db, Count: data.FakeCount{Num: 123, Err: nil}},
		{Name: "Error counting", DB: db, Count: data.FakeCount{Num: 0, Err: errors.New("SOME ERROR")}},
	}
	for _, testCase := range testCases {
		handler := ClustersHandler(testCase.DB, testCase.Count)
		w := httptest.NewRecorder()
		r, err := http.NewRequest("GET", "/1/clusters", &bytes.Buffer{})
		assert.NoErr(t, err)
		handler(w, r)
	}
}
