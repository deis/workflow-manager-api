package main

// Component type definition
type Component struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// Version type definition
type Version struct {
	Version string `json:"version"`
}

// ComponentVersion type definition
type ComponentVersion struct {
	Component
	Version
}

// Cluster type definition
type Cluster struct {
	Components []ComponentVersion `json:"components"`
}
