package handlers

import (
	"log"
	"net/http"

	"github.com/deis/workflow-manager-api/pkg/data"
	"github.com/deis/workflow-manager-api/pkg/swagger/models"
	"github.com/deis/workflow-manager-api/pkg/swagger/restapi/operations"
	"github.com/deis/workflow-manager/types"
	"github.com/go-swagger/go-swagger/httpkit/middleware"
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
func GetLatestVersions(params operations.GetComponentsByLatestReleaseParams, db *gorm.DB) middleware.Responder {
	reqStruct := params.Body

	componentAndTrainSlice := make([]data.ComponentAndTrain, len(reqStruct.Data))
	for i, d := range reqStruct.Data {
		componentAndTrainSlice[i] = data.ComponentAndTrain{
			ComponentName: d.Component.Name,
			Train:         d.Version.Train,
		}
	}

	componentVersions, err := data.GetLatestVersions(db, componentAndTrainSlice)
	if err != nil {
		log.Printf("data.GetLatestVersions error (%s)", err)
		return operations.NewGetComponentsByLatestReleaseDefault(http.StatusInternalServerError).WithPayload(&models.Error{Code: http.StatusInternalServerError, Message: "database error"})
	}
	ret := operations.GetComponentsByLatestReleaseOKBodyBody{Data: componentVersions}
	return operations.NewGetComponentsByLatestReleaseOK().WithPayload(ret)
}

// GetLatestVersionsForV2 is the handler for the POST /v2/versions/latest endpoint
func GetLatestVersionsForV2(params operations.GetComponentsByLatestReleaseForV2Params, db *gorm.DB) middleware.Responder {
	reqStruct := params.Body
	componentAndTrainSlice := make([]data.ComponentAndTrain, len(reqStruct.Data))
	for i, d := range reqStruct.Data {
		//this is just to make the existing workflow manager clients out in the beta world to work.
		if d.Version.Train == "" {
			d.Version.Train = "beta"
		}
		componentAndTrainSlice[i] = data.ComponentAndTrain{
			ComponentName: d.Component.Name,
			Train:         d.Version.Train,
		}
	}

	componentVersions, err := data.GetLatestVersions(db, componentAndTrainSlice)
	if err != nil {
		log.Printf("data.GetLatestVersions error (%s)", err)
		return operations.NewGetComponentsByLatestReleaseForV2Default(http.StatusInternalServerError).WithPayload(&models.Error{Code: http.StatusInternalServerError, Message: "database error"})
	}
	ret := operations.GetComponentsByLatestReleaseForV2OKBodyBody{Data: componentVersions}
	return operations.NewGetComponentsByLatestReleaseForV2OK().WithPayload(ret)
}
