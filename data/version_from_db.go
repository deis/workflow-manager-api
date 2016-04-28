package data

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/deis/workflow-manager/types"
)

// VersionFromDB fulfills the Version interface
type VersionFromDB struct{}

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
