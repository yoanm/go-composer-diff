// Package compdiff provides utilities for comparing PHP Composer lock files and identifying package differences.
package compdiff

import (
	"fmt"
	"os"

	depsdiff "github.com/yoanm/go-deps-diff"
	"github.com/yoanm/go-deps-diff/contract"
)

// FileDiff compares two Composer lock files from file paths and returns detailed information about package differences.
//
// It reads both the requirement (composer.json) and lock (composer.lock) files from disk
// and processes them to generate a comprehensive diff.
//
// Parameters:
//   - previous: FileInput containing paths to the previous composer.json and composer.lock files
//   - current: FileInput containing paths to the current composer.json and composer.lock files
//
// Returns:
//   - DiffMap: A map where keys are package names and values contain PackageChange information
//   - error: Non-nil if file reading or parsing fails
//
// Example:
//
//	changes, err := compdiff.FileDiff(
//	    &FileInput{Lock: "old/composer.lock", Requirement: "old/composer.json"},
//	    &FileInput{Lock: "new/composer.lock", Requirement: "new/composer.json"},
//	)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	for pkgName, change := range changes {
//	    fmt.Printf("%s: %s\n", pkgName, change.Operation.Name)
//	}
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

// Diff compares two sets of Composer requirement and lock file contents and returns detailed package differences.
//
// It takes parsed JSON content (not file paths) and processes them to identify all package additions,
// removals, and version changes. This function is useful when you already have the file contents
// in memory or from other sources.
//
// Parameters:
//   - previous: Input containing the previous composer.json and composer.lock file contents as byte slices
//   - current: Input containing the current composer.json and composer.lock file contents as byte slices
//
// Returns:
//   - DiffMap: A map where keys are package names and values contain PackageChange information
//   - error: Non-nil if JSON parsing or package map building fails
//
// Example:
//
//	previousReq, _ := os.ReadFile("old/composer.json")
//	previousLock, _ := os.ReadFile("old/composer.lock")
//	currentReq, _ := os.ReadFile("new/composer.json")
//	currentLock, _ := os.ReadFile("new/composer.lock")
//
//	changes, err := compdiff.Diff(
//	    &Input{Requirement: previousReq, Lock: previousLock},
//	    &Input{Requirement: currentReq, Lock: currentLock},
//	)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	for pkgName, change := range changes {
//	    fmt.Printf("%s: %s\n", pkgName, change.Operation.Name)
//	}
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
