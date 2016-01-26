package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func newServer() *httptest.Server {
	r := mux.NewRouter()
	// Routes consist of a path and a handler function.
	r.HandleFunc("/clusters", defaultHandler)
	r.HandleFunc("/versions", versionsHandler).Methods("GET")
	r.HandleFunc("/versions", versionsPostHandler).Methods("POST")
	r.HandleFunc("/clusters/{id}", clustersPostHandler).Methods("POST")
	return httptest.NewServer(r)
}

func TestGetClusters(t *testing.T) {
	server := newServer()
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
	id := "123"
	jsonData := `{"Components": [{"Name": "component-a", "Version": "1.0"}]}`
	server := newServer()
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
	if json[id].Components[0].Name != "component-a" {
		t.Error("unexpected component name from JSON response")
	}
	//TODO Why do we have to dereference "Version" twice?
	if json[id].Components[0].Version.Version != "1.0" {
		t.Error("unexpected component version from JSON response")
	}
}

func parseJSONClusters(r *http.Response) map[string]Cluster {
	rawJSON, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Print(err)
	}
	var rawJSONMap map[string]*json.RawMessage
	err = json.Unmarshal(rawJSON, &rawJSONMap)
	if err != nil {
		log.Print(err)
	}
	clusters := make(map[string]Cluster)
	for id := range rawJSONMap {
		var clusterObj Cluster
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
