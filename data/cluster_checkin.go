package data

import (
	stypes "github.com/jmoiron/sqlx/types"
)

// ClusterCheckinTable expresses the 'clusters_checkins' table schema. It's also the Gorm model
type ClusterCheckinTable struct {
	CheckinID string          `gorm:"primary_key;type:uuid;column:checkins_id"`
	ClusterID string          `gorm:"type:uuid;column:cluster_id"`      // indexed
	CreatedAt *Timestamp      `gorm:"type:timestamp;column:created_at"` // indexed
	Data      stypes.JSONText `gorm:"type:json;column:data"`
}

// TableName overrides the gorm-inferred table name
func (c ClusterCheckinTable) TableName() string {
	return "clusters_checkins"
}
