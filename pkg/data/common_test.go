package data

import (
	"github.com/deis/workflow-manager-api/pkg/swagger/models"
)

const (
	clusterID     = "testcluster"
	componentName = "testcomponent"
	version       = "testversion"
	train         = "stable"
	released      = "2006-01-02T15:04:05Z"
)

var (
	versionData = models.VersionData{
		Description: "release notes",
	}
	componentDescription = "this is a component"
	updateAvailable      = "yup"
)
