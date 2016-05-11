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

func upsertCluster(db *gorm.DB, id string, cluster ClusterStateful) (ClusterStateful, error) {
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
		return ClusterStateful{}, countDB.Error
	}
	var resDB *gorm.DB
	if numExisting == 0 {
		// no existing clusters, so create one
		createDB := db.Create(&clustersTable{ClusterID: id, Data: js})
		if createDB.Error != nil {
			return ClusterStateful{}, createDB.Error
		}
		resDB = createDB
	} else {
		updateDB := db.Save(&clustersTable{ClusterID: id, Data: js})
		if updateDB.Error != nil {
			return ClusterStateful{}, updateDB.Error
		}
		resDB = updateDB
	}
	if resDB.RowsAffected != 1 {
		return ClusterStateful{}, fmt.Errorf("%d rows were affected, but expected only 1", resDB.RowsAffected)
	}
	retCluster, err := GetCluster(db, id)
	if err != nil {
		return ClusterStateful{}, err
	}
	return retCluster, nil
}

// CheckInAndSetCluster creates or updates the cluster with the given ID.
// TODO: rename this function to better reflect what it does (https://github.com/deis/workflow-manager-api/issues/128)
func CheckInAndSetCluster(db *gorm.DB, id string, cluster ClusterStateful) (ClusterStateful, error) {
	txn := db.Begin()
	if txn.Error != nil {
		return ClusterStateful{}, txErr{orig: nil, err: txn.Error, op: "begin"}
	}
	ret, err := upsertCluster(txn, id, cluster)
	if err != nil {
		rbDB := txn.Rollback()
		if rbDB.Error != nil {
			return ClusterStateful{}, txErr{orig: err, err: rbDB.Error, op: "rollback"}
		}
		return ClusterStateful{}, err
	}
	comDB := txn.Commit()
	if comDB.Error != nil {
		return ClusterStateful{}, txErr{orig: nil, err: comDB.Error, op: "commit"}
	}
	return ret, nil
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
func FilterClustersByAge(db *gorm.DB, filter *ClusterAgeFilter) ([]ClusterStateful, error) {
	var rows []clustersTable
	execDB := db.Raw(`SELECT clusters.*
		FROM clusters, clusters_checkins
		WHERE clusters_checkins.cluster_id = clusters.cluster_id
		GROUP BY clusters_checkins.cluster_id, clusters.cluster_id
		HAVING MIN(clusters_checkins.created_at) > ?
		AND MIN(clusters_checkins.created_at) < ?
		AND MIN(clusters_checkins.created_at) > ?
		AND MAX(clusters_checkins.created_at) < ?`,
		Timestamp{Time: filter.CreatedAfter},
		Timestamp{Time: filter.CreatedBefore},
		Timestamp{Time: filter.CheckedInAfter},
		Timestamp{Time: filter.CheckedInBefore},
	).Find(&rows)
	if execDB.Error != nil {
		return nil, execDB.Error
	}

	clusters := make([]ClusterStateful, len(rows))
	for i, row := range rows {
		cluster, err := parseJSONCluster(row.Data)
		if err != nil {
			return nil, err
		}
		clusters[i] = cluster
	}
	return clusters, nil
}
