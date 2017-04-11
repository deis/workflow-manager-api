package data

import (
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq" // Pure Go Postgres driver for database/sql
)

// NewDB attempts to discover and connect to a postgres database
func NewDB() (*gorm.DB, error) {
	dataSourceName := "postgres://" + dBUser + ":" + dBPass + "@" + dBURL + "/" + dBName + "?sslmode=disable"
	db, err := gorm.Open("postgres", dataSourceName)
	if err != nil {
		log.Println("couldn't get a db connection!")
		return nil, err
	}
	if err := db.DB().Ping(); err != nil {
		log.Println("Failed to keep db connection alive")
		return nil, err
	}
	return db, nil
}
