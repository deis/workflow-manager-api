package data

import (
	"encoding/json"
	"testing"

	"github.com/arschles/assert"
	"github.com/deis/workflow-manager-api/pkg/swagger/models"
	sqlxTypes "github.com/jmoiron/sqlx/types"
)

func TestParseJSONComponentFail(t *testing.T) {
	jTxt := sqlxTypes.JSONText([]byte("hello world!"))
	_, err := parseJSONComponent(jTxt)
	if err == nil {
		t.Fatalf("returned error was nil")
	}
}

func TestParseJSONComponentSucc(t *testing.T) {
	cVer := models.ComponentVersion{
		Component:       &models.Component{Name: "test name", Description: &componentDescription},
		Version:         &models.Version{Train: "stable", Version: "test version", Released: "test release", Data: &versionData},
		UpdateAvailable: &updateAvailable,
	}
	b, err := json.Marshal(cVer)
	if err != nil {
		t.Fatalf("error marshalling ComponentVersion (%s)", err)
	}

	jTxt := sqlxTypes.JSONText(b)
	parsedCVer, err := parseJSONComponent(jTxt)
	if err != nil {
		t.Fatalf("returned error %s", err)
	}
	assert.Equal(t, parsedCVer, cVer, "component version")
}
