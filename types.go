package compdiff

// ReqFile represents the structure of a composer.json file (composer requirement).
type ReqFile struct {
	Require    map[string]string `json:"require,omitempty"`
	RequireDev map[string]string `json:"require-dev,omitempty"`
}

// LockFile represents the structure of a composer.lock file.
type LockFile struct {
	Packages    []Package `json:"packages,omitempty"`
	PackagesDev []Package `json:"packages-dev,omitempty"`
}

// Package represents a single package entry in composer.lock.
type Package struct {
	Name      string            `json:"name"`
	Version   string            `json:"version"`
	Source    *VersionReference `json:"source,omitempty"`
	Dist      *VersionReference `json:"dist,omitempty"`
	Support   *Support          `json:"support,omitempty"`
	Homepage  string            `json:"homepage,omitempty"`
	Abandoned interface{}       `json:"abandoned,omitempty"` // Can be bool or string
}

// VersionReference contains the reference (commit hash or tag).
type VersionReference struct {
	Reference string `json:"reference,omitempty"`
}

// Support contains links to documentation and support.
type Support struct {
	Wiki   string `json:"wiki,omitempty"`
	Docs   string `json:"docs,omitempty"`
	Source string `json:"source,omitempty"`
}
