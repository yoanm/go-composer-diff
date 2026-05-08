package compdiff_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/yoanm/go-deps-diff/contract"
	"github.com/yoanm/go-deps-diff/contract/semver"
	difftesting "github.com/yoanm/go-deps-diff/testing"

	compdiff "github.com/yoanm/go-composer-diff"
)

func TestBuildMapFromBytes_Error(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		reqData  []byte
		lockData []byte
		checkFn  func(err error) error
	}{
		{
			name:    "invalid json - req file",
			reqData: []byte(`{invalid}`),
			lockData: []byte(`{
				"packages": [
					{"name": "vendor/pkg", "version": "1.0.0"}
				]
			}`),
			checkFn: func(err error) error {
				if err.Error() != "parsing requirement file content: invalid JSON: invalid character 'i' looking for beginning of object key string" {
					return fmt.Errorf("unexpected error: %w", err)
				}

				return nil
			},
		},
		{
			name: "invalid json - lock file",
			reqData: []byte(`{
				"require": {"vendor/pkg": "^1.0"},
				"require-dev": {"vendor/test": "^1.0"}
			}`),
			lockData: []byte(`{invalid}`),
			checkFn: func(err error) error {
				if err.Error() != "parsing lock file content: invalid JSON: invalid character 'i' looking for beginning of object key string" {
					return fmt.Errorf("unexpected error: %w", err)
				}

				return nil
			},
		},
		{
			name:    "empty input - req file",
			reqData: []byte{},
			lockData: []byte(`{
				"packages": [
					{"name": "vendor/pkg", "version": "1.0.0"}
				]
			}`),
			checkFn: func(err error) error {
				if err.Error() != "parsing requirement file content: invalid format: empty input" {
					return fmt.Errorf("unexpected error: %w", err)
				}

				return nil
			},
		},
		{
			name: "empty input - lock file",
			reqData: []byte(`{
				"require": {"vendor/pkg": "^1.0"},
				"require-dev": {"vendor/test": "^1.0"}
			}`),
			lockData: []byte{},
			checkFn: func(err error) error {
				if err.Error() != "parsing lock file content: invalid format: empty input" {
					return fmt.Errorf("unexpected error: %w", err)
				}

				return nil
			},
		},
		{
			name:    "missing require arrays - req file",
			reqData: []byte(`{"other": "field"}`),
			lockData: []byte(`{
				"packages": [
					{"name": "vendor/pkg", "version": "1.0.0"}
				]
			}`),
			checkFn: func(err error) error {
				if err.Error() != "parsing requirement file content: invalid format: missing 'require' or 'require-dev' fields" {
					return fmt.Errorf("unexpected error: %w", err)
				}

				return nil
			},
		},
		{
			name: "missing require arrays - lock file",
			reqData: []byte(`{
				"require": {"vendor/pkg": "^1.0"},
				"require-dev": {"vendor/test": "^1.0"}
			}`),
			lockData: []byte(`{"other": "field"}`),
			checkFn: func(err error) error {
				if err.Error() != "parsing lock file content: invalid format: missing 'packages' or 'packages-dev' fields" {
					return fmt.Errorf("unexpected error: %w", err)
				}

				return nil
			},
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			_, err := compdiff.BuildMapFromBytes(testCase.reqData, testCase.lockData)
			if err == nil {
				t.Errorf("an error is expected")
			} else if err2 := testCase.checkFn(err); err2 != nil {
				t.Error(err2)
			}
		})
	}
}

func TestRootRequirementProperty(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		req      []byte
		lock     []byte
		expected bool
	}{
		{
			name: "Package not explicitly required",
			req:  []byte(`{"require": {}}`),
			lock: []byte(`{
				"packages": [{"name": "vendor/pkg", "version": "1.0.0", "source": {"reference": "abc"}}]
			}`),
			expected: false,
		},
		{
			name: "Package in require section",
			req:  []byte(`{"require": {"vendor/pkg": "^1.0"}}`),
			lock: []byte(`{
				"packages": [{"name": "vendor/pkg", "version": "1.0.0", "source": {"reference": "abc"}}]
			}`),
			expected: true,
		},
		{
			name: "Package in require-dev section",
			req:  []byte(`{"require-dev": {"vendor/pkg": "^1.0"}}`),
			lock: []byte(`{
				"packages-dev": [{"name": "vendor/pkg", "version": "1.0.0", "source": {"reference": "abc"}}]
			}`),
			expected: false,
		},
	}

	for _, testData := range tests {
		t.Run(testData.name, func(t *testing.T) {
			t.Parallel()

			pkgMap, err := compdiff.BuildMapFromBytes(testData.req, testData.lock)
			if err != nil {
				t.Fatal(fmt.Errorf("building map: %w", err))
			} else if len(pkgMap) != 1 {
				t.Fatal(fmt.Errorf("one and only one package is expected, got %d", len(pkgMap)))
			}

			pkg, pkgExists := pkgMap["vendor/pkg"]
			if !pkgExists {
				t.Fatal(errors.New("package 'vendor/pkg' is expected in the package map"))
			} else if pkg.IsRootRequirement() != testData.expected {
				t.Fatal(
					fmt.Errorf(
						"unexpected IsRootRequirement(): got %t, want %t",
						pkg.IsRootRequirement(),
						testData.expected,
					),
				)
			}
		})
	}
}

func TestRootDevRequirementProperty(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		req      []byte
		lock     []byte
		expected bool
	}{
		{
			name: "Package not explicitly required",
			req:  []byte(`{"require": {}}`),
			lock: []byte(`{
				"packages": [{"name": "vendor/pkg", "version": "1.0.0", "source": {"reference": "abc"}}]
			}`),
			expected: false,
		},
		{
			name: "Package in require-dev section",
			req:  []byte(`{"require-dev": {"vendor/pkg": "^1.0"}}`),
			lock: []byte(`{
				"packages-dev": [{"name": "vendor/pkg", "version": "1.0.0", "source": {"reference": "abc"}}]
			}`),
			expected: true,
		},
		{
			name: "Package in require section",
			req:  []byte(`{"require": {"vendor/pkg": "^1.0"}}`),
			lock: []byte(`{
				"packages": [{"name": "vendor/pkg", "version": "1.0.0", "source": {"reference": "abc"}}]
			}`),
			expected: false,
		},
	}

	for _, testData := range tests {
		t.Run(testData.name, func(t *testing.T) {
			t.Parallel()

			pkgMap, err := compdiff.BuildMapFromBytes(testData.req, testData.lock)
			if err != nil {
				t.Error(fmt.Errorf("building map: %w", err))

				return
			}

			pkg, pkgExists := pkgMap["vendor/pkg"]
			if !pkgExists {
				t.Fatal(errors.New("package 'vendor/pkg' is expected in the package map"))
			} else if pkg.IsRootDevRequirement() != testData.expected {
				t.Fatal(
					fmt.Errorf(
						"unexpected IsRootDevRequirement(): got %t, want %t",
						pkg.IsRootDevRequirement(),
						testData.expected,
					),
				)
			}
		})
	}
}

func TestDevOnlyProperty(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		req      []byte
		lock     []byte
		expected bool
	}{
		{
			name: "Package on packages section",
			req:  []byte(`{"require": {}}`),
			lock: []byte(`{
				"packages": [{"name": "vendor/pkg", "version": "1.0.0", "source": {"reference": "abc"}}]
			}`),
			expected: false,
		},
		{
			name: "Package on packages-dev section",
			req:  []byte(`{"require": {}}`),
			lock: []byte(`{
				"packages-dev": [{"name": "vendor/pkg", "version": "1.0.0", "source": {"reference": "abc"}}]
			}`),
			expected: true,
		},
		{
			// Special case !
			// Package is required as dev requirement but is actually present in the packages section. This may happen
			// if the package happens to be a dependency of another package defined in 'require' section.
			name: "Package on packages section BUT presents in require-dev section",
			req:  []byte(`{"require-dev": {"vendor/pkg": "^1.0"}}`),
			lock: []byte(`{
				"packages": [{"name": "vendor/pkg", "version": "1.0.0", "source": {"reference": "abc"}}]
			}`),
			expected: false,
		},
	}

	for _, testData := range tests {
		t.Run(testData.name, func(t *testing.T) {
			t.Parallel()

			pkgMap, err := compdiff.BuildMapFromBytes(testData.req, testData.lock)
			if err != nil {
				t.Error(fmt.Errorf("building map: %w", err))

				return
			}

			pkg, pkgExists := pkgMap["vendor/pkg"]
			if !pkgExists {
				t.Fatal(errors.New("package 'vendor/pkg' is expected in the package map"))
			} else if pkg.IsDevOnly() != testData.expected {
				t.Fatal(
					fmt.Errorf(
						"unexpected IsDevOnly(): got %t, want %t",
						pkg.IsDevOnly(),
						testData.expected,
					),
				)
			}
		})
	}
}

func TestIsAbandonedProperty(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		lock     []byte
		expected bool
	}{
		{
			name:     "abandoned false",
			lock:     []byte(`{"packages": [{"name": "vendor/pkg", "version": "1.0.0", "abandoned": false}]}`),
			expected: false,
		},
		{
			name:     "abandoned true",
			lock:     []byte(`{"packages": [{"name": "vendor/pkg", "version": "1.0.0", "abandoned": true}]}`),
			expected: true,
		},
		{
			name:     "abandoned string true",
			lock:     []byte(`{"packages": [{"name": "vendor/pkg", "version": "1.0.0", "abandoned": "true"}]}`),
			expected: true,
		},
		{
			name: "abandoned replacement string",
			lock: []byte(`{"packages": [
				{"name": "vendor/pkg", "version": "1.0.0", "abandoned": "https://example.com/replacement"}
			]}`),
			expected: true,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			pkgMap, err := compdiff.BuildMapFromBytes([]byte(`{"require": {}}`), testCase.lock)
			if err != nil {
				t.Fatal(fmt.Errorf("building map: %w", err))
			}

			if len(pkgMap) != 1 {
				t.Fatal(fmt.Errorf("one and only one package is expected, got %d", len(pkgMap)))
			}

			pkg, pkgExists := pkgMap["vendor/pkg"]
			if !pkgExists {
				t.Fatal(errors.New("package 'vendor/pkg' is expected in the package map"))
			} else if pkg.IsAbandoned() != testCase.expected {
				t.Fatalf("IsAbandoned() = %v, want %v", pkg.IsAbandoned(), testCase.expected)
			}
		})
	}
}

func TestLinkProperty(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		lock     []byte
		expected string
	}{
		{
			name: "wiki link priority",
			lock: []byte(`{"packages": [{"name": "vendor/pkg", "version": "1.0.0", "support": {
				"wiki": "https://wiki.example.com",
				"docs": "https://docs.example.com",
				"source": "https://source.example.com"
			}, "homepage": "https://homepage.example.com"}]}`),
			expected: "https://wiki.example.com",
		},
		{
			name: "docs link (no wiki)",
			lock: []byte(`{"packages": [{"name": "vendor/pkg", "version": "1.0.0", "support": {
				"docs": "https://docs.example.com",
				"source": "https://source.example.com"
			}, "homepage": "https://homepage.example.com"}]}`),
			expected: "https://docs.example.com",
		},
		{
			name: "source link (no wiki/docs)",
			lock: []byte(`{"packages": [{"name": "vendor/pkg", "version": "1.0.0", "support": {
				"source": "https://source.example.com"
			}, "homepage": "https://homepage.example.com"}]}`),
			expected: "https://source.example.com",
		},
		{
			name: "homepage link (no support)",
			lock: []byte(`{"packages": [
				{"name": "vendor/pkg", "version": "1.0.0", "homepage": "https://homepage.example.com"}
			]}`),
			expected: "https://homepage.example.com",
		},
		{
			name:     "no links",
			lock:     []byte(`{"packages": [{"name": "vendor/pkg", "version": "1.0.0"}]}`),
			expected: "",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			pkgMap, err := compdiff.BuildMapFromBytes([]byte(`{"require": {}}`), testCase.lock)
			if err != nil {
				t.Fatal(fmt.Errorf("building map: %w", err))
			}

			if len(pkgMap) != 1 {
				t.Fatal(fmt.Errorf("one and only one package is expected, got %d", len(pkgMap)))
			}

			pkg, pkgExists := pkgMap["vendor/pkg"]
			if !pkgExists {
				t.Fatal(errors.New("package 'vendor/pkg' is expected in the package map"))
			} else if pkg.GetLink() != testCase.expected {
				t.Fatalf("GetLink() = %v, want %v", pkg.GetLink(), testCase.expected)
			}
		})
	}
}

func TestVersionProperty(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		lock     []byte
		expected contract.PkgVersion
	}{
		{
			name: "Semver version",
			lock: []byte(`{"packages": [{"name": "vendor/pkg", "version": "1.2.3",
				"dist": {"reference": "abc123"}
			}]}`),
			expected: contract.PkgVersion{Raw: "1.2.3", Label: "1.2.3", Semver: &semver.Version{Major: 1, Minor: 2, Patch: 3, Extra: ""}},
		},
		{
			name: "Semver version with extra",
			lock: []byte(`{"packages": [{"name": "vendor/pkg", "version": "1.2.3+beta",
				"dist": {"reference": "abc123"}
			}]}`),
			expected: contract.PkgVersion{Raw: "1.2.3+beta", Label: "1.2.3+beta", Semver: &semver.Version{Major: 1, Minor: 2, Patch: 3, Extra: "+beta"}},
		},
		{
			name: "not semver - dist reference",
			lock: []byte(`{"packages": [{"name": "vendor/pkg", "version": "dev-master",
				"dist": {"reference": "abc123"}
			}]}`),
			expected: contract.PkgVersion{Raw: "abc123", Label: "dev-master#abc123", Semver: nil},
		},
		{
			name: "not semver - source reference (no dist)",
			lock: []byte(`{"packages": [{"name": "vendor/pkg", "version": "dev-master",
				"source": {"reference": "def456"}
			}]}`),
			expected: contract.PkgVersion{Raw: "def456", Label: "dev-master#def456", Semver: nil},
		},
		{
			name: "not semver - dist preferred over source",
			lock: []byte(`{"packages": [{"name": "vendor/pkg", "version": "dev-master",
				"source": {"reference": "nop"},
				"dist": {"reference": "abc123"}
			}]}`),
			expected: contract.PkgVersion{Raw: "abc123", Label: "dev-master#abc123", Semver: nil},
		},
		{
			name:     "not semver - no reference",
			lock:     []byte(`{"packages": [{"name": "vendor/pkg", "version": "dev-master"}]}`),
			expected: contract.PkgVersion{Raw: "dev-master", Label: "dev-master", Semver: nil},
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			pkgMap, err := compdiff.BuildMapFromBytes([]byte(`{"require": {}}`), testCase.lock)
			if err != nil {
				t.Fatal(fmt.Errorf("building map: %w", err))
			}

			if len(pkgMap) != 1 {
				t.Fatal(fmt.Errorf("one and only one package is expected, got %d", len(pkgMap)))
			}

			pkg, pkgExists := pkgMap["vendor/pkg"]
			if !pkgExists {
				t.Error("package 'vendor/pkg' is expected in the package map")
			} else if err2 := difftesting.ValidatePkgVersion(pkg.GetVersion(), testCase.expected); err2 != nil {
				t.Error(err2)
			}
		})
	}
}
