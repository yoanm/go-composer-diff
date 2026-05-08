package compdiff_test

import (
	"fmt"
	compdiff "github.com/yoanm/go-composer-diff"
	"testing"
)

func TestParseLock(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		data    []byte
		checkFn func(lock *compdiff.LockFile) bool
	}{
		{
			name: "valid simple lock",
			data: []byte(`{
				"packages": [
					{"name": "vendor/pkg", "version": "1.0.0"}
				]
			}`),
			checkFn: func(lock *compdiff.LockFile) bool {
				return lock != nil && len(lock.Packages) == 1 &&
					lock.Packages[0].Name == "vendor/pkg"
			},
		},
		{
			name: "valid lock with packages-dev",
			data: []byte(`{
				"packages": [],
				"packages-dev": [
					{"name": "vendor/test", "version": "1.0.0"}
				]
			}`),
			checkFn: func(lock *compdiff.LockFile) bool {
				return lock != nil && len(lock.PackagesDev) == 1 &&
					lock.PackagesDev[0].Name == "vendor/test"
			},
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			lock, err := compdiff.ParseLock(testCase.data)
			if err != nil {
				t.Errorf("ParseLock() error = %v", err)
			}

			if !testCase.checkFn(lock) {
				t.Errorf("ParseLock() result check failed")
			}
		})
	}
}

func TestParseLock_Error(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		data    []byte
		checkFn func(err error) error
	}{
		{
			name: "invalid json",
			data: []byte(`{invalid}`),
			checkFn: func(err error) error {
				if err.Error() != "invalid JSON: invalid character 'i' looking for beginning of object key string" {
					return fmt.Errorf("unexpected error: %w", err)
				}

				return nil
			},
		},
		{
			name: "empty input",
			data: []byte{},
			checkFn: func(err error) error {
				if err.Error() != "invalid format: empty input" {
					return fmt.Errorf("unexpected error: %w", err)
				}

				return nil
			},
		},
		{
			name: "missing packages arrays",
			data: []byte(`{"other": "field"}`),
			checkFn: func(err error) error {
				if err.Error() != "invalid format: missing 'packages' or 'packages-dev' fields" {
					return fmt.Errorf("unexpected error: %w", err)
				}

				return nil
			},
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			_, err := compdiff.ParseLock(testCase.data)
			if err == nil {
				t.Errorf("an error is expected")
			} else if err2 := testCase.checkFn(err); err2 != nil {
				t.Error(err2)
			}
		})
	}
}

func TestParseReq(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		data    []byte
		checkFn func(req *compdiff.ReqFile) bool
	}{
		{
			name: "valid compdiff.json",
			data: []byte(`{
				"require": {"vendor/pkg": "^1.0"},
				"require-dev": {"vendor/test": "^1.0"}
			}`),
			checkFn: func(req *compdiff.ReqFile) bool {
				return req != nil && req.Require["vendor/pkg"] == "^1.0" &&
					req.RequireDev["vendor/test"] == "^1.0"
			},
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			lock, err := compdiff.ParseReq(testCase.data)
			if err != nil {
				t.Errorf("ParseReq() error = %v", err)
			}

			if !testCase.checkFn(lock) {
				t.Errorf("ParseReq() result check failed")
			}
		})
	}
}

func TestParseReq_Error(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		data    []byte
		checkFn func(err error) error
	}{
		{
			name: "invalid json",
			data: []byte(`{invalid}`),
			checkFn: func(err error) error {
				if err.Error() != "invalid JSON: invalid character 'i' looking for beginning of object key string" {
					return fmt.Errorf("unexpected error: %w", err)
				}

				return nil
			},
		},
		{
			name: "empty input",
			data: []byte{},
			checkFn: func(err error) error {
				if err.Error() != "invalid format: empty input" {
					return fmt.Errorf("unexpected error: %w", err)
				}

				return nil
			},
		},
		{
			name: "missing require arrays",
			data: []byte(`{"other": "field"}`),
			checkFn: func(err error) error {
				if err.Error() != "invalid format: missing 'require' or 'require-dev' fields" {
					return fmt.Errorf("unexpected error: %w", err)
				}

				return nil
			},
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			_, err := compdiff.ParseReq(testCase.data)
			if err == nil {
				t.Errorf("an error is expected")
			} else if err2 := testCase.checkFn(err); err2 != nil {
				t.Error(err2)
			}
		})
	}
}
