package restapi

import (
	"log"
	"net/http"

	"github.com/deis/workflow-manager-api/config"
	"github.com/deis/workflow-manager-api/pkg/data"
	"github.com/deis/workflow-manager-api/pkg/handlers"
	"github.com/deis/workflow-manager-api/pkg/swagger/restapi/operations"
	errors "github.com/go-swagger/go-swagger/errors"
	httpkit "github.com/go-swagger/go-swagger/httpkit"
	middleware "github.com/go-swagger/go-swagger/httpkit/middleware"
	"github.com/jinzhu/gorm"
)

type GormDb struct {
	db *gorm.DB
}

// This file is safe to edit. Once it exists it will not be overwritten
func getDb(api *operations.WorkflowManagerAPI) *gorm.DB {
	for _, optsGroup := range api.CommandLineOptionsGroups {
		if optsGroup.ShortDescription == "deisUnitTests" {
			gormDb, ok := optsGroup.Options.(GormDb)
			if !ok {
				log.Fatalf("unable to cast to gorm db\n")
			}
			return gormDb.db
		}
	}
	db, err := data.NewDB()
	if err != nil {
		log.Fatalf("unable to create connection to DB (%s)", err)
	}
	if err := data.VerifyPersistentStorage(db); err != nil {
		log.Fatalf("unable to verify persistent storage\n%s", err)
	}
	return db
}

func configureFlags(api *operations.WorkflowManagerAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.WorkflowManagerAPI) http.Handler {

	db := getDb(api)
	db.LogMode(true)
	// configure the api here
	api.ServeError = errors.ServeError

	api.JSONConsumer = httpkit.JSONConsumer()

	api.JSONProducer = httpkit.JSONProducer()

	api.BasicAuth = func(user string, pass string) (interface{}, error) {
		expectedUser := config.Spec.DoctorAuthUser
		expectedPass := config.Spec.DoctorAuthPass
		if user == expectedUser && pass == expectedPass {
			return user, nil
		}
		return nil, errors.Unauthenticated("basic auth")
	}

	api.CreateClusterDetailsHandler = operations.CreateClusterDetailsHandlerFunc(func(params operations.CreateClusterDetailsParams) middleware.Responder {
		return handlers.ClusterCheckin(params, db)
	})

	api.CreateClusterDetailsForV2Handler = operations.CreateClusterDetailsForV2HandlerFunc(func(params operations.CreateClusterDetailsForV2Params) middleware.Responder {
		return handlers.ClusterCheckin(operations.CreateClusterDetailsParams{Body: params.Body}, db)
	})

	api.GetClusterByIDHandler = operations.GetClusterByIDHandlerFunc(func(params operations.GetClusterByIDParams) middleware.Responder {
		return handlers.GetCluster(params, db)
	})
	api.GetClustersByAgeHandler = operations.GetClustersByAgeHandlerFunc(func(params operations.GetClustersByAgeParams) middleware.Responder {
		return handlers.ClustersAge(params, db)
	})
	api.GetClustersCountHandler = operations.GetClustersCountHandlerFunc(func() middleware.Responder {
		return handlers.ClustersCount(db)
	})
	api.GetClusterCheckinsHandler = operations.GetClusterCheckinsHandlerFunc(func(params operations.GetClusterCheckinsParams) middleware.Responder {
		return handlers.ClusterCheckins(params, db)
	})
	api.GetPersistentClustersHandler = operations.GetPersistentClustersHandlerFunc(func(params operations.GetPersistentClustersParams) middleware.Responder {
		return handlers.PersistentClusters(params, db)
	})
	api.GetComponentByNameHandler = operations.GetComponentByNameHandlerFunc(func(params operations.GetComponentByNameParams) middleware.Responder {
		return handlers.GetComponentTrainVersions(params, db)
	})
	api.GetComponentByReleaseHandler = operations.GetComponentByReleaseHandlerFunc(func(params operations.GetComponentByReleaseParams) middleware.Responder {
		return handlers.GetVersion(params, db)
	})
	api.GetComponentsByLatestReleaseHandler = operations.GetComponentsByLatestReleaseHandlerFunc(func(params operations.GetComponentsByLatestReleaseParams) middleware.Responder {
		return handlers.GetLatestVersions(params, db)
	})
	api.GetComponentsByLatestReleaseForV2Handler = operations.GetComponentsByLatestReleaseForV2HandlerFunc(func(params operations.GetComponentsByLatestReleaseForV2Params) middleware.Responder {
		return handlers.GetLatestVersionsForV2(params, db)
	})
	api.GetDoctorInfoHandler = operations.GetDoctorInfoHandlerFunc(func(params operations.GetDoctorInfoParams, principal interface{}) middleware.Responder {
		return handlers.GetDoctor(params, db)
	})
	api.PublishComponentReleaseHandler = operations.PublishComponentReleaseHandlerFunc(func(params operations.PublishComponentReleaseParams) middleware.Responder {
		return handlers.PublishVersion(params, db)
	})
	api.PublishDoctorInfoHandler = operations.PublishDoctorInfoHandlerFunc(func(params operations.PublishDoctorInfoParams) middleware.Responder {
		return handlers.PublishDoctor(params, db)
	})
	api.PingHandler = operations.PingHandlerFunc(func() middleware.Responder {
		return handlers.Ping()
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
