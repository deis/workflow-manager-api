package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func getRoutes() *mux.Router {
	r := mux.NewRouter()
	// Routes consist of a path and a handler function.
	r.HandleFunc("/clusters", defaultHandler)
	r.HandleFunc("/versions", versionsHandler).Methods("GET")
	r.HandleFunc("/versions", versionsPostHandler).Methods("POST")
	r.HandleFunc("/clusters/{id}", clustersPostHandler).Methods("POST")
	return r
}

// handler echoes the HTTP request.
func defaultHandler(w http.ResponseWriter, r *http.Request) {
	js, err := json.Marshal(getAll())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(js, w)
}

// versionsHandler handles GET requests to "/versions"
func versionsHandler(w http.ResponseWriter, r *http.Request) {
	js, err := json.Marshal(getVersions())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(js, w)
}

func writeJSON(json []byte, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

// clustersPostHandler handles POST requests to /clusters/{id}
func clustersPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "expected application/json", http.StatusUnsupportedMediaType)
		return
	}
	id := mux.Vars(r)["id"]
	mu.Lock()
	cluster, _ := getCluster(id)
	err := json.NewDecoder(r.Body).Decode(&cluster)
	if err != nil {
		log.Print(err)
		return
	}
	setCluster(id, cluster)
	mu.Unlock()
	data, err := json.MarshalIndent(cluster, "", "  ")
	if err != nil {
		log.Fatalf("JSON marshaling failed: %s", err)
	}
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(data))
}

// versionsPostHandler handles POST requests to /versions
func versionsPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "expected application/json", http.StatusUnsupportedMediaType)
		return
	}
	mu.Lock()
	componentVersions := []ComponentVersion{}
	err := json.NewDecoder(r.Body).Decode(&componentVersions)
	if err != nil {
		log.Print(err)
		return
	}
	for _, version := range componentVersions {
		setLatestVersion(version)
	}
	mu.Unlock()
	data, err := json.MarshalIndent(componentVersions, "", "  ")
	if err != nil {
		log.Fatalf("JSON marshaling failed: %s", err)
	}
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(data))
}
