package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/deis/workflow-manager-api/pkg/data"
	"github.com/deis/workflow-manager/types"
	"github.com/jinzhu/gorm"
)

// SparseComponentInfo is the JSON compatible struct that holds limited data about a component
type SparseComponentInfo struct {
	Name string `json:"name"`
}

// SparseVersionInfo is the JSON compatible struct that holds limited data about a
// component version
type SparseVersionInfo struct {
	Train string `json:"train"`
}

// SparseComponentAndTrainInfo is the JSON compatible struct that holds a
// SparseComponentInfo and SparseVersionInfo
type SparseComponentAndTrainInfo struct {
	Component SparseComponentInfo `json:"component"`
	Version   SparseVersionInfo   `json:"version"`
}

// SparseComponentAndTrainInfoJSONWrapper is the JSON compatible struct that holds a slice of
// SparseComponentAndTrainInfo structs
type SparseComponentAndTrainInfoJSONWrapper struct {
	Data []SparseComponentAndTrainInfo `json:"data"`
}

// ComponentVersionsJSONWrapper is the JSON compatible struct that holds a slice of
// types.ComponentVersion structs
type ComponentVersionsJSONWrapper struct {
	Data []types.ComponentVersion `json:"data"`
}

// GetLatestVersions is the handler for the POST /{apiVersion}/versions/latest endpoint
func GetLatestVersions(db *gorm.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqStruct := new(SparseComponentAndTrainInfoJSONWrapper)
		if err := json.NewDecoder(r.Body).Decode(reqStruct); err != nil {
			http.Error(w, fmt.Sprintf("error decoding request body (%s)", err), http.StatusBadRequest)
			return
		}

		componentAndTrainSlice := make([]data.ComponentAndTrain, len(reqStruct.Data))
		for i, d := range reqStruct.Data {
			componentAndTrainSlice[i] = data.ComponentAndTrain{
				ComponentName: d.Component.Name,
				Train:         d.Version.Train,
			}
		}

		componentVersions, err := data.GetLatestVersions(db, componentAndTrainSlice)
		if err != nil {
			http.Error(w, "database error", http.StatusInternalServerError)
			return
		}

		ret := ComponentVersionsJSONWrapper{Data: componentVersions}
		if err := writeJSON(w, ret); err != nil {
			log.Printf("GetLatestVersions json marshal failed (%s)", err)
		}
	})
}
