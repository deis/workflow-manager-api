package data

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"

	"database/sql"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/deis/workflow-manager/components"
	"github.com/deis/workflow-manager/types"
	sqlxTypes "github.com/jmoiron/sqlx/types"
	_ "github.com/lib/pq" // Pure Go Postgres driver for database/sql
)

const (
	rDSRegionKey                             = "WORKFLOW_MANAGER_API_RDS_REGION"
	dBInstanceKey                            = "WORKFLOW_MANAGER_API_DBINSTANCE"
	dBUserKey                                = "WORKFLOW_MANAGER_API_DBUSER"
	dBPassKey                                = "WORKFLOW_MANAGER_API_DBPASS"
	clustersTableName                        = "clusters"
	clustersTableIDKey                       = "cluster_id"
	clustersTableFirstSeenKey                = "first_seen"
	clustersTableLastSeenKey                 = "last_seen"
	clustersTableDataKey                     = "data"
	clustersCheckinsTableName                = "clusters_checkins"
	clustersCheckinsTableIDKey               = "checkins_id"
	clustersCheckinsTableClusterIDKey        = "cluster_id"
	clustersCheckinsTableClusterCreatedAtKey = "created_at"
	clustersCheckinsTableDataKey             = "data"
	versionsTableName                        = "versions"
	versionsTableIDKey                       = "version_id"
	versionsTableComponentNameKey            = "component_name"
	versionsTableTrainKey                    = "train"
	versionsTableVersionKey                  = "version"
	versionsTableReleaseTimeStampKey         = "release_timestamp"
	versionsTableDataKey                     = "data"
)

var (
	rDSRegion                 = os.Getenv(rDSRegionKey)
	dBInstance                = os.Getenv(dBInstanceKey)
	dBUser                    = os.Getenv(dBUserKey)
	dBPass                    = os.Getenv(dBPassKey)
	mu                        sync.Mutex
	errInvalidDBRecordRequest = errors.New("invalid DB record request")
)

// ClustersTable type that expresses the `clusters` postgres table schema
type ClustersTable struct {
	clusterID string // PRIMARY KEY
	firstSeen *Timestamp
	lastSeen  *Timestamp
	data      sqlxTypes.JSONText
}

// ClustersCheckinsTable type that expresses the `clusters_checkins` postgres table schema
type ClustersCheckinsTable struct {
	checkinID string     // PRIMARY KEY, type uuid
	clusterID string     // indexed
	createdAt *Timestamp // indexed
	data      sqlxTypes.JSONText
}

// VersionsTable type that expresses the `deis_component_versions` postgres table schema
type VersionsTable struct {
	versionID        string // PRIMARY KEY
	componentName    string // indexed
	train            string // indexed
	version          string // indexed
	releaseTimestamp *Timestamp
	data             sqlxTypes.JSONText
}

// DB is an interface for managing a DB instance
type DB interface {
	Get() (*sql.DB, error)
}

// Cluster is an interface for managing a persistent cluster record
type Cluster interface {
	Get(*sql.DB, string) (types.Cluster, error)
	Set(*sql.DB, string, types.Cluster) (types.Cluster, error)
	Checkin(*sql.DB, string, types.Cluster) (sql.Result, error)
}

// ClusterFromDB fulfills the Cluster interface
type ClusterFromDB struct{}

// Get method for ClusterFromDB, the actual database/sql.DB implementation
func (c ClusterFromDB) Get(db *sql.DB, id string) (types.Cluster, error) {
	row := getDBRecord(db, clustersTableName, []string{clustersTableIDKey}, []string{id})
	rowResult := ClustersTable{}
	if err := row.Scan(&rowResult.clusterID, &rowResult.firstSeen, &rowResult.lastSeen, &rowResult.data); err != nil {
		return types.Cluster{}, err
	}
	cluster, err := components.ParseJSONCluster(rowResult.data)
	if err != nil {
		log.Println("error parsing cluster")
		return types.Cluster{}, err
	}
	cluster.FirstSeen = *rowResult.firstSeen.Time
	cluster.LastSeen = *rowResult.lastSeen.Time
	return cluster, nil
}

// Set method for ClusterFromDB, the actual database/sql.DB implementation
func (c ClusterFromDB) Set(db *sql.DB, id string, cluster types.Cluster) (types.Cluster, error) {
	var ret types.Cluster // return variable
	mu.Lock()
	js, err := json.Marshal(cluster)
	if err != nil {
		fmt.Println("error marshaling data")
	}
	row := getDBRecord(db, clustersTableName, []string{clustersTableIDKey}, []string{id})
	var result sql.Result
	// Register the "latest checkin" with the primary cluster record
	rowResult := ClustersTable{}
	if err := row.Scan(&rowResult.clusterID, &rowResult.firstSeen, &rowResult.lastSeen, &rowResult.data); err != nil {
		result, err = newClusterDBRecord(db, id, js)
		if err != nil {
			log.Println(err)
		}
	} else {
		result, err = updateClusterDBRecord(db, id, js)
		if err != nil {
			log.Println(err)
		}
	}
	affected, err := result.RowsAffected()
	if err != nil {
		log.Println("failed to get affected row count")
	}
	if affected == 0 {
		log.Println("no records updated")
	} else if affected == 1 {
		ret, err = c.Get(db, id)
		if err != nil {
			return types.Cluster{}, err
		}
	} else if affected > 1 {
		log.Println("updated more than one record with same ID value!")
	}
	mu.Unlock()
	return ret, nil
}

// Checkin method for ClusterFromDB, the actual database/sql.DB implementation
func (c ClusterFromDB) Checkin(db *sql.DB, id string, cluster types.Cluster) (sql.Result, error) {
	js, err := json.Marshal(cluster)
	if err != nil {
		fmt.Println("error marshaling data")
	}
	result, err := newClusterCheckinsDBRecord(db, id, js)
	if err != nil {
		log.Println("cluster checkin db record not created", err)
		return nil, err
	}
	return result, nil
}

// Version is an interface for managing a persistent cluster record
type Version interface {
	// Retrieve a single Version record from a DB matching the unique property values in a ComponentVersion struct
	Get(*sql.DB, types.ComponentVersion) (types.ComponentVersion, error)
	// Retrieve a list of Version records that match a given component + train
	Collection(db *sql.DB, train string, component string) ([]types.ComponentVersion, error)
	// Retrieve the most recent Version record that matches a given component + train
	Latest(db *sql.DB, train string, component string) (types.ComponentVersion, error)
	// Store/Update a single Version record into a DB
	Set(*sql.DB, types.ComponentVersion) (types.ComponentVersion, error)
}

// VersionFromDB fulfills the Version interface
type VersionFromDB struct{}

// Get method for VersionFromDB, the actual database/sql.DB implementation
func (c VersionFromDB) Get(db *sql.DB, cV types.ComponentVersion) (types.ComponentVersion, error) {
	row := getDBRecord(db, versionsTableName,
		[]string{versionsTableComponentNameKey, versionsTableTrainKey, versionsTableVersionKey},
		[]string{cV.Component.Name, cV.Version.Train, cV.Version.Version})
	rowResult := VersionsTable{}
	//TODO: sql.NullString is to pass tests, not for production
	var s sql.NullString
	if err := row.Scan(&s, &rowResult.componentName, &rowResult.train, &rowResult.version, &rowResult.releaseTimestamp, &rowResult.data); err != nil {
		return types.ComponentVersion{}, err
	}
	componentVersion := types.ComponentVersion{
		Component: types.Component{
			Name: rowResult.componentName,
		},
		Version: types.Version{
			Version:  rowResult.version,
			Released: rowResult.releaseTimestamp.String(),
			Train:    rowResult.train,
			Data:     rowResult.data,
		},
	}
	return componentVersion, nil
}

// Collection method for VersionFromDB, the actual database/sql.DB implementation
func (c VersionFromDB) Collection(db *sql.DB, train string, component string) ([]types.ComponentVersion, error) {
	rows, err := getDBRecords(db, versionsTableName,
		[]string{versionsTableTrainKey, versionsTableComponentNameKey},
		[]string{train, component})
	if err != nil {
		return []types.ComponentVersion{}, err
	}
	rowsResult := []VersionsTable{}
	var row VersionsTable
	defer rows.Close()
	for rows.Next() {
		//TODO: sql.NullString is to pass tests, not for production
		var s sql.NullString
		err = rows.Scan(&s, &row.componentName,
			&row.train, &row.version, &row.releaseTimestamp, &row.data)
		if err != nil {
			return []types.ComponentVersion{}, err
		}
		rowsResult = append(rowsResult, row)
	}
	componentVersions, err := parseDBVersions(rowsResult)
	if err != nil {
		log.Println("error parsing DB versions data")
		return []types.ComponentVersion{}, err
	}
	return componentVersions, nil
}

// Latest method for VersionFromDB, the actual database/sql.DB implementation
func (c VersionFromDB) Latest(db *sql.DB, train string, component string) (types.ComponentVersion, error) {
	// TODO: implement
	return types.ComponentVersion{}, nil
}

// Set method for VersionFromDB, the actual database/sql.DB implementation
func (c VersionFromDB) Set(db *sql.DB, componentVersion types.ComponentVersion) (types.ComponentVersion, error) {
	var ret types.ComponentVersion // return variable
	mu.Lock()
	row := getDBRecord(db, versionsTableName,
		[]string{versionsTableComponentNameKey, versionsTableTrainKey, versionsTableVersionKey},
		[]string{componentVersion.Component.Name, componentVersion.Version.Train, componentVersion.Version.Version})
	var result sql.Result
	rowResult := VersionsTable{}
	if err := row.Scan(&rowResult.versionID, &rowResult.componentName, &rowResult.train, &rowResult.version, &rowResult.releaseTimestamp, &rowResult.data); err != nil {
		result, err = newVersionDBRecord(db, componentVersion)
		if err != nil {
			log.Println(err)
		}
	} else {
		result, err = updateVersionDBRecord(db, componentVersion)
		if err != nil {
			log.Println(err)
		}
	}
	affected, err := result.RowsAffected()
	if err != nil {
		log.Println("failed to get affected row count")
	}
	if affected == 0 {
		log.Println("no records updated")
	} else if affected == 1 {
		ret, err = c.Get(db, componentVersion)
		if err != nil {
			return types.ComponentVersion{}, err
		}
	} else if affected > 1 {
		log.Println("updated more than one record with same ID value!")
	}
	mu.Unlock()
	return ret, nil
}

// Count is an interface for managing a record count
type Count interface {
	Get(db *sql.DB) (int, error)
}

// ClusterCount fulfills the Count interface
type ClusterCount struct{}

// Get method for ClusterCount
func (c ClusterCount) Get(db *sql.DB) (int, error) {
	count, err := getTableCount(db, clustersTableName)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// RDSDB fulfills the DB interface
type RDSDB struct{}

// Get method for RDSDB
func (r RDSDB) Get() (*sql.DB, error) {
	db, err := getRDSDB()
	if err != nil {
		return nil, err
	}
	return db, nil
}

// VerifyPersistentStorage is a high level interace for verifying storage abstractions
func VerifyPersistentStorage(dbGetter DB) (*sql.DB, error) {
	db, err := dbGetter.Get()
	if err != nil {
		return nil, err
	}
	if err := verifyVersionsTable(db); err != nil {
		log.Println("unable to verify " + versionsTableName + " table")
		return db, err
	}
	count, err := getTableCount(db, versionsTableName)
	if err != nil {
		log.Println("unable to get record count for " + versionsTableName + " table")
		return db, err
	}
	log.Println("counted " + strconv.Itoa(count) + " records for " + versionsTableName + " table")
	err = verifyClustersTable(db)
	if err != nil {
		log.Println("unable to verify " + clustersTableName + " table")
		return db, err
	}
	count, err = getTableCount(db, clustersTableName)
	if err != nil {
		log.Println("unable to get record count for " + clustersTableName + " table")
		return db, err
	}
	log.Println("counted " + strconv.Itoa(count) + " records for " + clustersTableName + " table")
	err = verifyClustersCheckinsTable(db)
	if err != nil {
		log.Println("unable to verify " + clustersCheckinsTableName + " table")
		return db, err
	}
	count, err = getTableCount(db, clustersCheckinsTableName)
	if err != nil {
		log.Println("unable to get record count for " + clustersCheckinsTableName + " table")
		return db, err
	}
	log.Println("counted " + strconv.Itoa(count) + " records for " + clustersCheckinsTableName + " table")
	return db, nil
}

// GetClusterCount is a high level interface for retrieving a simple cluster count
func GetClusterCount(db *sql.DB, c Count) (int, error) {
	count, err := c.Get(db)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// GetCluster is a high level interface for retrieving a cluster data record
func GetCluster(id string, db *sql.DB, c Cluster) (types.Cluster, error) {
	cluster, err := c.Get(db, id)
	if err != nil {
		return types.Cluster{}, err
	}
	return cluster, nil
}

// SetCluster is a high level interface for updating a cluster data record
func SetCluster(id string, cluster types.Cluster, db *sql.DB, c Cluster) (types.Cluster, error) {
	// Check in
	_, err := c.Checkin(db, id, cluster)
	if err != nil {
		return types.Cluster{}, err
	}
	// Update cluster record
	ret, err := c.Set(db, id, cluster)
	if err != nil {
		return types.Cluster{}, err
	}
	return ret, nil
}

// GetVersion is a high level interface for retrieving a version data record
func GetVersion(componentVersion types.ComponentVersion, db *sql.DB, v Version) (types.ComponentVersion, error) {
	componentVersion, err := v.Get(db, componentVersion)
	if err != nil {
		return types.ComponentVersion{}, err
	}
	return componentVersion, nil
}

// GetComponentTrainVersions is a high level interface for retrieving component versions for a given "train"
func GetComponentTrainVersions(train string, component string, db *sql.DB, v Version) ([]types.ComponentVersion, error) {
	componentVersions, err := v.Collection(db, train, component)
	if err != nil {
		return []types.ComponentVersion{}, err
	}
	return componentVersions, nil
}

// GetLatestComponentTrainVersion is a high level interface for retrieving the latest component version for a given "train"
func GetLatestComponentTrainVersion(train string, component string, db *sql.DB, v Version) (types.ComponentVersion, error) {
	componentVersion, err := v.Latest(db, train, component)
	if err != nil {
		return types.ComponentVersion{}, err
	}
	return componentVersion, nil
}

// SetVersion is a high level interface for updating a component version record
func SetVersion(componentVersion types.ComponentVersion, db *sql.DB, v Version) (types.ComponentVersion, error) {
	ret, err := v.Set(db, componentVersion)
	if err != nil {
		return types.ComponentVersion{}, err
	}
	return ret, nil
}

func getRDSSession() *rds.RDS {
	return rds.New(session.New(), &aws.Config{Region: aws.String(rDSRegion)})
}

func getRDSDB() (*sql.DB, error) {
	svc := getRDSSession()
	dbInstanceIdentifier := new(string)
	dbInstanceIdentifier = &dBInstance
	params := rds.DescribeDBInstancesInput{DBInstanceIdentifier: dbInstanceIdentifier}
	resp, err := svc.DescribeDBInstances(&params)
	if err != nil {
		return nil, err
	}
	if len(resp.DBInstances) > 1 {
		log.Printf("more than one database instance returned for %s, using the 1st one", dBInstance)
	}
	instance := resp.DBInstances[0]
	url := *instance.Endpoint.Address + ":" + strconv.FormatInt(*instance.Endpoint.Port, 10)
	dataSourceName := "postgres://" + dBUser + ":" + dBPass + "@" + url + "/" + *instance.DBName + "?sslmode=require"
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		log.Println("couldn't get a db connection!")
		return nil, err
	}
	if err := db.Ping(); err != nil {
		log.Println("Failed to keep db connection alive")
		return nil, err
	}
	return db, nil
}

func createClustersTable(db *sql.DB) (sql.Result, error) {
	return db.Exec(fmt.Sprintf(
		"CREATE TABLE IF NOT EXISTS %s ( %s uuid PRIMARY KEY, %s timestamp, %s timestamp DEFAULT current_timestamp, %s json )",
		clustersTableName,
		clustersTableIDKey,
		clustersTableFirstSeenKey,
		clustersTableLastSeenKey,
		clustersTableDataKey,
	))
}

func createClustersCheckinsTable(db *sql.DB) (sql.Result, error) {
	return db.Exec(fmt.Sprintf(
		"CREATE TABLE IF NOT EXISTS %s ( %s bigserial PRIMARY KEY, %s uuid, %s timestamp, %s json )",
		clustersCheckinsTableName,
		clustersCheckinsTableIDKey,
		clustersTableIDKey,
		clustersCheckinsTableClusterCreatedAtKey,
		clustersCheckinsTableDataKey,
	))
}

func createVersionsTable(db *sql.DB) (sql.Result, error) {
	query := fmt.Sprintf(
		"CREATE TABLE IF NOT EXISTS %s ( %s bigserial PRIMARY KEY, %s varchar(32), %s varchar(24), %s varchar(32), %s timestamp, %s json, unique (%s, %s, %s) )",
		versionsTableName,
		versionsTableIDKey,
		versionsTableComponentNameKey,
		versionsTableTrainKey,
		versionsTableVersionKey,
		versionsTableReleaseTimeStampKey,
		versionsTableDataKey,
		versionsTableComponentNameKey,
		versionsTableTrainKey,
		versionsTableVersionKey,
	)
	return db.Exec(query)
}

func verifyClustersTable(db *sql.DB) error {
	if _, err := createClustersTable(db); err != nil {
		log.Println("unable to verify clusters table exists")
		return err
	}
	return nil
}

func verifyClustersCheckinsTable(db *sql.DB) error {
	if _, err := createClustersCheckinsTable(db); err != nil {
		log.Println("unable to verify clusters table exists")
		return err
	}
	return nil
}

func verifyVersionsTable(db *sql.DB) error {
	if _, err := createVersionsTable(db); err != nil {
		log.Println("unable to verify versions table exists")
		return err
	}
	return nil
}

// getDBRecord is a convenience that executes a simple "SELECT *" SQL query against
// a passed-in db reference, accepting an arbitrary number of keys(table fields)/vals
// assumes a single record response
func getDBRecord(db *sql.DB, table string, keys []string, vals []string) *sql.Row {
	sliceEqualize(&keys, &vals)
	query := fmt.Sprintf("SELECT * FROM %s", table)
	for i, key := range keys {
		if i == 0 {
			query += fmt.Sprintf(" WHERE %s = '%s'", key, vals[i])
		} else {
			query += fmt.Sprintf(" AND %s = '%s'", key, vals[i])
		}
	}
	return db.QueryRow(query)
}

// getDBRecords is a convenience that executes a simple "SELECT *" SQL query against
// a passed-in db reference, accepting an arbitrary number of keys(table fields)/vals
func getDBRecords(db *sql.DB, table string, keys []string, vals []string) (*sql.Rows, error) {
	sliceEqualize(&keys, &vals)
	query := fmt.Sprintf("SELECT * FROM %s", table)
	for i, key := range keys {
		if i == 0 {
			query += fmt.Sprintf(" WHERE %s = '%s'", key, vals[i])
		} else {
			query += fmt.Sprintf(" AND %s = '%s'", key, vals[i])
		}
	}
	return db.Query(query)
}

// sliceEqualize is a convenience that ensures two slices of strings have equal lengths
// if not, the larger slice's elements that exceed the boundary of the smaller are stripped
func sliceEqualize(slice1 *[]string, slice2 *[]string) {
	if len(*slice1) != len(*slice2) {
		if len(*slice1) > len(*slice2) {
			*slice1 = (*slice1)[:len(*slice2)]
		} else {
			*slice2 = (*slice2)[:len(*slice1)]
		}
	}
}

func getAllRows(db *sql.DB, table string) (*sql.Rows, error) {
	return db.Query(fmt.Sprintf("SELECT * FROM %s", table))
}

func getTableCount(db *sql.DB, table string) (int, error) {
	rows, err := db.Query(fmt.Sprintf("SELECT COUNT(*) FROM %s", table))
	if err != nil {
		return 0, err
	}
	var count int
	for rows.Next() {
		err := rows.Scan(&count)
		if err != nil {
			return 0, err
		}
	}
	return count, nil
}

func newClusterDBRecord(db *sql.DB, id string, data []byte) (sql.Result, error) {
	insert := fmt.Sprintf("INSERT INTO %s (cluster_id, first_seen, last_seen, data) VALUES('%s', '%s', '%s', '%s')", clustersTableName, id, now(), now(), string(data))
	return db.Exec(insert)
}

func newVersionDBRecord(db *sql.DB, componentVersion types.ComponentVersion) (sql.Result, error) {
	insert := fmt.Sprintf(
		"INSERT INTO %s (%s, %s, %s, %s, %s) VALUES('%s', '%s', '%s', '%s', '%s')",
		versionsTableName,
		versionsTableComponentNameKey,
		versionsTableTrainKey,
		versionsTableVersionKey,
		versionsTableReleaseTimeStampKey,
		versionsTableDataKey,
		componentVersion.Component.Name,
		componentVersion.Version.Train,
		componentVersion.Version.Version,
		componentVersion.Version.Released,
		string(componentVersion.Version.Data[:]),
	)
	return db.Exec(insert)
}

func updateVersionDBRecord(db *sql.DB, componentVersion types.ComponentVersion) (sql.Result, error) {
	data, err := json.Marshal(componentVersion.Version.Data)
	if err != nil {
		return nil, err
	}
	update := fmt.Sprintf(
		"UPDATE %s SET %s='%s', %s='%s', %s='%s', %s='%s', %s='%s' WHERE %s='%s' AND %s='%s' AND %s='%s'",
		versionsTableName,
		versionsTableComponentNameKey,
		componentVersion.Component.Name,
		versionsTableTrainKey,
		componentVersion.Version.Train,
		versionsTableVersionKey,
		componentVersion.Version.Version,
		versionsTableReleaseTimeStampKey,
		componentVersion.Version.Released,
		versionsTableDataKey,
		string(data),
		versionsTableComponentNameKey,
		componentVersion.Component.Name,
		versionsTableTrainKey,
		componentVersion.Version.Train,
		versionsTableVersionKey,
		componentVersion.Version.Version,
	)
	return db.Exec(update)
}

func updateClusterDBRecord(db *sql.DB, id string, data []byte) (sql.Result, error) {
	update := fmt.Sprintf("UPDATE %s SET data='%s', last_seen='%s' WHERE cluster_id='%s'", clustersTableName, string(data), now(), id)
	return db.Exec(update)
}

func newClusterCheckinsDBRecord(db *sql.DB, id string, data []byte) (sql.Result, error) {
	update := fmt.Sprintf("INSERT INTO %s (data, created_at, cluster_id) VALUES('%s', '%s', '%s')", clustersCheckinsTableName, string(data), now(), id)
	return db.Exec(update)
}
