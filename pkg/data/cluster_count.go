package data

import (
	"database/sql"
	"fmt"

	"github.com/jinzhu/gorm"
)

// GetClusterCount returns the total number of clusters in the database
func GetClusterCount(db *gorm.DB) (int, error) {
	count := 0
	countDB := db.Model(&clustersTable{}).Count(&count)
	if countDB.Error != nil {
		return 0, countDB.Error
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
