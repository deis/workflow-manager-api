package data

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/deis/workflow-manager/types"
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

// GetLatestVersion gets the latest version from the DB for the given train & component
func GetLatestVersion(db *sql.DB, train string, component string) (types.ComponentVersion, error) {
	rows, err := getOrderedDBRecords(
		db,
		versionsTableName,
		[]string{versionsTableTrainKey, versionsTableComponentNameKey},
		[]string{train, component},
		newOrderBy(versionsTableReleaseTimeStampKey, "desc"),
	)
	if err != nil {
		return types.ComponentVersion{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var row VersionsTable
		//TODO: sql.NullString is to pass tests, not for production
		var s sql.NullString
		if err = rows.Scan(&s, &row.componentName, &row.train, &row.version, &row.releaseTimestamp, &row.data); err != nil {
			return types.ComponentVersion{}, err
		}
		cv, err := parseDBVersion(row)
		if err != nil {
			return types.ComponentVersion{}, err
		}
		return cv, nil
	}
	return types.ComponentVersion{}, errNoMoreRows{tableName: versionsTableName}
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

	rowsResult := []VersionsTable{}
	defer rows.Close()
	for rows.Next() {
		var row VersionsTable
		// note that we have to pass in a *sql.NullString as the first and last arg to ignore the
		// primary key and the final release timestamp returned from the MAX aggregate function
		// in the above SQL
		if err = rows.Scan(
			&sql.NullString{},
			&row.componentName,
			&row.train,
			&row.version,
			&row.releaseTimestamp,
			&row.data,
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

// GetVersionsList retrieves a list of version records from the DB that match a given train & component
func GetVersionsList(db *sql.DB, train string, component string) ([]types.ComponentVersion, error) {
	rows, err := getDBRecords(db, versionsTableName,
		[]string{versionsTableTrainKey, versionsTableComponentNameKey},
		[]string{train, component})
	if err != nil {
		return nil, err
	}
	rowsResult := []VersionsTable{}
	var row VersionsTable
	defer rows.Close()
	for rows.Next() {
		//TODO: sql.NullString is to pass tests, not for production
		var s sql.NullString
		err = rows.Scan(&s, &row.componentName,
			&row.train, &row.version, &row.releaseTimestamp, &row.data)
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
