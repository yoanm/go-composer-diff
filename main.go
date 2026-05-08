package compdiff

import (
	"fmt"

	depsdiff "github.com/yoanm/go-deps-diff"
	"github.com/yoanm/go-deps-diff/contract"
)

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
