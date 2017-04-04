package data

import (
	"fmt"
	"log"
	"strconv"

	"github.com/deis/workflow-manager-api/config"
	"github.com/jinzhu/gorm"
)

var (
	dBUser = config.Spec.DBUser
	dBPass = config.Spec.DBPass
	dBURL  = config.Spec.DBURL
	dBName = config.Spec.DBName
)

type errNoMoreRows struct {
	tableName string
}

func (e errNoMoreRows) Error() string {
	return fmt.Sprintf("no more rows available in the '%s' table", e.tableName)
}

// VerifyPersistentStorage is a high level interace for verifying storage abstractions
func VerifyPersistentStorage(db *gorm.DB) error {
	if _, err := createOrUpdateVersionsTable(db); err != nil {
		log.Println("unable to verify " + versionsTableName + " table")
		return err
	}
	count, err := getTableCount(db.DB(), versionsTableName)
	if err != nil {
		log.Println("unable to get record count for " + versionsTableName + " table")
		return err
	}
	log.Println("counted " + strconv.Itoa(count) + " records for " + versionsTableName + " table")
	err = verifyClustersTable(db.DB())
	if err != nil {
		log.Println("unable to verify " + clustersTableName + " table")
		return err
	}
	count, err = getTableCount(db.DB(), clustersTableName)
	if err != nil {
		log.Println("unable to get record count for " + clustersTableName + " table")
		return err
	}
	log.Println("counted " + strconv.Itoa(count) + " records for " + clustersTableName + " table")
	err = verifyClustersCheckinsTable(db.DB())
	if err != nil {
		log.Println("unable to verify " + clustersCheckinsTableName + " table")
		return err
	}
	count, err = getTableCount(db.DB(), clustersCheckinsTableName)
	if err != nil {
		log.Println("unable to get record count for " + clustersCheckinsTableName + " table")
		return err
	}
	err = verifyDoctorTable(db.DB())
	if err != nil {
		log.Println("unable to verify " + doctorTableName + " table")
		return err
	}
	count, err = getTableCount(db.DB(), doctorTableName)
	if err != nil {
		log.Println("unable to get record count for " + doctorTableName + " table")
		return err
	}
	log.Println("counted " + strconv.Itoa(count) + " records for " + clustersCheckinsTableName + " table")
	return nil
}
