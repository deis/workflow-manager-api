package data

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx/types"
)

const (
	clustersCheckinsTableName                = "clusters_checkins"
	clustersCheckinsTableIDKey               = "checkins_id"
	clustersCheckinsTableClusterIDKey        = "cluster_id"
	clustersCheckinsTableClusterCreatedAtKey = "created_at"
	clustersCheckinsTableDataKey             = "data"
)

// ClustersCheckinsTable type that expresses the `clusters_checkins` postgres table schema
type clustersCheckinsTable struct {
	CheckinID string         `gorm:"primary_key;type:bigserial;column_name:checkins_id"`
	ClusterID string         `gorm:"type:uuid;column_name:cluster_id;index"`
	CreatedAt *Timestamp     `gorm:"type:timestamp;column_name:created_at;index"`
	Data      types.JSONText `gorm:"type:json;column_name:data"`
}

func (c clustersCheckinsTable) TableName() string {
	return clustersCheckinsTableName
}

func createClustersCheckinsTable(db *sql.DB) (sql.Result, error) {
	return db.Exec(fmt.Sprintf(
		"CREATE TABLE IF NOT EXISTS %s ( %s bigserial PRIMARY KEY, %s uuid, %s timestamp, %s json )",
		clustersCheckinsTableName,
		clustersCheckinsTableIDKey,
		clustersTableIDKey,
		clustersCheckinsTableClusterCreatedAtKey,
		clustersCheckinsTableDataKey,
	))
}

func verifyClustersCheckinsTable(db *sql.DB) error {
	if _, err := createClustersCheckinsTable(db); err != nil {
		log.Println("unable to verify clusters table exists")
		return err
	}
	return nil
}
