package data

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com/deis/workflow-manager/components"
	"github.com/deis/workflow-manager/types"
	"github.com/jinzhu/gorm"
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

// GetCluster gets the cluster from the DB with the given cluster ID
func GetCluster(db *gorm.DB, id string) (types.Cluster, error) {
	ret := &clustersTable{}
	resDB := db.Where(&clustersTable{ClusterID: id}).First(ret)
	if resDB.Error != nil {
		return types.Cluster{}, resDB.Error
	}
	cluster, err := components.ParseJSONCluster(ret.Data)
	if err != nil {
		return types.Cluster{}, errParsingCluster{origErr: err}
	}
	return cluster, nil
}

// CheckInAndSetCluster checks the cluster with the given ID in, and then updates it
func CheckInAndSetCluster(db *gorm.DB, id string, cluster types.Cluster) (types.Cluster, error) {
	// Check in
	if _, err := CheckInCluster(db.DB(), id, cluster); err != nil {
		return types.Cluster{}, err
	}
	var ret types.Cluster // return variable
	js, err := json.Marshal(cluster)
	if err != nil {
		fmt.Println("error marshaling data")
	}
	row := getDBRecord(db.DB(), clustersTableName, []string{clustersTableIDKey}, []string{id})
	var result sql.Result
	// Register the "latest checkin" with the primary cluster record
	rowResult := clustersTable{}
	if scanErr := row.Scan(&rowResult.ClusterID, &rowResult.Data); scanErr != nil {
		result, err = newClusterDBRecord(db.DB(), id, js)
		if err != nil {
			log.Println(err)
		}
	} else {
		result, err = updateClusterDBRecord(db.DB(), id, js)
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
		ret, err = GetCluster(db, id)
		if err != nil {
			return types.Cluster{}, err
		}
	} else if affected > 1 {
		log.Println("updated more than one record with same ID value!")
	}
	return ret, nil
}

// CheckInCluster creates a new record in the cluster checkins DB to indicate that the cluster has checked in right now
func CheckInCluster(db *sql.DB, id string, cluster types.Cluster) (sql.Result, error) {
	js, err := json.Marshal(cluster)
	if err != nil {
		fmt.Println("error marshaling data")
	}
	result, err := newClusterCheckinsDBRecord(db, id, now(), js)
	if err != nil {
		log.Println("cluster checkin db record not created", err)
		return nil, err
	}
	return result, nil
}

// FilterClustersByAge returns a slice of clusters whose various time fields match the requirements
// in the given filter. Note that the filter's requirements are a conjunction, not a disjunction
func FilterClustersByAge(db *sql.DB, filter *ClusterAgeFilter) ([]types.Cluster, error) {
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

	clusters := []types.Cluster{}
	for rows.Next() {
		rowResult := clustersTable{}
		if err := rows.Scan(&rowResult.ClusterID, &rowResult.Data); err != nil {
			return nil, err
		}
		cluster, err := components.ParseJSONCluster(rowResult.Data)
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
