package data

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/deis/workflow-manager/types"
	"github.com/jinzhu/gorm"
)

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

// SetVersion adds or updates a single version record in the database
func SetVersion(db *gorm.DB, componentVersion types.ComponentVersion) (types.ComponentVersion, error) {
	// TODO: this read-modify-write should be done inside a transaction. Also, rename SetVersion to something else.
	// Both of these TODOs are captured in https://github.com/deis/workflow-manager-api/issues/90

	var ret types.ComponentVersion // return variable
	row := getDBRecord(db.DB(), versionsTableName,
		[]string{versionsTableComponentNameKey, versionsTableTrainKey, versionsTableVersionKey},
		[]string{componentVersion.Component.Name, componentVersion.Version.Train, componentVersion.Version.Version})
	var result sql.Result
	rowResult := versionsTable{}
	if err := row.Scan(&rowResult.VersionID, &rowResult.ComponentName, &rowResult.Train, &rowResult.Version, &rowResult.ReleaseTimestamp, &rowResult.Data); err != nil {
		result, err = newVersionDBRecord(db.DB(), componentVersion)
		if err != nil {
			log.Println(err)
			return types.ComponentVersion{}, err
		}
	} else {
		result, err = updateVersionDBRecord(db.DB(), componentVersion)
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

// GetLatestVersion gets the latest version from the DB for the given train & component
func GetLatestVersion(db *gorm.DB, train string, component string) (types.ComponentVersion, error) {
	resTable := new(versionsTable)
	query := versionsTable{ComponentName: component, Train: train}
	resDB := db.Where(query).Order("release_timestamp desc").First(resTable)
	if resDB.Error != nil {
		return types.ComponentVersion{}, resDB.Error
	}

	componentVersion, err := parseDBVersion(*resTable)
	if err != nil {
		return types.ComponentVersion{}, err
	}
	return componentVersion, nil
}

// GetLatestVersions fetches from the DB and returns the latest versions for each component/train pair
// given in ct. Returns an empty slice and non-nil error on any error communicating with the
// database or otherwise if the first returned value is not empty, it's guaranteed to:
//
// - Be the same length as ct
// - Have the same ordering as ct, with respect to the component name
func GetLatestVersions(db *sql.DB, ct []ComponentAndTrain) ([]types.ComponentVersion, error) {
	componentsList := []string{}
	listedComponents := make(map[string]struct{})
	trainsList := []string{}
	listedTrains := make(map[string]struct{})
	for _, c := range ct {
		if _, componentListed := listedComponents[c.ComponentName]; !componentListed {
			componentsList = append(componentsList, fmt.Sprintf("'%s'", c.ComponentName))
			listedComponents[c.ComponentName] = struct{}{}
		}
		if _, trainListed := listedTrains[c.Train]; !trainListed {
			trainsList = append(trainsList, fmt.Sprintf("'%s'", c.Train))
			listedTrains[c.Train] = struct{}{}
		}
	}
	query := fmt.Sprintf(
		"SELECT *, MAX(%s) FROM %s WHERE %s IN (%s) AND %s IN (%s) GROUP BY %s, %s",
		versionsTableReleaseTimeStampKey,
		versionsTableName,
		versionsTableComponentNameKey,
		strings.Join(componentsList, ","),
		versionsTableTrainKey,
		strings.Join(trainsList, ","),
		versionsTableComponentNameKey,
		versionsTableTrainKey,
	)
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}

	rowsResult := []versionsTable{}
	defer rows.Close()
	for rows.Next() {
		var row versionsTable
		// note that we have to pass in a *sql.NullString as the first and last arg to ignore the
		// primary key and the final release timestamp returned from the MAX aggregate function
		// in the above SQL
		if err = rows.Scan(
			&sql.NullString{},
			&row.ComponentName,
			&row.Train,
			&row.Version,
			&row.ReleaseTimestamp,
			&row.Data,
			&sql.NullString{},
		); err != nil {
			return nil, err
		}
		rowsResult = append(rowsResult, row)
	}
	componentVersions, err := parseDBVersions(rowsResult)
	if err != nil {
		return []types.ComponentVersion{}, err
	}
	return componentVersions, nil
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
func GetVersion(db *gorm.DB, cV types.ComponentVersion) (types.ComponentVersion, error) {
	resTable := new(versionsTable)
	resDB := db.Where(versionsTable{
		ComponentName: cV.Component.Name,
		Train:         cV.Version.Train,
		Version:       cV.Version.Version,
	}).First(resTable)
	if resDB.Error != nil {
		return types.ComponentVersion{}, resDB.Error
	}

	componentVersion, err := parseDBVersion(*resTable)
	if err != nil {
		return types.ComponentVersion{}, err
	}
	return componentVersion, nil
}

// GetVersionsList retrieves a list of version records from the DB that match a given train & component
func GetVersionsList(db *sql.DB, train string, component string) ([]types.ComponentVersion, error) {
	rows, err := getDBRecords(db, versionsTableName,
		[]string{versionsTableTrainKey, versionsTableComponentNameKey},
		[]string{train, component})
	if err != nil {
		return nil, err
	}
	rowsResult := []versionsTable{}
	var row versionsTable
	defer rows.Close()
	for rows.Next() {
		//TODO: sql.NullString is to pass tests, not for production
		var s sql.NullString
		err = rows.Scan(&s, &row.ComponentName,
			&row.Train, &row.Version, &row.ReleaseTimestamp, &row.Data)
		if err != nil {
			return nil, err
		}
		rowsResult = append(rowsResult, row)
	}
	componentVersions, err := parseDBVersions(rowsResult)
	if err != nil {
		log.Println("error parsing DB versions data")
		return nil, err
	}
	return componentVersions, nil
}
