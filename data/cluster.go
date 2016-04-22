package data

import (
	"database/sql"

	"github.com/deis/workflow-manager/types"
)

// Cluster is an interface for managing a persistent cluster record
type Cluster interface {
	Get(*sql.DB, string) (types.Cluster, error)
	Set(*sql.DB, string, types.Cluster) (types.Cluster, error)
	Checkin(*sql.DB, string, types.Cluster) (sql.Result, error)
	FilterByAge(*sql.DB, *ClusterAgeFilter) ([]types.Cluster, error)
}
