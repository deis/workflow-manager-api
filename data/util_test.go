package data

import (
	"encoding/json"
	"testing"

	"github.com/deis/workflow-manager/types"
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
	cVer := types.ComponentVersion{
		Component:       types.Component{Name: "test name", Description: "test description"},
		Version:         types.Version{Version: "test version", Released: "test release"},
		UpdateAvailable: "test update avail",
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
	if parsedCVer != cVer {
		t.Fatalf("old component version != new (%+v != %+v)", cVer, parsedCVer)
	}
}
