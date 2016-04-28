package data

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx/types"
)

// ClustersCheckinsTable type that expresses the `clusters_checkins` postgres table schema
type ClustersCheckinsTable struct {
	checkinID string     // PRIMARY KEY, type uuid
	clusterID string     // indexed
	createdAt *Timestamp // indexed
	data      types.JSONText
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
