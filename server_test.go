package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/arschles/assert"
	"github.com/deis/workflow-manager-api/data"
	"github.com/deis/workflow-manager/types"
)

func newServer(db data.DB, ver data.Version, counter data.Count, cluster data.Cluster) *httptest.Server {
	// Routes consist of a path and a handler function.
	return httptest.NewServer(getRoutes(db, ver, counter, cluster))
}

func urlPath(ver int, remainder ...string) string {
	return fmt.Sprintf("%d/%s", ver, strings.Join(remainder, "/"))
}

// tests the GET /{apiVersion}/versions/{component} endpoint
func TestGetVersions(t *testing.T) {
	t.Skip("TODO")
}

// tests the POST /{apiVersion}/versions/{component} endpoint
func TestPostVersions(t *testing.T) {
	t.Skip("TODO")
}

// tests the GET /{apiVersion}/clusters endpoint
func TestGetClusters(t *testing.T) {
	memDB, err := newMemDB()
	if err != nil {
		t.Fatalf("error creating new in-memory DB (%s)", err)
	}
	if err := data.VerifyPersistentStorage(memDB); err != nil {
		t.Fatalf("VerifyPersistentStorage (%s)", err)
	}
	server := newServer(memDB, data.VersionFromDB{}, data.ClusterCount{}, data.ClusterFromDB{})
	defer server.Close()
	resp, err := httpGet(server, urlPath(1, "clusters"))
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
	if err := data.VerifyPersistentStorage(memDB); err != nil {
		t.Fatalf("VerifyPersistentStorage (%s)", err)
	}
	clusterFromDB := data.ClusterFromDB{}
	srv := newServer(memDB, data.VersionFromDB{}, data.ClusterCount{}, clusterFromDB)
	defer srv.Close()
	id := "123"
	cluster := types.Cluster{ID: id, FirstSeen: time.Now(), LastSeen: time.Now().Add(1 * time.Minute), Components: nil}
	newCluster, err := data.SetCluster(id, cluster, memDB, clusterFromDB)
	assert.NoErr(t, err)
	resp, err := httpGet(srv, urlPath(1, "clusters", id))
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
	if err := data.VerifyPersistentStorage(memDB); err != nil {
		t.Fatalf("VerifyPersistentStorage (%s)", err)
	}
	id := "123"
	jsonData := `{"Components": [{"Component": {"Name": "component-a"}, "Version": {"Version": "1.0"}}]}`
	server := newServer(memDB, data.VersionFromDB{}, data.ClusterCount{}, data.ClusterFromDB{})
	defer server.Close()
	resp, err := httpPost(server, urlPath(1, "clusters", id), jsonData)
	if err != nil {
		t.Fatalf("POSTing to endpoint (%s)", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Received non-200 response: %d", resp.StatusCode)
	}
	resp, err = httpGet(server, urlPath(1, "clusters", id))
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
