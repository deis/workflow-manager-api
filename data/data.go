package data

import (
	"encoding/json"
	"fmt"
	"log"
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
	sqlxTypes "github.com/jmoiron/sqlx/types"
	_ "github.com/lib/pq" // Pure Go Postgres driver for database/sql
)

const (
	rDSRegionKey                  = "WORKFLOW_MANAGER_API_RDS_REGION"
	dBNameKey                     = "WORKFLOW_MANAGER_API_DBNAME"
	dBUserKey                     = "WORKFLOW_MANAGER_API_DBUSER"
	dBPassKey                     = "WORKFLOW_MANAGER_API_DBPASS"
	clustersTableName             = "clusters"
	clustersTableIDKey            = "cluster_id"
	clustersTableFirstSeenKey     = "first_seen"
	clustersTableLastSeenKey      = "last_seen"
	clustersTableDataKey          = "data"
	versionsTableName             = "versions"
	versionsTableComponentNameKey = "component_name"
	versionsTableLastUpdatedKey   = "last_updated"
	versionsTableDataKey          = "data"
)

var (
	rDSRegion = os.Getenv(rDSRegionKey)
	dBName    = os.Getenv(dBNameKey)
	dBUser    = os.Getenv(dBUserKey)
	dBPass    = os.Getenv(dBPassKey)
	mu        sync.Mutex
)

// ClustersTable type that expresses the `clusters` postgres table schema
type ClustersTable struct {
	clusterID string
	firstSeen time.Time
	lastSeen  time.Time
	data      sqlxTypes.JSONText
}

// VersionsTable type that expresses the `deis_component_versions` postgres table schema
type VersionsTable struct {
	componentName string
	lastUpdated   time.Time
	data          sqlxTypes.JSONText
}

// DB is an interface for managing a DB instance
type DB interface {
	Get() (*sql.DB, error)
}

// Cluster is an interface for managing a persistent cluster record
type Cluster interface {
	Get(*sql.DB, string) (types.Cluster, error)
	Set(*sql.DB, string, types.Cluster) (types.Cluster, error)
}

// ClusterFromDB fulfills the record interface
type ClusterFromDB struct{}

// Get method for ClusterFromDB
func (c ClusterFromDB) Get(db *sql.DB, id string) (types.Cluster, error) {
	row := getDBRecord(db, clustersTableName, clustersTableIDKey, id)
	rowResult := ClustersTable{}
	if err := row.Scan(&rowResult.clusterID, &rowResult.firstSeen, &rowResult.lastSeen, &rowResult.data); err != nil {
		return types.Cluster{}, err
	}
	cluster, err := components.ParseJSONCluster(rowResult.data)
	if err != nil {
		log.Println("error parsing cluster")
		return types.Cluster{}, err
	}
	cluster.FirstSeen = rowResult.firstSeen
	cluster.LastSeen = rowResult.lastSeen
	return cluster, nil
}

// Set method for ClusterFromDB
func (c ClusterFromDB) Set(db *sql.DB, id string, cluster types.Cluster) (types.Cluster, error) {
	var ret types.Cluster // return variable
	mu.Lock()
	js, err := json.Marshal(cluster)
	if err != nil {
		fmt.Println("error marshaling data")
	}
	row := getDBRecord(db, clustersTableName, clustersTableIDKey, id)
	var result sql.Result
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
		ret, err = GetCluster(id, RDSDB{}, ClusterFromDB{})
		if err != nil {
			return types.Cluster{}, err
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
func VerifyPersistentStorage() error {
	db, err := getRDSDB()
	if err != nil {
		log.Println("couldn't get a db connection")
		return err
	}
	err = verifyVersionsTable(db)
	if err != nil {
		log.Println("unable to verify versions table")
		return err
	}
	count, err := getTableCount(db, versionsTableName)
	if err != nil {
		log.Println("unable to get record count for versions table")
		return err
	}
	log.Println("counted " + strconv.Itoa(count) + " records for " + versionsTableName + " table")
	err = verifyClustersTable(db)
	if err != nil {
		log.Println("unable to verify clusters table")
		return err
	}
	count, err = getTableCount(db, clustersTableName)
	if err != nil {
		log.Println("unable to get record count for versions table")
		return err
	}
	log.Println("counted " + strconv.Itoa(count) + " records for " + clustersTableName + " table")
	return nil
}

// GetClusterCount is a high level interface for retrieving a simple cluster count
func GetClusterCount(d DB, c Count) (int, error) {
	db, err := d.Get()
	if err != nil {
		return 0, err
	}
	count, err := c.Get(db)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// GetCluster is a high level interface for retrieving a cluster data record
func GetCluster(id string, d DB, c Cluster) (types.Cluster, error) {
	db, err := d.Get()
	if err != nil {
		return types.Cluster{}, err
	}
	cluster, err := c.Get(db, id)
	if err != nil {
		return types.Cluster{}, err
	}
	return cluster, nil
}

// SetCluster is a high level interface for updating a cluster data record
func SetCluster(id string, cluster types.Cluster, d DB, c Cluster) (types.Cluster, error) {
	db, err := d.Get()
	if err != nil {
		return types.Cluster{}, err
	}
	ret, err := c.Set(db, id, cluster)
	if err != nil {
		return types.Cluster{}, err
	}
	return ret, nil
}

// GetVersion is a high level interface for retrieving a component version record
func GetVersion(component string) (types.ComponentVersion, bool) {
	db, err := getRDSDB()
	if err != nil {
		log.Fatal("couldn't get a db connection!")
	}
	row := getDBRecord(db, versionsTableName, versionsTableComponentNameKey, component)
	rowResult := VersionsTable{}
	if err := row.Scan(&rowResult.componentName, &rowResult.lastUpdated, &rowResult.data); err != nil {
		return types.ComponentVersion{}, false
	}
	componentVersion, err := components.ParseJSONComponent(rowResult.data)
	if err != nil {
		log.Println("error parsing component version")
		return types.ComponentVersion{}, false
	}
	return componentVersion, true
}

// SetVersion is a high level interface for updating a component version record
func SetVersion(component string, c types.ComponentVersion) types.ComponentVersion {
	var componentVersion types.ComponentVersion // return variable
	mu.Lock()
	js, err := json.Marshal(c)
	if err != nil {
		fmt.Println("error marshaling data")
	}
	db, err := getRDSDB()
	if err != nil {
		log.Fatal("couldn't get a db connection!")
	}
	row := getDBRecord(db, versionsTableName, versionsTableComponentNameKey, component)
	var result sql.Result
	rowResult := VersionsTable{}
	if err := row.Scan(&rowResult.componentName, &rowResult.lastUpdated, &rowResult.data); err != nil {
		result, err = newVersionDBRecord(db, component, js)
		if err != nil {
			log.Println(err)
		}
	} else {
		result, err = updateVersionDBRecord(db, component, js)
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
		componentVersion, ok = GetVersion(component)
		if !ok {
			log.Println("couldn't get component version after update")
		}
	} else if affected > 1 {
		log.Println("updated more than one record with same ID value!")
	}
	mu.Unlock()
	return componentVersion
}

func getRDSSession() *rds.RDS {
	return rds.New(session.New(), &aws.Config{Region: aws.String(rDSRegion)})
}

func getRDSDB() (*sql.DB, error) {
	svc := getRDSSession()
	dbInstanceIdentifier := new(string)
	dbInstanceIdentifier = &dBName
	params := rds.DescribeDBInstancesInput{DBInstanceIdentifier: dbInstanceIdentifier}
	resp, err := svc.DescribeDBInstances(&params)
	if err != nil {
		return nil, err
	}
	if len(resp.DBInstances) > 1 {
		log.Printf("more than one database instance returned for %s, using the 1st one\n", dBName)
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
	return db.Exec(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s ( %s uuid PRIMARY KEY, %s timestamp, %s timestamp DEFAULT current_timestamp, %s json )", clustersTableName, clustersTableIDKey, clustersTableFirstSeenKey, clustersTableLastSeenKey, clustersTableDataKey))
}

func createVersionsTable(db *sql.DB) (sql.Result, error) {
	return db.Exec(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s ( %s varchar(64) PRIMARY KEY, %s timestamp, %s json )", versionsTableName, versionsTableComponentNameKey, versionsTableLastUpdatedKey, versionsTableDataKey))
}

func verifyClustersTable(db *sql.DB) error {
	if _, err := createClustersTable(db); err != nil {
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

func getDBRecord(db *sql.DB, table string, key string, val string) *sql.Row {
	return db.QueryRow(fmt.Sprintf("SELECT * FROM %s WHERE %s = '%s'", table, key, val))
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
	now := time.Now().Format(time.RFC3339)
	insert := fmt.Sprintf("INSERT INTO %s (cluster_id, first_seen, last_seen, data) VALUES('%s', '%s', '%s', '%s')", clustersTableName, id, now, now, string(data))
	return db.Exec(insert)
}

func newVersionDBRecord(db *sql.DB, component string, data []byte) (sql.Result, error) {
	now := time.Now().Format(time.RFC3339)
	insert := fmt.Sprintf("INSERT INTO %s (component_name, last_updated, data) VALUES('%s', '%s', '%s')", versionsTableName, component, now, string(data))
	return db.Exec(insert)
}

func updateClusterDBRecord(db *sql.DB, id string, data []byte) (sql.Result, error) {
	now := time.Now().Format(time.RFC3339)
	update := fmt.Sprintf("UPDATE %s SET data='%s', last_seen='%s' WHERE cluster_id='%s'", clustersTableName, string(data), now, id)
	return db.Exec(update)
}

func updateVersionDBRecord(db *sql.DB, component string, data []byte) (sql.Result, error) {
	now := time.Now().Format(time.RFC3339)
	update := fmt.Sprintf("UPDATE %s SET data='%s', last_updated='%s' WHERE component_name='%s'", versionsTableName, string(data), now, component)
	return db.Exec(update)
}
