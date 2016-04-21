package data

import (
	"bytes"
	"encoding/json"
	"errors"

	"github.com/deis/workflow-manager/components"
	"github.com/deis/workflow-manager/types"
	"github.com/jinzhu/gorm"
	stypes "github.com/jmoiron/sqlx/types"
)

var (
	errNoRowsAffected = errors.New("no rows affected")
)

// ClusterTable expresses the 'clusters' table schema. It's also the Gorm model
type ClusterTable struct {
	ID   string          `gorm:"primary_key;type:uuid;column:cluster_id"` // PRIMARY KEY
	Data stypes.JSONText `gorm:"type:json;column:data"`
}

// TableName overrides the gorm-inferred table name
func (c ClustersTable) TableName() string {
	return "clusters"
}

func getClusterTableByID(db *gorm.DB, id string) (*ClusterTable, error) {
	cl := new(ClusterTable)
	byID := db.Where(&ClusterTable{ID: id}).First(cl)
	if byID.Error != nil {
		return nil, byID.Error
	}
	return cl, nil
}

// GetCluster gets the cluster with the specified ID
func GetCluster(db *gorm.DB, id string) (*types.Cluster, error) {
	cl, err := getClusterTableByID(db, id)
	if err != nil {
		return nil, err
	}
	cluster, err := components.ParseJSONCluster(cl.Data)
	if err != nil {
		return nil, err
	}

	return &cluster, nil
}

// CheckinAndSetCluster sets the new data on the cluster in the db, and checks it in
// TODO: implement checkin logic
func CheckinAndSetCluster(db *gorm.DB, id string, cluster *types.Cluster) (*types.Cluster, error) {
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(cluster); err != nil {
		return nil, err
	}
	data := stypes.JSONText(buf.Bytes())

	tx := db.Begin()
	clusterTable, err := getClusterTableByID(db, id)
	if err != nil {
		// assume the cluster was missing so create it
		created := db.Create(&ClusterTable{ID: cluster.ID, Data: data})
		if created.Error != nil {
			tx.Rollback()
			return nil, created.Error
		}
		if created.RowsAffected < 1 {
			tx.Rollback()
			return nil, errNoRowsAffected
		}
	} else {
		// assume the cluster was not missing so update it
		clusterTable.ID = cluster.ID
		clusterTable.Data = data
		updated := tx.Save(clusterTable)
		if updated.Error != nil {
			tx.Rollback()
			return nil, updated.Error
		}
		if updated.RowsAffected < 1 {
			return nil, errNoRowsAffected
		}
	}
	committed := tx.Commit()
	if committed.Error != nil {
		return nil, committed.Error
	}
	if committed.RowsAffected < 1 {
		return nil, errNoRowsAffected
	}
	return cluster, nil
}
