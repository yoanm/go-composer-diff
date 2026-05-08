package compdiff_test

import (
	"testing"

	"github.com/andreyvit/diff"

	"github.com/yoanm/go-deps-diff/contract"
	"github.com/yoanm/go-deps-diff/contract/semver"
	difftesting "github.com/yoanm/go-deps-diff/testing"

	compdiff "github.com/yoanm/go-composer-diff"
)

func TestIntegration_FileDiff_Errors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		previous      *compdiff.FileInput
		current       *compdiff.FileInput
		expectedError string
	}{
		{
			name:          "invalid previous req file",
			previous:      &compdiff.FileInput{Lock: "./testdata/simple-composer.lock", Requirement: "/do.not.exist/previous-req-file"},
			current:       &compdiff.FileInput{Lock: "./testdata/simple-composer.lock", Requirement: "./testdata/simple-composer.json"},
			expectedError: "reading previous files: reading requirement file: open /do.not.exist/previous-req-file: no such file or directory",
		},
		{
			name:          "invalid current req file",
			previous:      &compdiff.FileInput{Lock: "./testdata/simple-composer.lock", Requirement: "./testdata/simple-composer.json"},
			current:       &compdiff.FileInput{Lock: "./testdata/simple-composer.lock", Requirement: "/do.not.exist/current-req-file"},
			expectedError: "reading current files: reading requirement file: open /do.not.exist/current-req-file: no such file or directory",
		},
		{
			name:          "invalid previous lock file",
			previous:      &compdiff.FileInput{Lock: "/do.not.exist/previous-lock-file", Requirement: "./testdata/simple-composer.json"},
			current:       &compdiff.FileInput{Lock: "./testdata/simple-composer.lock", Requirement: "./testdata/simple-composer.json"},
			expectedError: "reading previous files: reading lock file: open /do.not.exist/previous-lock-file: no such file or directory",
		},
		{
			name:          "invalid current lock file",
			previous:      &compdiff.FileInput{Lock: "./testdata/simple-composer.lock", Requirement: "./testdata/simple-composer.json"},
			current:       &compdiff.FileInput{Lock: "/do.not.exist/current-lock-file", Requirement: "./testdata/simple-composer.json"},
			expectedError: "reading current files: reading lock file: open /do.not.exist/current-lock-file: no such file or directory",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			_, err := compdiff.FileDiff(testCase.previous, testCase.current)
			if err == nil {
				t.Fatal("an error is expected")
			} else if err.Error() != testCase.expectedError {
				t.Errorf("unexpected error: %s", diff.CharacterDiff(testCase.expectedError, err.Error()))
			}
		})
	}
}

func TestIntegration_Diff_Errors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		previousReqData  []byte
		previousLockData []byte
		currentReqData   []byte
		currentLockData  []byte
		expectedError    string
	}{
		{
			name:             "invalid json - previous req file",
			previousReqData:  []byte(`{invalid}`),
			previousLockData: []byte(`{"packages": [{"name": "vendor/pkg", "version": "1.0.0"}]}`),
			currentReqData:   []byte(`{"require": {"vendor/pkg": "^1.0"}}`),
			currentLockData:  []byte(`{"packages": [{"name": "vendor/pkg", "version": "1.0.0"}]}`),

			expectedError: "building previous package map: parsing requirement file content: invalid JSON: invalid character 'i' looking for beginning of object key string",
		},
		{
			name:             "invalid json - current req file",
			previousReqData:  []byte(`{"require": {"vendor/pkg": "^1.0"}}`),
			previousLockData: []byte(`{"packages": [{"name": "vendor/pkg", "version": "1.0.0"}]}`),
			currentReqData:   []byte(`{invalid}`),
			currentLockData:  []byte(`{"packages": [{"name": "vendor/pkg", "version": "1.0.0"}]}`),

			expectedError: "building current package map: parsing requirement file content: invalid JSON: invalid character 'i' looking for beginning of object key string",
		},
		{
			name:             "invalid json - previous lock file",
			previousReqData:  []byte(`{"require": {"vendor/pkg": "^1.0"}}`),
			previousLockData: []byte(`{invalid}`),
			currentReqData:   []byte(`{"require": {"vendor/pkg": "^1.0"}}`),
			currentLockData:  []byte(`{"packages": [{"name": "vendor/pkg", "version": "1.0.0"}]}`),

			expectedError: "building previous package map: parsing lock file content: invalid JSON: invalid character 'i' looking for beginning of object key string",
		},
		{
			name:             "invalid json - current lock file",
			previousReqData:  []byte(`{"require": {"vendor/pkg": "^1.0"}}`),
			previousLockData: []byte(`{"packages": [{"name": "vendor/pkg", "version": "1.0.0"}]}`),
			currentReqData:   []byte(`{"require": {"vendor/pkg": "^1.0"}}`),
			currentLockData:  []byte(`{invalid}`),

			expectedError: "building current package map: parsing lock file content: invalid JSON: invalid character 'i' looking for beginning of object key string",
		},
		{
			name:             "empty input - previous req file",
			previousReqData:  []byte{},
			previousLockData: []byte(`{"packages": [{"name": "vendor/pkg", "version": "1.0.0"}]}`),
			currentReqData:   []byte(`{"require": {"vendor/pkg": "^1.0"}}`),
			currentLockData:  []byte(`{"packages": [{"name": "vendor/pkg", "version": "1.0.0"}]}`),

			expectedError: "building previous package map: parsing requirement file content: invalid format: empty input",
		},
		{
			name:             "empty input - current req file",
			previousReqData:  []byte(`{"require": {"vendor/pkg": "^1.0"}}`),
			previousLockData: []byte(`{"packages": [{"name": "vendor/pkg", "version": "1.0.0"}]}`),
			currentReqData:   []byte{},
			currentLockData:  []byte(`{"packages": [{"name": "vendor/pkg", "version": "1.0.0"}]}`),

			expectedError: "building current package map: parsing requirement file content: invalid format: empty input",
		},
		{
			name:             "empty input - previous lock file",
			previousReqData:  []byte(`{"require": {"vendor/pkg": "^1.0"}}`),
			previousLockData: []byte{},
			currentReqData:   []byte(`{"require": {"vendor/pkg": "^1.0"}}`),
			currentLockData:  []byte(`{"packages": [{"name": "vendor/pkg", "version": "1.0.0"}]}`),
			expectedError:    "building previous package map: parsing lock file content: invalid format: empty input",
		},
		{
			name:             "empty input - current lock file",
			previousReqData:  []byte(`{"require": {"vendor/pkg": "^1.0"}}`),
			previousLockData: []byte(`{"packages": [{"name": "vendor/pkg", "version": "1.0.0"}]}`),
			currentReqData:   []byte(`{"require": {"vendor/pkg": "^1.0"}}`),
			currentLockData:  []byte{},
			expectedError:    "building current package map: parsing lock file content: invalid format: empty input",
		},
		{
			name:             "missing require arrays - previous req file",
			previousReqData:  []byte(`{"other": "field"}`),
			previousLockData: []byte(`{"packages": [{"name": "vendor/pkg", "version": "1.0.0"}]}`),
			currentReqData:   []byte(`{"require": {"vendor/pkg": "^1.0"}}`),
			currentLockData:  []byte(`{"packages": [{"name": "vendor/pkg", "version": "1.0.0"}]}`),

			expectedError: "building previous package map: parsing requirement file content: invalid format: missing 'require' or 'require-dev' fields",
		},
		{
			name:             "missing require arrays - current req file",
			previousReqData:  []byte(`{"require": {"vendor/pkg": "^1.0"}}`),
			previousLockData: []byte(`{"packages": [{"name": "vendor/pkg", "version": "1.0.0"}]}`),
			currentReqData:   []byte(`{"other": "field"}`),
			currentLockData:  []byte(`{"packages": [{"name": "vendor/pkg", "version": "1.0.0"}]}`),

			expectedError: "building current package map: parsing requirement file content: invalid format: missing 'require' or 'require-dev' fields",
		},
		{
			name:             "missing require arrays - previous lock file",
			previousReqData:  []byte(`{"require": {"vendor/pkg": "^1.0"}}`),
			previousLockData: []byte(`{"other": "field"}`),
			currentReqData:   []byte(`{"require": {"vendor/pkg": "^1.0"}}`),
			currentLockData:  []byte(`{"packages": [{"name": "vendor/pkg", "version": "1.0.0"}]}`),

			expectedError: "building previous package map: parsing lock file content: invalid format: missing 'packages' or 'packages-dev' fields",
		},
		{
			name:             "missing require arrays - current lock file",
			previousReqData:  []byte(`{"require": {"vendor/pkg": "^1.0"}}`),
			previousLockData: []byte(`{"packages": [{"name": "vendor/pkg", "version": "1.0.0"}]}`),
			currentReqData:   []byte(`{"require": {"vendor/pkg": "^1.0"}}`),
			currentLockData:  []byte(`{"other": "field"}`),

			expectedError: "building current package map: parsing lock file content: invalid format: missing 'packages' or 'packages-dev' fields",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			_, err := compdiff.Diff(
				&compdiff.Input{
					Lock:        testCase.previousLockData,
					Requirement: testCase.previousReqData,
				},
				&compdiff.Input{
					Lock:        testCase.currentLockData,
					Requirement: testCase.currentReqData,
				},
			)
			if err == nil {
				t.Fatal("an error is expected")
			} else if err.Error() != testCase.expectedError {
				t.Errorf("unexpected error: %s", diff.CharacterDiff(testCase.expectedError, err.Error()))
			}
		})
	}
}

func TestIntegration_Diff_OriginalDataset(t *testing.T) {
	t.Parallel()

	out, err := compdiff.FileDiff(
		&compdiff.FileInput{
			Lock:        "./testdata/original_dataset-PREVIOUS-composer.lock",
			Requirement: "./testdata/original_dataset-PREVIOUS-composer.json",
		},
		&compdiff.FileInput{
			Lock:        "./testdata/original_dataset-CURRENT-composer.lock",
			Requirement: "./testdata/original_dataset-CURRENT-composer.json",
		},
	)
	if err != nil {
		t.Errorf("FileDiff() error = %v", err)

		return
	}

	for _, err := range difftesting.ValidateChanges(out, integrationOriginalDatasetExpectation) {
		t.Error(err)
	}
}

var integrationOriginalDatasetExpectation = contract.DiffMap{
	"sebastian/diff": { // sebastian/diff	4.0.4 	↘️‼️️ 	3.0.3
		Package: &difftesting.TestPkgWrapper{
			Name:               "sebastian/diff",
			Abandoned:          false,
			Version:            contract.PkgVersion{Raw: "3.0.3", Label: "3.0.3", Semver: &semver.Version{Major: 3, Minor: 0, Patch: 3, Extra: ""}},
			Link:               "https://github.com/sebastianbergmann/diff/tree/3.0.3",
			DevOnly:            false,
			RootRequirement:    false,
			RootDevRequirement: true,
		},
		Operation: contract.Operation{
			Name:       contract.DowngradeOperation,
			SemverType: contract.SemverMajorUpdate,
		},
		PreviousVersion: contract.PkgVersion{Raw: "4.0.4", Label: "4.0.4", Semver: &semver.Version{Major: 4, Minor: 0, Patch: 4, Extra: ""}},
	},
	"symfony/asset": { // symfony/asset	v4.4.27 	↗️️ 	v5.4.21
		Package: &difftesting.TestPkgWrapper{
			Name:               "symfony/asset",
			Abandoned:          false,
			Version:            contract.PkgVersion{Raw: "v5.4.21", Label: "v5.4.21", Semver: &semver.Version{Major: 5, Minor: 4, Patch: 21, Extra: ""}},
			Link:               "https://github.com/symfony/asset/tree/v5.4.21",
			DevOnly:            false,
			RootRequirement:    false,
			RootDevRequirement: true,
		},
		Operation: contract.Operation{
			Name:       contract.UpgradeOperation,
			SemverType: contract.SemverMajorUpdate,
		},
		PreviousVersion: contract.PkgVersion{Raw: "v4.4.27", Label: "v4.4.27", Semver: &semver.Version{Major: 4, Minor: 4, Patch: 27, Extra: ""}},
	},
	"yoanm/jsonrpc-server-sdk": { // yoanm/jsonrpc-server-sdk	dev-master#dcd886d❗ 	➡️ 	v1.3.0
		Package: &difftesting.TestPkgWrapper{
			Name:               "yoanm/jsonrpc-server-sdk",
			Abandoned:          false,
			Version:            contract.PkgVersion{Raw: "v1.3.0", Label: "v1.3.0", Semver: &semver.Version{Major: 1, Minor: 3, Patch: 0, Extra: ""}},
			Link:               "https://github.com/yoanm/php-jsonrpc-server-sdk/tree/v1.3.0",
			DevOnly:            false,
			RootRequirement:    true,
			RootDevRequirement: false,
		},
		Operation: contract.Operation{
			Name:       contract.UnknownUpdateOperation,
			SemverType: contract.SemverUnknownUpdate,
		},
		PreviousVersion: contract.PkgVersion{Raw: "dcd886d0ae9246129ec8fbf5e082eff1fc3c49ea", Label: "dev-master#dcd886d", Semver: nil},
	},
	"yoanm/jsonrpc-server-doc-sdk": { // yoanm/jsonrpc-server-doc-sdk	➕️ 	dev-master#a0febcc❗
		Package: &difftesting.TestPkgWrapper{
			Name:               "yoanm/jsonrpc-server-doc-sdk",
			Abandoned:          false,
			Version:            contract.PkgVersion{Raw: "a0febcca883a64c71ed5c97d9e2bacc46a26ff30", Label: "dev-master#a0febcc", Semver: nil},
			Link:               "https://github.com/yoanm/php-jsonrpc-server-doc-sdk/tree/v0.3.0",
			DevOnly:            false,
			RootRequirement:    true,
			RootDevRequirement: false,
		},
		Operation: contract.Operation{
			Name:       contract.AdditionOperation,
			SemverType: contract.SemverNoUpdate,
		},
	},
	"behat/gherkin": { // behat/gherkin	v4.8.0 	↘️‼️️ 	v4.7.0
		Package: &difftesting.TestPkgWrapper{
			Name:               "behat/gherkin",
			Abandoned:          false,
			Version:            contract.PkgVersion{Raw: "v4.7.0", Label: "v4.7.0", Semver: &semver.Version{Major: 4, Minor: 7, Patch: 0, Extra: ""}},
			Link:               "https://github.com/Behat/Gherkin/tree/v4.7.0",
			DevOnly:            false,
			RootRequirement:    false,
			RootDevRequirement: true,
		},
		Operation: contract.Operation{
			Name:       contract.DowngradeOperation,
			SemverType: contract.SemverMinorUpdate,
		},
		PreviousVersion: contract.PkgVersion{Raw: "v4.8.0", Label: "v4.8.0", Semver: &semver.Version{Major: 4, Minor: 8, Patch: 0, Extra: ""}},
	},
	"symfony/deprecation-contracts": { // symfony/deprecation-contracts	v2.2.0 	↗️️ 	v2.5.2
		Package: &difftesting.TestPkgWrapper{
			Name:               "symfony/deprecation-contracts",
			Abandoned:          false,
			Version:            contract.PkgVersion{Raw: "v2.5.2", Label: "v2.5.2", Semver: &semver.Version{Major: 2, Minor: 5, Patch: 2, Extra: ""}},
			Link:               "https://github.com/symfony/deprecation-contracts/tree/v2.5.2",
			DevOnly:            false,
			RootRequirement:    true,
			RootDevRequirement: false,
		},
		Operation: contract.Operation{
			Name:       contract.UpgradeOperation,
			SemverType: contract.SemverMinorUpdate,
		},
		PreviousVersion: contract.PkgVersion{Raw: "v2.2.0", Label: "v2.2.0", Semver: &semver.Version{Major: 2, Minor: 2, Patch: 0, Extra: ""}},
	},
	"symfony/polyfill-ctype": { // symfony/polyfill-ctype	v1.23.0 	↗️️ 	v1.27.0
		Package: &difftesting.TestPkgWrapper{
			Name:               "symfony/polyfill-ctype",
			Abandoned:          false,
			Version:            contract.PkgVersion{Raw: "v1.27.0", Label: "v1.27.0", Semver: &semver.Version{Major: 1, Minor: 27, Patch: 0, Extra: ""}},
			Link:               "https://github.com/symfony/polyfill-ctype/tree/v1.27.0",
			DevOnly:            false,
			RootRequirement:    false,
			RootDevRequirement: false,
		},
		Operation: contract.Operation{
			Name:       contract.UpgradeOperation,
			SemverType: contract.SemverMinorUpdate,
		},
		PreviousVersion: contract.PkgVersion{Raw: "v1.23.0", Label: "v1.23.0", Semver: &semver.Version{Major: 1, Minor: 23, Patch: 0, Extra: ""}},
	},
	"symfony/polyfill-php80": { // symfony/polyfill-php80	v1.23.1 	↗️️ 	v1.27.0
		Package: &difftesting.TestPkgWrapper{
			Name:               "symfony/polyfill-php80",
			Abandoned:          false,
			Version:            contract.PkgVersion{Raw: "v1.27.0", Label: "v1.27.0", Semver: &semver.Version{Major: 1, Minor: 27, Patch: 0, Extra: ""}},
			Link:               "https://github.com/symfony/polyfill-php80/tree/v1.27.0",
			DevOnly:            false,
			RootRequirement:    false,
			RootDevRequirement: false,
		},
		Operation: contract.Operation{
			Name:       contract.UpgradeOperation,
			SemverType: contract.SemverMinorUpdate,
		},
		PreviousVersion: contract.PkgVersion{Raw: "v1.23.1", Label: "v1.23.1", Semver: &semver.Version{Major: 1, Minor: 23, Patch: 1, Extra: ""}},
	},
	"phpstan/phpstan": { // phpstan/phpstan	0.12.96 	↗️️ 	0.12.100
		Package: &difftesting.TestPkgWrapper{
			Name:               "phpstan/phpstan",
			Abandoned:          false,
			Version:            contract.PkgVersion{Raw: "0.12.100", Label: "0.12.100", Semver: &semver.Version{Major: 0, Minor: 12, Patch: 100, Extra: ""}},
			Link:               "https://github.com/phpstan/phpstan/tree/0.12.100",
			DevOnly:            false,
			RootRequirement:    false,
			RootDevRequirement: true,
		},
		Operation: contract.Operation{
			Name:       contract.UpgradeOperation,
			SemverType: contract.SemverPatchUpdate,
		},
		PreviousVersion: contract.PkgVersion{Raw: "0.12.96", Label: "0.12.96", Semver: &semver.Version{Major: 0, Minor: 12, Patch: 96, Extra: ""}},
	},
	"sebastian/code-unit": { // sebastian/code-unit	1.0.8 	↘️‼️️ 	1.0.7
		Package: &difftesting.TestPkgWrapper{
			Name:               "sebastian/code-unit",
			Abandoned:          false,
			Version:            contract.PkgVersion{Raw: "1.0.7", Label: "1.0.7", Semver: &semver.Version{Major: 1, Minor: 0, Patch: 7, Extra: ""}},
			Link:               "https://github.com/sebastianbergmann/code-unit/tree/1.0.7",
			DevOnly:            false,
			RootRequirement:    false,
			RootDevRequirement: true,
		},
		Operation: contract.Operation{
			Name:       contract.DowngradeOperation,
			SemverType: contract.SemverPatchUpdate,
		},
		PreviousVersion: contract.PkgVersion{Raw: "1.0.8", Label: "1.0.8", Semver: &semver.Version{Major: 1, Minor: 0, Patch: 8, Extra: ""}},
	},
	"symfony/cache-contracts": { // symfony/cache-contracts	v1.1.1 	↗️️ 	v1.1.13
		Package: &difftesting.TestPkgWrapper{
			Name:               "symfony/cache-contracts",
			Abandoned:          false,
			Version:            contract.PkgVersion{Raw: "v1.1.13", Label: "v1.1.13", Semver: &semver.Version{Major: 1, Minor: 1, Patch: 13, Extra: ""}},
			Link:               "https://github.com/symfony/cache-contracts/tree/v1.1.13",
			DevOnly:            false,
			RootRequirement:    true,
			RootDevRequirement: false,
		},
		Operation: contract.Operation{
			Name:       contract.UpgradeOperation,
			SemverType: contract.SemverPatchUpdate,
		},
		PreviousVersion: contract.PkgVersion{Raw: "v1.1.1", Label: "v1.1.1", Semver: &semver.Version{Major: 1, Minor: 1, Patch: 1, Extra: ""}},
	},
	"psr/cache": { // ➕ 	psr/cache 	3.0.0
		Package: &difftesting.TestPkgWrapper{
			Name:               "psr/cache",
			Abandoned:          false,
			Version:            contract.PkgVersion{Raw: "3.0.0", Label: "3.0.0", Semver: &semver.Version{Major: 3, Minor: 0, Patch: 0, Extra: ""}},
			Link:               "https://github.com/php-fig/cache/tree/3.0.0",
			DevOnly:            false,
			RootRequirement:    false,
			RootDevRequirement: false,
		},
		Operation: contract.Operation{
			Name:       contract.AdditionOperation,
			SemverType: contract.SemverNoUpdate,
		},
	},
	"psr/container": { // ➕ 	psr/container 	1.1.2
		Package: &difftesting.TestPkgWrapper{
			Name:               "psr/container",
			Abandoned:          false,
			Version:            contract.PkgVersion{Raw: "1.1.2", Label: "1.1.2", Semver: &semver.Version{Major: 1, Minor: 1, Patch: 2, Extra: ""}},
			Link:               "https://github.com/php-fig/container/tree/1.1.2",
			DevOnly:            false,
			RootRequirement:    true,
			RootDevRequirement: false,
		},
		Operation: contract.Operation{
			Name:       contract.AdditionOperation,
			SemverType: contract.SemverNoUpdate,
		},
	},
	"symfony/console": { // ➕ 	symfony/console 	v5.4.21
		Package: &difftesting.TestPkgWrapper{
			Name:               "symfony/console",
			Abandoned:          false,
			Version:            contract.PkgVersion{Raw: "v5.4.21", Label: "v5.4.21", Semver: &semver.Version{Major: 5, Minor: 4, Patch: 21, Extra: ""}},
			Link:               "https://github.com/symfony/console/tree/v5.4.21",
			DevOnly:            false,
			RootRequirement:    true,
			RootDevRequirement: false,
		},
		Operation: contract.Operation{
			Name:       contract.AdditionOperation,
			SemverType: contract.SemverNoUpdate,
		},
	},
	"symfony/polyfill-intl-grapheme": { // ➕ 	symfony/polyfill-intl-grapheme 	v1.27.0
		Package: &difftesting.TestPkgWrapper{
			Name:               "symfony/polyfill-intl-grapheme",
			Abandoned:          false,
			Version:            contract.PkgVersion{Raw: "v1.27.0", Label: "v1.27.0", Semver: &semver.Version{Major: 1, Minor: 27, Patch: 0, Extra: ""}},
			Link:               "https://github.com/symfony/polyfill-intl-grapheme/tree/v1.27.0",
			DevOnly:            false,
			RootRequirement:    false,
			RootDevRequirement: false,
		},
		Operation: contract.Operation{
			Name:       contract.AdditionOperation,
			SemverType: contract.SemverNoUpdate,
		},
	},
	"symfony/polyfill-intl-normalizer": { // ➕ 	symfony/polyfill-intl-normalizer 	v1.27.0
		Package: &difftesting.TestPkgWrapper{
			Name:               "symfony/polyfill-intl-normalizer",
			Abandoned:          false,
			Version:            contract.PkgVersion{Raw: "v1.27.0", Label: "v1.27.0", Semver: &semver.Version{Major: 1, Minor: 27, Patch: 0, Extra: ""}},
			Link:               "https://github.com/symfony/polyfill-intl-normalizer/tree/v1.27.0",
			DevOnly:            false,
			RootRequirement:    false,
			RootDevRequirement: false,
		},
		Operation: contract.Operation{
			Name:       contract.AdditionOperation,
			SemverType: contract.SemverNoUpdate,
		},
	},
	"symfony/polyfill-mbstring": { // ➕ 	symfony/polyfill-mbstring 	v1.27.0
		Package: &difftesting.TestPkgWrapper{
			Name:               "symfony/polyfill-mbstring",
			Abandoned:          false,
			Version:            contract.PkgVersion{Raw: "v1.27.0", Label: "v1.27.0", Semver: &semver.Version{Major: 1, Minor: 27, Patch: 0, Extra: ""}},
			Link:               "https://github.com/symfony/polyfill-mbstring/tree/v1.27.0",
			DevOnly:            false,
			RootRequirement:    false,
			RootDevRequirement: false,
		},
		Operation: contract.Operation{
			Name:       contract.AdditionOperation,
			SemverType: contract.SemverNoUpdate,
		},
	},
	"symfony/polyfill-php73": { // ➕ 	symfony/polyfill-php73 	v1.27.0
		Package: &difftesting.TestPkgWrapper{
			Name:               "symfony/polyfill-php73",
			Abandoned:          false,
			Version:            contract.PkgVersion{Raw: "v1.27.0", Label: "v1.27.0", Semver: &semver.Version{Major: 1, Minor: 27, Patch: 0, Extra: ""}},
			Link:               "https://github.com/symfony/polyfill-php73/tree/v1.27.0",
			DevOnly:            false,
			RootRequirement:    false,
			RootDevRequirement: false,
		},
		Operation: contract.Operation{
			Name:       contract.AdditionOperation,
			SemverType: contract.SemverNoUpdate,
		},
	},
	"symfony/service-contracts": { // ➕ 	symfony/service-contracts 	v2.5.2
		Package: &difftesting.TestPkgWrapper{
			Name:               "symfony/service-contracts",
			Abandoned:          false,
			Version:            contract.PkgVersion{Raw: "v2.5.2", Label: "v2.5.2", Semver: &semver.Version{Major: 2, Minor: 5, Patch: 2, Extra: ""}},
			Link:               "https://github.com/symfony/service-contracts/tree/v2.5.2",
			DevOnly:            false,
			RootRequirement:    false,
			RootDevRequirement: false,
		},
		Operation: contract.Operation{
			Name:       contract.AdditionOperation,
			SemverType: contract.SemverNoUpdate,
		},
	},
	"symfony/string": { // ➕ 	symfony/string 	v6.2.7
		Package: &difftesting.TestPkgWrapper{
			Name:               "symfony/string",
			Abandoned:          false,
			Version:            contract.PkgVersion{Raw: "v6.2.7", Label: "v6.2.7", Semver: &semver.Version{Major: 6, Minor: 2, Patch: 7, Extra: ""}},
			Link:               "https://github.com/symfony/string/tree/v6.2.7",
			DevOnly:            false,
			RootRequirement:    false,
			RootDevRequirement: false,
		},
		Operation: contract.Operation{
			Name:       contract.AdditionOperation,
			SemverType: contract.SemverNoUpdate,
		},
	},
	"doctrine/instantiator": { // ➖ 	doctrine/instantiator 	1.4.0
		Package: &difftesting.TestPkgWrapper{
			Name:               "doctrine/instantiator",
			Abandoned:          false,
			Version:            contract.PkgVersion{Raw: "1.4.0", Label: "1.4.0", Semver: &semver.Version{Major: 1, Minor: 4, Patch: 0, Extra: ""}},
			Link:               "https://github.com/doctrine/instantiator/tree/1.4.0",
			DevOnly:            false,
			RootRequirement:    false,
			RootDevRequirement: false,
		},
		Operation: contract.Operation{
			Name:       contract.RemovalOperation,
			SemverType: contract.SemverNoUpdate,
		},
	},
	"myclabs/deep-copy": { // ➖ 	myclabs/deep-copy 	1.10.2
		Package: &difftesting.TestPkgWrapper{
			Name:               "myclabs/deep-copy",
			Abandoned:          false,
			Version:            contract.PkgVersion{Raw: "1.10.2", Label: "1.10.2", Semver: &semver.Version{Major: 1, Minor: 10, Patch: 2, Extra: ""}},
			Link:               "https://github.com/myclabs/DeepCopy/tree/1.10.2",
			DevOnly:            false,
			RootRequirement:    false,
			RootDevRequirement: false,
		},
		Operation: contract.Operation{
			Name:       contract.RemovalOperation,
			SemverType: contract.SemverNoUpdate,
		},
	},
	"nikic/php-parser": { // ➖ 	nikic/php-parser 	v4.12.0
		Package: &difftesting.TestPkgWrapper{
			Name:               "nikic/php-parser",
			Abandoned:          false,
			Version:            contract.PkgVersion{Raw: "v4.12.0", Label: "v4.12.0", Semver: &semver.Version{Major: 4, Minor: 12, Patch: 0, Extra: ""}},
			Link:               "https://github.com/nikic/PHP-Parser/tree/v4.12.0",
			DevOnly:            false,
			RootRequirement:    false,
			RootDevRequirement: false,
		},
		Operation: contract.Operation{
			Name:       contract.RemovalOperation,
			SemverType: contract.SemverNoUpdate,
		},
	},
	"phar-io/manifest": { // ➖ 	phar-io/manifest 	2.0.3
		Package: &difftesting.TestPkgWrapper{
			Name:               "phar-io/manifest",
			Abandoned:          false,
			Version:            contract.PkgVersion{Raw: "2.0.3", Label: "2.0.3", Semver: &semver.Version{Major: 2, Minor: 0, Patch: 3, Extra: ""}},
			Link:               "https://github.com/phar-io/manifest/tree/2.0.3",
			DevOnly:            false,
			RootRequirement:    false,
			RootDevRequirement: false,
		},
		Operation: contract.Operation{
			Name:       contract.RemovalOperation,
			SemverType: contract.SemverNoUpdate,
		},
	},
	"phar-io/version": { // ➖ 	phar-io/version 	3.1.0
		Package: &difftesting.TestPkgWrapper{
			Name:               "phar-io/version",
			Abandoned:          false,
			Version:            contract.PkgVersion{Raw: "3.1.0", Label: "3.1.0", Semver: &semver.Version{Major: 3, Minor: 1, Patch: 0, Extra: ""}},
			Link:               "https://github.com/phar-io/version/tree/3.1.0",
			DevOnly:            false,
			RootRequirement:    false,
			RootDevRequirement: false,
		},
		Operation: contract.Operation{
			Name:       contract.RemovalOperation,
			SemverType: contract.SemverNoUpdate,
		},
	},
	"phpdocumentor/reflection-common": { // ➖ 	phpdocumentor/reflection-common 	2.2.0
		Package: &difftesting.TestPkgWrapper{
			Name:               "phpdocumentor/reflection-common",
			Abandoned:          false,
			Version:            contract.PkgVersion{Raw: "2.2.0", Label: "2.2.0", Semver: &semver.Version{Major: 2, Minor: 2, Patch: 0, Extra: ""}},
			Link:               "https://github.com/phpDocumentor/ReflectionCommon/tree/2.x",
			DevOnly:            false,
			RootRequirement:    false,
			RootDevRequirement: false,
		},
		Operation: contract.Operation{
			Name:       contract.RemovalOperation,
			SemverType: contract.SemverNoUpdate,
		},
	},
	"phpdocumentor/reflection-docblock": { // ➖ 	phpdocumentor/reflection-docblock 	5.2.2
		Package: &difftesting.TestPkgWrapper{
			Name:               "phpdocumentor/reflection-docblock",
			Abandoned:          false,
			Version:            contract.PkgVersion{Raw: "5.2.2", Label: "5.2.2", Semver: &semver.Version{Major: 5, Minor: 2, Patch: 2, Extra: ""}},
			Link:               "https://github.com/phpDocumentor/ReflectionDocBlock/tree/master",
			DevOnly:            false,
			RootRequirement:    false,
			RootDevRequirement: false,
		},
		Operation: contract.Operation{
			Name:       contract.RemovalOperation,
			SemverType: contract.SemverNoUpdate,
		},
	},
	"phpdocumentor/type-resolver": { // ➖ 	phpdocumentor/type-resolver 	1.4.0
		Package: &difftesting.TestPkgWrapper{
			Name:               "phpdocumentor/type-resolver",
			Abandoned:          false,
			Version:            contract.PkgVersion{Raw: "1.4.0", Label: "1.4.0", Semver: &semver.Version{Major: 1, Minor: 4, Patch: 0, Extra: ""}},
			Link:               "https://github.com/phpDocumentor/TypeResolver/tree/1.4.0",
			DevOnly:            false,
			RootRequirement:    false,
			RootDevRequirement: false,
		},
		Operation: contract.Operation{
			Name:       contract.RemovalOperation,
			SemverType: contract.SemverNoUpdate,
		},
	},
	"phpspec/prophecy": { // ➖ 	phpspec/prophecy 	1.13.0
		Package: &difftesting.TestPkgWrapper{
			Name:               "phpspec/prophecy",
			Abandoned:          false,
			Version:            contract.PkgVersion{Raw: "1.13.0", Label: "1.13.0", Semver: &semver.Version{Major: 1, Minor: 13, Patch: 0, Extra: ""}},
			Link:               "https://github.com/phpspec/prophecy/tree/1.13.0",
			DevOnly:            false,
			RootRequirement:    false,
			RootDevRequirement: false,
		},
		Operation: contract.Operation{
			Name:       contract.RemovalOperation,
			SemverType: contract.SemverNoUpdate,
		},
	},
	"phpunit/php-code-coverage": { // ➖ 	phpunit/php-code-coverage 	9.2.6
		Package: &difftesting.TestPkgWrapper{
			Name:               "phpunit/php-code-coverage",
			Abandoned:          false,
			Version:            contract.PkgVersion{Raw: "9.2.6", Label: "9.2.6", Semver: &semver.Version{Major: 9, Minor: 2, Patch: 6, Extra: ""}},
			Link:               "https://github.com/sebastianbergmann/php-code-coverage/tree/9.2.6",
			DevOnly:            false,
			RootRequirement:    false,
			RootDevRequirement: false,
		},
		Operation: contract.Operation{
			Name:       contract.RemovalOperation,
			SemverType: contract.SemverNoUpdate,
		},
	},
	"phpunit/php-file-iterator": { // ➖ 	phpunit/php-file-iterator 	3.0.5
		Package: &difftesting.TestPkgWrapper{
			Name:               "phpunit/php-file-iterator",
			Abandoned:          false,
			Version:            contract.PkgVersion{Raw: "3.0.5", Label: "3.0.5", Semver: &semver.Version{Major: 3, Minor: 0, Patch: 5, Extra: ""}},
			Link:               "https://github.com/sebastianbergmann/php-file-iterator/tree/3.0.5",
			DevOnly:            false,
			RootRequirement:    false,
			RootDevRequirement: false,
		},
		Operation: contract.Operation{
			Name:       contract.RemovalOperation,
			SemverType: contract.SemverNoUpdate,
		},
	},
	"phpunit/php-invoker": { // ➖ 	phpunit/php-invoker 	3.1.1
		Package: &difftesting.TestPkgWrapper{
			Name:               "phpunit/php-invoker",
			Abandoned:          false,
			Version:            contract.PkgVersion{Raw: "3.1.1", Label: "3.1.1", Semver: &semver.Version{Major: 3, Minor: 1, Patch: 1, Extra: ""}},
			Link:               "https://github.com/sebastianbergmann/php-invoker/tree/3.1.1",
			DevOnly:            false,
			RootRequirement:    false,
			RootDevRequirement: false,
		},
		Operation: contract.Operation{
			Name:       contract.RemovalOperation,
			SemverType: contract.SemverNoUpdate,
		},
	},
	"phpunit/php-text-template": { // ➖ 	phpunit/php-text-template 	2.0.4
		Package: &difftesting.TestPkgWrapper{
			Name:               "phpunit/php-text-template",
			Abandoned:          false,
			Version:            contract.PkgVersion{Raw: "2.0.4", Label: "2.0.4", Semver: &semver.Version{Major: 2, Minor: 0, Patch: 4, Extra: ""}},
			Link:               "https://github.com/sebastianbergmann/php-text-template/tree/2.0.4",
			DevOnly:            false,
			RootRequirement:    false,
			RootDevRequirement: false,
		},
		Operation: contract.Operation{
			Name:       contract.RemovalOperation,
			SemverType: contract.SemverNoUpdate,
		},
	},
	"phpunit/php-timer": { // ➖ 	phpunit/php-timer 	5.0.3
		Package: &difftesting.TestPkgWrapper{
			Name:               "phpunit/php-timer",
			Abandoned:          false,
			Version:            contract.PkgVersion{Raw: "5.0.3", Label: "5.0.3", Semver: &semver.Version{Major: 5, Minor: 0, Patch: 3, Extra: ""}},
			Link:               "https://github.com/sebastianbergmann/php-timer/tree/5.0.3",
			DevOnly:            false,
			RootRequirement:    false,
			RootDevRequirement: false,
		},
		Operation: contract.Operation{
			Name:       contract.RemovalOperation,
			SemverType: contract.SemverNoUpdate,
		},
	},
	"phpunit/phpunit": { // ➖ 	phpunit/phpunit 	9.3.0
		Package: &difftesting.TestPkgWrapper{
			Name:               "phpunit/phpunit",
			Abandoned:          false,
			Version:            contract.PkgVersion{Raw: "9.3.0", Label: "9.3.0", Semver: &semver.Version{Major: 9, Minor: 3, Patch: 0, Extra: ""}},
			Link:               "https://github.com/sebastianbergmann/phpunit/tree/9.3",
			DevOnly:            false,
			RootRequirement:    false,
			RootDevRequirement: true,
		},
		Operation: contract.Operation{
			Name:       contract.RemovalOperation,
			SemverType: contract.SemverNoUpdate,
		},
	},
	"sebastian/code-unit-reverse-lookup": { // ➖ 	sebastian/code-unit-reverse-lookup 	2.0.3
		Package: &difftesting.TestPkgWrapper{
			Name:               "sebastian/code-unit-reverse-lookup",
			Abandoned:          false,
			Version:            contract.PkgVersion{Raw: "2.0.3", Label: "2.0.3", Semver: &semver.Version{Major: 2, Minor: 0, Patch: 3, Extra: ""}},
			Link:               "https://github.com/sebastianbergmann/code-unit-reverse-lookup/tree/2.0.3",
			DevOnly:            false,
			RootRequirement:    false,
			RootDevRequirement: false,
		},
		Operation: contract.Operation{
			Name:       contract.RemovalOperation,
			SemverType: contract.SemverNoUpdate,
		},
	},
	"sebastian/comparator": { // ➖ 	sebastian/comparator 	4.0.6
		Package: &difftesting.TestPkgWrapper{
			Name:               "sebastian/comparator",
			Abandoned:          false,
			Version:            contract.PkgVersion{Raw: "4.0.6", Label: "4.0.6", Semver: &semver.Version{Major: 4, Minor: 0, Patch: 6, Extra: ""}},
			Link:               "https://github.com/sebastianbergmann/comparator/tree/4.0.6",
			DevOnly:            false,
			RootRequirement:    false,
			RootDevRequirement: false,
		},
		Operation: contract.Operation{
			Name:       contract.RemovalOperation,
			SemverType: contract.SemverNoUpdate,
		},
	},
	"sebastian/complexity": { // ➖ 	sebastian/complexity 	2.0.2
		Package: &difftesting.TestPkgWrapper{
			Name:               "sebastian/complexity",
			Abandoned:          false,
			Version:            contract.PkgVersion{Raw: "2.0.2", Label: "2.0.2", Semver: &semver.Version{Major: 2, Minor: 0, Patch: 2, Extra: ""}},
			Link:               "https://github.com/sebastianbergmann/complexity/tree/2.0.2",
			DevOnly:            false,
			RootRequirement:    false,
			RootDevRequirement: false,
		},
		Operation: contract.Operation{
			Name:       contract.RemovalOperation,
			SemverType: contract.SemverNoUpdate,
		},
	},
	"sebastian/environment": { // ➖ 	sebastian/environment 	5.1.3
		Package: &difftesting.TestPkgWrapper{
			Name:               "sebastian/environment",
			Abandoned:          false,
			Version:            contract.PkgVersion{Raw: "5.1.3", Label: "5.1.3", Semver: &semver.Version{Major: 5, Minor: 1, Patch: 3, Extra: ""}},
			Link:               "https://github.com/sebastianbergmann/environment/tree/5.1.3",
			DevOnly:            false,
			RootRequirement:    false,
			RootDevRequirement: false,
		},
		Operation: contract.Operation{
			Name:       contract.RemovalOperation,
			SemverType: contract.SemverNoUpdate,
		},
	},
	"sebastian/exporter": { // ➖ 	sebastian/exporter 	4.0.3
		Package: &difftesting.TestPkgWrapper{
			Name:               "sebastian/exporter",
			Abandoned:          false,
			Version:            contract.PkgVersion{Raw: "4.0.3", Label: "4.0.3", Semver: &semver.Version{Major: 4, Minor: 0, Patch: 3, Extra: ""}},
			Link:               "https://github.com/sebastianbergmann/exporter/tree/4.0.3",
			DevOnly:            false,
			RootRequirement:    false,
			RootDevRequirement: false,
		},
		Operation: contract.Operation{
			Name:       contract.RemovalOperation,
			SemverType: contract.SemverNoUpdate,
		},
	},
	"sebastian/global-state": { // ➖ 	sebastian/global-state 	5.0.3
		Package: &difftesting.TestPkgWrapper{
			Name:               "sebastian/global-state",
			Abandoned:          false,
			Version:            contract.PkgVersion{Raw: "5.0.3", Label: "5.0.3", Semver: &semver.Version{Major: 5, Minor: 0, Patch: 3, Extra: ""}},
			Link:               "https://github.com/sebastianbergmann/global-state/tree/5.0.3",
			DevOnly:            false,
			RootRequirement:    false,
			RootDevRequirement: false,
		},
		Operation: contract.Operation{
			Name:       contract.RemovalOperation,
			SemverType: contract.SemverNoUpdate,
		},
	},
	"sebastian/lines-of-code": { // ➖ 	sebastian/lines-of-code 	1.0.3
		Package: &difftesting.TestPkgWrapper{
			Name:               "sebastian/lines-of-code",
			Abandoned:          false,
			Version:            contract.PkgVersion{Raw: "1.0.3", Label: "1.0.3", Semver: &semver.Version{Major: 1, Minor: 0, Patch: 3, Extra: ""}},
			Link:               "https://github.com/sebastianbergmann/lines-of-code/tree/1.0.3",
			DevOnly:            false,
			RootRequirement:    false,
			RootDevRequirement: false,
		},
		Operation: contract.Operation{
			Name:       contract.RemovalOperation,
			SemverType: contract.SemverNoUpdate,
		},
	},
	"sebastian/object-enumerator": { // ➖ 	sebastian/object-enumerator 	4.0.4
		Package: &difftesting.TestPkgWrapper{
			Name:               "sebastian/object-enumerator",
			Abandoned:          false,
			Version:            contract.PkgVersion{Raw: "4.0.4", Label: "4.0.4", Semver: &semver.Version{Major: 4, Minor: 0, Patch: 4, Extra: ""}},
			Link:               "https://github.com/sebastianbergmann/object-enumerator/tree/4.0.4",
			DevOnly:            false,
			RootRequirement:    false,
			RootDevRequirement: false,
		},
		Operation: contract.Operation{
			Name:       contract.RemovalOperation,
			SemverType: contract.SemverNoUpdate,
		},
	},
	"sebastian/object-reflector": { // ➖ 	sebastian/object-reflector 	2.0.4
		Package: &difftesting.TestPkgWrapper{
			Name:               "sebastian/object-reflector",
			Abandoned:          false,
			Version:            contract.PkgVersion{Raw: "2.0.4", Label: "2.0.4", Semver: &semver.Version{Major: 2, Minor: 0, Patch: 4, Extra: ""}},
			Link:               "https://github.com/sebastianbergmann/object-reflector/tree/2.0.4",
			DevOnly:            false,
			RootRequirement:    false,
			RootDevRequirement: false,
		},
		Operation: contract.Operation{
			Name:       contract.RemovalOperation,
			SemverType: contract.SemverNoUpdate,
		},
	},
	"sebastian/recursion-context": { // ➖ 	sebastian/recursion-context 	4.0.4
		Package: &difftesting.TestPkgWrapper{
			Name:               "sebastian/recursion-context",
			Abandoned:          false,
			Version:            contract.PkgVersion{Raw: "4.0.4", Label: "4.0.4", Semver: &semver.Version{Major: 4, Minor: 0, Patch: 4, Extra: ""}},
			Link:               "https://github.com/sebastianbergmann/recursion-context/tree/4.0.4",
			DevOnly:            false,
			RootRequirement:    false,
			RootDevRequirement: false,
		},
		Operation: contract.Operation{
			Name:       contract.RemovalOperation,
			SemverType: contract.SemverNoUpdate,
		},
	},
	"sebastian/resource-operations": { // ➖ 	sebastian/resource-operations 	3.0.3
		Package: &difftesting.TestPkgWrapper{
			Name:               "sebastian/resource-operations",
			Abandoned:          false,
			Version:            contract.PkgVersion{Raw: "3.0.3", Label: "3.0.3", Semver: &semver.Version{Major: 3, Minor: 0, Patch: 3, Extra: ""}},
			Link:               "https://github.com/sebastianbergmann/resource-operations/tree/3.0.3",
			DevOnly:            false,
			RootRequirement:    false,
			RootDevRequirement: false,
		},
		Operation: contract.Operation{
			Name:       contract.RemovalOperation,
			SemverType: contract.SemverNoUpdate,
		},
	},
	"sebastian/type": { // ➖ 	sebastian/type 	2.3.4
		Package: &difftesting.TestPkgWrapper{
			Name:               "sebastian/type",
			Abandoned:          false,
			Version:            contract.PkgVersion{Raw: "2.3.4", Label: "2.3.4", Semver: &semver.Version{Major: 2, Minor: 3, Patch: 4, Extra: ""}},
			Link:               "https://github.com/sebastianbergmann/type/tree/2.3.4",
			DevOnly:            false,
			RootRequirement:    false,
			RootDevRequirement: false,
		},
		Operation: contract.Operation{
			Name:       contract.RemovalOperation,
			SemverType: contract.SemverNoUpdate,
		},
	},
	"sebastian/version": { // ➖ 	sebastian/version 	3.0.2
		Package: &difftesting.TestPkgWrapper{
			Name:               "sebastian/version",
			Abandoned:          false,
			Version:            contract.PkgVersion{Raw: "3.0.2", Label: "3.0.2", Semver: &semver.Version{Major: 3, Minor: 0, Patch: 2, Extra: ""}},
			Link:               "https://github.com/sebastianbergmann/version/tree/3.0.2",
			DevOnly:            false,
			RootRequirement:    false,
			RootDevRequirement: false,
		},
		Operation: contract.Operation{
			Name:       contract.RemovalOperation,
			SemverType: contract.SemverNoUpdate,
		},
	},
	"theseer/tokenizer": { // ➖ 	theseer/tokenizer 	1.2.1
		Package: &difftesting.TestPkgWrapper{
			Name:               "theseer/tokenizer",
			Abandoned:          false,
			Version:            contract.PkgVersion{Raw: "1.2.1", Label: "1.2.1", Semver: &semver.Version{Major: 1, Minor: 2, Patch: 1, Extra: ""}},
			Link:               "https://github.com/theseer/tokenizer/tree/1.2.1",
			DevOnly:            false,
			RootRequirement:    false,
			RootDevRequirement: false,
		},
		Operation: contract.Operation{
			Name:       contract.RemovalOperation,
			SemverType: contract.SemverNoUpdate,
		},
	},
	"twig/twig": { // ➖ 	twig/twig 	v1.44.4
		Package: &difftesting.TestPkgWrapper{
			Name:               "twig/twig",
			Abandoned:          false,
			Version:            contract.PkgVersion{Raw: "v1.44.4", Label: "v1.44.4", Semver: &semver.Version{Major: 1, Minor: 44, Patch: 4, Extra: ""}},
			Link:               "https://github.com/twigphp/Twig/tree/v1.44.4",
			DevOnly:            false,
			RootRequirement:    false,
			RootDevRequirement: false,
		},
		Operation: contract.Operation{
			Name:       contract.RemovalOperation,
			SemverType: contract.SemverNoUpdate,
		},
	},
	"webmozart/assert": { // ➖ 	webmozart/assert 	1.10.0
		Package: &difftesting.TestPkgWrapper{
			Name:               "webmozart/assert",
			Abandoned:          false,
			Version:            contract.PkgVersion{Raw: "1.10.0", Label: "1.10.0", Semver: &semver.Version{Major: 1, Minor: 10, Patch: 0, Extra: ""}},
			Link:               "https://github.com/webmozarts/assert/tree/1.10.0",
			DevOnly:            false,
			RootRequirement:    false,
			RootDevRequirement: false,
		},
		Operation: contract.Operation{
			Name:       contract.RemovalOperation,
			SemverType: contract.SemverNoUpdate,
		},
	},
	"yoanm/init-php-repository": { // ➖ 	yoanm/init-php-repository 	dev-master#02c0922❗
		Package: &difftesting.TestPkgWrapper{
			Name:               "yoanm/init-php-repository",
			Abandoned:          false,
			Version:            contract.PkgVersion{Raw: "02c0922c4691e02b606c7cfe4cf8978233b1e978", Label: "dev-master#02c0922", Semver: nil},
			Link:               "https://github.com/yoanm/initPhpRepository/tree/v0.2.0",
			DevOnly:            false,
			RootRequirement:    false,
			RootDevRequirement: true,
		},
		Operation: contract.Operation{
			Name:       contract.RemovalOperation,
			SemverType: contract.SemverNoUpdate,
		},
	},
	"squizlabs/php_codesniffer": { // ➖ 	squizlabs/php_codesniffer 	3.6.2
		Package: &difftesting.TestPkgWrapper{
			Name:               "squizlabs/php_codesniffer",
			Abandoned:          false,
			Version:            contract.PkgVersion{Raw: "2acf168", Label: "2.9.x-dev#2acf168", Semver: nil},
			Link:               "https://github.com/squizlabs/PHP_CodeSniffer/wiki",
			DevOnly:            true,
			RootRequirement:    false,
			RootDevRequirement: true,
		},
		Operation: contract.Operation{
			Name:       contract.NoChangeOperation,
			SemverType: contract.SemverNoUpdate,
		},
	},
}
