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
	"github.com/jinzhu/gorm"
)

// ClustersCount route handler
func ClustersCount(db *gorm.DB) http.Handler {
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
func GetCluster(db *gorm.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		cluster, err := data.GetCluster(db, id)
		if err != nil {
			log.Printf("data.GetCluster error (%s)", err)
			http.NotFound(w, r)
			return
		}
		if err := writeJSON(w, cluster); err != nil {
			log.Printf("GetCluster json marshal failed (%s)", err)
		}
	})
}

// ClusterCheckin route handler
func ClusterCheckin(db *gorm.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		cluster := data.ClusterStateful{}
		err := json.NewDecoder(r.Body).Decode(&cluster)
		if err != nil {
			log.Printf("Error decoding POST body JSON data (%s)", err)
			return
		}
		var result data.ClusterStateful
		result, err = data.CheckInAndSetCluster(db, id, cluster)
		if err != nil {
			log.Printf("data.SetCluster error (%s)", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		if err := writeJSON(w, result); err != nil {
			log.Printf("ClusterCheckin json marshal error (%s)", err)
		}
	})
}

// GetVersion route handler
func GetVersion(db *gorm.DB) http.Handler {
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
		if err := writeJSON(w, componentVersion); err != nil {
			log.Printf("GetVersion json marshal failed (%s)", err)
		}
	})
}

// GetComponentTrainVersions route handler
func GetComponentTrainVersions(db *sql.DB) http.Handler {
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
		if err := writeJSON(w, componentVersions); err != nil {
			log.Printf("GetComponentTrainVersions json marshal failed (%s)", err)
		}
	})
}

// PublishVersion route handler
func PublishVersion(db *gorm.DB) http.Handler {
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
		if err := writeJSON(w, result); err != nil {
			log.Printf("PublishVersion json marshal error (%s)", err)
		}
	})
}

// writeJSON is a helper function for writing HTTP JSON data
func writeJSON(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"error":"%s","error_type":"json"}`, err)))
		return err
	}
	return nil
}

// writePlainText is a helper function for writing HTTP text data
func writePlainText(text string, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(text))
}
