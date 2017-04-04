package config

import "github.com/kelseyhightower/envconfig"

// Specification config struct
type Specification struct {
	DoctorAuthUser string `envconfig:"DOCTOR_AUTH_USER" required:"true"`
	DoctorAuthPass string `envconfig:"DOCTOR_AUTH_PASS" required:"true"`
	DBUser         string `envconfig:"DBUSER" required:"true"`
	DBPass         string `envconfig:"DBPASS" required:"true"`
	DBURL          string `envconfig:"DBURL" required:"true"`
	DBName         string `envconfig:"DBNAME" required:"true"`
}

// Spec is an exportable variable that contains workflow manager config data
var Spec Specification

func init() {
	envconfig.Process("WORKFLOW_MANAGER_API", &Spec)
}
