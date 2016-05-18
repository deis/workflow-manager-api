package data

import (
	"encoding/json"
	"github.com/deis/workflow-manager-api/pkg/swagger/models"
	"github.com/deis/workflow-manager/types"
	sqlxTypes "github.com/jmoiron/sqlx/types"
	"log"
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
func parseJSONCluster(rawJSON []byte) (models.Cluster, error) {
	var cluster models.Cluster
	if err := json.Unmarshal(rawJSON, &cluster); err != nil {
		log.Print(err)
		return models.Cluster{}, err
	}
	return cluster, nil
}
