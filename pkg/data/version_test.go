package data

import (
	"fmt"
	"testing"
	"time"

	"github.com/arschles/assert"
	"github.com/deis/workflow-manager-api/pkg/swagger/models"
)

func testComponentVersion() *models.ComponentVersion {
	return &models.ComponentVersion{
		Component: &models.Component{
			Name:        componentName,
			Description: &componentDescription,
		},
		Version: &models.Version{
			Version:  version,
			Released: released,
			Train:    train,
			Data:     &versionData,
		},
		UpdateAvailable: &updateAvailable,
	}
}

func TestVersionRoundTrip(t *testing.T) {
	sqliteDB, err := NewMemDB()
	assert.NoErr(t, err)
	assert.NoErr(t, VerifyPersistentStorage(sqliteDB))
	componentVersion := testComponentVersion()
	cVerNoExist, err := GetVersion(sqliteDB, *componentVersion)
	assert.True(t, err != nil, "error not returned but expected")
	assert.Equal(t, cVerNoExist, models.ComponentVersion{}, "component version")
	cVerSet, err := UpsertVersion(sqliteDB, *componentVersion)
	assert.NoErr(t, err)
	assert.Equal(t, cVerSet.Component.Name, componentVersion.Component.Name, "component name")
	assert.Equal(t, cVerSet.Version.Version, componentVersion.Version.Version, "version string")
	assert.Equal(t, cVerSet.Version.Released, componentVersion.Version.Released, "released string")
	assert.Equal(t, cVerSet.Version.Train, componentVersion.Version.Train, "version train")
	getCVer, err := GetVersion(sqliteDB, *componentVersion)
	assert.NoErr(t, err)
	assert.Equal(t, getCVer.Component.Name, componentVersion.Component.Name, "component name")
	assert.Equal(t, getCVer.Version.Version, componentVersion.Version.Version, "version string")
	assert.Equal(t, getCVer.Version.Released, componentVersion.Version.Released, "released string")
	assert.Equal(t, getCVer.Version.Train, componentVersion.Version.Train, "version train")
}

func TestGetLatestVersion(t *testing.T) {
	sqliteDB, err := NewMemDB()
	assert.NoErr(t, err)
	assert.NoErr(t, VerifyPersistentStorage(sqliteDB))

	const numCVs = 4
	const latestCVIdx = 2
	componentVersions := make([]models.ComponentVersion, numCVs)
	for i := 0; i < numCVs; i++ {
		cv := testComponentVersion()
		cv.Component.Name = componentName
		desc := fmt.Sprintf("description%d", i)
		cv.Component.Description = &desc
		cv.Version.Train = train
		cv.Version.Version = fmt.Sprintf("testversion%d", i)
		cv.Version.Released = time.Now().Add(time.Duration(i) * time.Hour).Format(released)
		cv.Version.Data = &models.VersionData{
			Description: fmt.Sprintf("data%d", i),
		}
		if i == latestCVIdx {
			cv.Version.Released = time.Now().Add(time.Duration(numCVs+1) * time.Hour).Format(released)
		}
		if _, setErr := UpsertVersion(sqliteDB, *cv); setErr != nil {
			t.Fatalf("error setting component version %d (%s)", i, setErr)
		}
		componentVersions[i] = *cv
	}
	cv, err := GetLatestVersion(sqliteDB, train, componentName)
	assert.NoErr(t, err)
	exCV := componentVersions[latestCVIdx]
	assert.Equal(t, cv.Component.Name, exCV.Component.Name, "component name")
	// since the versions table doesn't store a description now, make sure it comes back empty
	assert.Nil(t, cv.Component.Description, "component name")

	assert.Equal(t, cv.Version.Train, exCV.Version.Train, "component version")
	assert.Equal(t, cv.Version.Version, exCV.Version.Version, "component version")
	assert.Equal(t, cv.Version.Released, exCV.Version.Released, "component release time")
	assert.Equal(t, cv.Version.Data, exCV.Version.Data, "component version data")
}

func TestGetLatestVersions(t *testing.T) {
	memDB, err := NewMemDB()
	assert.NoErr(t, err)
	assert.NoErr(t, VerifyPersistentStorage(memDB))
	componentNames := []string{"component1", "component2", "component3"}
	trains := []string{"train1", "train2", "train3"}
	componentAndTrainSlice := make([]ComponentAndTrain, len(componentNames)*len(trains))
	// mapping from component/train (string representation) to release time
	releaseTimes := make(map[string]time.Time)

	const numForEach = 3
	for i, componentName := range componentNames {
		for j, train := range trains {
			ct := ComponentAndTrain{
				ComponentName: componentName,
				Train:         train,
			}
			offset := i * len(componentNames)
			idx := offset + j
			componentAndTrainSlice[idx] = ct
			for n := 0; n < numForEach; n++ {
				cv := testComponentVersion()
				cv.Component.Name = componentName
				cv.Version.Train = train

				//specify a different version for each component version in this name/train
				cv.Version.Version = fmt.Sprintf("version%d-%d", idx, n)
				cvReleaseTime := time.Now().Add(time.Duration((idx+1)*(n+1)) * time.Hour)
				cv.Version.Released = cvReleaseTime.Format(released)

				// record the latest release time for each component
				ct := componentAndTrainFromComponentVersion(cv)
				if releaseTimes[ct.String()].Before(cvReleaseTime) {
					releaseTimes[ct.String()] = cvReleaseTime
				}
				if _, setErr := UpsertVersion(memDB, *cv); setErr != nil {
					t.Fatalf("error setting component version %d (%s)", idx, setErr)
				}
			}
		}
	}

	// now add a component & train that we aren't gonna filter for
	cv := testComponentVersion()
	cv.Component.Name = "invalid"
	cv.Version.Train = "invalid"
	cv.Version.Version = "invalid"
	cv.Version.Released = time.Now().Format(released)
	_, setErr := UpsertVersion(memDB, *cv)
	assert.NoErr(t, setErr)

	componentVersions, err := GetLatestVersions(memDB, componentAndTrainSlice)
	assert.NoErr(t, err)
	assert.Equal(t, len(componentVersions), len(releaseTimes), "number of returned components")
	for _, componentVersion := range componentVersions {
		ct := componentAndTrainFromComponentVersion(componentVersion)
		exReleaseTime := releaseTimes[ct.String()].Format(released)
		if exReleaseTime != componentVersion.Version.Released {
			t.Errorf(
				"expected release time (%s) doesn't match actual (%s) for component %s",
				exReleaseTime,
				componentVersion.Version.Released,
				componentVersion.Component.Name,
			)
		}
	}
}
