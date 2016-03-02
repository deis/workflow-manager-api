package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"database/sql"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/deis/workflow-manager/components"
	"github.com/deis/workflow-manager/types"
	"github.com/gorilla/mux"
	sqlxTypes "github.com/jmoiron/sqlx/types"
	_ "github.com/lib/pq"
)

// package-level constants
const (
	listenPort                = "8443"
	dBNameKey                 = "WORKFLOW_MANAGER_API_DBNAME"
	clustersTableName         = "clusters"
	clustersTableIDKey        = "cluster_id"
	clustersTableFirstSeenKey = "first_seen"
	clustersTableLastSeenKey  = "last_seen"
	clustersTableDataKey      = "data"
)

var (
	dbName         = os.Getenv(dBNameKey)
	memoClusters   = make(map[string]types.Cluster)
	latestVersions = make(map[string]types.Version)
	mu             sync.Mutex
)

// ClusterTable type that expresses the `clusters` postgres table schema
type ClusterTable struct {
	clusterID string
	firstSeen time.Time
	lastSeen  time.Time
	data      sqlxTypes.JSONText
}

// DeisComponentVersion type that expresses the `deis_component_versions` postgres table schema
type DeisComponentVersion struct {
	name     string
	version  string
	released time.Time
}

func getRDSSession() *rds.RDS {
	return rds.New(session.New(), &aws.Config{Region: aws.String("us-west-1")})
}

func getDB() (*sql.DB, error) {
	svc := getRDSSession()
	dbInstanceIdentifier := new(string)
	dbInstanceIdentifier = &dbName
	params := rds.DescribeDBInstancesInput{DBInstanceIdentifier: dbInstanceIdentifier}
	resp, err := svc.DescribeDBInstances(&params)
	if err != nil {
		return nil, err
	}
	if len(resp.DBInstances) > 1 {
		log.Printf("more than one database instance returned for %s, using the 1st one\n", dbName)
	}
	instance := resp.DBInstances[0]
	url := *instance.Endpoint.Address + ":" + strconv.FormatInt(*instance.Endpoint.Port, 10)
	dataSourceName := "postgres://golang:golangadmin@" + url + "/" + *instance.DBName + "?sslmode=require"
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		fmt.Println("Failed to keep connection alive")
		return nil, err
	}
	return db, nil
}

func createClustersTable(db *sql.DB) (sql.Result, error) {
	return db.Exec(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s ( %s uuid PRIMARY KEY, %s timestamp, %s timestamp DEFAULT current_timestamp, %s json )", clustersTableName, clustersTableIDKey, clustersTableFirstSeenKey, clustersTableLastSeenKey, clustersTableDataKey))
}

func getClustersDB(db *sql.DB) (*sql.Rows, error) {
	if _, err := createClustersTable(db); err != nil {
		log.Println("unable to verify clusters table exists")
	}
	return getAllRows(db, clustersTableName)
}

func getAllRows(db *sql.DB, table string) (*sql.Rows, error) {
	return db.Query("SELECT * FROM Clusters")
}

func newClusterDB(db *sql.DB, id string, data []byte) (sql.Result, error) {
	now := time.Now().Format(time.RFC3339)
	insert := fmt.Sprintf("INSERT INTO %s (cluster_id, first_seen, last_seen, data) VALUES('%s', '%s', '%s' '%s')", clustersTableName, id, now, now, string(data))
	return db.Exec(insert)
}

func setClusterDB(db *sql.DB, id string, data []byte) (sql.Result, error) {
	now := time.Now().Format(time.RFC3339)
	update := fmt.Sprintf("UPDATE %s SET data='%s', last_seen='%s' WHERE cluster_id='%s'", clustersTableName, string(data), now, id)
	return db.Exec(update)
}

func getClusterDB(db *sql.DB, id string) *sql.Row {
	return db.QueryRow(fmt.Sprintf("SELECT * FROM %s WHERE %s = '%s'", clustersTableName, clustersTableIDKey, id))
}

// Main opens up a TLS listening port
func main() {
	r := getRoutes()
	// Bind to a port and pass our router in
	err := http.ListenAndServeTLS(":"+listenPort, "server.pem", "server.key", r)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

// get all cluster data
func getAll() map[string]types.Cluster {
	return memoClusters
}

// get all versions data
func getVersions() map[string]types.Version {
	return latestVersions
}

// make a new Cluster struct
func newCluster() types.Cluster {
	return types.Cluster{}
}

// make a new ComponentVersion struct
func newComponentVersion() types.ComponentVersion {
	return types.ComponentVersion{}
}

// get a cluster record, returns a new Cluster that the caller can optionally use
func getCluster(id string) (types.Cluster, bool) {
	db, err := getDB()
	if err != nil {
		log.Fatal("couldn't get a db connection!")
	}
	row := getClusterDB(db, id)
	rowResult := ClusterTable{}
	if err := row.Scan(&rowResult.clusterID, &rowResult.firstSeen, &rowResult.lastSeen, &rowResult.data); err != nil {
		return types.Cluster{}, false
	}
	cluster, err := components.ParseJSONCluster(rowResult.data)
	if err != nil {
		log.Println("error parsing cluster")
		return types.Cluster{}, false
	}
	cluster.FirstSeen = rowResult.firstSeen
	cluster.LastSeen = rowResult.lastSeen
	return cluster, true
}

// get a component version record, returns a new ComponentVersion that the caller can optionally use
func getComponentVersion(name string) (types.ComponentVersion, bool) {
	version, ok := latestVersions[name]
	if !ok {
		return newComponentVersion(), false
	}
	componentVersion := types.ComponentVersion{Component: types.Component{Name: name}, Version: version}
	return componentVersion, true
}

// cluster record set'er
func setCluster(id string, c types.Cluster) types.Cluster {
	var cluster types.Cluster // return variable
	mu.Lock()
	js, err := json.Marshal(c)
	if err != nil {
		fmt.Println("error marshaling data")
	}
	db, err := getDB()
	if err != nil {
		log.Fatal("couldn't get a db connection!")
	}
	row := getClusterDB(db, id)
	var result sql.Result
	rowResult := ClusterTable{}
	if err := row.Scan(&rowResult.clusterID, &rowResult.firstSeen, &rowResult.lastSeen, &rowResult.data); err != nil {
		result, err = newClusterDB(db, id, js)
		if err != nil {
			log.Println(err)
		}
	} else {
		result, err = setClusterDB(db, id, js)
		if err != nil {
			log.Println(err)
		}
	}
	affected, err := result.RowsAffected()
	if err != nil {
		log.Println("failed to get affected row count")
	}
	var ok bool
	if affected == 0 {
		log.Println("no records updated")
	} else if affected == 1 {
		cluster, ok = getCluster(id)
		if !ok {
			log.Println("couldn't get cluster after update")
		}
	} else if affected > 1 {
		log.Println("updated more than one record with same ID value!")
	}
	mu.Unlock()
	return cluster
}

// component version record set'er
func setLatestVersion(cV types.ComponentVersion) types.Version {
	latestVersions[cV.Component.Name] = cV.Version
	return latestVersions[cV.Component.Name]
}

func getRoutes() *mux.Router {
	r := mux.NewRouter()
	// Routes consist of a path and a handler function.
	r.HandleFunc("/clusters", defaultHandler)
	r.HandleFunc("/versions", versionsHandler).Methods("GET")
	r.HandleFunc("/versions", versionsPostHandler).Methods("POST")
	r.HandleFunc("/clusters/{id}", clustersHandler).Methods("GET")
	r.HandleFunc("/clusters/{id}", clustersPostHandler).Methods("POST")
	return r
}

// handler echoes the HTTP request.
func defaultHandler(w http.ResponseWriter, r *http.Request) {
	js, err := json.Marshal(getAll())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(js, w)
}

// versionsHandler handles GET requests to "/versions"
func versionsHandler(w http.ResponseWriter, r *http.Request) {
	js, err := json.Marshal(getVersions())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(js, w)
}

func writeJSON(json []byte, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

// clustersHandler handles GET requests to "/clusters/{id}"
func clustersHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	mu.Lock()
	cluster, ok := getCluster(id)
	if !ok {
		http.NotFound(w, r)
		return
	}
	js, err := json.Marshal(cluster)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(js, w)
	mu.Unlock()
}

// clustersPostHandler handles POST requests to /clusters/{id}
func clustersPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "expected application/json", http.StatusUnsupportedMediaType)
		return
	}
	id := mux.Vars(r)["id"]
	cluster := types.Cluster{}
	err := json.NewDecoder(r.Body).Decode(&cluster)
	if err != nil {
		log.Print(err)
		return
	}
	data, err := json.MarshalIndent(setCluster(id, cluster), "", "  ")
	if err != nil {
		log.Fatalf("JSON marshaling failed: %s", err)
	}
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(data))
}

// versionsPostHandler handles POST requests to /versions
func versionsPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "expected application/json", http.StatusUnsupportedMediaType)
		return
	}
	mu.Lock()
	componentVersions := []types.ComponentVersion{}
	err := json.NewDecoder(r.Body).Decode(&componentVersions)
	if err != nil {
		log.Print(err)
		return
	}
	for _, version := range componentVersions {
		setLatestVersion(version)
	}
	mu.Unlock()
	data, err := json.MarshalIndent(componentVersions, "", "  ")
	if err != nil {
		log.Fatalf("JSON marshaling failed: %s", err)
	}
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(data))
}
