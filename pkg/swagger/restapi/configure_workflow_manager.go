package restapi

import (
	"net/http"

	errors "github.com/go-swagger/go-swagger/errors"
	httpkit "github.com/go-swagger/go-swagger/httpkit"
	middleware "github.com/go-swagger/go-swagger/httpkit/middleware"

	"github.com/deis/workflow-manager-api/pkg/swagger/restapi/operations"
)

// This file is safe to edit. Once it exists it will not be overwritten

func configureFlags(api *operations.WorkflowManagerAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.WorkflowManagerAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	api.JSONConsumer = httpkit.JSONConsumer()

	api.JSONProducer = httpkit.JSONProducer()

	api.CreateClusterDetailsHandler = operations.CreateClusterDetailsHandlerFunc(func(params operations.CreateClusterDetailsParams) middleware.Responder {
		return middleware.NotImplemented("operation .CreateClusterDetails has not yet been implemented")
	})
	api.GetClusterByIDHandler = operations.GetClusterByIDHandlerFunc(func(params operations.GetClusterByIDParams) middleware.Responder {
		return middleware.NotImplemented("operation .GetClusterByID has not yet been implemented")
	})
	api.GetClustersByAgeHandler = operations.GetClustersByAgeHandlerFunc(func(params operations.GetClustersByAgeParams) middleware.Responder {
		return middleware.NotImplemented("operation .GetClustersByAge has not yet been implemented")
	})
	api.GetClustersCountHandler = operations.GetClustersCountHandlerFunc(func() middleware.Responder {
		return middleware.NotImplemented("operation .GetClustersCount has not yet been implemented")
	})
	api.GetComponentByLatestReleaseHandler = operations.GetComponentByLatestReleaseHandlerFunc(func(params operations.GetComponentByLatestReleaseParams) middleware.Responder {
		return middleware.NotImplemented("operation .GetComponentByLatestRelease has not yet been implemented")
	})
	api.GetComponentByNameHandler = operations.GetComponentByNameHandlerFunc(func(params operations.GetComponentByNameParams) middleware.Responder {
		return middleware.NotImplemented("operation .GetComponentByName has not yet been implemented")
	})
	api.GetComponentByReleaseHandler = operations.GetComponentByReleaseHandlerFunc(func(params operations.GetComponentByReleaseParams) middleware.Responder {
		return middleware.NotImplemented("operation .GetComponentByRelease has not yet been implemented")
	})
	api.GetComponentsByLatestReleaseHandler = operations.GetComponentsByLatestReleaseHandlerFunc(func(params operations.GetComponentsByLatestReleaseParams) middleware.Responder {
		return middleware.NotImplemented("operation .GetComponentsByLatestRelease has not yet been implemented")
	})
	api.PublishComponentReleaseHandler = operations.PublishComponentReleaseHandlerFunc(func(params operations.PublishComponentReleaseParams) middleware.Responder {
		return middleware.NotImplemented("operation .PublishComponentRelease has not yet been implemented")
	})

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
