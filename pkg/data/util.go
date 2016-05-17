package data

import (
	"encoding/json"
	"log"

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

// parseJSONCluster converts a JSON representation of a cluster
// to a ClusterStateful type
func parseJSONCluster(rawJSON []byte) (ClusterStateful, error) {
	var cluster ClusterStateful
	if err := json.Unmarshal(rawJSON, &cluster); err != nil {
		log.Print(err)
		return ClusterStateful{}, err
	}
	return cluster, nil
}
