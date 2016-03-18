package data

import (
	"database/sql"
)

// FakeCount is an implementation of the Count interface, to be used for tests
type FakeCount struct {
	Num int
	Err error
}

func (f FakeCount) Get(*sql.DB) (int, error) {
	return f.Num, f.Err
}
