package data

// DBType is a fmt.Stringer that indicates the type of database the workflow manager API should
// connect to.
type DBType string

// String is the fmt.Stringer interface implementation
func (d DBType) String() string {
	return string(d)
}

const (
	// RDSDBType indicates that the workflow manager API should connect to an RDS instance
	// defined by the Amazon AWS API
	RDSDBType DBType = "rds"
	// InClusterPostgresDBType indicates that the workflow manager API should connect to a postgres
	// database that's hosted in the same cluster, under the same namespace
	InClusterPostgresDBType DBType = "incluster"
)
