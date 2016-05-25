package data

import (
	"encoding/json"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/pborman/uuid"
)

var (
	validClusterAgeFilters = []ClusterAgeFilter{
		ClusterAgeFilter{
			CheckedInBefore: timeFuture().Add(2 * time.Hour),
			CheckedInAfter:  timePast(),
			CreatedAfter:    timePast().Add(-1 * time.Hour),
			CreatedBefore:   timeFuture(),
		},
		ClusterAgeFilter{
			CheckedInBefore: timeFuture().Add(3 * time.Hour),
			CheckedInAfter:  timePast(),
			CreatedAfter:    timePast().Add(-1 * time.Hour),
			CreatedBefore:   timeFuture(),
		},
	}
	validCreatedAtTime = timeNow()
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
	ret := []filterAndCreatedAt{}

	// valid cluster filters and created time
	for _, validClusterAgeFilter := range validClusterAgeFilters {
		ret = append(ret, filterAndCreatedAt{
			filter:    validClusterAgeFilter,
			createdAt: validCreatedAtTime,
			expected:  true,
		})
	}

	// valid cluster filters and created time that's guaranteed to return no clusters
	for _, validClusterAgeFilter := range validClusterAgeFilters {
		ret = append(ret, filterAndCreatedAt{
			filter:    validClusterAgeFilter,
			createdAt: validClusterAgeFilter.CreatedBefore.Add(20 * time.Hour),
			expected:  false,
		})
	}

	// invalid cluster filters and a created time that's 1 hour before the filter's created before time
	for _, filter := range invalidClusterAgeFilters() {
		ret = append(ret, filterAndCreatedAt{
			filter:    filter,
			createdAt: filter.CreatedBefore.Add(-1 * time.Hour),
			expected:  false,
		})
	}
	return ret
}

func createAndCheckinClusters(db *gorm.DB, totalNumClusters, filterNum int, fca filterAndCreatedAt) error {
	// create and check in clusters in the DB
	for clusterNum := 0; clusterNum < totalNumClusters; clusterNum++ {
		cluster := testCluster()
		cluster.ID = uuid.New()
		clusterJSON, marshalErr := json.Marshal(cluster)
		if marshalErr != nil {
			return fmt.Errorf("error JSON serializing cluster %d for filter %d (%s)", clusterNum, filterNum, marshalErr)
		}
		createDB := db.Model(&clustersTable{}).Create(&clustersTable{ClusterID: cluster.ID, Data: string(clusterJSON)})
		if createDB.Error != nil {
			return fmt.Errorf("error creating cluster %s for filter %d in DB (%s)", cluster.ID, filterNum, createDB.Error)
		}
		if cErr := CheckInCluster(db, cluster.ID, fca.createdAt, cluster); cErr != nil {
			return fmt.Errorf("error creating checkin for cluster %d, filter %d (%s)", clusterNum, filterNum, cErr)
		}
	}
	return nil
}

func TestClusterFilterByAge(t *testing.T) {
	const numClusters = 4

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

			db, err := newDB()
			if err != nil {
				t.Errorf("error creating new DB in case %d (%s)", i, err)
				return
			}

			if ccErr := createAndCheckinClusters(db, numClusters, i, fca); ccErr != nil {
				t.Errorf("Error creating and checking in clusters (%s)", ccErr)
				return
			}

			// filter the cluster & test results
			filteredClusters, filterErr := FilterClustersByAge(db, &fca.filter)
			if filterErr != nil {
				t.Errorf("error filtering for case %d (%s)", i, filterErr)
				return
			}
			if fca.expected && len(filteredClusters) != numClusters {
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
			filteredClusters, filterErr = FilterClustersByAge(db, &fca.filter)
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
