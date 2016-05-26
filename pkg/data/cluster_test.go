package data

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/arschles/assert"
	"github.com/deis/workflow-manager-api/pkg/swagger/models"
	"github.com/jinzhu/gorm"
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
func newDB() (*gorm.DB, error) {
	db, err := NewMemDB()
	if err != nil {
		return nil, err
	}
	if err := VerifyPersistentStorage(db); err != nil {
		return nil, err
	}
	return db, nil
}

func testCluster() models.Cluster {
	cluster := models.Cluster{}
	cluster.ID = clusterID
	cluster.Components = []*models.ComponentVersion{testComponentVersion()}
	return cluster
}

func TestClusterRoundTrip(t *testing.T) {
	sqliteDB, err := newDB()
	assert.NoErr(t, err)
	cluster, err := GetCluster(sqliteDB, clusterID)
	assert.True(t, err != nil, "error not returned when expected")
	assert.Equal(t, cluster, models.Cluster{}, "returned cluster")
	expectedCluster := testCluster()
	// the first time we invoke .CheckInAndSetCluster() it will create a new record
	newCluster, err := UpsertCluster(sqliteDB, clusterID, expectedCluster)
	assert.NoErr(t, err)
	assert.Equal(t, newCluster.ID, expectedCluster.ID, "cluster ID property")
	assert.Equal(t, newCluster.Components[0].Component.Description, expectedCluster.Components[0].Component.Description, "cluster component description property")
	// modify the cluster object
	desc := "new description"
	expectedCluster.Components[0].Component.Description = &desc
	// the next time we invoke .CheckInAndSetCluster() it should update the existing record we just created
	updatedCluster, err := UpsertCluster(sqliteDB, clusterID, expectedCluster)
	assert.NoErr(t, err)
	assert.Equal(t, updatedCluster.Components[0].Component.Description, expectedCluster.Components[0].Component.Description, "cluster component description property")
	getCluster, err := GetCluster(sqliteDB, clusterID)
	assert.NoErr(t, err)
	assert.Equal(t, getCluster, updatedCluster, "cluster")
}

func TestClusterFromDBCheckin(t *testing.T) {
	sqliteDB, err := newDB()
	assert.NoErr(t, err)
	assert.NoErr(t, CheckInCluster(sqliteDB, clusterID, time.Now(), testCluster()))
}
