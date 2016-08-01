package config

import "github.com/kelseyhightower/envconfig"

// Specification config struct
type Specification struct {
	DoctorAuthUser string `envconfig:"DOCTOR_AUTH_USER" required:"true"`
	DoctorAuthPass string `envconfig:"DOCTOR_AUTH_PASS" required:"true"`
	DBInstance     string `envconfig:"DBINSTANCE"`
	DBUser         string `envconfig:"DBUSER"`
	DBPass         string `envconfig:"DBPASS"`
	RDSRegion      string `envconfig:"RDS_REGION"`
}

// Spec is an exportable variable that contains workflow manager config data
var Spec Specification

func init() {
	envconfig.Process("WORKFLOW_MANAGER_API", &Spec)
}
