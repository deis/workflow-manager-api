package data

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq" // Pure Go Postgres driver for database/sql
)

func NewPostgresDB(host string, port int, username, password, dbName string) (*gorm.DB, error) {
	dbString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=require", username, password, host, port, dbName)
	return gorm.Open("postgres", dbString)
}
