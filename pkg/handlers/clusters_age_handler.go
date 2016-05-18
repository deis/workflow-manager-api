package handlers

import (
	"log"
	"net/http"

	"github.com/deis/workflow-manager-api/pkg/data"
	"github.com/jinzhu/gorm"
)

// ClustersAge is the handler for the GET /{apiVersion}/clusters/age endpoint
func ClustersAge(db *gorm.DB) http.Handler {
	type clustersAgeResp struct {
		Data []data.ClusterStateful `json:"data"`
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clusterAgeFilter, err := parseAgeQueryKeys(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		clusters, err := data.FilterClustersByAge(db, clusterAgeFilter)
		if err != nil {
			log.Printf("Error filtering clusters by age (%s)", err)
			http.Error(w, "database error", http.StatusInternalServerError)
			return
		}
		ret := clustersAgeResp{Data: clusters}
		if err := writeJSON(w, ret); err != nil {
			log.Printf("ClustersAge json marshal error (%s)", err)
		}
	})
}
