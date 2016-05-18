package data

import (
	"fmt"
	"strings"
	"time"

	"github.com/deis/workflow-manager-api/rest"
)

type keyAndTime struct {
	key  string
	time time.Time
}

func (k keyAndTime) String() string {
	return fmt.Sprintf("%s (%s)", k.key, k.time)
}

// ErrImpossibleFilter is the error returned when a caller tries to create a new ClusterAgeFilter
// with parameters that would create a filter that is guaranteed to produce no results. One such
// "impossible" query is a "created before" time is after a "created " time that is after  (i.e.
// an impossible filter). See the documentation on ClusterAgeFilter for examples of impossible
// filters
type ErrImpossibleFilter struct {
	vals   []keyAndTime
	reason string
}

// Error is the error interface implementation
func (e ErrImpossibleFilter) Error() string {
	strs := make([]string, len(e.vals))
	for i := 0; i < len(e.vals); i++ {
		strs[i] = e.vals[i].String()
	}
	return fmt.Sprintf("impossible filter for keys/times (%s): %s", strings.Join(strs, ", "), e.reason)
}

// ClusterAgeFilter is the struct used to filter on cluster ages. It represents the conjunction
// of all of its fields. For example:
//
//  created_time<=CreatedBefore
//  AND
//  created_time>=CreatedAfter
//  AND
//  checked_in_time<=CheckedInBefore
//  AND
//  checked_in_time>=CheckedInAfter
type ClusterAgeFilter struct {
	CheckedInBefore time.Time
	CheckedInAfter  time.Time
	CreatedBefore   time.Time
	CreatedAfter    time.Time
}

// NewClusterAgeFilter returns a new ClusterAgeFilter if the given times can result in a valid
// query that would return clusters. If not, returns nil and an ErrImpossibleFilter error
func NewClusterAgeFilter(
	checkedInBefore,
	checkedInAfter,
	createdBefore,
	createdAfter time.Time,
) (*ClusterAgeFilter, error) {
	candidate := ClusterAgeFilter{
		CheckedInBefore: checkedInBefore,
		CheckedInAfter:  checkedInAfter,
		CreatedBefore:   createdBefore,
		CreatedAfter:    createdAfter,
	}
	if err := candidate.checkValid(); err != nil {
		return nil, err
	}
	return &candidate, nil
}

func (c ClusterAgeFilter) String() string {
	return fmt.Sprintf(
		"created_before %s, created_after %s, checked_in_before %s, checked_in_after %s",
		c.CreatedBefore,
		c.CreatedAfter,
		c.CheckedInBefore,
		c.CheckedInAfter,
	)
}

func (c ClusterAgeFilter) checkValid() error {
	if c.CreatedBefore.After(c.CheckedInBefore) {
		// you can't have clusters that were checked in before they were created
		return ErrImpossibleFilter{
			vals: []keyAndTime{
				keyAndTime{key: rest.CreatedBeforeQueryStringKey, time: c.CreatedBefore},
				keyAndTime{key: rest.CheckedInBeforeQueryStringKey, time: c.CheckedInBefore},
			},
			reason: fmt.Sprintf(
				"%s needs to be greater than or equal to %s",
				rest.CheckedInBeforeQueryStringKey,
				rest.CreatedBeforeQueryStringKey,
			),
		}
	} else if c.CheckedInAfter.After(c.CheckedInBefore) || c.CheckedInAfter.Equal(c.CheckedInBefore) {
		// you can't have clusters that were checked in before time T-1
		// and at the same time checked in after time T+1
		return ErrImpossibleFilter{
			vals: []keyAndTime{
				keyAndTime{key: rest.CheckedInBeforeQueryStringKey, time: c.CheckedInBefore},
				keyAndTime{key: rest.CheckedInAfterQueryStringKey, time: c.CheckedInAfter},
			},
			reason: fmt.Sprintf(
				"%s needs to be greater than %s",
				rest.CheckedInBeforeQueryStringKey,
				rest.CheckedInAfterQueryStringKey,
			),
		}
	} else if c.CreatedAfter.After(c.CreatedBefore) || c.CreatedAfter.Equal(c.CreatedBefore) {
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
	} else if c.CheckedInBefore.Before(c.CreatedAfter) || c.CheckedInBefore.Equal(c.CreatedAfter) {
		// you can't have clusters that were checked in before time T-1
		// and at the same time created after time T+1
		return ErrImpossibleFilter{
			vals: []keyAndTime{
				keyAndTime{key: rest.CheckedInBeforeQueryStringKey, time: c.CheckedInBefore},
				keyAndTime{key: rest.CreatedAfterQueryStringKey, time: c.CreatedAfter},
			},
			reason: fmt.Sprintf(
				"%s needs to be after %s",
				rest.CheckedInBeforeQueryStringKey,
				rest.CreatedAfterQueryStringKey,
			),
		}
	}
	return nil
}
