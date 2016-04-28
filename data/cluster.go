package data

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com/deis/workflow-manager/components"
	"github.com/deis/workflow-manager/types"
)

// Cluster is an interface for managing a persistent cluster record
type Cluster interface {
	Set(*sql.DB, string, types.Cluster) (types.Cluster, error)
}

func updateClusterDBRecord(db *sql.DB, id string, data []byte) (sql.Result, error) {
	update := fmt.Sprintf("UPDATE %s SET data='%s' WHERE cluster_id='%s'", clustersTableName, string(data), id)
	return db.Exec(update)
}

// GetCluster gets the cluster from the DB with the given cluster ID
func GetCluster(db *sql.DB, id string) (types.Cluster, error) {
	row := getDBRecord(db, clustersTableName, []string{clustersTableIDKey}, []string{id})
	rowResult := ClustersTable{}
	if err := row.Scan(&rowResult.clusterID, &rowResult.data); err != nil {
		return types.Cluster{}, err
	}
	cluster, err := components.ParseJSONCluster(rowResult.data)
	if err != nil {
		log.Println("error parsing cluster")
		return types.Cluster{}, err
	}
	return cluster, nil
}

// SetCluster is a high level interface for updating a cluster data record
func SetCluster(id string, cluster types.Cluster, db *sql.DB, c Cluster) (types.Cluster, error) {
	// Check in
	if _, err := CheckInCluster(db, id, cluster); err != nil {
		return types.Cluster{}, err
	}
	// Update cluster record
	ret, err := c.Set(db, id, cluster)
	if err != nil {
		return types.Cluster{}, err
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
	query := fmt.Sprintf(`SELECT DISTINCT clusters.*
		FROM clusters, clusters_checkins
		WHERE clusters_checkins.cluster_id = clusters.cluster_id
		GROUP BY clusters_checkins.cluster_id
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
		rowResult := ClustersTable{}
		if err := rows.Scan(&rowResult.clusterID, &rowResult.data); err != nil {
			return nil, err
		}
		cluster, err := components.ParseJSONCluster(rowResult.data)
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
