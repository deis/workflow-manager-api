package main

import (
	"log"
	"net/http"
	"sync"
)

// package-level constants
const (
	listenPort = "8443"
)

var memoClusters = make(map[string]Cluster)
var latestVersions = make(map[string]Version)
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
func getAll() map[string]Cluster {
	return memoClusters
}

// get all versions data
func getVersions() map[string]Version {
	return latestVersions
}

// make a new Cluster struct
func newCluster() Cluster {
	return Cluster{}
}

// make a new ComponentVersion struct
func newComponentVersion() ComponentVersion {
	return ComponentVersion{}
}

// get a cluster record, returns a new Cluster that the caller can optionally use
func getCluster(id string) (Cluster, bool) {
	cluster, ok := memoClusters[id]
	if !ok {
		return newCluster(), false
	}
	return cluster, true
}

// get a component version record, returns a new ComponentVersion that the caller can optionally use
func getComponentVersion(name string) (ComponentVersion, bool) {
	version, ok := latestVersions[name]
	if !ok {
		return newComponentVersion(), false
	}
	componentVersion := ComponentVersion{Component: Component{Name: name}, Version: version}
	return componentVersion, true
}

// cluster record set'er
func setCluster(id string, c Cluster) Cluster {
	memoClusters[id] = c
	return memoClusters[id]
}

// component version record set'er
func setLatestVersion(cV ComponentVersion) Version {
	latestVersions[cV.Name] = cV.Version
	return latestVersions[cV.Name]
}
