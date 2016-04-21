package data

import (
	"testing"

	"github.com/arschles/assert"
	"github.com/deis/workflow-manager/types"
)

func testCluster() types.Cluster {
	return types.Cluster{
		ID:         clusterID,
		Components: []types.ComponentVersion{testComponentVersion()},
	}
}

func TestClusterFromDBRoundTrip(t *testing.T) {
	memDB, err := NewMemDB()
	assert.NoErr(t, err)
	sqliteDB, err := memDB.Get()
	assert.NoErr(t, err)
	_, err = VerifyPersistentStorage(memDB)
	assert.NoErr(t, err)
	c := ClusterFromDB{}
	cluster, err := c.Get(sqliteDB, clusterID)
	assert.True(t, err != nil, "error not returned when expected")
	assert.Equal(t, cluster, types.Cluster{}, "returned cluster")
	expectedCluster := testCluster()
	// the first time we invoke .Set() it will create a new record
	newCluster, err := c.Set(sqliteDB, clusterID, expectedCluster)
	assert.NoErr(t, err)
	assert.Equal(t, newCluster.ID, expectedCluster.ID, "cluster ID property")
	assert.Equal(t, newCluster.Components[0].Component.Description, expectedCluster.Components[0].Component.Description, "cluster component description property")
	// modify the cluster object
	expectedCluster.Components[0].Component.Description = "new description"
	// the next time we invoke .Set() it should update the existing record we just created
	updatedCluster, err := c.Set(sqliteDB, clusterID, expectedCluster)
	assert.NoErr(t, err)
	assert.Equal(t, updatedCluster.Components[0].Component.Description, expectedCluster.Components[0].Component.Description, "cluster component description property")
	getCluster, err := c.Get(sqliteDB, clusterID)
	assert.NoErr(t, err)
	assert.Equal(t, getCluster, updatedCluster, "cluster")
}

func TestClusterFromDBCheckin(t *testing.T) {
	memDB, err := NewMemDB()
	assert.NoErr(t, err)
	sqliteDB, err := memDB.Get()
	assert.NoErr(t, err)
	db, err := VerifyPersistentStorage(memDB)
	assert.NotNil(t, db, "db")
	assert.NoErr(t, err)
	c := ClusterFromDB{}
	res, err := c.Checkin(sqliteDB, clusterID, testCluster())
	assert.NoErr(t, err)
	rowsAffected, err := res.RowsAffected()
	assert.NoErr(t, err)
	assert.Equal(t, rowsAffected, int64(1), "number of rows affected")
}
