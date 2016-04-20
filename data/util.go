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
		cver, err := parseDBVersion(version)
		if err != nil {
			return nil, err
		}
		componentVersions[i] = cver
	}
	return componentVersions, nil
}

func parseDBVersion(version VersionsTable) (types.ComponentVersion, error) {
	versionData, err := version.data.MarshalJSON()
	if err != nil {
		return types.ComponentVersion{}, err
	}
	return types.ComponentVersion{
		Component: types.Component{
			Name: version.componentName,
		},
		Version: types.Version{
			Train:    version.train,
			Version:  version.version,
			Released: version.releaseTimestamp.String(),
			Data:     versionData,
		},
	}, nil
}
