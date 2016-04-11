package data

import (
	"fmt"
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
		},
		UpdateAvailable: updateAvailable,
	}
}

func cVerEqual(v1, v2 types.ComponentVersion) error {
	if v1.Component.Name != v2.Component.Name {
		return fmt.Errorf("component name %s != %s", v1.Component.Name, v2.Component.Name)
	}
	if v1.Component.Description != v2.Component.Description {
		return fmt.Errorf("component description %s != %s", v1.Component.Description, v2.Component.Description)
	}
	if v1.Version.Version != v2.Version.Version {
		return fmt.Errorf("version %s != %s", v1.Version.Version, v2.Version.Version)
	}
	return nil
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
	cVerNoExist, err := ver.Get(sqliteDB, componentName)
	assert.True(t, err != nil, "error not returned but expected")
	assert.Equal(t, cVerNoExist, types.ComponentVersion{}, "component version")
	expectedCVer := testComponentVersion()
	cVerSet, err := ver.Set(sqliteDB, componentName, expectedCVer)
	assert.NoErr(t, err)
	assert.NoErr(t, cVerEqual(cVerSet, expectedCVer))
	getCVer, err := ver.Get(sqliteDB, componentName)
	assert.NoErr(t, err)
	assert.NoErr(t, cVerEqual(getCVer, expectedCVer))
}
