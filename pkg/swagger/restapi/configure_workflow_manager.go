package restapi

import (
	"log"
	"net/http"
	"os"

	"github.com/deis/workflow-manager-api/pkg/data"
	"github.com/deis/workflow-manager-api/pkg/handlers"
	"github.com/deis/workflow-manager-api/pkg/swagger/restapi/operations"
	errors "github.com/go-swagger/go-swagger/errors"
	httpkit "github.com/go-swagger/go-swagger/httpkit"
	middleware "github.com/go-swagger/go-swagger/httpkit/middleware"
	"github.com/jinzhu/gorm"
)

const (
	rDSRegionKey  = "WORKFLOW_MANAGER_API_RDS_REGION"
	dBUserKey     = "WORKFLOW_MANAGER_API_DBUSER"
	dBPassKey     = "WORKFLOW_MANAGER_API_DBPASS"
	dBInstanceKey = "WORKFLOW_MANAGER_API_DBINSTANCE"
	dBFlavor      = "postgres"
)

var (
	rdsRegion  = os.Getenv(rDSRegionKey)
	dBInstance = os.Getenv(dBInstanceKey)
	dBUser     = os.Getenv(dBUserKey)
	dBPass     = os.Getenv(dBPassKey)
)

type GormDb struct {
	db *gorm.DB
}

// This file is safe to edit. Once it exists it will not be overwritten
func getDb(api *operations.WorkflowManagerAPI, d data.DB) *gorm.DB {
	for _, optsGroup := range api.CommandLineOptionsGroups {
		if optsGroup.ShortDescription == "deisUnitTests" {
			gormDb, ok := optsGroup.Options.(GormDb)
			if !ok {
				log.Fatalf("unable to cast to gorm db\n")
			}
			return gormDb.db
		}
	}
	db, err := d.Get()
	if err != nil {
		log.Fatalf("unable to create connection to RDS DB (%s)", err)
	}
	if err := data.VerifyPersistentStorage(db); err != nil {
		log.Fatalf("unable to verify persistent storage\n%s", err)
	}
	return db
}

func configureFlags(api *operations.WorkflowManagerAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

// extend the arity of this function to accept as a 2nd arg a *sql.DB
func configureAPI(api *operations.WorkflowManagerAPI) http.Handler {
	db := data.NewRDSDB(rdsRegion, dBUser, dBPass, dBFlavor, dBInstance)
	rdsDB := getDb(api, db)
	rdsDB.LogMode(true)
	// configure the api here
	api.ServeError = errors.ServeError

	api.JSONConsumer = httpkit.JSONConsumer()

	api.JSONProducer = httpkit.JSONProducer()

	api.CreateClusterDetailsHandler = operations.CreateClusterDetailsHandlerFunc(func(params operations.CreateClusterDetailsParams) middleware.Responder {
		return handlers.ClusterCheckin(params, rdsDB)
	})

	api.CreateClusterDetailsForV2Handler = operations.CreateClusterDetailsForV2HandlerFunc(func(params operations.CreateClusterDetailsForV2Params) middleware.Responder {
		return handlers.ClusterCheckin(operations.CreateClusterDetailsParams{Body: params.Body}, rdsDB)
	})

	api.GetClusterByIDHandler = operations.GetClusterByIDHandlerFunc(func(params operations.GetClusterByIDParams) middleware.Responder {
		return handlers.GetCluster(params, rdsDB)
	})
	api.GetClustersByAgeHandler = operations.GetClustersByAgeHandlerFunc(func(params operations.GetClustersByAgeParams) middleware.Responder {
		return handlers.ClustersAge(params, rdsDB)
	})
	api.GetClustersCountHandler = operations.GetClustersCountHandlerFunc(func() middleware.Responder {
		return handlers.ClustersCount(rdsDB)
	})
	api.GetComponentByNameHandler = operations.GetComponentByNameHandlerFunc(func(params operations.GetComponentByNameParams) middleware.Responder {
		return handlers.GetComponentTrainVersions(params, rdsDB)
	})
	api.GetComponentByReleaseHandler = operations.GetComponentByReleaseHandlerFunc(func(params operations.GetComponentByReleaseParams) middleware.Responder {
		return handlers.GetVersion(params, rdsDB)
	})
	api.GetComponentsByLatestReleaseHandler = operations.GetComponentsByLatestReleaseHandlerFunc(func(params operations.GetComponentsByLatestReleaseParams) middleware.Responder {
		return handlers.GetLatestVersions(params, rdsDB)
	})
	api.GetComponentsByLatestReleaseForV2Handler = operations.GetComponentsByLatestReleaseForV2HandlerFunc(func(params operations.GetComponentsByLatestReleaseForV2Params) middleware.Responder {
		return handlers.GetLatestVersionsForV2(params, rdsDB)
	})
	api.GetDoctorInfoHandler = operations.GetDoctorInfoHandlerFunc(func(params operations.GetDoctorInfoParams) middleware.Responder {
		return handlers.GetDoctor(params, rdsDB)
	})
	api.PublishComponentReleaseHandler = operations.PublishComponentReleaseHandlerFunc(func(params operations.PublishComponentReleaseParams) middleware.Responder {
		return handlers.PublishVersion(params, rdsDB)
	})
	api.PublishDoctorInfoHandler = operations.PublishDoctorInfoHandlerFunc(func(params operations.PublishDoctorInfoParams) middleware.Responder {
		return handlers.PublishDoctor(params, rdsDB)
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
