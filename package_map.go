package compdiff

import (
	"fmt"

	"github.com/yoanm/go-deps-diff/contract"
)

func BuildMapFromBytes(reqContent, lockContent []byte) (contract.PackageMap, error) {
	reqData, err := ParseReq(reqContent)
	if err != nil {
		return nil, fmt.Errorf("parsing requirement file content: %w", err)
	}

	lockData, err := ParseLock(lockContent)
	if err != nil {
		return nil, fmt.Errorf("parsing lock file content: %w", err)
	}

	return BuildMap(reqData, lockData)
}

// BuildMap creates an efficient lookup map for composer packages.
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
