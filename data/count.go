package data

import (
	"database/sql"
	"fmt"
)

// Count is an interface for managing a record count
type Count interface {
	Get(db *sql.DB) (int, error)
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

// GetClusterCount is a high level interface for retrieving a simple cluster count
func GetClusterCount(db *sql.DB, c Count) (int, error) {
	count, err := c.Get(db)
	if err != nil {
		return 0, err
	}
	return count, nil
}
