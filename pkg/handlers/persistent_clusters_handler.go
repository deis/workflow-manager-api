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

// PersistentClusters is the handler for the GET /{apiVersion}/clusters/persistent endpoint
func PersistentClusters(params operations.GetPersistentClustersParams, db *gorm.DB) middleware.Responder {
	persistentClustersFilter, err := parsePersistentClusterQueryKeys(params)
	if err != nil {
		return operations.NewGetPersistentClustersDefault(http.StatusBadRequest).WithPayload(&models.Error{Code: http.StatusBadRequest, Message: err.Error()})
	}

	checkins, err := data.FilterPersistentClusters(db, persistentClustersFilter)
	if err != nil {
		log.Printf("Error filtering persistent clusters (%s)", err)
		return operations.NewGetPersistentClustersDefault(http.StatusInternalServerError).WithPayload(&models.Error{Code: http.StatusInternalServerError, Message: err.Error()})
	}
	numResults := int64(len(checkins))
	clustersCount := models.ClustersCount{Count: &numResults, Data: checkins}
	return operations.NewGetClusterCheckinsOK().WithPayload(&clustersCount)
}
