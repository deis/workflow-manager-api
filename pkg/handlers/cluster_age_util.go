package handlers

import (
	"fmt"
	"time"

	"github.com/deis/workflow-manager-api/pkg/data"
	"github.com/deis/workflow-manager-api/pkg/swagger/restapi/operations"
	"github.com/deis/workflow-manager-api/rest"
)

type errInvalidTimeFmt struct {
	key string
	err error
}

func (e errInvalidTimeFmt) Error() string {
	return fmt.Sprintf("%s is an invalid timestamp (%s)", e.key, e.err)
}

func parseAgeQueryKeys(params operations.GetClustersByAgeParams) (*data.ClusterAgeFilter, error) {
	checkedInBefore, err := time.Parse(
		data.StdTimestampFmt,
		params.CheckedInBefore.String(),
	)
	if err != nil {
		return nil, errInvalidTimeFmt{key: rest.CheckedInBeforeQueryStringKey, err: err}
	}

	checkedInAfter, err := time.Parse(
		data.StdTimestampFmt,
		params.CheckedInAfter.String(),
	)
	if err != nil {
		return nil, errInvalidTimeFmt{key: rest.CheckedInAfterQueryStringKey, err: err}
	}

	createdBefore, err := time.Parse(
		data.StdTimestampFmt,
		params.CreatedBefore.String(),
	)
	if err != nil {
		return nil, errInvalidTimeFmt{key: rest.CreatedBeforeQueryStringKey, err: err}
	}

	createdAfter, err := time.Parse(
		data.StdTimestampFmt,
		params.CreatedAfter.String(),
	)
	if err != nil {
		return nil, errInvalidTimeFmt{key: rest.CreatedAfterQueryStringKey, err: err}
	}

	return data.NewClusterAgeFilter(checkedInBefore, checkedInAfter, createdBefore, createdAfter)
}
