package main

import (
	"log"
	"net/http"
	"os"

	"github.com/deis/workflow-manager-api/data"
	"github.com/deis/workflow-manager-api/handlers"
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
	db := data.RDSDB{}
	if err := data.VerifyPersistentStorage(db); err != nil {
		log.Fatal("unable to verify persistent storage")
	}
	r := getRoutes(db, data.VersionFromDB{}, data.ClusterCount{}, data.ClusterFromDB{})
	if err := http.ListenAndServeTLS(":"+listenPort, "server.pem", "server.key", r); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func getRoutes(db data.DB, version data.Version, count data.Count, cluster data.Cluster) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/{apiVersion}/versions/{component}", handlers.VersionsGetHandler(db, version)).Methods("GET")
	r.HandleFunc("/{apiVersion}/versions/{component}", handlers.VersionsPostHandler(db, version)).Methods("POST")
	r.HandleFunc("/{apiVersion}/clusters", handlers.ClustersHandler(db, count)).Methods("GET")
	r.HandleFunc("/{apiVersion}/clusters/{id}", handlers.ClustersGetHandler(db, cluster)).Methods("GET")
	r.HandleFunc("/{apiVersion}/clusters/{id}", handlers.ClustersPostHandler(db, cluster)).Methods("POST")
	return r
}
