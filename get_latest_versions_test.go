package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/arschles/assert"
	"github.com/deis/workflow-manager-api/data"
	"github.com/deis/workflow-manager-api/handlers"
	"github.com/deis/workflow-manager/types"
)

func TestGetLatestVersions(t *testing.T) {
	memDB, err := data.NewMemDB()
	assert.NoErr(t, err)
	assert.NoErr(t, data.VerifyPersistentStorage(memDB))
	assert.NoErr(t, err)
	srv := newServer(memDB)
	defer srv.Close()
	const numComponentVersions = 6
	components := make(map[string]types.ComponentVersion)
	for i := 0; i < numComponentVersions; i++ {
		name := fmt.Sprintf("component%d", i)
		train := fmt.Sprintf("train%d", i)
		releaseTime1 := time.Now().Add(time.Duration(i+1) * time.Hour)
		releaseTime2 := time.Now().Add(time.Duration((i+1)*2) * time.Hour)
		cv1 := types.ComponentVersion{
			Component: types.Component{Name: name},
			Version: types.Version{
				Train:    train,
				Released: releaseTime1.Format(releaseTimeFormat),
				Version:  fmt.Sprintf("version%d-1", i),
			},
		}
		cv2 := types.ComponentVersion{
			Component: types.Component{Name: name},
			Version: types.Version{
				Train:    train,
				Released: releaseTime2.Format(releaseTimeFormat),
				Version:  fmt.Sprintf("version%d-2", i),
			},
		}

		if _, err := data.SetVersion(memDB, cv1); err != nil {
			t.Fatalf("Error setting component %d (%s)", i, err)
		}
		if _, err := data.SetVersion(memDB, cv2); err != nil {
			t.Fatalf("Error setting component %d (%s)", i, err)
		}
		components[cv2.Component.Name] = cv2
	}

	postBody := handlers.SparseComponentAndTrainInfoJSONWrapper{
		Data: make([]handlers.SparseComponentAndTrainInfo, len(components)),
	}
	i := 0
	for _, component := range components {
		postBody.Data[i] = handlers.SparseComponentAndTrainInfo{
			Component: handlers.SparseComponentInfo{Name: component.Component.Name},
			Version:   handlers.SparseVersionInfo{Train: component.Version.Train},
		}
		i++
	}
	postBodyReader := new(bytes.Buffer)
	assert.NoErr(t, json.NewEncoder(postBodyReader).Encode(postBody))

	resp, err := httpPost(srv, urlPath(2, "versions", "latest"), string(postBodyReader.Bytes()))
	assert.NoErr(t, err)
	defer resp.Body.Close()
	assert.Equal(t, resp.StatusCode, http.StatusOK, "response code")
	respStruct := new(handlers.ComponentVersionsJSONWrapper)
	assert.NoErr(t, json.NewDecoder(resp.Body).Decode(respStruct))
	assert.Equal(t, len(respStruct.Data), len(components), "length of response")
	for _, component := range respStruct.Data {
		expected, ok := components[component.Component.Name]
		assert.True(t, ok, "%s not found in the component/train response body", component.Component.Name)
		assert.Equal(t, component.Version.Version, expected.Version.Version, "component version")
		assert.Equal(t, component.Version.Train, expected.Version.Train, "component train")
		assert.Equal(t, component.Version.Released, expected.Version.Released, "component released time")
	}
}
