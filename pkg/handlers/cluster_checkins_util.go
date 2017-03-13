package handlers

import (
	"time"

	"github.com/deis/workflow-manager-api/pkg/data"
	"github.com/deis/workflow-manager-api/pkg/swagger/restapi/operations"
	"github.com/deis/workflow-manager-api/rest"
)

func parseCheckinsQueryKeys(params operations.GetClusterCheckinsParams) (*data.ClusterCheckinsFilter, error) {
	createdAfter, err := time.Parse(
		data.StdTimestampFmt,
		params.CreatedAfter.String(),
	)
	if err != nil {
		return nil, errInvalidTimeFmt{key: rest.CreatedAfterQueryStringKey, err: err}
	}

	createdBefore, err := time.Parse(
		data.StdTimestampFmt,
		params.CreatedBefore.String(),
	)
	if err != nil {
		return nil, errInvalidTimeFmt{key: rest.CreatedBeforeQueryStringKey, err: err}
	}

	return data.NewClusterCheckinsFilter(createdAfter, createdBefore)
}
