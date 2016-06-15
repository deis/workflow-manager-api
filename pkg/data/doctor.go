package data

import (
	"github.com/deis/workflow-manager-api/pkg/swagger/models"
	"github.com/jinzhu/gorm"
)

// UpsertVersion adds or updates a single version record in the database
func AddDoctroInfo(db *gorm.DB, doctorInfo models.DoctorInfo) (models.DoctorInfo, error) {
	//toDo: add entry to db table
	return nil, nil
}
