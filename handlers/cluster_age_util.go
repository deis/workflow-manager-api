package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/deis/workflow-manager-api/data"
	"github.com/deis/workflow-manager-api/rest"
)

type errInvalidTimeFmt struct {
	key string
	err error
}

func (e errInvalidTimeFmt) Error() string {
	return fmt.Sprintf("%s is an invalid timestamp (%s)", e.key, e.err)
}

func parseAgeQueryKeys(r *http.Request) (*data.ClusterAgeFilter, error) {
	checkedInBefore, err := time.Parse(
		data.StdTimestampFmt,
		r.URL.Query().Get(rest.CheckedInBeforeQueryStringKey),
	)
	if err != nil {
		return nil, errInvalidTimeFmt{key: rest.CheckedInBeforeQueryStringKey, err: err}
	}

	checkedInAfter, err := time.Parse(
		data.StdTimestampFmt,
		r.URL.Query().Get(rest.CheckedInAfterQueryStringKey),
	)
	if err != nil {
		return nil, errInvalidTimeFmt{key: rest.CheckedInAfterQueryStringKey, err: err}
	}

	createdBefore, err := time.Parse(
		data.StdTimestampFmt,
		r.URL.Query().Get(rest.CreatedBeforeQueryStringKey),
	)
	if err != nil {
		return nil, errInvalidTimeFmt{key: rest.CreatedBeforeQueryStringKey, err: err}
	}

	createdAfter, err := time.Parse(
		data.StdTimestampFmt,
		r.URL.Query().Get(rest.CreatedAfterQueryStringKey),
	)
	if err != nil {
		return nil, errInvalidTimeFmt{key: rest.CreatedAfterQueryStringKey, err: err}
	}

	return data.NewClusterAgeFilter(checkedInBefore, checkedInAfter, createdBefore, createdAfter)
}
