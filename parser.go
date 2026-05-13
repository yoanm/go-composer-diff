package compdiff

import (
	"encoding/json"
)

// ParseLock parses a composer.lock file from JSON bytes and returns a structured LockFile.
//
// It validates that the JSON is well-formed and contains the required fields
// (either 'packages' or 'packages-dev' or both).
//
// Parameters:
//   - data: The raw JSON content of a composer.lock file as bytes
//
// Returns:
//   - *LockFile: Pointer to the parsed lock file structure
//   - error: InvalidJSONError if JSON is malformed, InvalidFormatError if structure is invalid, or nil on success
//
// Example:
//
//	lockContent, err := os.ReadFile("composer.lock")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	lock, err := ParseLock(lockContent)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Installed %d packages\n", len(lock.Packages))
func ParseLock(data []byte) (*LockFile, error) {
	if len(data) == 0 {
		return nil, &InvalidFormatError{message: "empty input"}
	}

	var lock LockFile

	if err := json.Unmarshal(data, &lock); err != nil {
		return nil, &InvalidJSONError{err: err}
	}

	if lock.Packages == nil && lock.PackagesDev == nil {
		return nil, &InvalidFormatError{message: "missing 'packages' or 'packages-dev' fields"}
	}

	return &lock, nil
}

// ParseReq parses a composer.json file from JSON bytes and returns a structured ReqFile.
//
// It validates that the JSON is well-formed and contains the required fields
// (either 'require' or 'require-dev' or both).
//
// Parameters:
//   - data: The raw JSON content of a composer.json file as bytes
//
// Returns:
//   - *ReqFile: Pointer to the parsed requirement file structure
//   - error: InvalidJSONError if JSON is malformed, InvalidFormatError if structure is invalid, or nil on success
//
// Example:
//
//	reqContent, err := os.ReadFile("composer.json")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	req, err := ParseReq(reqContent)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Direct requirements: %d\n", len(req.Require))
func ParseReq(data []byte) (*ReqFile, error) {
	if len(data) == 0 {
		return nil, &InvalidFormatError{message: "empty input"}
	}

	var result ReqFile

	if err := json.Unmarshal(data, &result); err != nil {
		return nil, &InvalidJSONError{err: err}
	}

	if result.Require == nil && result.RequireDev == nil {
		return nil, &InvalidFormatError{message: "missing 'require' or 'require-dev' fields"}
	}

	return &result, nil
}

// InvalidJSONError indicates that the input JSON content is malformed and cannot be parsed.
//
// This error is returned by ParseLock or ParseReq when the JSON is syntactically invalid.
type InvalidJSONError struct {
	err error
}

// Error returns a string representation of the InvalidJSONError.
func (e InvalidJSONError) Error() string {
	return "invalid JSON: " + e.err.Error()
}

// InvalidFormatError indicates that the JSON is syntactically valid but doesn't match
// the expected Composer file structure.
//
// This error is returned when the JSON is missing required fields or has an unexpected schema.
type InvalidFormatError struct {
	message string
}

// Error returns a string representation of the InvalidFormatError.
func (e InvalidFormatError) Error() string {
	return "invalid format: " + e.message
}
