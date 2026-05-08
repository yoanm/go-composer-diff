package compdiff

import (
	"fmt"
	"os"

	depsdiff "github.com/yoanm/go-deps-diff"
	"github.com/yoanm/go-deps-diff/contract"
)

func FileDiff(previous, current *FileInput) (contract.DiffMap, error) {
	previousReq, err := os.ReadFile(previous.Requirement)
	if err != nil {
		return nil, fmt.Errorf("reading previous requirement file: %w", err)
	}

	currentReq, err := os.ReadFile(current.Requirement)
	if err != nil {
		return nil, fmt.Errorf("reading current requirement file: %w", err)
	}

	previousLock, err := os.ReadFile(previous.Lock)
	if err != nil {
		return nil, fmt.Errorf("reading previous lock file: %w", err)
	}

	currentLock, err := os.ReadFile(current.Lock)
	if err != nil {
		return nil, fmt.Errorf("reading current lock file: %w", err)
	}

	return Diff(
		&Input{Lock: previousLock, Requirement: previousReq},
		&Input{Lock: currentLock, Requirement: currentReq},
	)
}

func Diff(previous, current *Input) (contract.DiffMap, error) {
	previousMap, err := BuildMapFromBytes(previous.Requirement, previous.Lock)
	if err != nil {
		return nil, fmt.Errorf("building previous package map: %w", err)
	}

	currentMap, err := BuildMapFromBytes(current.Requirement, current.Lock)
	if err != nil {
		return nil, fmt.Errorf("building current package map: %w", err)
	}

	//nolint:wrapcheck // Diff is the main purpose, so we don't want to wrap the error it can return
	return depsdiff.Diff(previousMap, currentMap)
}
