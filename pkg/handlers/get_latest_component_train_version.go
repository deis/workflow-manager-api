package handlers

import (
	"fmt"
	"net/http"

	"github.com/deis/workflow-manager-api/pkg/data"
	"github.com/deis/workflow-manager-api/pkg/swagger/models"
	"github.com/deis/workflow-manager-api/pkg/swagger/restapi/operations"
	"github.com/go-swagger/go-swagger/httpkit/middleware"
	"github.com/jinzhu/gorm"
)

// GetLatestComponentTrainVersion returns the response for the
// GET /:apiVersion/versions/:train/:component/latest endpoint
func GetLatestComponentTrainVersion(params operations.GetComponentByLatestReleaseParams, db *gorm.DB) middleware.Responder {
	train := params.Train
	component := params.Component
	if train == "" {
		return operations.NewGetComponentByLatestReleaseDefault(http.StatusBadRequest).WithPayload(&models.Error{Code: http.StatusBadRequest, Message: "train is required"})
	}
	if component == "" {
		return operations.NewGetComponentByLatestReleaseDefault(http.StatusBadRequest).WithPayload(&models.Error{Code: http.StatusBadRequest, Message: "component is required"})
	}
	cv, err := data.GetLatestVersion(db, train, component)
	if err != nil {
		return operations.NewGetComponentByLatestReleaseDefault(http.StatusInternalServerError).WithPayload(&models.Error{Code: http.StatusInternalServerError, Message: fmt.Sprintf("error getting component (%s)", err)})
	}
	return operations.NewGetComponentByLatestReleaseOK().WithPayload(&cv)
}
