# compdiff

Package compdiff provides utilities for comparing PHP Composer lock files and identifying package differences.

## Functions

### func [BuildMap](/package_map.go#L73)

`func BuildMap(reqData *ReqFile, lockData *LockFile) (contract.PackageMap, error)`

BuildMap creates an efficient lookup map for Composer packages from structured ReqFile and LockFile.

This function takes already-parsed ReqFile and LockFile structures and builds a PackageMap.
Both regular and dev packages are included in the map. Use this when you have already parsed
the files, or use BuildMapFromBytes to parse and build in one step.

Parameters:

```diff
- reqData: Parsed composer.json structure containing explicit package requirements
- lockData: Parsed composer.lock structure containing installed package information
```

Returns:

```diff
- contract.PackageMap: A map of package names to their wrapper objects
- error: Always returns nil (errors would have occurred during parsing)
```

Example:

```go
req, _ := ParseReq(reqContent)
lock, _ := ParseLock(lockContent)
packageMap, _ := compdiff.BuildMap(req, lock)
for name, pkg := range packageMap {
    fmt.Printf("%s (dev: %v): %s\n", name, pkg.IsDevOnly(), pkg.GetVersion().Raw)
}
```

### func [BuildMapFromBytes](/package_map.go#L33)

`func BuildMapFromBytes(reqContent, lockContent []byte) (contract.PackageMap, error)`

BuildMapFromBytes creates an efficient lookup map for Composer packages from raw JSON content.

This function parses the composer.json (requirement) and composer.lock (lock) files
and builds a PackageMap for use with the depsdiff comparison functions.

Parameters:

```diff
- reqContent: The raw JSON content of a composer.json file
- lockContent: The raw JSON content of a composer.lock file
```

Returns:

```diff
- contract.PackageMap: A map of package names to their wrapper objects
- error: Non-nil if parsing or building fails
```

Example:

```go
reqContent, _ := os.ReadFile("composer.json")
lockContent, _ := os.ReadFile("composer.lock")
packageMap, err := compdiff.BuildMapFromBytes(reqContent, lockContent)
if err != nil {
    log.Fatal(err)
}
for name, pkg := range packageMap {
    fmt.Printf("%s: %s\n", name, pkg.GetVersion().Label)
}
```

### func [Diff](/main.go#L89)

`func Diff(previous, current *Input) (contract.DiffMap, error)`

Diff compares two sets of Composer requirement and lock file contents and returns detailed package differences.

It takes parsed JSON content (not file paths) and processes them to identify all package additions,
removals, and version changes. This function is useful when you already have the file contents
in memory or from other sources.

Parameters:

```diff
- previous: Input containing the previous composer.json and composer.lock file contents as byte slices
- current: Input containing the current composer.json and composer.lock file contents as byte slices
```

Returns:

```diff
- DiffMap: A map where keys are package names and values contain PackageChange information
- error: Non-nil if JSON parsing or package map building fails
```

Example:

```go
previousReq, _ := os.ReadFile("old/composer.json")
previousLock, _ := os.ReadFile("old/composer.lock")
currentReq, _ := os.ReadFile("new/composer.json")
currentLock, _ := os.ReadFile("new/composer.lock")

changes, err := compdiff.Diff(
    &Input{Requirement: previousReq, Lock: previousLock},
    &Input{Requirement: currentReq, Lock: currentLock},
)
if err != nil {
    log.Fatal(err)
}
for pkgName, change := range changes {
    fmt.Printf("%s: %s\n", pkgName, change.Operation.Name)
}
```

### func [FileDiff](/main.go#L37)

`func FileDiff(previous, current *FileInput) (contract.DiffMap, error)`

FileDiff compares two Composer lock files from file paths and returns detailed information about package differences.

It reads both the requirement (composer.json) and lock (composer.lock) files from disk
and processes them to generate a comprehensive diff.

Parameters:

```diff
- previous: FileInput containing paths to the previous composer.json and composer.lock files
- current: FileInput containing paths to the current composer.json and composer.lock files
```

Returns:

```diff
- DiffMap: A map where keys are package names and values contain PackageChange information
- error: Non-nil if file reading or parsing fails
```

Example:

```go
changes, err := compdiff.FileDiff(
    &FileInput{Lock: "old/composer.lock", Requirement: "old/composer.json"},
    &FileInput{Lock: "new/composer.lock", Requirement: "new/composer.json"},
)
if err != nil {
    log.Fatal(err)
}
for pkgName, change := range changes {
    fmt.Printf("%s: %s\n", pkgName, change.Operation.Name)
}
```

## Types

### type [ComposerPackageWrapper](/wrapper_types.go#L10)

`type ComposerPackageWrapper struct { ... }`

ComposerPackageWrapper wraps Composer package information and implements the contract.PkgWrapper interface.

It provides access to package metadata including name, version, status, and various classifications
(root requirement, dev-only, abandoned status, etc.). This type is the main way to interact with
package information when working with the diff results.

#### func (*ComposerPackageWrapper) [GetLink](/wrapper_types.go#L42)

`func (w *ComposerPackageWrapper) GetLink() string`

GetLink returns a URL link to the package (e.g., homepage or documentation).

Returns an empty string if no link is available.

#### func (*ComposerPackageWrapper) [GetName](/wrapper_types.go#L24)

`func (w *ComposerPackageWrapper) GetName() string`

GetName returns the package name (e.g., "vendor/package-name").

#### func (*ComposerPackageWrapper) [GetVersion](/wrapper_types.go#L35)

`func (w *ComposerPackageWrapper) GetVersion() contract.PkgVersion`

GetVersion returns the package version information including raw version, label,
and parsed semantic version if applicable.

#### func (*ComposerPackageWrapper) [IsAbandoned](/wrapper_types.go#L29)

`func (w *ComposerPackageWrapper) IsAbandoned() bool`

IsAbandoned returns whether this package has been marked as abandoned by its maintainer.

#### func (*ComposerPackageWrapper) [IsDevOnly](/wrapper_types.go#L51)

`func (w *ComposerPackageWrapper) IsDevOnly() bool`

IsDevOnly returns whether this package is only required for development (not in production).

A package is dev-only if it appears only in the packages-dev section of composer.lock.
However, a dev-only package can be a root requirement, a root dev requirement or a transitive dependency. If
you are looking for requirement vs dev requirement see IsRootRequirement and IsRootDevRequirement.

#### func (*ComposerPackageWrapper) [IsRootDevRequirement](/wrapper_types.go#L65)

`func (w *ComposerPackageWrapper) IsRootDevRequirement() bool`

IsRootDevRequirement returns whether this package is explicitly required for development in composer.json.

Returns true if the package is listed in the "require-dev" section of composer.json.

#### func (*ComposerPackageWrapper) [IsRootRequirement](/wrapper_types.go#L58)

`func (w *ComposerPackageWrapper) IsRootRequirement() bool`

IsRootRequirement returns whether this package is explicitly required in composer.json.

Returns true if the package is listed in the "require" section of composer.json.

### type [FileInput](/input.go#L19)

`type FileInput struct { ... }`

FileInput holds the file paths to Composer requirement and lock files.

This is used when you want to compare lock files from disk locations.
Use Input if you already have file contents in memory.

### type [Input](/input.go#L7)

`type Input struct { ... }`

Input holds the byte contents of Composer requirement and lock files.

This is used when you have the file contents in memory (e.g., read from disk or network).
Use FileInput if you have file paths instead.

### type [InvalidFormatError](/parser.go#L105)

`type InvalidFormatError struct { ... }`

InvalidFormatError indicates that the JSON is syntactically valid but doesn't match
the expected Composer file structure.

This error is returned when the JSON is missing required fields or has an unexpected schema.

#### func (InvalidFormatError) [Error](/parser.go#L110)

`func (e InvalidFormatError) Error() string`

Error returns a string representation of the InvalidFormatError.

### type [InvalidJSONError](/parser.go#L92)

`type InvalidJSONError struct { ... }`

InvalidJSONError indicates that the input JSON content is malformed and cannot be parsed.

This error is returned by ParseLock or ParseReq when the JSON is syntactically invalid.

#### func (InvalidJSONError) [Error](/parser.go#L97)

`func (e InvalidJSONError) Error() string`

Error returns a string representation of the InvalidJSONError.

### type [LockFile](/types.go#L17)

`type LockFile struct { ... }`

LockFile represents the structure of a composer.lock file.

It contains the complete list of installed packages and their versions, separated
into regular packages and development-only packages. The lock file is the authoritative
source for the exact versions and references used in a project.

#### func [ParseLock](/parser.go#L30)

`func ParseLock(data []byte) (*LockFile, error)`

ParseLock parses a composer.lock file from JSON bytes and returns a structured LockFile.

It validates that the JSON is well-formed and contains the required fields
(either 'packages' or 'packages-dev' or both).

Parameters:

```go
- data: The raw JSON content of a composer.lock file as bytes
```

Returns:

```diff
- *LockFile: Pointer to the parsed lock file structure
- error: InvalidJSONError if JSON is malformed, InvalidFormatError if structure is invalid, or nil on success
```

Example:

```go
lockContent, err := os.ReadFile("composer.lock")
if err != nil {
    log.Fatal(err)
}
lock, err := ParseLock(lockContent)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Installed %d packages\n", len(lock.Packages))
```

### type [Package](/types.go#L27)

`type Package struct { ... }`

Package represents a single package entry in composer.lock.

It contains all metadata about an installed package including its version,
source information (e.g., repository reference), distribution information,
documentation links, and whether it's been abandoned.

### type [ReqFile](/types.go#L7)

`type ReqFile struct { ... }`

ReqFile represents the structure of a composer.json file (composer requirement).

It contains the explicit package requirements defined by the project, separated into
regular dependencies and development-only dependencies.

#### func [ParseReq](/parser.go#L71)

`func ParseReq(data []byte) (*ReqFile, error)`

ParseReq parses a composer.json file from JSON bytes and returns a structured ReqFile.

It validates that the JSON is well-formed and contains the required fields
(either 'require' or 'require-dev' or both).

Parameters:

```go
- data: The raw JSON content of a composer.json file as bytes
```

Returns:

```diff
- *ReqFile: Pointer to the parsed requirement file structure
- error: InvalidJSONError if JSON is malformed, InvalidFormatError if structure is invalid, or nil on success
```

Example:

```go
reqContent, err := os.ReadFile("composer.json")
if err != nil {
    log.Fatal(err)
}
req, err := ParseReq(reqContent)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Direct requirements: %d\n", len(req.Require))
```

### type [Support](/types.go#L48)

`type Support struct { ... }`

Support contains links to documentation and support.

It provides URLs to various documentation and support resources for a package.

### type [VersionReference](/types.go#L41)

`type VersionReference struct { ... }`

VersionReference contains the reference (commit hash or tag).

It points to the specific commit or tag in the source repository that corresponds
to the package version.

---
Readme created from Go doc with [goreadme](https://github.com/posener/goreadme)
