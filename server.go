package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/deis/workflow-manager-api/data"
	"github.com/deis/workflow-manager-api/handlers"
	"github.com/deis/workflow-manager-api/rest"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

const (
	listenPortKey = "WORKFLOW_MANAGER_API_PORT"
)

var (
	listenPort = os.Getenv(listenPortKey)
)

// Main opens up a TLS listening port
func main() {
	db, err := data.VerifyPersistentStorage(data.RDSDB{})
	if err != nil {
		log.Fatalf("unable to verify persistent storage\n%s", err)
	}
	r := getRoutes(db, data.VersionFromDB{}, data.ClusterCount{}, data.ClusterFromDB{})
	if err := http.ListenAndServe(":"+listenPort, r); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func getRoutes(db *sql.DB, version data.Version, count data.Count, cluster data.Cluster) *mux.Router {
	r := mux.NewRouter()
	r.Handle("/{apiVersion}/versions/latest", handlers.GetLatestVersions(db, version)).Methods("POST").
		Headers(handlers.ContentTypeHeaderKey, handlers.JSONContentType)
	r.Handle("/{apiVersion}/versions/{train}/{component}", handlers.GetComponentTrainVersions(db, version)).Methods("GET")
	// Note: the following route must go before the route that ends with {version}, so that
	// Gorilla mux always routes the static "latest" route to the appropriate handler, and "latest"
	// doesn't get interpreted as a {version}
	r.Handle("/{apiVersion}/versions/{train}/{component}/latest", handlers.GetLatestComponentTrainVersion(db, version)).Methods("GET")
	r.Handle("/{apiVersion}/versions/{train}/{component}/{version}", handlers.GetVersion(db, version)).Methods("GET")
	r.Handle("/{apiVersion}/versions/{train}/{component}/{version}", handlers.PublishVersion(db, version)).Methods("POST").
		Headers(handlers.ContentTypeHeaderKey, handlers.JSONContentType)
	r.Handle("/{apiVersion}/clusters/count", handlers.ClustersCount(db, count)).Methods("GET")
	r.Handle("/{apiVersion}/clusters/age", handlers.ClustersAge(db, cluster)).Methods("GET").
		Queries(
			rest.CheckedInBeforeQueryStringKey,
			"",
			rest.CheckedInAfterQueryStringKey,
			"",
			rest.CreatedBeforeQueryStringKey,
			"",
			rest.CreatedAfterQueryStringKey,
			"",
		)
	r.Handle("/{apiVersion}/clusters/{id}", handlers.GetCluster(db, cluster)).Methods("GET")
	r.Handle("/{apiVersion}/clusters/{id}", handlers.ClusterCheckin(db, cluster)).Methods("POST").
		Headers(handlers.ContentTypeHeaderKey, handlers.JSONContentType)
	return r
}
