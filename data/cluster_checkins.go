package data

import (
	"database/sql"
	"fmt"
)

func newClusterCheckinsDBRecord(db *sql.DB, id string, createdAt *Timestamp, data []byte) (sql.Result, error) {
	update := fmt.Sprintf(
		"INSERT INTO %s (data, created_at, cluster_id) VALUES('%s', '%s', '%s')",
		clustersCheckinsTableName,
		string(data),
		createdAt.String(),
		id,
	)
	return db.Exec(update)
}
