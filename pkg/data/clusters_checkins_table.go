package data

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

const (
	clustersCheckinsTableName                = "clusters_checkins"
	clustersCheckinsTableIDKey               = "checkins_id"
	clustersCheckinsTableClusterIDKey        = "cluster_id"
	clustersCheckinsTableClusterCreatedAtKey = "created_at"
	clustersCheckinsTableDataKey             = "data"
)

// ClustersCheckinsTable type that expresses the `clusters_checkins` postgres table schema
type clustersCheckinsTable struct {
	CheckinsID string `gorm:"primary_key;type:bigserial;column_name:checkins_id"`
	ClusterID  string `gorm:"type:uuid;column_name:cluster_id;index"`
	CreatedAt  string `gorm:"type:timestamp;column_name:created_at;index"`
	Data       string `gorm:"type:json;column_name:data"`
}

func newClustersCheckinsTable(checkinID, clusterID string, createdAt time.Time, clusterData []byte) clustersCheckinsTable {
	return clustersCheckinsTable{
		CheckinsID: checkinID,
		ClusterID:  clusterID,
		CreatedAt:  Timestamp{Time: createdAt}.String(),
		Data:       string(clusterData),
	}
}

// BeforeSave is the gorm callback for saving a new cluster checkin
func (c clustersCheckinsTable) BeforeSave() error {
	_, err := newTimestampFromStr(c.CreatedAt)
	return err
}

// BeforeUpdate is the gorm callback for updating a cluster checkin
func (c clustersCheckinsTable) BeforeUpdate() error {
	_, err := newTimestampFromStr(c.CreatedAt)
	return err
}

func (c clustersCheckinsTable) createdAtTime() (time.Time, error) {
	return time.Parse(StdTimestampFmt, c.CreatedAt)
}

func (c clustersCheckinsTable) TableName() string {
	return clustersCheckinsTableName
}

func createClustersCheckinsTable(db *sql.DB) (sql.Result, error) {
	return db.Exec(fmt.Sprintf(
		"CREATE TABLE IF NOT EXISTS %s ( %s bigserial PRIMARY KEY, %s uuid, %s timestamp, %s json )",
		clustersCheckinsTableName,
		clustersCheckinsTableIDKey,
		clustersTableIDKey,
		clustersCheckinsTableClusterCreatedAtKey,
		clustersCheckinsTableDataKey,
	))
}

func verifyClustersCheckinsTable(db *sql.DB) error {
	if _, err := createClustersCheckinsTable(db); err != nil {
		log.Println("unable to verify clusters table exists")
		return err
	}
	return nil
}
