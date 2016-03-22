package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/arschles/assert"
)

const (
	jsonStr   = `{"this":"is json"}`
	plainText = "this is plain text"
)

func TestWriteJSON(t *testing.T) {
	json := []byte(jsonStr)
	w := httptest.NewRecorder()
	writeJSON(json, w)
	assert.Equal(t, w.Code, http.StatusOK, "response code")
	assert.Equal(t, w.HeaderMap.Get("Content-Type"), "application/json", "content type header")
	assert.True(t, w.Body != nil, "response body was nil")
	assert.Equal(t, string(w.Body.Bytes()), jsonStr, "response body")
}

func TestWritePlainText(t *testing.T) {
	w := httptest.NewRecorder()
	writePlainText(plainText, w)
	assert.Equal(t, w.Code, http.StatusOK, "response code")
	assert.Equal(t, w.HeaderMap.Get("Content-Type"), "text/plain", "content type header")
	assert.True(t, w.Body != nil, "response body was nil")
	assert.Equal(t, string(w.Body.Bytes()), plainText, "response body")
}
