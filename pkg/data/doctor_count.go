package data

import "github.com/jinzhu/gorm"

// GetDoctorCount returns the total number of doctor reports in the database
func GetDoctorCount(db *gorm.DB) (int, error) {
	count := 0
	countDB := db.Model(&doctorTable{}).Count(&count)
	if countDB.Error != nil {
		return 0, countDB.Error
	}
	return count, nil
}
