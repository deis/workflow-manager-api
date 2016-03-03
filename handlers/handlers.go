package handlers

// handler echoes the HTTP request.
import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/deis/workflow-manager-api/data"
	"github.com/deis/workflow-manager/types"
	"github.com/gorilla/mux"
)

// ClustersHandler route handler
func ClustersHandler(d data.DB, c data.Count) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		count, err := data.GetClusterCount(d, c)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writePlainText(strconv.Itoa(count), w)
	}
}

// ClustersGetHandler route handler
func ClustersGetHandler(d data.DB, c data.Cluster) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		cluster, err := data.GetCluster(id, d, c)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		js, err := json.Marshal(cluster)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writeJSON(js, w)
	}
}

// ClustersPostHandler route handler
func ClustersPostHandler(d data.DB, c data.Cluster) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			http.Error(w, "expected application/json", http.StatusUnsupportedMediaType)
			return
		}
		id := mux.Vars(r)["id"]
		cluster := types.Cluster{}
		err := json.NewDecoder(r.Body).Decode(&cluster)
		if err != nil {
			log.Print(err)
			return
		}
		var result types.Cluster
		result, err = data.SetCluster(id, cluster, d, c)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		data, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			log.Fatalf("JSON marshaling failed: %s", err)
		}
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(data))
	}
}

// VersionsGetHandler handles GET requests to "/versions/{component}"
func VersionsGetHandler(w http.ResponseWriter, r *http.Request) {
	component := mux.Vars(r)["component"]
	componentVersion, ok := data.GetVersion(component)
	if !ok {
		http.NotFound(w, r)
		return
	}
	js, err := json.Marshal(componentVersion)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(js, w)
}

// VersionsPostHandler handles POST requests to /versions/{component}
func VersionsPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "expected application/json", http.StatusUnsupportedMediaType)
		return
	}
	component := mux.Vars(r)["component"]
	componentVersion := types.ComponentVersion{}
	err := json.NewDecoder(r.Body).Decode(&componentVersion)
	if err != nil {
		log.Print(err)
		return
	}
	data, err := json.MarshalIndent(data.SetVersion(component, componentVersion), "", "  ")
	if err != nil {
		log.Fatalf("JSON marshaling failed: %s", err)
	}
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(data))
}

// writeJSON is a helper function for writing HTTP JSON data
func writeJSON(json []byte, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

// writePlainText is a helper function for writing HTTP text data
func writePlainText(text string, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(text))
}
