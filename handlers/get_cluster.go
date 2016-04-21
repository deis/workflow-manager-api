package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/deis/workflow-manager-api/data"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

// GetCluster route handler
func GetCluster(db *gorm.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		cluster, err := data.GetCluster(db, id)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		if err := json.NewEncoder(w).Encode(cluster); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}
