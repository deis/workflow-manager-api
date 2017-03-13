package data

import (
	"fmt"
	"time"

	"github.com/deis/workflow-manager-api/rest"
)

// PersistentClustersFilter is the struct used to filter on persistent clusters at a given time.
// For example:
//
//  MIN(created_at) > epoch
//  AND
//  MIN(created_at) > timestamp
//  AND
//  COUNT(checkins) > 1
//  AND MAX(created_at) > timestamp - 24 hours
type PersistentClustersFilter struct {
	Epoch             time.Time
	Timestamp         time.Time
	RelativeYesterday time.Time
}

// NewPersistentClustersFilter returns a new PersistentClustersFilter if the given times can result in a valid
// query that would return clusters. If not, returns nil and an ErrImpossibleFilter error
func NewPersistentClustersFilter(
	epoch,
	timestamp time.Time,
) (*PersistentClustersFilter, error) {
	candidate := PersistentClustersFilter{
		Epoch:             epoch,
		Timestamp:         timestamp,
		RelativeYesterday: timestamp.AddDate(0, 0, -1),
	}
	if err := candidate.checkValid(); err != nil {
		return nil, err
	}
	return &candidate, nil
}

func (c PersistentClustersFilter) String() string {
	return fmt.Sprintf(
		"epoch %s, timestamp %s",
		c.Epoch,
		c.Timestamp,
	)
}

func (c PersistentClustersFilter) checkValid() error {
	if c.Epoch.After(c.Timestamp) || c.Epoch.Equal(c.Timestamp) {
		// you can't have persistent clusters at a time T with an epoch time T+1
		return ErrImpossibleFilter{
			vals: []keyAndTime{
				keyAndTime{key: rest.EpochQueryStringKey, time: c.Epoch},
				keyAndTime{key: rest.TimestampQueryStringKey, time: c.Timestamp},
			},
			reason: fmt.Sprintf(
				"%s needs to be less than %s",
				rest.EpochQueryStringKey,
				rest.TimestampQueryStringKey,
			),
		}
	}
	return nil
}
