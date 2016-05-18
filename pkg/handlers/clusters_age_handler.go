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

// ClustersAge is the handler for the GET /{apiVersion}/clusters/age endpoint
func ClustersAge(params operations.GetClustersByAgeParams, db *gorm.DB) middleware.Responder {
	clusterAgeFilter, err := parseAgeQueryKeys(params)
	if err != nil {
		return operations.NewGetClustersByAgeDefault(http.StatusBadRequest).WithPayload(&models.Error{Code: http.StatusBadRequest, Message: err.Error()})
	}

	clusters, err := data.FilterClustersByAge(db, clusterAgeFilter)
	if err != nil {
		log.Printf("Error filtering clusters by age (%s)", err)
		return operations.NewGetClustersByAgeDefault(http.StatusInternalServerError).WithPayload(&models.Error{Code: http.StatusInternalServerError, Message: err.Error()})
	}
	return operations.NewGetClustersByAgeOK().WithPayload(operations.GetClustersByAgeOKBodyBody{Data: clusters})
}
