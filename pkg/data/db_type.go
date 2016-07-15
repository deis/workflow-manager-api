package data

import (
	"fmt"
)

const (
	// RDSDBType indicates that the workflow manager API should connect to an RDS instance
	// defined by the Amazon AWS API
	RDSDBType DBType = "rds"
	// InClusterPostgresDBType indicates that the workflow manager API should connect to a postgres
	// database that's hosted in the same cluster, under the same namespace
	InClusterPostgresDBType DBType = "incluster"
)

// ErrInvalidDBType is the error returned when a string couldn't be parsed into a DBType
type ErrInvalidDBType string

// Error is the error interface implementation
func (e ErrInvalidDBType) Error() string {
	return fmt.Sprintf("%s is an invalid DB type", string(e))
}

// DBType is a fmt.Stringer that indicates the type of database the workflow manager API should
// connect to.
type DBType string

// DBTypeFromString parses a DBType from s. Returns an empty DBType and an ErrInvalidDBType
// error if s doesn't represent a known DB type
func DBTypeFromString(s string) (DBType, error) {
	switch s {
	case RDSDBType.String():
		return RDSDBType, nil
	case InClusterPostgresDBType.String():
		return InClusterPostgresDBType, nil
	default:
		return DBType(""), ErrInvalidDBType(s)
	}
}

// String is the fmt.Stringer interface implementation
func (d DBType) String() string {
	return string(d)
}
