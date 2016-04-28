package data

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx/types"
)

// ClustersTable type that expresses the `clusters` postgres table schema
type ClustersTable struct {
	clusterID string // PRIMARY KEY
	data      types.JSONText
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
