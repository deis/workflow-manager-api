package handlers

import (
	"log"
	"net/http"

	"github.com/deis/workflow-manager-api/pkg/data"
	"github.com/deis/workflow-manager-api/pkg/swagger/models"
	"github.com/deis/workflow-manager-api/pkg/swagger/restapi/operations"
	"github.com/go-swagger/go-swagger/httpkit/middleware"
	"github.com/jinzhu/gorm"
)

// ClusterCheckins is the handler for the GET /{apiVersion}/clusters/checkins endpoint
func ClusterCheckins(params operations.GetClusterCheckinsParams, db *gorm.DB) middleware.Responder {
	clusterCheckinsFilter, err := parseCheckinsQueryKeys(params)
	if err != nil {
		return operations.NewGetClusterCheckinsDefault(http.StatusBadRequest).WithPayload(&models.Error{Code: http.StatusBadRequest, Message: err.Error()})
	}

	checkins, err := data.FilterClusterCheckins(db, clusterCheckinsFilter)
	if err != nil {
		log.Printf("Error filtering cluster checkins (%s)", err)
		return operations.NewGetClusterCheckinsDefault(http.StatusInternalServerError).WithPayload(&models.Error{Code: http.StatusInternalServerError, Message: err.Error()})
	}
	numResults := int64(len(checkins))
	clustersCount := models.ClustersCount{Count: &numResults, Data: checkins}
	return operations.NewGetClusterCheckinsOK().WithPayload(&clustersCount)
}
