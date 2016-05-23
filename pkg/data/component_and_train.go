package data

import (
	"fmt"

	"github.com/deis/workflow-manager-api/pkg/swagger/models"
)

// ComponentAndTrain represents a component and its train. It is used in functions such as
// the Version interface's MultiLatest function, where the caller must specify each component
// and its train at once.
type ComponentAndTrain struct {
	ComponentName string
	Train         string
}

func componentAndTrainFromComponentVersion(cv *models.ComponentVersion) *ComponentAndTrain {
	return &ComponentAndTrain{
		ComponentName: cv.Component.Name,
		Train:         cv.Version.Train,
	}
}

func (c ComponentAndTrain) String() string {
	return fmt.Sprintf("%s (%s)", c.ComponentName, c.Train)
}
