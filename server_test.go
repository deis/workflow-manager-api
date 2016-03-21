package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/deis/workflow-manager-api/data"
	"github.com/deis/workflow-manager-api/handlers"
	"github.com/deis/workflow-manager/types"
	"github.com/gorilla/mux"
)

func newServer(db data.DB, ver data.Version, counter data.Count, cluster data.Cluster) *httptest.Server {
	// Routes consist of a path and a handler function.
	r := mux.NewRouter()
	r.HandleFunc("/{apiVersion}/versions/{component}", handlers.VersionsGetHandler(db, ver)).Methods("GET")
	r.HandleFunc("/{apiVersion}/versions/{component}", handlers.VersionsPostHandler(db, ver)).Methods("POST")
	r.HandleFunc("/{apiVersion}/clusters", handlers.ClustersHandler(db, counter)).Methods("GET")
	r.HandleFunc("/{apiVersion}/clusters/{id}", handlers.ClustersGetHandler(db, cluster)).Methods("GET")
	r.HandleFunc("/{apiVersion}/clusters/{id}", handlers.ClustersPostHandler(db, cluster)).Methods("POST")
	return httptest.NewServer(r)
}

func TestGetClusters(t *testing.T) {
	memDB, err := data.NewMemDB()
	if err != nil {
		t.Fatalf("error creating new in-memory DB (%s)", err)
	}
	server := newServer(memDB, data.FakeVersion{}, data.FakeCount{}, data.FakeCluster{})
	defer server.Close()
	resp, err := httpGet(server, "/versions")
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Received non-200 response: %d\n", resp.StatusCode)
	}
}

func TestPostClusters(t *testing.T) {
	memDB, err := data.NewMemDB()
	if err != nil {
		t.Fatalf("error creating new in-memory DB (%s)", err)
	}
	id := "123"
	jsonData := `{"Components": [{"Component": {"Name": "component-a"}, "Version": {"Version": "1.0"}}]}`
	server := newServer(memDB, data.FakeVersion{}, data.FakeCount{}, data.FakeCluster{})
	defer server.Close()
	resp, err := httpPost(server.URL+"/clusters/"+id, jsonData)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Received non-200 response: %d\n", resp.StatusCode)
	}
	resp, err = httpGet(server, "/clusters")
	defer resp.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Received non-200 response: %d\n", resp.StatusCode)
	}
	json := parseJSONClusters(resp)
	if json[id].Components[0].Component.Name != "component-a" {
		t.Error("unexpected component name from JSON response")
	}
	//TODO Why do we have to dereference "Version" twice?
	if json[id].Components[0].Version.Version != "1.0" {
		t.Error("unexpected component version from JSON response")
	}
}

func parseJSONClusters(r *http.Response) map[string]types.Cluster {
	rawJSON, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Print(err)
	}
	var rawJSONMap map[string]*json.RawMessage
	err = json.Unmarshal(rawJSON, &rawJSONMap)
	if err != nil {
		log.Print(err)
	}
	clusters := make(map[string]types.Cluster)
	for id := range rawJSONMap {
		var clusterObj types.Cluster
		err = json.Unmarshal(*rawJSONMap[id], &clusterObj)
		if err != nil {
			log.Print(err)
		}
		clusters[id] = clusterObj
	}
	return clusters
}

func httpGet(s *httptest.Server, route string) (*http.Response, error) {
	return http.Get(s.URL + route)
}

func httpPost(url string, json string) (*http.Response, error) {
	jsonStr := []byte(json)
	return http.Post(url, "application/json", bytes.NewBuffer(jsonStr))
}
