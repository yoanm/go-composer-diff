package compdiff

import (
	"fmt"

	"github.com/yoanm/go-deps-diff/contract"
)

// BuildMapFromBytes creates an efficient lookup map for Composer packages from raw JSON content.
//
// This function parses the composer.json (requirement) and composer.lock (lock) files
// and builds a PackageMap for use with the depsdiff comparison functions.
//
// Parameters:
//   - reqContent: The raw JSON content of a composer.json file
//   - lockContent: The raw JSON content of a composer.lock file
//
// Returns:
//   - contract.PackageMap: A map of package names to their wrapper objects
//   - error: Non-nil if parsing or building fails
//
// Example:
//
//	reqContent, _ := os.ReadFile("composer.json")
//	lockContent, _ := os.ReadFile("composer.lock")
//	packageMap, err := compdiff.BuildMapFromBytes(reqContent, lockContent)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	for name, pkg := range packageMap {
//	    fmt.Printf("%s: %s\n", name, pkg.GetVersion().Label)
//	}
func BuildMapFromBytes(reqContent, lockContent []byte) (contract.PackageMap, error) {
	var (
		reqData  *ReqFile
		lockData *LockFile
		err      error
	)

	if reqData, err = ParseReq(reqContent); err != nil {
		return nil, fmt.Errorf("parsing requirement file content: %w", err)
	}

	if lockData, err = ParseLock(lockContent); err != nil {
		return nil, fmt.Errorf("parsing lock file content: %w", err)
	}

	return BuildMap(reqData, lockData)
}

// BuildMap creates an efficient lookup map for Composer packages from structured ReqFile and LockFile.
//
// This function takes already-parsed ReqFile and LockFile structures and builds a PackageMap.
// Both regular and dev packages are included in the map. Use this when you have already parsed
// the files, or use BuildMapFromBytes to parse and build in one step.
//
// Parameters:
//   - reqData: Parsed composer.json structure containing explicit package requirements
//   - lockData: Parsed composer.lock structure containing installed package information
//
// Returns:
//   - contract.PackageMap: A map of package names to their wrapper objects
//   - error: Always returns nil (errors would have occurred during parsing)
//
// Example:
//
//	req, _ := ParseReq(reqContent)
//	lock, _ := ParseLock(lockContent)
//	packageMap, _ := compdiff.BuildMap(req, lock)
//	for name, pkg := range packageMap {
//	    fmt.Printf("%s (dev: %v): %s\n", name, pkg.IsDevOnly(), pkg.GetVersion().Raw)
//	}
func BuildMap(reqData *ReqFile, lockData *LockFile) (contract.PackageMap, error) {
	packageMap := make(map[string]contract.PkgWrapper)

	// Add regular packages
	if lockData.Packages != nil {
		for i := range lockData.Packages {
			pkg := &lockData.Packages[i]
			packageMap[pkg.Name] = createWrapper(pkg, reqData, false)
		}
	}

	// Add dev packages
	if lockData.PackagesDev != nil {
		for i := range lockData.PackagesDev {
			pkg := &lockData.PackagesDev[i]
			// Avoid overwriting regular packages with dev packages (unlikely but just in case)
			if _, exists := packageMap[pkg.Name]; !exists {
				packageMap[pkg.Name] = createWrapper(pkg, reqData, true)
			}
		}
	}

	return packageMap, nil
}
