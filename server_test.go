package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

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
	memDB, err := data.NewMemDB()
	if err != nil {
		t.Fatalf("error creating new in-memory DB (%s)", err)
	}
	server := newServer(memDB, data.VersionFromDB{}, data.ClusterCount{}, data.ClusterFromDB{})
	defer server.Close()
	resp, err := httpGet(server, urlPath(1, "clusters"))
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Received non-200 response: %d\n", resp.StatusCode)
	}
}

// tests the GET /{apiVersion}/clusters/{id} endpoint
func TestGetClusterByID(t *testing.T) {
	t.Skip("TODO")
}

// tests the POST {apiVersion}/clusters/{id} endpoint
func TestPostClusters(t *testing.T) {
	memDB, err := data.NewMemDB()
	if err != nil {
		t.Fatalf("error creating new in-memory DB (%s)", err)
	}
	id := "123"
	jsonData := `{"Components": [{"Component": {"Name": "component-a"}, "Version": {"Version": "1.0"}}]}`
	server := newServer(memDB, data.VersionFromDB{}, data.ClusterCount{}, data.ClusterFromDB{})
	defer server.Close()
	resp, err := httpPost(server, urlPath(1, "clusters", id), jsonData)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Received non-200 response: %d\n", resp.StatusCode)
	}
	resp, err = httpGet(server, urlPath(1, "clusters", id))
	defer resp.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Received non-200 response: %d", resp.StatusCode)
	}
	clusterMap := make(map[string]types.Cluster)
	// if err := json.NewDecoder(resp.Body).Decode(&clusterMap); err != nil {
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading response body (%s)", err)
	}
	t.Logf("got response body %s", string(respBody))
	if err := json.Unmarshal(respBody, &clusterMap); err != nil {
		t.Fatalf("parsing clusters map (%s)", err)
	}
	if clusterMap[id].Components[0].Component.Name != "component-a" {
		t.Error("unexpected component name from JSON response")
	}
	//TODO Why do we have to dereference "Version" twice?
	if clusterMap[id].Components[0].Version.Version != "1.0" {
		t.Error("unexpected component version from JSON response")
	}
}

func parseJSONClusters(r *http.Response) (map[string]types.Cluster, error) {
	rawJSONMap := make(map[string]*json.RawMessage)
	if err := json.NewDecoder(r.Body).Decode(&rawJSONMap); err != nil {
		return nil, err
	}
	log.Printf("received raw json map %+v", rawJSONMap)

	clusters := make(map[string]types.Cluster)
	for id := range rawJSONMap {
		var clusterObj types.Cluster
		if rawJSONMap[id] == nil {
			return nil, fmt.Errorf("id %s is nil", id)
		}
		if err := json.Unmarshal(*rawJSONMap[id], &clusterObj); err != nil {
			log.Print(err)
			continue
		}
		clusters[id] = clusterObj
	}
	return clusters, nil
}

func httpGet(s *httptest.Server, route string) (*http.Response, error) {
	return http.Get(s.URL + "/" + route)
}

func httpPost(s *httptest.Server, route string, json string) (*http.Response, error) {
	fullURL := s.URL + "/" + route
	return http.Post(fullURL, "application/json", bytes.NewBuffer([]byte(json)))
}
