package compdiff

import (
	"encoding/json"
)

// ParseLock parses a composer.lock file from JSON bytes.
func ParseLock(data []byte) (*LockFile, error) {
	if len(data) == 0 {
		return nil, &InvalidFormatError{message: "empty input"}
	}

	var lock LockFile

	err := json.Unmarshal(data, &lock)
	if err != nil {
		return nil, &InvalidJSONError{err: err}
	}

	if lock.Packages == nil && lock.PackagesDev == nil {
		return nil, &InvalidFormatError{message: "missing 'packages' or 'packages-dev' fields"}
	}

	return &lock, nil
}

// ParseReq parses a composer.json file from JSON bytes.
func ParseReq(data []byte) (*ReqFile, error) {
	if len(data) == 0 {
		return nil, &InvalidFormatError{message: "empty input"}
	}

	var result ReqFile

	err := json.Unmarshal(data, &result)
	if err != nil {
		return nil, &InvalidJSONError{err: err}
	}

	if result.Require == nil && result.RequireDev == nil {
		return nil, &InvalidFormatError{message: "missing 'require' or 'require-dev' fields"}
	}

	return &result, nil
}

// InvalidJSONError indicates the input JSON is malformed.
type InvalidJSONError struct {
	err error
}

func (e InvalidJSONError) Error() string {
	return "invalid JSON: " + e.err.Error()
}

// InvalidFormatError indicates the JSON is valid but doesn't match expected structure.
type InvalidFormatError struct {
	message string
}

func (e InvalidFormatError) Error() string {
	return "invalid format: " + e.message
}
