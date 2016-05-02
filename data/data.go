package data

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/jinzhu/gorm"
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
func VerifyPersistentStorage(db *gorm.DB) error {
	if _, err := createOrUpdateVersionsTable(db); err != nil {
		log.Println("unable to verify " + versionsTableName + " table")
		return err
	}
	count, err := getTableCount(db.DB(), versionsTableName)
	if err != nil {
		log.Println("unable to get record count for " + versionsTableName + " table")
		return err
	}
	log.Println("counted " + strconv.Itoa(count) + " records for " + versionsTableName + " table")
	err = verifyClustersTable(db.DB())
	if err != nil {
		log.Println("unable to verify " + clustersTableName + " table")
		return err
	}
	count, err = getTableCount(db.DB(), clustersTableName)
	if err != nil {
		log.Println("unable to get record count for " + clustersTableName + " table")
		return err
	}
	log.Println("counted " + strconv.Itoa(count) + " records for " + clustersTableName + " table")
	err = verifyClustersCheckinsTable(db.DB())
	if err != nil {
		log.Println("unable to verify " + clustersCheckinsTableName + " table")
		return err
	}
	count, err = getTableCount(db.DB(), clustersCheckinsTableName)
	if err != nil {
		log.Println("unable to get record count for " + clustersCheckinsTableName + " table")
		return err
	}
	log.Println("counted " + strconv.Itoa(count) + " records for " + clustersCheckinsTableName + " table")
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
