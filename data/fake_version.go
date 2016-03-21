package data

import (
	"database/sql"

	"github.com/deis/workflow-manager/types"
)

// FakeVersion is an implementation of Version, for use in testing
type FakeVersion struct {
	GetComponentVersion types.ComponentVersion
	GetErr              error
	SetComponentVersion types.ComponentVersion
	SetErr              error
}

// Get is the interface implementation
func (f FakeVersion) Get(*sql.DB, string) (types.ComponentVersion, error) {
	return f.GetComponentVersion, f.GetErr
}

// Set is the interface implementation
func (f FakeVersion) Set(*sql.DB, string, types.ComponentVersion) (types.ComponentVersion, error) {
	return f.SetComponentVersion, f.SetErr
}
