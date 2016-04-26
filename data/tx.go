package data

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

type errCommitTx struct {
	origErr error
}

func (e errCommitTx) Error() string {
	return fmt.Sprintf("Error committing transaction (%s)", e.origErr)
}

type errRollbackTx struct {
	txErr   error
	origErr error
}

func (e errRollbackTx) Error() string {
	return fmt.Sprintf("Error rolling back transaction after error (%s) (%s)", e.origErr, e.txErr)
}

// inTx runs fn inside a transaction. if fn failed, tries to rollback. otherwise tries to commit
//
// on rollback, returns nil and fn's error if the rollback succeeds. otehrwise returns nil and errRollbackTx
// on commit, returns errCommitTx if the commit fails, otherwise returns the resultant DB and nil
func inTx(db *gorm.DB, fn func(*gorm.DB) error) (*gorm.DB, error) {
	tx := db.Begin()
	err := fn(tx)
	if err != nil {
		rbDB := tx.Rollback()
		if rbDB.Error != nil {
			return nil, errRollbackTx{txErr: rbDB.Error, origErr: err}
		}
		return nil, err
	}
	cDB := tx.Commit()
	if cDB.Error != nil {
		return nil, errCommitTx{origErr: cDB.Error}
	}
	return cDB, nil
}
