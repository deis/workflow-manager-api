package data

import (
	"fmt"
	"testing"
	"time"

	"github.com/arschles/assert"
	"github.com/deis/workflow-manager/types"
)

func testComponentVersion() types.ComponentVersion {
	return types.ComponentVersion{
		Component: types.Component{
			Name:        componentName,
			Description: componentDescription,
		},
		Version: types.Version{
			Version:  version,
			Released: released,
			Train:    train,
			Data:     versionData,
		},
		UpdateAvailable: updateAvailable,
	}
}

func TestVersionFromDBRoundTrip(t *testing.T) {
	memDB, err := NewMemDB()
	assert.NoErr(t, err)
	sqliteDB, err := memDB.Get()
	assert.NoErr(t, err)
	db, err := VerifyPersistentStorage(memDB)
	assert.NotNil(t, db, "db")
	assert.NoErr(t, err)
	ver := VersionFromDB{}
	componentVersion := testComponentVersion()
	cVerNoExist, err := ver.Get(sqliteDB, componentVersion)
	assert.True(t, err != nil, "error not returned but expected")
	assert.Equal(t, cVerNoExist, types.ComponentVersion{}, "component version")
	cVerSet, err := ver.Set(sqliteDB, componentVersion)
	assert.NoErr(t, err)
	assert.Equal(t, cVerSet.Component.Name, componentVersion.Component.Name, "component name")
	assert.Equal(t, cVerSet.Version.Version, componentVersion.Version.Version, "version string")
	assert.Equal(t, cVerSet.Version.Released, componentVersion.Version.Released, "released string")
	assert.Equal(t, cVerSet.Version.Train, componentVersion.Version.Train, "version train")
	getCVer, err := ver.Get(sqliteDB, componentVersion)
	assert.NoErr(t, err)
	assert.Equal(t, getCVer.Component.Name, componentVersion.Component.Name, "component name")
	assert.Equal(t, getCVer.Version.Version, componentVersion.Version.Version, "version string")
	assert.Equal(t, getCVer.Version.Released, componentVersion.Version.Released, "released string")
	assert.Equal(t, getCVer.Version.Train, componentVersion.Version.Train, "version train")
}

func TestVersionFromDBLatest(t *testing.T) {
	memDB, err := NewMemDB()
	assert.NoErr(t, err)
	sqliteDB, err := memDB.Get()
	assert.NoErr(t, err)
	db, err := VerifyPersistentStorage(memDB)
	assert.NotNil(t, db, "db")
	assert.NoErr(t, err)
	ver := VersionFromDB{}

	const numCVs = 4
	const latestCVIdx = 2
	componentVersions := make([]types.ComponentVersion, numCVs)
	for i := 0; i < numCVs; i++ {
		cv := testComponentVersion()
		cv.Component.Name = componentName
		cv.Component.Description = fmt.Sprintf("description%d", i)
		cv.Version.Train = train
		cv.Version.Version = fmt.Sprintf("testversion%d", i)
		cv.Version.Released = time.Now().Add(time.Duration(i) * time.Hour).Format(released)
		cv.Version.Data = []byte(fmt.Sprintf("data%d", i))
		if i == latestCVIdx {
			cv.Version.Released = time.Now().Add(time.Duration(numCVs+1) * time.Hour).Format(released)
		}
		if _, setErr := ver.Set(sqliteDB, cv); setErr != nil {
			t.Fatalf("error setting component version %d (%s)", i, setErr)
		}
		componentVersions[i] = cv
	}
	cv, err := ver.Latest(sqliteDB, train, componentName)
	assert.NoErr(t, err)
	exCV := componentVersions[latestCVIdx]
	assert.Equal(t, cv.Component.Name, exCV.Component.Name, "component name")
	// since the versions table doesn't store a description now, make sure it comes back empty
	assert.Equal(t, cv.Component.Description, "", "component name")

	assert.Equal(t, cv.Version.Train, exCV.Version.Train, "component version")
	assert.Equal(t, cv.Version.Version, exCV.Version.Version, "component version")
	assert.Equal(t, cv.Version.Released, exCV.Version.Released, "component release time")
	assert.Equal(t, string(cv.Version.Data), string(exCV.Version.Data), "component version data")
}
