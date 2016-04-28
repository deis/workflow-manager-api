package data

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/deis/workflow-manager/types"
)

// VersionFromDB fulfills the Version interface
type VersionFromDB struct{}

// Get method for VersionFromDB, the actual database/sql.DB implementation
func (c VersionFromDB) Get(db *sql.DB, cV types.ComponentVersion) (types.ComponentVersion, error) {
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

// Collection method for VersionFromDB, the actual database/sql.DB implementation
func (c VersionFromDB) Collection(db *sql.DB, train string, component string) ([]types.ComponentVersion, error) {
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

// Latest method for VersionFromDB, the actual database/sql.DB implementation
func (c VersionFromDB) Latest(db *sql.DB, train string, component string) (types.ComponentVersion, error) {
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

// Set method for VersionFromDB, the actual database/sql.DB implementation
func (c VersionFromDB) Set(db *sql.DB, componentVersion types.ComponentVersion) (types.ComponentVersion, error) {
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
		ret, err = c.Get(db, componentVersion)
		if err != nil {
			return types.ComponentVersion{}, err
		}
	} else if affected > 1 {
		log.Println("updated more than one record with same ID value!")
	}
	return ret, nil
}

// MultiLatest is the Version interface implementation
func (c VersionFromDB) MultiLatest(db *sql.DB, ct []ComponentAndTrain) ([]types.ComponentVersion, error) {
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
