package handlers

// handler echoes the HTTP request.
import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/deis/workflow-manager-api/data"
	"github.com/deis/workflow-manager/types"
	"github.com/gorilla/mux"
)

// ClustersCount route handler
func ClustersCount(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count, err := data.GetClusterCount(db)
		if err != nil {
			log.Printf("data.GetClusterCount error (%s)", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writePlainText(strconv.Itoa(count), w)
	})
}

// GetCluster route handler
func GetCluster(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		cluster, err := data.GetCluster(db, id)
		if err != nil {
			log.Printf("data.GetCluster error (%s)", err)
			http.NotFound(w, r)
			return
		}
		js, err := json.Marshal(cluster)
		if err != nil {
			log.Printf("JSON marshaling failed (%s)", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writeJSON(js, w)
	})
}

// ClusterCheckin route handler
func ClusterCheckin(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		cluster := types.Cluster{}
		err := json.NewDecoder(r.Body).Decode(&cluster)
		if err != nil {
			log.Printf("Error decoding POST body JSON data (%s)", err)
			return
		}
		var result types.Cluster
		result, err = data.CheckInAndSetCluster(db, id, cluster)
		if err != nil {
			log.Printf("data.SetCluster error (%s)", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		data, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			log.Printf("JSON marshaling failed (%s)", err)
			http.Error(w, fmt.Sprintf("JSON marshaling failed (%s)", err), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(data))
	})
}

// GetVersion route handler
func GetVersion(db *sql.DB, v data.Version) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		routeParams := mux.Vars(r)
		train := routeParams["train"]
		component := routeParams["component"]
		version := routeParams["version"]
		params := types.ComponentVersion{
			Component: types.Component{Name: component},
			Version:   types.Version{Train: train, Version: version},
		}
		componentVersion, err := data.GetVersion(db, params)
		if err != nil {
			log.Printf("data.GetVersion error (%s)", err)
			http.NotFound(w, r)
			return
		}
		js, err := json.Marshal(componentVersion)
		if err != nil {
			log.Printf("JSON marshaling failed (%s)", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writeJSON(js, w)
	})
}

// GetComponentTrainVersions route handler
func GetComponentTrainVersions(db *sql.DB, v data.Version) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		routeParams := mux.Vars(r)
		train := routeParams["train"]
		component := routeParams["component"]
		componentVersions, err := data.GetVersionsList(db, train, component)
		if err != nil {
			log.Printf("data.GetComponentTrainVersions error (%s)", err)
			http.NotFound(w, r)
			return
		}
		js, err := json.Marshal(componentVersions)
		if err != nil {
			log.Printf("JSON marshaling failed (%s)", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writeJSON(js, w)
	})
}

// PublishVersion route handler
func PublishVersion(db *sql.DB, v data.Version) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		componentVersion := types.ComponentVersion{}
		err := json.NewDecoder(r.Body).Decode(&componentVersion)
		if err != nil {
			log.Printf("Error decoding POST body JSON data (%s)", err)
			return
		}
		//TODO: validate request body parameter values for "component", "train", and "version"
		// match the values passed in with the URL
		routeParams := mux.Vars(r)
		componentVersion.Component.Name = routeParams["component"]
		componentVersion.Version.Train = routeParams["train"]
		componentVersion.Version.Version = routeParams["version"]
		result, err := data.SetVersion(db, componentVersion)
		if err != nil {
			log.Printf("data.SetVersion error (%s)", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		data, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			log.Printf("JSON marshaling failed (%s)", err)
		}
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(data))
	})
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
