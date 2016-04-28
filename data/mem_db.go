// +build testonly

package data

import (
	"database/sql"
	// this import registers "sqlite3" as a name you can provide to sql.Open
	_ "github.com/mxk/go-sqlite/sqlite3"
)

// NOTE: we are using a build tag to conditionally build this file (see the comment at the top
// of this file) because we want to use it in tests, but also want to build with CGO_ENABLED=0.
// When CGO_ENABLED is set to 0, C source files are not allowed in the build, and since
// sqlite uses C source files, we can't build this.
//
// However, we can run tests with CGO_ENABLED=1, so as long as we run tests with -tags testonly,
// this file will be included and we we're good to go.
//
// Why not put this in a test file (i.e. mem_db_test.go), you ask? In short, symbols in test
// files are not exported outside of the package they live in, so the in-memory DB logic would
// have to be copied everywhere it's needed. Right now, that's in this package as well as 'main'.
// Not the worst thing, but this is cleaner.
//
// For more information on build tags, see
// http://dave.cheney.net/2013/10/12/how-to-use-conditional-compilation-with-the-go-build-tool

const (
	sqlite3Str = "sqlite3"
	memStr     = ":memory:"
)

// NewMemDB returns a DB implementation that stores all data in-memory. The initial database
// is empty, and is best used for testing.
func NewMemDB() (*sql.DB, error) {
	db, err := sql.Open(sqlite3Str, memStr)
	if err != nil {
		return nil, err
	}
	return db, nil
}
