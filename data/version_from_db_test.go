package data

import (
	"testing"

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
	memDB, err := newMemDB()
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
