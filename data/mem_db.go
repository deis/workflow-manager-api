package data

import (
	"database/sql"

	"github.com/mxk/go-sqlite/sqlite3"
)

const (
	memStr = ":memory:"
)

type memDB struct {
	db *sql.DB
}

func (m memDB) Get() (*sql.DB, error) {
	return m.db, nil
}

func NewMemDB() (*sql.DB, error) {
	db, err := sqlite3.Open(memStr)
	if err != nil {
		return nil, err
	}
	return &memDB{db:db}, nil
}
