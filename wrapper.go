package compdiff

import (
	"fmt"
	"log/slog"

	"github.com/yoanm/go-deps-diff/contract"
	"github.com/yoanm/go-deps-diff/semver"
)

func createWrapper(pkg *Package, reqData *ReqFile, isDevOnly bool) *ComposerPackageWrapper {
	_, isRootReq := reqData.Require[pkg.Name]
	_, isRootDevReq := reqData.RequireDev[pkg.Name]

	return &ComposerPackageWrapper{
		name:                 pkg.Name,
		isAbandoned:          isAbandonedPkg(pkg),
		version:              parsePkgVersion(pkg),
		link:                 getPkgLink(pkg),
		isDevOnly:            isDevOnly,
		isRootRequirement:    isRootReq,
		isRootDevRequirement: isRootDevReq,
	}
}

// isAbandonedPkg safely extracts the abandoned status from the composer package
// Returns true if the field is explicitly set to true (boolean or string "true").
func isAbandonedPkg(pkg *Package) bool {
	if pkg.Abandoned == nil {
		return false
	}

	switch v := pkg.Abandoned.(type) {
	case bool:
		return v
	case string:
		// Non-empty string is considered abandoned (e.g. an url to the replacement), "true" is explicitly true
		return v == "true" || v != ""
	default:
		return false
	}
}

const shortRefLength = 7

// parsePkgVersion parses a version string into a PkgVersionOld.
func parsePkgVersion(pkg *Package) contract.PkgVersion {
	var (
		sVersion *contract.Semver
		err      error
	)

	if sVersion, err = semver.Parse(pkg.Version); err != nil {
		if semver.IsValid(pkg.Version) { // Something went wrong during parsing then (unlikely to happen though)
			slog.Warn(fmt.Errorf("error while parsing semver version %s: %w", pkg.Version, err).Error())
		} else if ref := getPkgRef(pkg); ref != "" { // Not semver - check if there is a commit reference
			shortRef := ref
			if len(shortRef) > shortRefLength {
				shortRef = shortRef[:shortRefLength]
			}

			return contract.PkgVersion{
				Raw:    ref,
				Label:  pkg.Version + "#" + shortRef,
				Semver: nil,
			}
		}
	}

	return contract.PkgVersion{
		Raw:    pkg.Version,
		Label:  pkg.Version,
		Semver: sVersion,
	}
}

// getPkgRef extracts the commit hash from a package
// Prefers dist.reference over source.reference.
func getPkgRef(pkg *Package) string {
	switch {
	case pkg.Dist != nil && pkg.Dist.Reference != "":
		return pkg.Dist.Reference
	case pkg.Source != nil && pkg.Source.Reference != "":
		return pkg.Source.Reference
	default:
		return ""
	}
}

// getPkgLink extracts the best available link from a package
// Priority: wiki -> docs -> source -> homepage.
func getPkgLink(pkg *Package) string {
	switch {
	case pkg.Support != nil:
		switch {
		case pkg.Support.Wiki != "":
			return pkg.Support.Wiki
		case pkg.Support.Docs != "":
			return pkg.Support.Docs
		case pkg.Support.Source != "":
			return pkg.Support.Source
		}
	case pkg.Homepage != "":
		return pkg.Homepage
	}

	return ""
}
