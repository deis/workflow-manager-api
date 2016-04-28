package data

import (
	"database/sql"
)

// ClusterCount fulfills the Count interface
type ClusterCount struct{}

// Get method for ClusterCount
func (c ClusterCount) Get(db *sql.DB) (int, error) {
	count, err := getTableCount(db, clustersTableName)
	if err != nil {
		return 0, err
	}
	return count, nil
}
