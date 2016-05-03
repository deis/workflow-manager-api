package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/arschles/assert"
	"github.com/deis/workflow-manager-api/data"
	"github.com/deis/workflow-manager-api/rest"
	"github.com/deis/workflow-manager/types"
	"github.com/pborman/uuid"
)

var (
	nowTime    = time.Now()
	futureTime = nowTime.Add(1 * time.Hour)
	pastTime   = nowTime.Add(-1 * time.Hour)
)

func timeFuture() time.Time {
	return futureTime
}

func timePast() time.Time {
	return pastTime
}

func timeNow() time.Time {
	return nowTime
}

func TestFilterByClusterAge(t *testing.T) {
	filter := data.ClusterAgeFilter{
		CheckedInBefore: timeFuture().Add(2 * time.Hour),
		CheckedInAfter:  timePast(),
		CreatedAfter:    timePast().Add(-1 * time.Hour),
		CreatedBefore:   timeFuture(),
	}
	memDB, err := data.NewMemDB()
	assert.NoErr(t, err)
	assert.NoErr(t, data.VerifyPersistentStorage(memDB))
	cluster := types.Cluster{ID: uuid.New()}
	srv := newServer(memDB)
	defer srv.Close()
	_, setErr := data.CheckInAndSetCluster(memDB, cluster.ID, cluster)
	assert.NoErr(t, setErr)
	_, checkInErr := data.CheckInCluster(memDB.DB(), cluster.ID, cluster)
	assert.NoErr(t, checkInErr)
	queryPairsMap := map[string]string{
		rest.CheckedInBeforeQueryStringKey: filter.CheckedInBefore.Format(data.StdTimestampFmt),
		rest.CheckedInAfterQueryStringKey:  filter.CheckedInAfter.Format(data.StdTimestampFmt),
		rest.CreatedBeforeQueryStringKey:   filter.CreatedBefore.Format(data.StdTimestampFmt),
		rest.CreatedAfterQueryStringKey:    filter.CreatedAfter.Format(data.StdTimestampFmt),
	}
	queryPairs := make([]string, len(queryPairsMap))
	i := 0
	for k, v := range queryPairsMap {
		queryPairs[i] = fmt.Sprintf("%s=%s", k, v)
		i++
	}

	route := urlPath(2, "clusters", "age")
	route += fmt.Sprintf("?%s", strings.Join(queryPairs, "&"))
	resp, err := httpGet(srv, route)
	assert.NoErr(t, err)
	defer resp.Body.Close()
	assert.Equal(t, resp.StatusCode, http.StatusOK, "response code")
	var respEnvelope struct {
		Data []types.Cluster `json:"data"`
	}
	assert.NoErr(t, json.NewDecoder(resp.Body).Decode(&respEnvelope))
	assert.Equal(t, len(respEnvelope.Data), 1, "length of the clusters list")
	assert.Equal(t, respEnvelope.Data[0].ID, cluster.ID, "returned cluster ID")
}
