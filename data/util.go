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
