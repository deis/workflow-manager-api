package data

import (
	"encoding/json"
	"fmt"

	"github.com/deis/workflow-manager-api/pkg/swagger/models"
	"github.com/jinzhu/gorm"
)

// GetDoctor gets the doctorInfo from the DB with the given report ID
func GetDoctor(db *gorm.DB, id string) (models.DoctorInfo, error) {
	ret := &doctorTable{}
	resDB := db.Where(&doctorTable{ReportID: id}).First(ret)
	if resDB.Error != nil {
		return models.DoctorInfo{}, resDB.Error
	}
	doctor, err := parseJSONDoctor([]byte(ret.Data))
	if err != nil {
		return models.DoctorInfo{}, errParsingCluster{origErr: err}
	}
	return doctor, nil
}

func upsertDoctor(db *gorm.DB, id string, doctor models.DoctorInfo) (models.DoctorInfo, error) {
	js, err := json.Marshal(doctor)
	if err != nil {
		return models.DoctorInfo{}, err
	}
	var numExisting int
	query := doctorTable{ReportID: id}
	countDB := db.Model(&doctorTable{}).Where(&query).Count(&numExisting)
	if countDB.Error != nil {
		return models.DoctorInfo{}, countDB.Error
	}
	var resDB *gorm.DB
	if numExisting == 0 {
		// no existing clusters, so create one
		createDB := db.Create(&doctorTable{ReportID: id, Data: string(js)})
		if createDB.Error != nil {
			return models.DoctorInfo{}, createDB.Error
		}
		resDB = createDB
	} else {
		updateDB := db.Save(&doctorTable{ReportID: id, Data: string(js)})
		if updateDB.Error != nil {
			return models.DoctorInfo{}, updateDB.Error
		}
		resDB = updateDB
	}
	if resDB.RowsAffected != 1 {
		return models.DoctorInfo{}, fmt.Errorf("%d rows were affected, but expected only 1", resDB.RowsAffected)
	}
	retDoctor, err := GetDoctor(db, id)
	if err != nil {
		return models.DoctorInfo{}, err
	}
	return retDoctor, nil
}

// UpsertDoctor creates or updates the doctor with the given ID.
func UpsertDoctor(db *gorm.DB, id string, doctor models.DoctorInfo) (models.DoctorInfo, error) {
	txn := db.Begin()
	if txn.Error != nil {
		return models.DoctorInfo{}, txErr{orig: nil, err: txn.Error, op: "begin"}
	}
	ret, err := upsertDoctor(txn, id, doctor)
	if err != nil {
		rbDB := txn.Rollback()
		if rbDB.Error != nil {
			return models.DoctorInfo{}, txErr{orig: err, err: rbDB.Error, op: "rollback"}
		}
		return models.DoctorInfo{}, err
	}
	comDB := txn.Commit()
	if comDB.Error != nil {
		return models.DoctorInfo{}, txErr{orig: nil, err: comDB.Error, op: "commit"}
	}
	return ret, nil
}
