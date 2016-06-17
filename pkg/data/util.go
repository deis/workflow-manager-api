package data

import (
	"encoding/json"

	"github.com/deis/workflow-manager-api/pkg/swagger/models"
	sqlxTypes "github.com/jmoiron/sqlx/types"
)

func parseJSONComponent(jTxt sqlxTypes.JSONText) (models.ComponentVersion, error) {
	ret := new(models.ComponentVersion)
	if err := json.Unmarshal(jTxt, ret); err != nil {
		return models.ComponentVersion{}, err
	}
	return *ret, nil
}

// parseJSONCluster converts a JSON representation of a cluster
// to a ClusterStateful type
func parseJSONCluster(rawJSON []byte) (models.Cluster, error) {
	var cluster models.Cluster
	if err := json.Unmarshal(rawJSON, &cluster); err != nil {
		return models.Cluster{}, err
	}
	return cluster, nil
}

func parseJSONDoctor(rawJSON []byte) (models.DoctorInfo, error) {
	var doctor models.DoctorInfo
	if err := json.Unmarshal(rawJSON, &doctor); err != nil {
		return models.DoctorInfo{}, err
	}
	return doctor, nil
}
