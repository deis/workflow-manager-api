package data

import (
	"database/sql"
)

// DB is an interface for managing a DB instance
type DB interface {
	Get() (*sql.DB, error)
}
