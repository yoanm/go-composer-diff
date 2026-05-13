package compdiff

// ReqFile represents the structure of a composer.json file (composer requirement).
//
// It contains the explicit package requirements defined by the project, separated into
// regular dependencies and development-only dependencies.
type ReqFile struct {
	Require    map[string]string `json:"require,omitempty"`
	RequireDev map[string]string `json:"require-dev,omitempty"`
}

// LockFile represents the structure of a composer.lock file.
//
// It contains the complete list of installed packages and their versions, separated
// into regular packages and development-only packages. The lock file is the authoritative
// source for the exact versions and references used in a project.
type LockFile struct {
	Packages    []Package `json:"packages,omitempty"`
	PackagesDev []Package `json:"packages-dev,omitempty"`
}

// Package represents a single package entry in composer.lock.
//
// It contains all metadata about an installed package including its version,
// source information (e.g., repository reference), distribution information,
// documentation links, and whether it's been abandoned.
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
//
// It points to the specific commit or tag in the source repository that corresponds
// to the package version.
type VersionReference struct {
	Reference string `json:"reference,omitempty"`
}

// Support contains links to documentation and support.
//
// It provides URLs to various documentation and support resources for a package.
type Support struct {
	Wiki   string `json:"wiki,omitempty"`
	Docs   string `json:"docs,omitempty"`
	Source string `json:"source,omitempty"`
}
