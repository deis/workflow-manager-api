package data

import (
	"testing"
	"time"

	"github.com/arschles/assert"
	"github.com/deis/workflow-manager/types"
)

func testCluster() types.Cluster {
	return types.Cluster{
		ID:         clusterID,
		FirstSeen:  time.Now(),
		LastSeen:   time.Now().Add(1 * time.Hour),
		Components: []types.ComponentVersion{testComponentVersion()},
	}
}

func TestClusterFromDBRoundTrip(t *testing.T) {
	db, err := newMemDB()
	assert.NoErr(t, err)
	sqliteDB, err := db.Get()
	assert.NoErr(t, err)

	assert.NoErr(t, VerifyPersistentStorage(db))
	c := ClusterFromDB{}
	cluster, err := c.Get(sqliteDB, clusterID)
	assert.True(t, err != nil, "error not returned when expected")
	assert.Equal(t, cluster, types.Cluster{}, "returned cluster")
	expectedCluster := testCluster()
	setCluster, err := c.Set(sqliteDB, clusterID, expectedCluster)
	assert.NoErr(t, err)
	assert.Equal(t, setCluster.ID, expectedCluster.ID, "cluster")
	getCluster, err := c.Get(sqliteDB, clusterID)
	assert.NoErr(t, err)
	assert.Equal(t, getCluster, setCluster, "cluster")
}

func TestClusterFromDBCheckin(t *testing.T) {
	db, err := newMemDB()
	assert.NoErr(t, err)
	sqliteDB, err := db.Get()
	assert.NoErr(t, err)

	assert.NoErr(t, VerifyPersistentStorage(db))
	c := ClusterFromDB{}
	res, err := c.Checkin(sqliteDB, clusterID, testCluster())
	assert.NoErr(t, err)
	rowsAffected, err := res.RowsAffected()
	assert.NoErr(t, err)
	assert.Equal(t, rowsAffected, int64(1), "number of rows affected")
}
