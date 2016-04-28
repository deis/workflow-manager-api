package data

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/deis/workflow-manager/components"
	"github.com/deis/workflow-manager/types"
)

// Cluster is an interface for managing a persistent cluster record
type Cluster interface {
	Set(*sql.DB, string, types.Cluster) (types.Cluster, error)
	Checkin(*sql.DB, string, types.Cluster) (sql.Result, error)
	FilterByAge(*sql.DB, *ClusterAgeFilter) ([]types.Cluster, error)
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

func newClusterDBRecord(db *sql.DB, id string, data []byte) (sql.Result, error) {
	insert := fmt.Sprintf("INSERT INTO %s (cluster_id, data) VALUES('%s', '%s')", clustersTableName, id, string(data))
	return db.Exec(insert)
}
