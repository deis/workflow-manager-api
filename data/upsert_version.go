package data

import (
	"github.com/deis/workflow-manager/types"
	"github.com/jinzhu/gorm"
)

func upsertVersion(db *gorm.DB, queryExisting versionsTable, setNew versionsTable) (*types.ComponentVersion, error) {
	var count int
	countDB := db.Where(&queryExisting).Count(&count)
	if countDB.Error != nil {
		return nil, countDB.Error
	}
	if count == 0 {
		createDB := db.Create(&setNew)
		if createDB.Error != nil {
			return nil, createDB.Error
		}
	} else {
		saveDB := db.Save(&setNew)
		if saveDB.Error != nil {
			return nil, saveDB.Error
		}
	}
	var ret versionsTable
	queryDB := db.Where(&queryExisting).First(&ret)
	if queryDB.Error != nil {
		return nil, queryDB.Error
	}
	cv, err := parseDBVersion(ret)
	if err != nil {
		return nil, err
	}
	return &cv, nil
}
