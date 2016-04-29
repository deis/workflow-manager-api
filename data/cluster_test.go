package data

import (
	"database/sql"
	"errors"
	"fmt"
	"testing"

	"github.com/arschles/assert"
	"github.com/deis/workflow-manager/types"
)

var (
	errInvalidCreatedTime = errors.New("no valid created time is possible")
)

type errInvalidNumSet struct {
	num int
}

func (e errInvalidNumSet) Error() string {
	return fmt.Sprintf("invalid num set: %d", e.num)
}

// creates a new DB and calls VerifyPersistentStorage on it to set it up
func newDB() (*sql.DB, error) {
	db, err := NewMemDB()
	if err != nil {
		return nil, err
	}
	if err := VerifyPersistentStorage(db); err != nil {
		return nil, err
	}
	return db, nil
}

func testCluster() types.Cluster {
	return types.Cluster{
		ID:         clusterID,
		Components: []types.ComponentVersion{testComponentVersion()},
	}
}

func TestClusterFromDBRoundTrip(t *testing.T) {
	sqliteDB, err := newDB()
	assert.NoErr(t, err)
	cluster, err := GetCluster(sqliteDB, clusterID)
	assert.True(t, err != nil, "error not returned when expected")
	assert.Equal(t, cluster, types.Cluster{}, "returned cluster")
	expectedCluster := testCluster()
	// the first time we invoke .CheckInAndSetCluster() it will create a new record
	newCluster, err := CheckInAndSetCluster(sqliteDB, clusterID, expectedCluster)
	assert.NoErr(t, err)
	assert.Equal(t, newCluster.ID, expectedCluster.ID, "cluster ID property")
	assert.Equal(t, newCluster.Components[0].Component.Description, expectedCluster.Components[0].Component.Description, "cluster component description property")
	// modify the cluster object
	expectedCluster.Components[0].Component.Description = "new description"
	// the next time we invoke .CheckInAndSetCluster() it should update the existing record we just created
	updatedCluster, err := CheckInAndSetCluster(sqliteDB, clusterID, expectedCluster)
	assert.NoErr(t, err)
	assert.Equal(t, updatedCluster.Components[0].Component.Description, expectedCluster.Components[0].Component.Description, "cluster component description property")
	getCluster, err := GetCluster(sqliteDB, clusterID)
	assert.NoErr(t, err)
	assert.Equal(t, getCluster, updatedCluster, "cluster")
}

func TestClusterFromDBCheckin(t *testing.T) {
	sqliteDB, err := newDB()
	assert.NoErr(t, err)
	res, err := CheckInCluster(sqliteDB, clusterID, testCluster())
	assert.NoErr(t, err)
	rowsAffected, err := res.RowsAffected()
	assert.NoErr(t, err)
	assert.Equal(t, rowsAffected, int64(1), "number of rows affected")
}
