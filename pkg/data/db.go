package data

import (
	"github.com/jinzhu/gorm"
)

// DB is an interface for managing a database
type DB interface {
	Get() (*gorm.DB, error)
}
