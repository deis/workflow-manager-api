package data

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/deis/workflow-manager/types"
	"github.com/jinzhu/gorm"
)

var (
	errNoRowsAffected = errors.New("No rows affected")
)

func updateClusterDBRecord(db *sql.DB, id string, data []byte) (sql.Result, error) {
	update := fmt.Sprintf("UPDATE %s SET data='%s' WHERE cluster_id='%s'", clustersTableName, string(data), id)
	return db.Exec(update)
}

type errParsingCluster struct {
	origErr error
}

func (e errParsingCluster) Error() string {
	return fmt.Sprintf("Error parsing cluster (%s)", e.origErr)
}

// ClusterStateful definition
// This is a wrapper around a cluster object to include properties for use in stateful contexts
type ClusterStateful struct {
	// FirstSeen and/or LastSeen suggests a Cluster object in a lifecycle context,
	// i.e., for use in business logic which needs to determine a cluster's "freshness" or "staleness"
	FirstSeen time.Time `json:"firstSeen"`
	LastSeen  time.Time `json:"lastSeen"`
	types.Cluster
}

// GetCluster gets the cluster from the DB with the given cluster ID
func GetCluster(db *gorm.DB, id string) (ClusterStateful, error) {
	ret := &clustersTable{}
	resDB := db.Where(&clustersTable{ClusterID: id}).First(ret)
	if resDB.Error != nil {
		return ClusterStateful{}, resDB.Error
	}
	cluster, err := parseJSONCluster(ret.Data)
	if err != nil {
		return ClusterStateful{}, errParsingCluster{origErr: err}
	}
	return cluster, nil
}

// CheckInAndSetCluster checks the cluster with the given ID in, and then updates it
func CheckInAndSetCluster(db *gorm.DB, id string, cluster ClusterStateful) (ClusterStateful, error) {
	// Check in
	if err := CheckInCluster(db, id, time.Now(), cluster); err != nil {
		return ClusterStateful{}, err
	}
	js, err := json.Marshal(cluster)
	if err != nil {
		return ClusterStateful{}, err
	}
	var numExisting int
	query := clustersTable{ClusterID: id}
	countDB := db.Model(&clustersTable{}).Where(&query).Count(&numExisting)
	if countDB.Error != nil {
		log.Printf("ERROR1 (%s)", countDB.Error)
		return ClusterStateful{}, countDB.Error
	}
	var resDB *gorm.DB
	if numExisting == 0 {
		// no existing clusters, so create one
		createDB := db.Create(&clustersTable{ClusterID: id, Data: js})
		if createDB.Error != nil {
			log.Printf("ERROR2 (%s)", createDB.Error)
			return ClusterStateful{}, createDB.Error
		}
		resDB = createDB
	} else {
		updateDB := db.Save(&clustersTable{ClusterID: id, Data: js})
		if updateDB.Error != nil {
			log.Printf("ERROR3 (%s)", updateDB.Error)
			return ClusterStateful{}, updateDB.Error
		}
		resDB = updateDB
	}
	if resDB.RowsAffected != 1 {
		return ClusterStateful{}, fmt.Errorf("%d rows were affected, but expected only 1", resDB.RowsAffected)
	}
	retCluster, err := GetCluster(db, id)
	if err != nil {
		log.Printf("ERROR4 (%s)", err)
		return ClusterStateful{}, err
	}
	return retCluster, nil
}

// CheckInCluster creates a new record in the cluster checkins DB to indicate that the cluster has checked in right now
func CheckInCluster(db *gorm.DB, id string, checkinTime time.Time, cluster ClusterStateful) error {
	js, err := json.Marshal(cluster)
	if err != nil {
		fmt.Println("error marshaling data")
	}
	record := newClustersCheckinsTable("", id, checkinTime, js)
	createdDB := db.Create(&record)
	if createdDB.Error != nil {
		log.Println("cluster checkin db record not created", createdDB.Error)
		return createdDB.Error
	}
	if createdDB.RowsAffected != 1 {
		return errNoRowsAffected
	}
	return nil
}

// FilterClustersByAge returns a slice of clusters whose various time fields match the requirements
// in the given filter. Note that the filter's requirements are a conjunction, not a disjunction
func FilterClustersByAge(db *sql.DB, filter *ClusterAgeFilter) ([]ClusterStateful, error) {
	query := fmt.Sprintf(`SELECT clusters.*
		FROM clusters, clusters_checkins
		WHERE clusters_checkins.cluster_id = clusters.cluster_id
		GROUP BY clusters_checkins.cluster_id, clusters.cluster_id
		HAVING MIN(clusters_checkins.created_at) > '%s'
		AND MIN(clusters_checkins.created_at) < '%s'
		AND MIN(clusters_checkins.created_at) > '%s'
		AND MAX(clusters_checkins.created_at) < '%s'`,
		filter.createdAfterTimestamp(),
		filter.createdBeforeTimestamp(),
		filter.checkedInAfterTimestamp(),
		filter.checkedInBeforeTimestamp(),
	)

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}

	clusters := []ClusterStateful{}
	for rows.Next() {
		rowResult := clustersTable{}
		if err := rows.Scan(&rowResult.ClusterID, &rowResult.Data); err != nil {
			return nil, err
		}
		cluster, err := parseJSONCluster(rowResult.Data)
		if err != nil {
			return nil, err
		}
		clusters = append(clusters, cluster)
	}
	return clusters, nil
}

func newClusterDBRecord(db *sql.DB, id string, data []byte) (sql.Result, error) {
	insert := fmt.Sprintf("INSERT INTO %s (cluster_id, data) VALUES('%s', '%s')", clustersTableName, id, string(data))
	return db.Exec(insert)
}
