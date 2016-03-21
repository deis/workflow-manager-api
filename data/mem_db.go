package data

import (
	"database/sql"

	// this import registers "sqlite3" as a name you can provide to sql.Open
	_ "github.com/mxk/go-sqlite/sqlite3"
)

const (
	sqlite3Str = "sqlite3"
	memStr     = ":memory:"
)

type memDB struct {
	db *sql.DB
}

func (m memDB) Get() (*sql.DB, error) {
	return m.db, nil
}

// NewMemDB returns a DB implementation that stores all data in-memory. The initial database
// is empty, and is best used for testing.
func NewMemDB() (DB, error) {
	db, err := sql.Open(sqlite3Str, memStr)
	if err != nil {
		return nil, err
	}
	return &memDB{db: db}, nil
}
