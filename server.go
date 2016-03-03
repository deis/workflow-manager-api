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
	err := data.VerifyPersistentStorage()
	if err != nil {
		log.Fatal("unable to verify persistent storage")
	}
	r := getRoutes()
	err = http.ListenAndServeTLS(":"+listenPort, "server.pem", "server.key", r)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func getRoutes() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/versions/{component}", handlers.VersionsGetHandler).Methods("GET")
	r.HandleFunc("/versions/{component}", handlers.VersionsPostHandler).Methods("POST")
	r.HandleFunc("/clusters", handlers.ClustersHandler(data.RDSDB{}, data.ClusterCount{})).Methods("GET")
	r.HandleFunc("/clusters/{id}", handlers.ClustersGetHandler(data.RDSDB{}, data.ClusterFromDB{})).Methods("GET")
	r.HandleFunc("/clusters/{id}", handlers.ClustersPostHandler(data.RDSDB{}, data.ClusterFromDB{})).Methods("POST")
	return r
}
