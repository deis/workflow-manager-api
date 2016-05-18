package handlers

// handler echoes the HTTP request.
import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/deis/workflow-manager-api/pkg/data"
	"github.com/deis/workflow-manager-api/pkg/swagger/models"
	"github.com/deis/workflow-manager-api/pkg/swagger/restapi/operations"
	"github.com/go-swagger/go-swagger/httpkit/middleware"
	"github.com/jinzhu/gorm"
)

// ClustersCount route handler
func ClustersCount(db *gorm.DB) middleware.Responder {
	count, err := data.GetClusterCount(db)
	if err != nil {
		log.Printf("data.GetClusterCount error (%s)", err)
		return operations.NewGetClustersCountDefault(http.StatusInternalServerError).WithPayload(&models.Error{Code: http.StatusInternalServerError, Message: err.Error()})
	}
	return operations.NewGetClustersCountOK().WithPayload(int64(count))
}

// GetCluster route handler
func GetCluster(params operations.GetClusterByIDParams, db *gorm.DB) middleware.Responder {
	id := params.ID
	cluster, err := data.GetCluster(db, id)
	if err != nil {
		log.Printf("data.GetCluster error (%s)", err)
		return operations.NewGetClusterByIDDefault(http.StatusNotFound).WithPayload(&models.Error{Code: http.StatusNotFound, Message: "404 cluster not found"})
	}
	return operations.NewGetClusterByIDOK().WithPayload(&cluster)
}

// ClusterCheckin route handler
func ClusterCheckin(params operations.CreateClusterDetailsParams, db *gorm.DB) middleware.Responder {
	cluster := *params.Body
	id := cluster.ID
	var result models.Cluster
	result, err := data.UpsertCluster(db, id, cluster)
	if err != nil {
		log.Printf("data.SetCluster error (%s)", err)
		return operations.NewCreateClusterDetailsDefault(http.StatusInternalServerError).WithPayload(&models.Error{Code: http.StatusInternalServerError, Message: err.Error()})
	}
	return operations.NewCreateClusterDetailsOK().WithPayload(&result)
}

// GetVersion route handler
func GetVersion(params operations.GetComponentByReleaseParams, db *gorm.DB) middleware.Responder {
	train := params.Train
	component := params.Component
	version := params.Release
	componentVersion := models.ComponentVersion{
		Component: &models.Component{Name: component},
		Version:   &models.Version{Train: train, Version: version},
	}
	componentVersion, err := data.GetVersion(db, componentVersion)
	if err != nil {
		log.Printf("data.GetVersion error (%s)", err)
		return operations.NewGetComponentByReleaseDefault(http.StatusNotFound).WithPayload(&models.Error{Code: http.StatusNotFound, Message: "404 release not found"})
	}
	return operations.NewGetComponentByReleaseOK().WithPayload(&componentVersion)
}

// GetComponentTrainVersions route handler
func GetComponentTrainVersions(params operations.GetComponentByNameParams, db *gorm.DB) middleware.Responder {
	train := params.Train
	component := params.Component
	componentVersions, err := data.GetVersionsList(db, train, component)
	if err != nil {
		log.Printf("data.GetComponentTrainVersions error (%s)", err)
		return operations.NewGetComponentByNameDefault(http.StatusNotFound).WithPayload(&models.Error{Code: http.StatusNotFound, Message: "404 component not found"})
	}
	return operations.NewGetComponentByNameOK().WithPayload(operations.GetComponentByNameOKBodyBody{Data: componentVersions})
}

// PublishVersion route handler
func PublishVersion(params operations.PublishComponentReleaseParams, db *gorm.DB) middleware.Responder {
	componentVersion := *params.Body
	//TODO: validate request body parameter values for "component", "train", and "version"
	// match the values passed in with the URL
	componentVersion.Component.Name = params.Component
	componentVersion.Version.Train = params.Train
	componentVersion.Version.Version = params.Release
	result, err := data.UpsertVersion(db, componentVersion)
	if err != nil {
		log.Printf("data.SetVersion error (%s)", err)
		operations.NewPublishComponentReleaseDefault(http.StatusInternalServerError).WithPayload(&models.Error{Code: http.StatusInternalServerError, Message: err.Error()})
	}
	return operations.NewPublishComponentReleaseOK().WithPayload(&result)
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
