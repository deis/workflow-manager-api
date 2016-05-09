package data

import (
	"database/sql"
	"fmt"

	"github.com/jinzhu/gorm"
)

const (
	versionsTableName                = "versions"
	versionsTableIDKey               = "version_id"
	versionsTableComponentNameKey    = "component_name"
	versionsTableTrainKey            = "train"
	versionsTableVersionKey          = "version"
	versionsTableReleaseTimeStampKey = "release_timestamp"
	versionsTableDataKey             = "data"
)

// VersionsTable type that expresses the `deis_component_versions` postgres table schema
type versionsTable struct {
	VersionID        string    `gorm:"primary_key;type:uuid;column:version_id"`
	ComponentName    string    `gorm:"column:component_name;index;unique"`
	Train            string    `gorm:"column:train;index;unique"`
	Version          string    `gorm:"column:version;index;unique"`
	ReleaseTimestamp Timestamp `gorm:"column:release_timestamp;type:timestamp"`
	Data             string    `gorm:"column:data;type:json"`
}

func (v versionsTable) TableName() string {
	return versionsTableName
}

func createOrUpdateVersionsTable(db *gorm.DB) (sql.Result, error) {
	query := fmt.Sprintf(
		"CREATE TABLE IF NOT EXISTS %s ( %s bigserial PRIMARY KEY, %s varchar(32), %s varchar(24), %s varchar(32), %s timestamp, %s json, unique (%s, %s, %s) )",
		versionsTableName,
		versionsTableIDKey,
		versionsTableComponentNameKey,
		versionsTableTrainKey,
		versionsTableVersionKey,
		versionsTableReleaseTimeStampKey,
		versionsTableDataKey,
		versionsTableComponentNameKey,
		versionsTableTrainKey,
		versionsTableVersionKey,
	)
	return db.DB().Exec(query)
}
