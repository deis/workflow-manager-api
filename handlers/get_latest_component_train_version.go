package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/deis/workflow-manager-api/data"
	"github.com/gorilla/mux"
)

// GetLatestComponentTrainVersion returns the handler for the
// GET /:apiVersion/versions/:train/:component/latest endpoint
func GetLatestComponentTrainVersion(db *sql.DB, ver data.Version) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		train := vars[TrainPathKey]
		component := vars[ComponentPathKey]
		if train == "" {
			http.Error(w, "train is required", http.StatusBadRequest)
			return
		}
		if component == "" {
			http.Error(w, "component is required", http.StatusBadRequest)
			return
		}
		cv, err := data.GetLatestVersion(db, train, component)
		if err != nil {
			http.Error(w, fmt.Sprintf("error getting component (%s)", err), http.StatusInternalServerError)
			return
		}
		if json.NewEncoder(w).Encode(cv); err != nil {
			http.Error(w, fmt.Sprintf("rrror getting latest component version (%s)", err), http.StatusInternalServerError)
			return
		}
	})
}
