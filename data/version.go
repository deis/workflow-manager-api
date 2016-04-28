package data

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com/deis/workflow-manager/types"
)

// Version is an interface for managing a persistent cluster record
type Version interface {
	// Retrieve a list of Version records that match a given component + train
	Collection(db *sql.DB, train string, component string) ([]types.ComponentVersion, error)
	// Retrieve the most recent Version record that matches a given component + train
	Latest(db *sql.DB, train string, component string) (types.ComponentVersion, error)
	// MultiLatest fetches from the DB and returns the latest release for each component/train pair
	// given in ct. Returns an empty slice and non-nil error on any error communicating with the
	// database or otherwise if the first returned value is not empty, it's guaranteed to:
	//
	// - Be the same length as ct
	// - Have the same ordering as ct, with respect to the component name
	MultiLatest(db *sql.DB, ct []ComponentAndTrain) ([]types.ComponentVersion, error)
}

func updateVersionDBRecord(db *sql.DB, componentVersion types.ComponentVersion) (sql.Result, error) {
	data, err := json.Marshal(componentVersion.Version.Data)
	if err != nil {
		log.Printf("JSON marshaling failed (%s)", err)
		return nil, err
	}
	update := fmt.Sprintf(
		"UPDATE %s SET %s='%s', %s='%s', %s='%s', %s='%s', %s='%s' WHERE %s='%s' AND %s='%s' AND %s='%s'",
		versionsTableName,
		versionsTableComponentNameKey,
		componentVersion.Component.Name,
		versionsTableTrainKey,
		componentVersion.Version.Train,
		versionsTableVersionKey,
		componentVersion.Version.Version,
		versionsTableReleaseTimeStampKey,
		componentVersion.Version.Released,
		versionsTableDataKey,
		string(data),
		versionsTableComponentNameKey,
		componentVersion.Component.Name,
		versionsTableTrainKey,
		componentVersion.Version.Train,
		versionsTableVersionKey,
		componentVersion.Version.Version,
	)
	return db.Exec(update)
}

// GetComponentTrainVersions is a high level interface for retrieving component versions for a given "train"
func GetComponentTrainVersions(train string, component string, db *sql.DB, v Version) ([]types.ComponentVersion, error) {
	componentVersions, err := v.Collection(db, train, component)
	if err != nil {
		return nil, err
	}
	return componentVersions, nil
}

// SetVersion adds or updates a single version record in the database
func SetVersion(db *sql.DB, componentVersion types.ComponentVersion) (types.ComponentVersion, error) {
	// TODO: this read-modify-write should be done inside a transaction. Also, rename SetVersion to something else.
	// Both of these TODOs are captured in https://github.com/deis/workflow-manager-api/issues/90

	var ret types.ComponentVersion // return variable
	row := getDBRecord(db, versionsTableName,
		[]string{versionsTableComponentNameKey, versionsTableTrainKey, versionsTableVersionKey},
		[]string{componentVersion.Component.Name, componentVersion.Version.Train, componentVersion.Version.Version})
	var result sql.Result
	rowResult := VersionsTable{}
	if err := row.Scan(&rowResult.versionID, &rowResult.componentName, &rowResult.train, &rowResult.version, &rowResult.releaseTimestamp, &rowResult.data); err != nil {
		result, err = newVersionDBRecord(db, componentVersion)
		if err != nil {
			log.Println(err)
			return types.ComponentVersion{}, err
		}
	} else {
		result, err = updateVersionDBRecord(db, componentVersion)
		if err != nil {
			log.Println(err)
			return types.ComponentVersion{}, err
		}
	}
	affected, err := result.RowsAffected()
	if err != nil {
		log.Println("failed to get affected row count")
	}
	if affected == 0 {
		log.Println("no records updated")
	} else if affected == 1 {
		ret, err = GetVersion(db, componentVersion)
		if err != nil {
			return types.ComponentVersion{}, err
		}
	} else if affected > 1 {
		log.Println("updated more than one record with same ID value!")
	}
	return ret, nil
}

// GetLatestComponentTrainVersion is a high level interface for retrieving the latest component version for a given "train"
func GetLatestComponentTrainVersion(train string, component string, db *sql.DB, v Version) (types.ComponentVersion, error) {
	componentVersion, err := v.Latest(db, train, component)
	if err != nil {
		return types.ComponentVersion{}, err
	}
	return componentVersion, nil
}

func newVersionDBRecord(db *sql.DB, componentVersion types.ComponentVersion) (sql.Result, error) {
	data, err := json.Marshal(componentVersion.Version.Data)
	if err != nil {
		log.Printf("JSON marshaling failed (%s)", err)
		return nil, err
	}
	insert := fmt.Sprintf(
		"INSERT INTO %s (%s, %s, %s, %s, %s) VALUES('%s', '%s', '%s', '%s', '%s')",
		versionsTableName,
		versionsTableComponentNameKey,
		versionsTableTrainKey,
		versionsTableVersionKey,
		versionsTableReleaseTimeStampKey,
		versionsTableDataKey,
		componentVersion.Component.Name,
		componentVersion.Version.Train,
		componentVersion.Version.Version,
		componentVersion.Version.Released,
		string(data),
	)
	return db.Exec(insert)
}

// GetVersion gets a single version record from a DB matching the unique property values in a ComponentVersion struct
func GetVersion(db *sql.DB, cV types.ComponentVersion) (types.ComponentVersion, error) {
	row := getDBRecord(db, versionsTableName,
		[]string{versionsTableComponentNameKey, versionsTableTrainKey, versionsTableVersionKey},
		[]string{cV.Component.Name, cV.Version.Train, cV.Version.Version})
	rowResult := VersionsTable{}
	//TODO: sql.NullString is to pass tests, not for production
	var s sql.NullString
	if err := row.Scan(&s, &rowResult.componentName, &rowResult.train, &rowResult.version, &rowResult.releaseTimestamp, &rowResult.data); err != nil {
		return types.ComponentVersion{}, err
	}
	componentVersion, err := parseDBVersion(rowResult)
	if err != nil {
		return types.ComponentVersion{}, err
	}
	return componentVersion, nil
}
