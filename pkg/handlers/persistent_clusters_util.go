package handlers

import (
	"time"

	"github.com/deis/workflow-manager-api/pkg/data"
	"github.com/deis/workflow-manager-api/pkg/swagger/restapi/operations"
	"github.com/deis/workflow-manager-api/rest"
)

func parsePersistentClusterQueryKeys(params operations.GetPersistentClustersParams) (*data.PersistentClustersFilter, error) {
	epoch, err := time.Parse(
		data.StdTimestampFmt,
		params.Epoch.String(),
	)
	if err != nil {
		return nil, errInvalidTimeFmt{key: rest.EpochQueryStringKey, err: err}
	}

	timestamp, err := time.Parse(
		data.StdTimestampFmt,
		params.Timestamp.String(),
	)
	if err != nil {
		return nil, errInvalidTimeFmt{key: rest.TimestampQueryStringKey, err: err}
	}

	return data.NewPersistentClustersFilter(epoch, timestamp)
}
