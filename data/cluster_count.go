package data

import (
	"database/sql"
	"fmt"
)

// GetClusterCount returns the total number of clusters in the database
func GetClusterCount(db *sql.DB) (int, error) {
	count, err := getTableCount(db, clustersTableName)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func getTableCount(db *sql.DB, table string) (int, error) {
	rows, err := db.Query(fmt.Sprintf("SELECT COUNT(*) FROM %s", table))
	if err != nil {
		return 0, err
	}
	var count int
	for rows.Next() {
		err := rows.Scan(&count)
		if err != nil {
			return 0, err
		}
	}
	return count, nil
}
