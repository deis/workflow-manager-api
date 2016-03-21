package data

import (
	"database/sql"

	"github.com/deis/workflow-manager/types"
)

// FakeCluster is a Cluster implementation, to be used for testing
type FakeCluster struct {
	GetCluster    types.Cluster
	GetErr        error
	SetCluster    types.Cluster
	SetErr        error
	CheckinResult sql.Result
	CheckinErr    error
}

// Get is the interface implementation
func (f FakeCluster) Get(*sql.DB, string) (types.Cluster, error) {
	return f.GetCluster, f.GetErr
}

// Set is the interface implementation
func (f FakeCluster) Set(*sql.DB, string, types.Cluster) (types.Cluster, error) {
	return f.SetCluster, f.SetErr
}

// Checkin is the interface implementation
func (f FakeCluster) Checkin(*sql.DB, string, types.Cluster) (sql.Result, error) {
	return f.CheckinResult, f.CheckinErr
}
