package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/arschles/assert"
	"github.com/deis/workflow-manager-api/data"
	"github.com/deis/workflow-manager/types"
)

const (
	componentName = "testcomponent"
	clusterID     = "testcluster"
)

func newServer(d data.DB, ver data.Version, counter data.Count, cluster data.Cluster) *httptest.Server {
	db, _ := d.Get()
	// Routes consist of a path and a handler function.
	return httptest.NewServer(getRoutes(db, ver, counter, cluster))
}

func urlPath(ver int, remainder ...string) string {
	return fmt.Sprintf("%d/%s", ver, strings.Join(remainder, "/"))
}

// tests the GET /{apiVersion}/versions/{train}/{component}/{version} endpoint
func TestGetVersion(t *testing.T) {
	memDB, err := newMemDB()
	assert.NoErr(t, err)
	db, err := data.VerifyPersistentStorage(memDB)
	assert.NoErr(t, err)
	versionFromDB := data.VersionFromDB{}
	srv := newServer(memDB, versionFromDB, data.ClusterCount{}, data.ClusterFromDB{})
	defer srv.Close()
	componentVer := types.ComponentVersion{
		Component: types.Component{Name: componentName},
		Version:   types.Version{Train: "beta", Version: "2.0.0-beta-2", Released: "2016-03-31T23:54:39Z", Data: []byte(`{"description": "release details"}`)},
	}
	_, err = data.SetVersion(componentVer, db, versionFromDB)
	assert.NoErr(t, err)
	resp, err := httpGet(srv, urlPath(1, "versions", componentVer.Version.Train, componentVer.Component.Name, componentVer.Version.Version))
	assert.NoErr(t, err)
	defer resp.Body.Close()
	assert.Equal(t, resp.StatusCode, http.StatusOK, "response code")
	decodedVer := new(types.ComponentVersion)
	assert.NoErr(t, json.NewDecoder(resp.Body).Decode(decodedVer))
	assert.Equal(t, *decodedVer, componentVer, "component version")
}

// tests the GET /{apiVersion}/versions/{train}/{component} endpoint
func TestGetComponentTrainVersions(t *testing.T) {
	memDB, err := newMemDB()
	assert.NoErr(t, err)
	db, err := data.VerifyPersistentStorage(memDB)
	assert.NoErr(t, err)
	versionFromDB := data.VersionFromDB{}
	srv := newServer(memDB, versionFromDB, data.ClusterCount{}, data.ClusterFromDB{})
	defer srv.Close()
	componentVers := []types.ComponentVersion{}
	componentVer1 := types.ComponentVersion{
		Component: types.Component{Name: componentName},
		Version:   types.Version{Train: "beta", Version: "2.0.0-beta-1", Released: "2016-03-30T23:54:39Z", Data: []byte(`{"description": "release details"}`)},
	}
	componentVer2 := types.ComponentVersion{
		Component: types.Component{Name: componentName},
		Version:   types.Version{Train: "beta", Version: "2.0.0-beta-2", Released: "2016-03-31T23:54:39Z", Data: []byte(`{"description": "release details"}`)},
	}
	componentVers = append(componentVers, componentVer1)
	componentVers = append(componentVers, componentVer2)
	_, err = data.SetVersion(componentVers[0], db, versionFromDB)
	assert.NoErr(t, err)
	_, err = data.SetVersion(componentVers[1], db, versionFromDB)
	assert.NoErr(t, err)
	resp, err := httpGet(srv, urlPath(1, "versions", componentVer1.Version.Train, componentVer1.Component.Name))
	assert.NoErr(t, err)
	defer resp.Body.Close()
	assert.Equal(t, resp.StatusCode, http.StatusOK, "response code")
	decodedVer := new([]types.ComponentVersion)
	assert.NoErr(t, json.NewDecoder(resp.Body).Decode(decodedVer))
	assert.Equal(t, *decodedVer, componentVers, "component versions")
}

// tests the POST /{apiVersion}/versions/{train}/{component}/{version} endpoint
func TestPostVersions(t *testing.T) {
	memDB, err := newMemDB()
	assert.NoErr(t, err)
	db, err := data.VerifyPersistentStorage(memDB)
	assert.NoErr(t, err)
	versionFromDB := data.VersionFromDB{}
	srv := newServer(memDB, versionFromDB, data.ClusterCount{}, data.ClusterFromDB{})
	defer srv.Close()
	train := "beta"
	version := "2.0.0-beta-2"
	componentVer := types.ComponentVersion{
		Component: types.Component{Name: componentName},
		Version:   types.Version{Train: train, Version: version, Released: "2016-03-31T23:54:39Z", Data: []byte(`{"description": "release details"}`)},
	}
	body := new(bytes.Buffer)
	assert.NoErr(t, json.NewEncoder(body).Encode(componentVer))
	resp, err := httpPost(srv, urlPath(2, "versions", train, componentName, version), string(body.Bytes()))
	assert.NoErr(t, err)
	defer resp.Body.Close()
	assert.Equal(t, resp.StatusCode, http.StatusOK, "response code")
	retComponentVersion := new(types.ComponentVersion)
	assert.NoErr(t, json.NewDecoder(resp.Body).Decode(retComponentVersion))
	// TODO: version data property not traveling and returning as expected
	assert.Equal(t, *retComponentVersion, componentVer, "component version")
	fetchedComponentVersion, err := data.GetVersion(componentVer, db, versionFromDB)
	assert.NoErr(t, err)
	assert.Equal(t, fetchedComponentVersion, componentVer, "component version")
}

// tests the GET /{apiVersion}/clusters/count endpoint
func TestGetClusters(t *testing.T) {
	memDB, err := newMemDB()
	if err != nil {
		t.Fatalf("error creating new in-memory DB (%s)", err)
	}
	db, err := data.VerifyPersistentStorage(memDB)
	assert.NotNil(t, db, "db")
	if err != nil {
		log.Fatalf("VerifyPersistentStorage (%s)", err)
	}
	server := newServer(memDB, data.VersionFromDB{}, data.ClusterCount{}, data.ClusterFromDB{})
	defer server.Close()
	resp, err := httpGet(server, urlPath(1, "clusters", "count"))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatalf("Received non-200 response: %d\n", resp.StatusCode)
	}
}

// tests the GET /{apiVersion}/clusters/{id} endpoint
func TestGetClusterByID(t *testing.T) {
	memDB, err := newMemDB()
	assert.NoErr(t, err)
	db, err := data.VerifyPersistentStorage(memDB)
	if err != nil {
		log.Fatalf("VerifyPersistentStorage (%s)", err)
	}
	clusterFromDB := data.ClusterFromDB{}
	srv := newServer(memDB, data.VersionFromDB{}, data.ClusterCount{}, clusterFromDB)
	defer srv.Close()
	cluster := types.Cluster{ID: clusterID, FirstSeen: time.Now(), LastSeen: time.Now().Add(1 * time.Minute), Components: nil}
	newCluster, err := data.SetCluster(clusterID, cluster, db, clusterFromDB)
	assert.NoErr(t, err)
	resp, err := httpGet(srv, urlPath(1, "clusters", clusterID))
	assert.NoErr(t, err)
	defer resp.Body.Close()
	assert.Equal(t, resp.StatusCode, 200, "response code")
	decodedCluster := new(types.Cluster)
	assert.NoErr(t, json.NewDecoder(resp.Body).Decode(decodedCluster))
	assert.Equal(t, *decodedCluster, newCluster, "returned cluster")
}

// tests the POST {apiVersion}/clusters/{id} endpoint
func TestPostClusters(t *testing.T) {
	memDB, err := newMemDB()
	if err != nil {
		t.Fatalf("error creating new in-memory DB (%s)", err)
	}
	db, err := data.VerifyPersistentStorage(memDB)
	assert.NotNil(t, db, "db")
	if err != nil {
		log.Fatalf("VerifyPersistentStorage (%s)", err)
	}
	jsonData := `{"Components": [{"Component": {"Name": "component-a"}, "Version": {"Version": "1.0"}}]}`
	server := newServer(memDB, data.VersionFromDB{}, data.ClusterCount{}, data.ClusterFromDB{})
	defer server.Close()
	resp, err := httpPost(server, urlPath(1, "clusters", clusterID), jsonData)
	if err != nil {
		t.Fatalf("POSTing to endpoint (%s)", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Received non-200 response: %d", resp.StatusCode)
	}
	resp, err = httpGet(server, urlPath(1, "clusters", clusterID))
	defer resp.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Received non-200 response: %d", resp.StatusCode)
	}
	cluster := new(types.Cluster)
	if err := json.NewDecoder(resp.Body).Decode(cluster); err != nil {
		t.Fatalf("error reading response body (%s)", err)
	}
	if len(cluster.Components) <= 0 {
		t.Fatalf("no components returned")
	}
	if cluster.Components[0].Component.Name != "component-a" {
		t.Error("unexpected component name from JSON response")
	}
	// Note that we have to dereference "Version" twice because cluster.Components[0].Version
	// is itself a types.Version, which has both a "Released" and "Version" field
	if cluster.Components[0].Version.Version != "1.0" {
		t.Error("unexpected component version from JSON response")
	}
}

func httpGet(s *httptest.Server, route string) (*http.Response, error) {
	return http.Get(s.URL + "/" + route)
}

func httpPost(s *httptest.Server, route string, json string) (*http.Response, error) {
	fullURL := s.URL + "/" + route
	return http.Post(fullURL, "application/json", bytes.NewBuffer([]byte(json)))
}
