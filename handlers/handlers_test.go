package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/arschles/assert"
)

const (
	plainText = "this is plain text"
)

func TestWriteJSON(t *testing.T) {
	json := map[string]string{"this": "is json"}
	const expected = `{"this":"is json"}`
	w := httptest.NewRecorder()
	assert.NoErr(t, writeJSON(w, json))
	assert.Equal(t, w.Code, http.StatusOK, "response code")
	assert.Equal(t, w.HeaderMap.Get("Content-Type"), "application/json", "content type header")
	assert.True(t, w.Body != nil, "response body was nil")
	assert.Equal(t, strings.TrimSpace(string(w.Body.Bytes())), expected, "response body")
}

func TestWriteJSONFail(t *testing.T) {
	toEncode := map[string]interface{}{
		"a": "b",
		// funcs cannot be encoded as json. see https://godoc.org/encoding/json#Marshal
		"c": func(i int) int { return i },
	}
	w := httptest.NewRecorder()
	err := writeJSON(w, toEncode)
	switch e := err.(type) {
	case *json.UnsupportedTypeError:
	default:
		t.Fatalf("JSON encoding returned err (%s), expected an UnsupportedTypeError", e)
	}

	assert.Equal(t, w.Code, http.StatusInternalServerError, "response code")
	assert.Equal(t, w.HeaderMap.Get("Content-Type"), "application/json", "content type header")
	assert.True(t, w.Body != nil, "response body was nil")
	var errBody struct {
		Error     string `json:"error"`
		ErrorType string `json:"error_type"`
	}
	assert.NoErr(t, json.NewDecoder(w.Body).Decode(&errBody))
	assert.Equal(t, errBody.Error, err.Error(), "error string")
	assert.Equal(t, errBody.ErrorType, "json", "error type")
}

func TestWritePlainText(t *testing.T) {
	w := httptest.NewRecorder()
	writePlainText(plainText, w)
	assert.Equal(t, w.Code, http.StatusOK, "response code")
	assert.Equal(t, w.HeaderMap.Get("Content-Type"), "text/plain", "content type header")
	assert.True(t, w.Body != nil, "response body was nil")
	assert.Equal(t, string(w.Body.Bytes()), plainText, "response body")
}
