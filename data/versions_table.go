package data

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx/types"
)

// VersionsTable type that expresses the `deis_component_versions` postgres table schema
type VersionsTable struct {
	versionID        string // PRIMARY KEY
	componentName    string // indexed
	train            string // indexed
	version          string // indexed
	releaseTimestamp *Timestamp
	data             types.JSONText
}

func createVersionsTable(db *sql.DB) (sql.Result, error) {
	query := fmt.Sprintf(
		"CREATE TABLE IF NOT EXISTS %s ( %s bigserial PRIMARY KEY, %s varchar(32), %s varchar(24), %s varchar(32), %s timestamp, %s json, unique (%s, %s, %s) )",
		versionsTableName,
		versionsTableIDKey,
		versionsTableComponentNameKey,
		versionsTableTrainKey,
		versionsTableVersionKey,
		versionsTableReleaseTimeStampKey,
		versionsTableDataKey,
		versionsTableComponentNameKey,
		versionsTableTrainKey,
		versionsTableVersionKey,
	)
	return db.Exec(query)
}

func verifyVersionsTable(db *sql.DB) error {
	if _, err := createVersionsTable(db); err != nil {
		log.Println("unable to verify versions table exists")
		return err
	}
	return nil
}
