package data

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/deis/workflow-manager-api/pkg/swagger/models"
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

// GetCluster gets the cluster from the DB with the given cluster ID
func GetCluster(db *gorm.DB, id string) (models.Cluster, error) {
	ret := &clustersTable{}
	resDB := db.Where(&clustersTable{ClusterID: id}).First(ret)
	if resDB.Error != nil {
		return models.Cluster{}, resDB.Error
	}
	cluster, err := parseJSONCluster([]byte(ret.Data))
	if err != nil {
		return models.Cluster{}, errParsingCluster{origErr: err}
	}
	return cluster, nil
}

func upsertCluster(db *gorm.DB, id string, cluster models.Cluster) (models.Cluster, error) {
	// Check in
	if err := CheckInCluster(db, id, time.Now(), cluster); err != nil {
		return models.Cluster{}, err
	}
	js, err := json.Marshal(cluster)
	if err != nil {
		return models.Cluster{}, err
	}
	var numExisting int
	query := clustersTable{ClusterID: id}
	countDB := db.Model(&clustersTable{}).Where(&query).Count(&numExisting)
	if countDB.Error != nil {
		return models.Cluster{}, countDB.Error
	}
	var resDB *gorm.DB
	if numExisting == 0 {
		// no existing clusters, so create one
		createDB := db.Create(&clustersTable{ClusterID: id, Data: string(js)})
		if createDB.Error != nil {
			return models.Cluster{}, createDB.Error
		}
		resDB = createDB
	} else {
		updateDB := db.Save(&clustersTable{ClusterID: id, Data: string(js)})
		if updateDB.Error != nil {
			return models.Cluster{}, updateDB.Error
		}
		resDB = updateDB
	}
	if resDB.RowsAffected != 1 {
		return models.Cluster{}, fmt.Errorf("%d rows were affected, but expected only 1", resDB.RowsAffected)
	}
	retCluster, err := GetCluster(db, id)
	if err != nil {
		return models.Cluster{}, err
	}
	return retCluster, nil
}

// UpsertCluster creates or updates the cluster with the given ID.
func UpsertCluster(db *gorm.DB, id string, cluster models.Cluster) (models.Cluster, error) {
	txn := db.Begin()
	if txn.Error != nil {
		return models.Cluster{}, txErr{orig: nil, err: txn.Error, op: "begin"}
	}
	ret, err := upsertCluster(txn, id, cluster)
	if err != nil {
		rbDB := txn.Rollback()
		if rbDB.Error != nil {
			return models.Cluster{}, txErr{orig: err, err: rbDB.Error, op: "rollback"}
		}
		return models.Cluster{}, err
	}
	comDB := txn.Commit()
	if comDB.Error != nil {
		return models.Cluster{}, txErr{orig: nil, err: comDB.Error, op: "commit"}
	}
	return ret, nil
}

// CheckInCluster creates a new record in the cluster checkins DB to indicate that the cluster has checked in right now
func CheckInCluster(db *gorm.DB, id string, checkinTime time.Time, cluster models.Cluster) error {
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
func FilterClustersByAge(db *gorm.DB, filter *ClusterAgeFilter) ([]*models.Cluster, error) {
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

	clusters, err := makeClusters(rows)
	if err != nil {
		return nil, err
	}
	return clusters, nil
}

// FilterClusterCheckins returns a slice of clusters whose various time fields match the requirements
// in the given filter.
func FilterClusterCheckins(db *gorm.DB, filter *ClusterCheckinsFilter) ([]*models.ClusterCheckin, error) {
	var rows []clustersCheckinsFilterResponse
	execDB := db.Raw(`SELECT cluster_id,
		MIN(created_at) AS first_seen,
		MAX(created_at) AS last_seen,
		AGE(MAX(created_at), MIN(created_at)) AS cluster_age,
		AGE(MAX(created_at), NOW()) AS last_checkin,
		COUNT(1) AS checkins
		FROM   clusters_checkins
		GROUP  BY cluster_id
		HAVING MIN(created_at) > ? AND MIN(created_at) < ?
		ORDER  BY first_seen DESC;`,
		Timestamp{Time: filter.CreatedAfter},
		Timestamp{Time: filter.CreatedBefore},
	).Find(&rows)

	if execDB.Error != nil {
		return nil, execDB.Error
	}

	checkins, err := makeClusterCheckins(rows)
	if err != nil {
		return nil, err
	}
	return checkins, nil
}

// FilterPersistentClusters returns a slice of clusters whose various time fields match the requirements
// in the given filter.
func FilterPersistentClusters(db *gorm.DB, filter *PersistentClustersFilter) ([]*models.ClusterCheckin, error) {
	var rows []clustersCheckinsFilterResponse
	execDB := db.Raw(`SELECT cluster_id,
        MIN(created_at) AS first_seen,
        MAX(created_at) AS last_seen,
        AGE(MAX(created_at), MIN(created_at)) AS cluster_age,
        COUNT(1) AS checkins
        FROM clusters_checkins
        GROUP  BY cluster_id
        HAVING MIN(created_at) > ? AND MIN(created_at) < ?
        AND COUNT(1) > 1 AND MAX(created_at) > ?
		ORDER  BY first_seen ASC`,
		Timestamp{Time: filter.Epoch},
		Timestamp{Time: filter.Timestamp},
		Timestamp{Time: filter.RelativeYesterday},
	).Find(&rows)
	if execDB.Error != nil {
		return nil, execDB.Error
	}

	checkins, err := makeClusterCheckins(rows)
	if err != nil {
		return nil, err
	}
	return checkins, nil
}

func makeClusters(rows []clustersTable) ([]*models.Cluster, error) {
	clusters := make([]*models.Cluster, len(rows))
	for i, row := range rows {
		cluster, err := parseJSONCluster([]byte(row.Data))
		if err != nil {
			return nil, err
		}
		clusters[i] = &cluster
	}
	return clusters, nil
}

func makeClusterCheckins(rows []clustersCheckinsFilterResponse) ([]*models.ClusterCheckin, error) {
	clusterCheckins := make([]*models.ClusterCheckin, len(rows))
	for i, row := range rows {
		c, err := strconv.ParseInt(row.Checkins, 10, 64)
		if err != nil {
			return nil, err
		}
		clusterCheckin := models.ClusterCheckin{
			ClusterID:  row.ClusterID,
			FirstSeen:  row.FirstSeen,
			LastSeen:   row.LastSeen,
			ClusterAge: row.ClusterAge,
			Checkins:   c,
		}
		if row.LastCheckin != "" {
			clusterCheckin.LastCheckin = row.LastCheckin
		}

		clusterCheckins[i] = &clusterCheckin
	}
	return clusterCheckins, nil
}
