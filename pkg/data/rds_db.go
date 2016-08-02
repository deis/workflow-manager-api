package data

import (
	"log"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/deis/workflow-manager-api/config"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq" // Pure Go Postgres driver for database/sql
)

var (
	rDSRegion = config.Spec.RDSRegion
)

// NewRDSDB attempts to discover and connect to a postgres database managed by Amazon RDS
func NewRDSDB() (*gorm.DB, error) {
	db, err := getRDSDB()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func getRDSSession() *rds.RDS {
	return rds.New(session.New(), &aws.Config{Region: aws.String(rDSRegion)})
}

func getRDSDB() (*gorm.DB, error) {
	svc := getRDSSession()
	dbInstanceIdentifier := new(string)
	dbInstanceIdentifier = &dBInstance
	params := rds.DescribeDBInstancesInput{DBInstanceIdentifier: dbInstanceIdentifier}
	resp, err := svc.DescribeDBInstances(&params)
	if err != nil {
		return nil, err
	}
	if len(resp.DBInstances) > 1 {
		log.Printf("more than one database instance returned for %s, using the 1st one", dBInstance)
	}
	instance := resp.DBInstances[0]
	url := *instance.Endpoint.Address + ":" + strconv.FormatInt(*instance.Endpoint.Port, 10)
	dataSourceName := "postgres://" + dBUser + ":" + dBPass + "@" + url + "/" + *instance.DBName + "?sslmode=require"
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
