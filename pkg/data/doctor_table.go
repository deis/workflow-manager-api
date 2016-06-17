package data

import (
	"database/sql"
	"fmt"
	"log"
)

const (
	doctorTableName    = "doctors"
	doctorTableIDKey   = "report_id"
	doctorTableDataKey = "data"
)

// doctorTable type that expresses the `DcotorInfo` postgres table schema
type doctorTable struct {
	ReportID string `gorm:"primary_key;type:uuid;column:report_id"` // PRIMARY KEY
	Data     string `gorm:"type:json;column:data"`
}

//TableName return doctor table name
func (d doctorTable) TableName() string {
	return doctorTableName
}

func createDoctorTable(db *sql.DB) (sql.Result, error) {
	return db.Exec(fmt.Sprintf(
		"CREATE TABLE IF NOT EXISTS %s ( %s uuid PRIMARY KEY, %s json )",
		doctorTableName,
		doctorTableIDKey,
		doctorTableDataKey,
	))
}

func verifyDoctorTable(db *sql.DB) error {
	if _, err := createDoctorTable(db); err != nil {
		log.Println("unable to verify clusters table exists")
		return err
	}
	return nil
}
