package compdiff

import (
	"fmt"
	"os"

	depsdiff "github.com/yoanm/go-deps-diff"
	"github.com/yoanm/go-deps-diff/contract"
)

func FileDiff(previous, current *FileInput) (contract.DiffMap, error) {
	var (
		previousReq, previousLock []byte
		currentReq, currentLock   []byte
		err                       error
	)

	if previousReq, previousLock, err = readFiles(previous); err != nil {
		return nil, fmt.Errorf("reading previous files: %w", err)
	}

	if currentReq, currentLock, err = readFiles(current); err != nil {
		return nil, fmt.Errorf("reading current files: %w", err)
	}

	return Diff(
		&Input{Lock: previousLock, Requirement: previousReq},
		&Input{Lock: currentLock, Requirement: currentReq},
	)
}

func Diff(previous, current *Input) (contract.DiffMap, error) {
	var (
		previousMap, currentMap contract.PackageMap
		err                     error
	)

	if previousMap, err = BuildMapFromBytes(previous.Requirement, previous.Lock); err != nil {
		return nil, fmt.Errorf("building previous package map: %w", err)
	}
	if currentMap, err = BuildMapFromBytes(current.Requirement, current.Lock); err != nil {
		return nil, fmt.Errorf("building current package map: %w", err)
	}

	return depsdiff.Diff(previousMap, currentMap), nil
}

func readFiles(files *FileInput) ([]byte, []byte, error) {
	var (
		lockContent, reqContent []byte
		err                     error
	)

	if reqContent, err = os.ReadFile(files.Requirement); err != nil {
		return nil, nil, fmt.Errorf("reading requirement file: %w", err)
	}

	if lockContent, err = os.ReadFile(files.Lock); err != nil {
		return nil, nil, fmt.Errorf("reading lock file: %w", err)
	}

	return reqContent, lockContent, nil
}
