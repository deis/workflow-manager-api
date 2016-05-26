package data

import (
	"database/sql"
	"fmt"
	"log"
)

const (
	clustersTableName    = "clusters"
	clustersTableIDKey   = "cluster_id"
	clustersTableDataKey = "data"
)

// ClustersTable type that expresses the `clusters` postgres table schema
type clustersTable struct {
	ClusterID string `gorm:"primary_key;type:uuid;column:cluster_id"` // PRIMARY KEY
	Data      string `gorm:"type:json;column:data"`
}

func (c clustersTable) TableName() string {
	return clustersTableName
}

func createClustersTable(db *sql.DB) (sql.Result, error) {
	return db.Exec(fmt.Sprintf(
		"CREATE TABLE IF NOT EXISTS %s ( %s uuid PRIMARY KEY, %s json )",
		clustersTableName,
		clustersTableIDKey,
		clustersTableDataKey,
	))
}

func verifyClustersTable(db *sql.DB) error {
	if _, err := createClustersTable(db); err != nil {
		log.Println("unable to verify clusters table exists")
		return err
	}
	return nil
}
