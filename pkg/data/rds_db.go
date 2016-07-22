package data

import (
	"log"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq" // Pure Go Postgres driver for database/sql
)

// rdsDB is an implementation of the DB interface
type rdsDB struct {
	config   *aws.Config
	user     string
	pass     string
	flavor   string
	instance *string
}

// NewRDSDB is a constructor that returns an instance of a DB
// that can connect to an Amazon RDS instance
func NewRDSDB(region string, user string, pass string, flavor string, instance string) DB {
	return &rdsDB{
		config:   &aws.Config{Region: aws.String(region)},
		user:     user,
		pass:     pass,
		flavor:   flavor,
		instance: &instance,
	}
}

func (r rdsDB) Get() (*gorm.DB, error) {
	svc := rds.New(session.New(), r.config)
	params := rds.DescribeDBInstancesInput{DBInstanceIdentifier: r.instance}
	resp, err := svc.DescribeDBInstances(&params)
	if err != nil {
		return nil, err
	}
	if len(resp.DBInstances) > 1 {
		log.Printf("more than one database instance returned for %s, using the 1st one", *r.instance)
	}
	instance := resp.DBInstances[0]
	url := *instance.Endpoint.Address + ":" + strconv.FormatInt(*instance.Endpoint.Port, 10)
	dataSourceName := r.flavor + "://" + r.user + ":" + r.pass + "@" + url + "/" + *instance.DBName + "?sslmode=require"
	db, err := gorm.Open(r.flavor, dataSourceName)
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
