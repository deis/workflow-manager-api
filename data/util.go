package data

import (
	"encoding/json"

	"github.com/deis/workflow-manager/types"
	sqlxTypes "github.com/jmoiron/sqlx/types"
)

func parseJSONComponent(jTxt sqlxTypes.JSONText) (types.ComponentVersion, error) {
	ret := new(types.ComponentVersion)
	if err := json.Unmarshal(jTxt, ret); err != nil {
		return types.ComponentVersion{}, err
	}
	return *ret, nil
}

func parseDBVersions(versions []VersionsTable) ([]types.ComponentVersion, error) {
	componentVersions := make([]types.ComponentVersion, len(versions))
	for i, version := range versions {
		versionData, err := version.data.MarshalJSON()
		if err != nil {
			return nil, err
		}
		component := types.Component{Name: version.componentName}
		version := types.Version{Train: version.train, Version: version.version, Released: version.releaseTimestamp.String(), Data: versionData}
		componentVersions[i] = types.ComponentVersion{Component: component, Version: version}
	}
	return componentVersions, nil
}
