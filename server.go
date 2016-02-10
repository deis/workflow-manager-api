package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/deis/workflow-manager/types"
)

// package-level constants
const (
	listenPort = "8443"
)

var memoClusters = make(map[string]types.Cluster)
var latestVersions = make(map[string]types.Version)
var mu sync.Mutex

// Main opens up a TLS listening port
func main() {
	r := getRoutes()
	// Bind to a port and pass our router in
	err := http.ListenAndServeTLS(":"+listenPort, "server.pem", "server.key", r)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

// get all cluster data
func getAll() map[string]types.Cluster {
	return memoClusters
}

// get all versions data
func getVersions() map[string]types.Version {
	return latestVersions
}

// make a new Cluster struct
func newCluster() types.Cluster {
	return types.Cluster{}
}

// make a new ComponentVersion struct
func newComponentVersion() types.ComponentVersion {
	return types.ComponentVersion{}
}

// get a cluster record, returns a new Cluster that the caller can optionally use
func getCluster(id string) (types.Cluster, bool) {
	cluster, ok := memoClusters[id]
	if !ok {
		return newCluster(), false
	}
	return cluster, true
}

// get a component version record, returns a new ComponentVersion that the caller can optionally use
func getComponentVersion(name string) (types.ComponentVersion, bool) {
	version, ok := latestVersions[name]
	if !ok {
		return newComponentVersion(), false
	}
	componentVersion := types.ComponentVersion{Component: types.Component{Name: name}, Version: version}
	return componentVersion, true
}

// cluster record set'er
func setCluster(id string, c types.Cluster) types.Cluster {
	memoClusters[id] = c
	return memoClusters[id]
}

// component version record set'er
func setLatestVersion(cV types.ComponentVersion) types.Version {
	latestVersions[cV.Name] = cV.Version
	return latestVersions[cV.Name]
}
