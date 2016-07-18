package data

import (
	"github.com/kelseyhightower/envconfig"
)

const (
	envconfigAppName = "wfm"
)

// InClusterPostgresConfig is the envconfig-compatible struct for connecting to an in-cluster postgres DB
type InClusterPostgresConfig struct {
	Host         string `envconfig:"DEIS_WFMPOSTGRES_SERVICE_HOST" required:"true"`
	Post         int    `envconfig:"DEIS_WFMPOSTGRES_SERVICE_PORT" required:"true"`
	UsernameFile string `envconfig:"DEIS_WFMPOSTGRES_USERNAME_FILE" default:"/var/run/secrets/postgres/auth/username"`
	PasswordFile string `envconfig:"DEIS_WFMPOSTGRES_PASSWORD_FILE" default:"/var/run/secrets/postgres/auth/password"`
	DBName       string `envconfig:"DEIS_WFMPOSTGRES_DB_NAME" default:"wfm"`
}

// GetInClusterPostgresConfig parses and returns the configuration for the in-cluster postgres
// database. returns a non-nil error if any parse issues.
func GetInClusterPostgresConfig() (*InClusterPostgresConfig, error) {
	ret := new(InClusterPostgresConfig)
	if err := envconfig.Process(envconfigAppName, ret); err != nil {
		return nil, err
	}
	return ret, nil
}
