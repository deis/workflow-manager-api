package data

import (
	"fmt"
	"time"

	"github.com/deis/workflow-manager-api/rest"
)

// ClusterCheckinsFilter is the struct used to filter on cluster checkins. It represents the conjunction
// of all of its fields. For example:
//
//  MIN(created_at) > created_after
//  AND
//  MIN(created_at) > created_before
type ClusterCheckinsFilter struct {
	CreatedAfter  time.Time
	CreatedBefore time.Time
}

// clustersCheckinsFilterResponse type that represents a `clusters_checkins` filter result
type clustersCheckinsFilterResponse struct {
	ClusterID   string `gorm:"type:uuid;column_name:cluster_id;index"`
	FirstSeen   string `gorm:"type:timestamp;column_name:first_seen"`
	LastSeen    string `gorm:"type:timestamp;column_name:last_seen"`
	ClusterAge  string `gorm:"type:interval;column_name:cluster_age"`
	LastCheckin string `gorm:"type:interval;column_name:last_checkin"`
	Checkins    string `gorm:"type:int;column_name:checkins"`
}

// NewClusterCheckinsFilter returns a new ClusterCheckinsFilter if the given times can result in a valid
// query that would return clusters. If not, returns nil and an ErrImpossibleFilter error
func NewClusterCheckinsFilter(
	createdAfter,
	createdBefore time.Time,
) (*ClusterCheckinsFilter, error) {
	candidate := ClusterCheckinsFilter{
		CreatedAfter:  createdAfter,
		CreatedBefore: createdBefore,
	}
	if err := candidate.checkValid(); err != nil {
		return nil, err
	}
	return &candidate, nil
}

func (c ClusterCheckinsFilter) String() string {
	return fmt.Sprintf(
		"created_after %s, created_before %s",
		c.CreatedAfter,
		c.CreatedBefore,
	)
}

func (c ClusterCheckinsFilter) checkValid() error {
	if c.CreatedAfter.After(c.CreatedBefore) || c.CreatedAfter.Equal(c.CreatedBefore) {
		// you can't have clusters that were created after time T+1
		// and at the same time created before time T-1
		return ErrImpossibleFilter{
			vals: []keyAndTime{
				keyAndTime{key: rest.CreatedAfterQueryStringKey, time: c.CreatedAfter},
				keyAndTime{key: rest.CreatedBeforeQueryStringKey, time: c.CreatedBefore},
			},
			reason: fmt.Sprintf(
				"%s needs to be greater than %s",
				rest.CreatedBeforeQueryStringKey,
				rest.CreatedAfterQueryStringKey,
			),
		}
	}
	return nil
}
