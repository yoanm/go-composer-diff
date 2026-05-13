package compdiff

import "github.com/yoanm/go-deps-diff/contract"

// ComposerPackageWrapper wraps Composer package information and implements the contract.PkgWrapper interface.
//
// It provides access to package metadata including name, version, status, and various classifications
// (root requirement, dev-only, abandoned status, etc.). This type is the main way to interact with
// package information when working with the diff results.
type ComposerPackageWrapper struct {
	name                 string
	isAbandoned          bool
	version              contract.PkgVersion
	link                 string
	isDevOnly            bool // true if only in lock file "packages-dev" section (dev-only dependency)
	isRootRequirement    bool // true if exists in requirement file "require" section
	isRootDevRequirement bool // true if exists in requirement file "require-dev" section
}

// Ensure that *ComposerPackageWrapper implements contract.PkgWrapper.
var _ contract.PkgWrapper = (*ComposerPackageWrapper)(nil)

// GetName returns the package name (e.g., "vendor/package-name").
func (w *ComposerPackageWrapper) GetName() string {
	return w.name
}

// IsAbandoned returns whether this package has been marked as abandoned by its maintainer.
func (w *ComposerPackageWrapper) IsAbandoned() bool {
	return w.isAbandoned
}

// GetVersion returns the package version information including raw version, label,
// and parsed semantic version if applicable.
func (w *ComposerPackageWrapper) GetVersion() contract.PkgVersion {
	return w.version
}

// GetLink returns a URL link to the package (e.g., homepage or documentation).
//
// Returns an empty string if no link is available.
func (w *ComposerPackageWrapper) GetLink() string {
	return w.link
}

// IsDevOnly returns whether this package is only required for development (not in production).
//
// A package is dev-only if it appears only in the packages-dev section of composer.lock.
// However, a dev-only package can be a root requirement, a root dev requirement or a transitive dependency. If
// you are looking for requirement vs dev requirement see IsRootRequirement and IsRootDevRequirement
func (w *ComposerPackageWrapper) IsDevOnly() bool {
	return w.isDevOnly
}

// IsRootRequirement returns whether this package is explicitly required in composer.json.
//
// Returns true if the package is listed in the "require" section of composer.json.
func (w *ComposerPackageWrapper) IsRootRequirement() bool {
	return w.isRootRequirement
}

// IsRootDevRequirement returns whether this package is explicitly required for development in composer.json.
//
// Returns true if the package is listed in the "require-dev" section of composer.json.
func (w *ComposerPackageWrapper) IsRootDevRequirement() bool {
	return w.isRootDevRequirement
}
