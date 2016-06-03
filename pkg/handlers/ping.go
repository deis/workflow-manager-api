package handlers

import (
	"github.com/deis/workflow-manager-api/pkg/swagger/restapi/operations"
	"github.com/go-swagger/go-swagger/httpkit/middleware"
)

// Ping is the handler for the ping endpoint
func Ping() middleware.Responder {
	return operations.NewPingOK()
}
