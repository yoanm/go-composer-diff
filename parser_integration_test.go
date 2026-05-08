package compdiff_test

import (
	"errors"
	"fmt"
	compdiff "github.com/yoanm/go-composer-diff"
	"testing"

	"github.com/yoanm/go-deps-diff/contract"
)

func TestAbandonedNullAndNumber(t *testing.T) {
	t.Parallel()

	req := []byte(`{"require": {}}`)
	lock := []byte(`{
  "packages": [
    {"name":"vendor/nullabandoned","version":"1.0.0","abandoned":null},
    {"name":"vendor/numabandoned","version":"1.0.0","abandoned":1}
  ]
}`)

	packageMap, err := compdiff.BuildMapFromBytes(req, lock)
	if err != nil {
		t.Fatalf("BuildMapFromBytes failed: %v", err)
	}

	if packageMap["vendor/nullabandoned"].IsAbandoned() {
		t.Errorf("vendor/nullabandoned should not be considered abandoned")
	}

	if packageMap["vendor/numabandoned"].IsAbandoned() {
		t.Errorf("vendor/numabandoned should not be considered abandoned for numeric value")
	}
}

func TestParseLock_SupportTypeMismatch(t *testing.T) {
	t.Parallel()

	// When support is a string instead of object, unmarshal should fail
	data := []byte(`{"packages": [{"name":"vendor/bad","version":"1.0.0","support":"https://wiki.example"}]}`)
	if _, err := compdiff.ParseLock(data); err == nil {
		t.Fatalf("expected ParseLock to fail on support type mismatch, but it succeeded")
	}
}

func TestParseLock_WithoutOptionalFields(t *testing.T) {
	t.Parallel()

	// Packages without optional fields should parse fine
	data := []byte(`{
  "packages": [
    {"name":"vendor/minimal","version":"1.2.3"}
  ]
}`)

	lock, err := compdiff.ParseLock(data)
	if err != nil {
		t.Fatalf("ParseLock failed: %v", err)
	}

	if len(lock.Packages) != 1 {
		t.Fatalf("expected 1 package, got %d", len(lock.Packages))
	}
	// BuildMap should not panic
	req := []byte(`{"require": {}}`)
	if _, err2 := compdiff.BuildMapFromBytes(req, data); err2 != nil {
		t.Fatalf("BuildMapFromBytes failed: %v", err2)
	}
}

func TestBuildMap_LinkAndAbandoned(t *testing.T) {
	t.Parallel()

	req := []byte(`{"require": {}}`)
	lock := []byte(`{
		"packages": [
			{
				"name":"vendor/pkg1",
				"version":"1.0.0",
				"support":{
					"wiki":"https://wiki.example/1",
					"docs":"https://docs.example/1",
					"source":"https://example.com/1"
				}
			},
			{
				"name":"vendor/pkg2",
				"version":"1.0.0",
				"support":{"docs":"https://docs.example/2","source":"https://example.com/2"}
			},
			{"name":"vendor/pkg3","version":"1.0.0","support":{"source":"https://example.com/3"}},
			{"name":"vendor/pkg4","version":"1.0.0","homepage":"https://home.example/4"},
			{"name":"vendor/pkg5","version":"1.0.0","abandoned":true},
			{"name":"vendor/pkg6","version":"1.0.0","abandoned":"https://packagist.org/packages/vendor/pkg6"},
			{"name":"vendor/pkg7","version":"1.0.0","abandoned":false}
		]
	}`)

	packageMap, err := compdiff.BuildMapFromBytes(req, lock)
	if err != nil {
		t.Fatalf("BuildMapFromBytes failed: %v", err)
	}

	if err2 := validatePackageLinkAndIsAbandoned(packageMap, "vendor/pkg1", false, "https://wiki.example/1"); err2 != nil {
		t.Error(err2)
	}

	if err2 := validatePackageLinkAndIsAbandoned(packageMap, "vendor/pkg2", false, "https://docs.example/2"); err2 != nil {
		t.Error(err2)
	}

	if err2 := validatePackageLinkAndIsAbandoned(packageMap, "vendor/pkg3", false, "https://example.com/3"); err2 != nil {
		t.Error(err2)
	}

	if err2 := validatePackageLinkAndIsAbandoned(packageMap, "vendor/pkg4", false, "https://home.example/4"); err2 != nil {
		t.Error(err2)
	}

	if err2 := validatePackageLinkAndIsAbandoned(packageMap, "vendor/pkg5", true, ""); err2 != nil {
		t.Error(err2)
	}

	if err2 := validatePackageLinkAndIsAbandoned(packageMap, "vendor/pkg6", true, ""); err2 != nil {
		t.Error(err2)
	}

	if err2 := validatePackageLinkAndIsAbandoned(packageMap, "vendor/pkg7", false, ""); err2 != nil {
		t.Error(err2)
	}
}

func validatePackageLinkAndIsAbandoned(
	pkgMap contract.PackageMap,
	targetPkgName string,
	expectedIsAbandoned bool,
	expectedLink string,
) error {
	pkg, ok := pkgMap[targetPkgName]
	if !ok {
		return errors.New("package is missing from map")
	}

	if expectedIsAbandoned != pkg.IsAbandoned() {
		return fmt.Errorf("invalid isAbandoned. Got %v, want %v", pkg.IsAbandoned(), expectedIsAbandoned)
	}

	if expectedLink != pkg.GetLink() {
		return fmt.Errorf("link mismatch: got %q, want %q", pkg.GetLink(), expectedLink)
	}

	return nil
}

func TestBuildMap_RootAndDevRequirements(t *testing.T) {
	t.Parallel()

	req := []byte(`{
		"require": {"vendor/root": "^1.0"},
		"require-dev": {"vendor/dev": "^1.0", "vendor/dev2": "^1.0"}
	}`)
	lock := []byte(`{
		"packages": [
			{"name":"vendor/root","version":"1.0.0"},
			{"name":"vendor/dev2","version":"1.0.0"}
		],
		"packages-dev": [
			{"name":"vendor/dev","version":"1.0.0"}
		]
	}`)

	packageMap, err := compdiff.BuildMapFromBytes(req, lock)
	if err != nil {
		t.Fatalf("BuildMapFromBytes failed: %v", err)
	}

	if err2 := validatePackageIsRequirementIsDevOnly(packageMap, "vendor/root", true, false, false); err2 != nil {
		t.Error(err2)
	}

	if err2 := validatePackageIsRequirementIsDevOnly(packageMap, "vendor/dev", false, true, true); err2 != nil {
		t.Error(err2)
	}

	if err2 := validatePackageIsRequirementIsDevOnly(packageMap, "vendor/dev2", false, true, false); err2 != nil {
		t.Error(err2)
	}
}

func validatePackageIsRequirementIsDevOnly(
	pkgMap contract.PackageMap,
	targetPkgName string,
	expectedIsRootRequirement bool,
	expectedIsRootDevRequirement bool,
	expectedIsDevOnly bool,
) error {
	pkg, ok := pkgMap[targetPkgName]
	if !ok {
		return errors.New("package is missing from map")
	}

	if expectedIsRootRequirement != pkg.IsRootRequirement() {
		return fmt.Errorf(
			"invalid isRootRequirement. Got %v, want %v",
			pkg.IsRootRequirement(),
			expectedIsRootRequirement,
		)
	}

	if expectedIsRootDevRequirement != pkg.IsRootDevRequirement() {
		return fmt.Errorf(
			"invalid isRootDevRequirement. Got %v, want %v",
			pkg.IsRootDevRequirement(),
			expectedIsRootDevRequirement,
		)
	}

	if expectedIsDevOnly != pkg.IsDevOnly() {
		return fmt.Errorf("invalid isDevOnly. Got %v, want %v", pkg.IsDevOnly(), expectedIsDevOnly)
	}

	return nil
}

func TestVersionHandling(t *testing.T) {
	t.Parallel()

	req := []byte(`{"require": {}}`)

	cases := []struct {
		name        string
		lock        string
		expectRaw   string
		expectLabel string
	}{
		{
			name:        "semver_valid_ignores_refs",
			lock:        `{"packages":[{"name":"vendor/my-pkg","version":"2.0.0","dist":{"reference":"deadbeef","url":"https://example.org"}}]}`,
			expectRaw:   "2.0.0",
			expectLabel: "2.0.0",
		},
		{
			name:        "non_semver_uses_dist_reference",
			lock:        `{"packages":[{"name":"vendor/my-pkg","version":"dev-master","dist":{"reference":"abcd1234dead","url":"https://example.org"}}]}`,
			expectRaw:   "abcd1234dead",
			expectLabel: "dev-master#abcd123",
		},
		{
			name:        "non_semver_uses_source_reference_when_dist_absent",
			lock:        `{"packages":[{"name":"vendor/my-pkg","version":"dev-branch","source":{"reference":"1234567890abcdef","url":"https://example.org"}}]}`,
			expectRaw:   "1234567890abcdef",
			expectLabel: "dev-branch#1234567",
		},
		{
			name:        "no_reference_available",
			lock:        `{"packages":[{"name":"vendor/my-pkg","version":"v1.2.3-alpha"}]}`,
			expectRaw:   "v1.2.3-alpha",
			expectLabel: "v1.2.3-alpha",
		},
		{
			name:        "prefers_dist_over_source",
			lock:        `{"packages":[{"name":"vendor/my-pkg","version":"dev-feature","dist":{"reference":"distref12","url":"https://example.org"},"source":{"reference":"sourceref34","url":"https://example.org"}}]}`,
			expectRaw:   "distref12",
			expectLabel: "dev-feature#distref",
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			pm, err := compdiff.BuildMapFromBytes(req, []byte(testCase.lock))
			if err != nil {
				t.Fatalf("BuildMapFromBytes failed: %v", err)
			}

			w, ok := pm["vendor/my-pkg"]
			if !ok {
				t.Fatal("package vendor/my-pkg missing from map")
			}

			ver := w.GetVersion()
			if ver.Raw != testCase.expectRaw {
				t.Errorf("%s: unexpected Raw. got=%q want=%q", testCase.name, ver.Raw, testCase.expectRaw)
			}

			if ver.Label != testCase.expectLabel {
				t.Errorf("%s: unexpected Label. got=%q want=%q", testCase.name, ver.Label, testCase.expectLabel)
			}
		})
	}
}
