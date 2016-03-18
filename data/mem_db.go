package data

import (
	"database/sql"

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

func NewMemDB() (DB, error) {
	db, err := sql.Open(sqlite3Str, memStr)
	if err != nil {
		return nil, err
	}
	return &memDB{db: db}, nil
}
