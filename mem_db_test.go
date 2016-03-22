package main

import (
	"database/sql"

	"github.com/deis/workflow-manager-api/data"
	// this import registers "sqlite3" as a name you can provide to sql.Open
	_ "github.com/mxk/go-sqlite/sqlite3"
)

const (
	sqlite3Str = "sqlite3"
	memStr     = ":memory:"
)

// memDB is a sqlite in-memory data.DB implementation. note that it's in a *_test.go file since
// it can't be built without cgo. Since the production binary is built with CGO_ENABLED=0, cgo
// is off, but *_test.go files are omitted.
type memDB struct {
	db *sql.DB
}

func (m memDB) Get() (*sql.DB, error) {
	return m.db, nil
}

// newMemDB returns a DB implementation that stores all data in-memory. The initial database
// is empty, and is best used for testing.
func newMemDB() (data.DB, error) {
	db, err := sql.Open(sqlite3Str, memStr)
	if err != nil {
		return nil, err
	}
	return &memDB{db: db}, nil
}
