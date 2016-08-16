package data

import (
	"database/sql"
	"encoding/json"
	"log"

	"github.com/deis/workflow-manager-api/pkg/swagger/models"
	"github.com/jinzhu/gorm"
)

// UpsertVersion adds or updates a single version record in the database
func UpsertVersion(db *gorm.DB, componentVersion models.ComponentVersion) (models.ComponentVersion, error) {
	releaseTimestamp, err := newTimestampFromStr(componentVersion.Version.Released)
	if err != nil {
		return models.ComponentVersion{}, err
	}

	js, err := json.Marshal(componentVersion.Version.Data)
	if err != nil {
		return models.ComponentVersion{}, err
	}

	// the query used to find the original version
	queryVsn := versionsTable{
		ComponentName: componentVersion.Component.Name,
		Train:         componentVersion.Version.Train,
		Version:       componentVersion.Version.Version,
	}

	// the new version
	newVsn := versionsTable{
		ComponentName:    componentVersion.Component.Name,
		Train:            componentVersion.Version.Train,
		Version:          componentVersion.Version.Version,
		ReleaseTimestamp: releaseTimestamp,
		Data:             string(js),
	}

	tx := db.Begin()
	cvPtr, err := upsertVersion(tx, queryVsn, newVsn)
	if err != nil {
		rollbackDB := tx.Rollback()
		if rollbackDB.Error != nil {
			return models.ComponentVersion{}, txErr{op: "rollback", orig: err, err: rollbackDB.Error}
		}
		return models.ComponentVersion{}, err
	}
	commitDB := tx.Commit()
	if commitDB.Error != nil {
		log.Println("6")
		return models.ComponentVersion{}, txErr{op: "commit", orig: nil, err: commitDB.Error}
	}
	return *cvPtr, nil
}

// GetLatestVersion gets the latest version from the DB for the given train & component
func GetLatestVersion(db *gorm.DB, train string, component string) (models.ComponentVersion, error) {
	resTable := new(versionsTable)
	query := versionsTable{ComponentName: component, Train: train}
	resDB := db.Where(query).Order("release_timestamp desc").First(resTable)
	if resDB.Error != nil {
		return models.ComponentVersion{}, resDB.Error
	}

	componentVersion, err := parseDBVersion(*resTable)
	if err != nil {
		return models.ComponentVersion{}, err
	}
	return componentVersion, nil
}

// GetLatestVersions fetches from the DB and returns the latest versions for each component/train pair
// given in ct. Returns an empty slice and non-nil error on any error communicating with the
// database or otherwise if the first returned value is not empty, it's guaranteed to:
//
// - Be the same length as ct
// - Have the same ordering as ct, with respect to the component name
func GetLatestVersions(db *gorm.DB, ct []ComponentAndTrain) ([]*models.ComponentVersion, error) {
	componentsList := []string{}
	listedComponents := make(map[string]struct{})
	trainsList := []string{}
	listedTrains := make(map[string]struct{})
	for _, c := range ct {
		if _, componentListed := listedComponents[c.ComponentName]; !componentListed {
			componentsList = append(componentsList, c.ComponentName)
			listedComponents[c.ComponentName] = struct{}{}
		}
		if _, trainListed := listedTrains[c.Train]; !trainListed {
			trainsList = append(trainsList, c.Train)
			listedTrains[c.Train] = struct{}{}
		}
	}
	rows, err := db.Raw("select * from versions as ver where ver.component_name IN (?) AND ver.train IN (?) AND release_timestamp = (select MAX(release_timestamp) from versions as ver1 where ver1.component_name = ver.component_name AND ver1.train = ver.train)", componentsList, trainsList).
		Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	rowsResult := []versionsTable{}
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
		); err != nil {
			return nil, err
		}
		rowsResult = append(rowsResult, row)
	}
	if rErr := rows.Err(); rErr != nil {
		return nil, rErr
	}

	componentVersions, err := parseDBVersions(rowsResult)
	if err != nil {
		return []*models.ComponentVersion{}, err
	}
	return componentVersions, nil
}

// GetVersion gets a single version record from a DB matching the unique property values in a ComponentVersion struct
func GetVersion(db *gorm.DB, cV models.ComponentVersion) (models.ComponentVersion, error) {
	resTable := new(versionsTable)
	resDB := db.Where(versionsTable{
		ComponentName: cV.Component.Name,
		Train:         cV.Version.Train,
		Version:       cV.Version.Version,
	}).First(resTable)
	if resDB.Error != nil {
		return models.ComponentVersion{}, resDB.Error
	}

	componentVersion, err := parseDBVersion(*resTable)
	if err != nil {
		return models.ComponentVersion{}, err
	}
	return componentVersion, nil
}

// GetVersionsList retrieves a list of version records from the DB that match a given train & component
func GetVersionsList(db *gorm.DB, train string, component string) ([]*models.ComponentVersion, error) {
	var rowsResult []versionsTable
	resDB := db.Where(&versionsTable{Train: train, ComponentName: component}).Find(&rowsResult)
	if resDB.Error != nil {
		return nil, resDB.Error
	}
	componentVersions, err := parseDBVersions(rowsResult)
	if err != nil {
		log.Println("error parsing DB versions data")
		return nil, err
	}
	return componentVersions, nil
}

func parseDBVersions(versions []versionsTable) ([]*models.ComponentVersion, error) {
	componentVersions := make([]*models.ComponentVersion, len(versions))
	for i, version := range versions {
		cver, err := parseDBVersion(version)
		if err != nil {
			return nil, err
		}
		componentVersions[i] = &cver
	}
	return componentVersions, nil
}

func parseDBVersion(version versionsTable) (models.ComponentVersion, error) {
	data := models.VersionData{}
	if err := json.Unmarshal([]byte(version.Data), &data); err != nil {
		return models.ComponentVersion{}, err
	}
	return models.ComponentVersion{
		Component: &models.Component{
			Name: version.ComponentName,
		},
		Version: &models.Version{
			Train:    version.Train,
			Version:  version.Version,
			Released: version.ReleaseTimestamp.String(),
			Data:     &data,
		},
	}, nil
}
