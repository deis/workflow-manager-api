package data

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"database/sql"
)

const (
	dBInstanceKey                            = "WORKFLOW_MANAGER_API_DBINSTANCE"
	dBUserKey                                = "WORKFLOW_MANAGER_API_DBUSER"
	dBPassKey                                = "WORKFLOW_MANAGER_API_DBPASS"
	clustersTableName                        = "clusters"
	clustersTableIDKey                       = "cluster_id"
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
	dBInstance = os.Getenv(dBInstanceKey)
	dBUser     = os.Getenv(dBUserKey)
	dBPass     = os.Getenv(dBPassKey)
)

type errNoMoreRows struct {
	tableName string
}

func (e errNoMoreRows) Error() string {
	return fmt.Sprintf("no more rows available in the '%s' table", e.tableName)
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

func getOrderedDBRecords(db *sql.DB, table string, keys, vals []string, ordering *orderBy) (*sql.Rows, error) {
	sliceEqualize(&keys, &vals)
	query := fmt.Sprintf("SELECT * FROM %s", table)
	for i, key := range keys {
		if i == 0 {
			query += fmt.Sprintf(" WHERE %s = '%s'", key, vals[i])
		} else {
			query += fmt.Sprintf(" AND %s = '%s'", key, vals[i])
		}
	}
	if ordering != nil {
		query += fmt.Sprintf(" %s", ordering.String())
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
