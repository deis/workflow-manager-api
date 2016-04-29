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

func parseDBVersions(versions []versionsTable) ([]types.ComponentVersion, error) {
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

func parseDBVersion(version versionsTable) (types.ComponentVersion, error) {
	data := make(map[string]interface{})
	if err := json.Unmarshal(version.Data, &data); err != nil {
		return types.ComponentVersion{}, err
	}
	return types.ComponentVersion{
		Component: types.Component{
			Name: version.ComponentName,
		},
		Version: types.Version{
			Train:    version.Train,
			Version:  version.Version,
			Released: version.ReleaseTimestamp.String(),
			Data:     data,
		},
	}, nil
}
