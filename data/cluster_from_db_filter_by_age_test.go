package data

import (
	"database/sql"
	"encoding/json"
	"sync"
	"testing"
	"time"

	"github.com/deis/workflow-manager/types"
	"github.com/pborman/uuid"
)

var (
	validClusterAgeFilter = ClusterAgeFilter{
		CheckedInBefore: timeFuture().Add(2 * time.Hour),
		CheckedInAfter:  timePast(),
		CreatedAfter:    timePast().Add(-1 * time.Hour),
		CreatedBefore:   timeFuture(),
	}
	validCreatedAtTime = timeNow()

	invalidCreatedAtTimes = []time.Time{
		// after the checked_in_before constraint
		validClusterAgeFilter.CheckedInBefore.Add(1 * time.Hour),
		// before the checked_in_after constraint
		validClusterAgeFilter.CheckedInAfter.Add(-1 * time.Hour),
		// before the created_after constraint
		validClusterAgeFilter.CreatedAfter.Add(-1 * time.Hour),
		// after the created_before constraint
		validClusterAgeFilter.CreatedBefore.Add(1 * time.Hour),
	}
)

func invalidClusterAgeFilters() []ClusterAgeFilter {
	return []ClusterAgeFilter{
		// created_before > checked_in_before
		ClusterAgeFilter{
			CreatedBefore:   timeNow(),
			CreatedAfter:    timePast(),
			CheckedInBefore: timePast(),
			CheckedInAfter:  timePast().Add(-1 * time.Hour),
		},
		// checked_in_after > checked_in_before
		ClusterAgeFilter{
			CreatedBefore:   timeFuture(),
			CreatedAfter:    timeNow(),
			CheckedInBefore: timePast(),
			CheckedInAfter:  timeNow(),
		},
		// checked_in_after == checked_in_before
		ClusterAgeFilter{
			CreatedBefore:   timeNow(),
			CreatedAfter:    timePast(),
			CheckedInBefore: timeFuture(),
			CheckedInAfter:  timeFuture(),
		},
		// created_after > created_before
		ClusterAgeFilter{
			CreatedBefore:   timePast(),
			CreatedAfter:    timeNow(),
			CheckedInBefore: timeFuture().Add(1 * time.Hour),
			CheckedInAfter:  timeFuture(),
		},
		// created_after == created_before
		ClusterAgeFilter{
			CreatedBefore:   timePast(),
			CreatedAfter:    timePast(),
			CheckedInBefore: timeFuture(),
			CheckedInAfter:  timeNow(),
		},
		// checked_in_before < created_after
		ClusterAgeFilter{
			CreatedBefore:   timeFuture(),
			CreatedAfter:    timeNow(),
			CheckedInBefore: timePast(),
			CheckedInAfter:  timePast().Add(-1 * time.Hour),
		},
		// checked_in_before == created_after
		ClusterAgeFilter{
			CreatedBefore:   timeFuture(),
			CreatedAfter:    timeNow(),
			CheckedInBefore: timeNow(),
			CheckedInAfter:  timePast().Add(-1 * time.Hour),
		},
	}
}

type filterAndCreatedAt struct {
	filter    ClusterAgeFilter // the filter to run queries with
	createdAt time.Time        // the time to set in the created_at field of the cluster checkin
	expected  bool             // whether or not results are expected to come back for the filter & created at time
}

func createFilters() []filterAndCreatedAt {
	ret := []filterAndCreatedAt{
		filterAndCreatedAt{
			filter:    validClusterAgeFilter,
			createdAt: validCreatedAtTime,
			expected:  true,
		},
	}
	for _, invalidCreatedAtTime := range invalidCreatedAtTimes {
		ret = append(ret, filterAndCreatedAt{
			filter:    validClusterAgeFilter,
			createdAt: invalidCreatedAtTime,
			expected:  false,
		})
	}
	for _, filter := range invalidClusterAgeFilters() {
		ret = append(ret, filterAndCreatedAt{
			filter:    filter,
			createdAt: validCreatedAtTime,
			expected:  false,
		})
	}
	return ret
}

func createCheckin(db *sql.DB, cl types.Cluster, createdAt *Timestamp) error {
	data, err := json.Marshal(cl)
	if err != nil {
		return err
	}
	if _, err = newClusterCheckinsDBRecord(db, cl.ID, createdAt, data); err != nil {
		return err
	}
	return nil
}

func TestClusterFromDBFilterByAge(t *testing.T) {
	cluster := testCluster()
	cluster.ID = uuid.New()

	// generate all combinations of filters
	filters := createFilters()

	var wg sync.WaitGroup
	for i, fca := range filters {
		wg.Add(1)
		go func(i int, fca filterAndCreatedAt) {
			defer wg.Done()
			// sanity check the filter so we don't get false negatives later on in the query
			validErr := fca.filter.checkValid()
			if validErr != nil && fca.expected {
				t.Errorf("expected results for case %d but filter is invalid (%s)", i, validErr)
				return
			}

			c := ClusterFromDB{}
			db, err := newDB()
			if err != nil {
				t.Errorf("error creating new DB in case %d (%s)", i, err)
				return
			}

			// create cluster in the DB
			if _, setErr := c.Set(db, cluster.ID, cluster); setErr != nil {
				t.Errorf("error creating cluster %s in DB (%s)", cluster.ID, setErr)
				return
			}

			// check in the cluster
			if cErr := createCheckin(db, cluster, &Timestamp{Time: &fca.createdAt}); cErr != nil {
				t.Errorf("error creating checkin for case %d (%s)", i, cErr)
				return
			}

			// filter the cluster & test results
			filteredClusters, filterErr := c.FilterByAge(db, &fca.filter)
			if filterErr != nil {
				t.Errorf("error filtering for case %d (%s)", i, filterErr)
				return
			}
			if fca.expected && len(filteredClusters) != 1 {
				t.Errorf("expected %d filtered cluster(s) on filter %d, got %d", 1, i, len(filteredClusters))
				return
			} else if !fca.expected && len(filteredClusters) > 0 {
				t.Errorf("expected 0 filtered clusters for filter %d, got %d", i, len(filteredClusters))
				return
			}

			db, err = newDB()
			if err != nil {
				t.Errorf("error creating new DB in case %d (%s)", i, err)
				return
			}
			filteredClusters, filterErr = c.FilterByAge(db, &fca.filter)
			if filterErr != nil {
				t.Errorf("error filtering for case %d (%s)", i, filterErr)
				return
			}
			if len(filteredClusters) != 0 {
				t.Errorf("expected 0 filtered clusters, got %d for case %d", len(filteredClusters), i)
				return
			}
		}(i, fca)
	}
	wg.Wait()
}
