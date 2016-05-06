package data

import (
	"testing"

	"github.com/arschles/assert"
	"github.com/pborman/uuid"
)

func TestGetClusterCount(t *testing.T) {
	db, err := newDB()
	count, err := GetClusterCount(db)
	assert.NoErr(t, err)
	assert.Equal(t, count, 0, "count")
	d1 := db.Create(&clustersTable{ClusterID: uuid.New(), Data: []byte("{}")})
	assert.NoErr(t, d1.Error)
	count, err = GetClusterCount(d1)
	assert.NoErr(t, err)
	assert.Equal(t, count, 1, "count")
}
