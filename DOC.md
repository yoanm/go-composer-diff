# compdiff

## Functions

### func [BuildMap](/wrapper.go#L26)

`func BuildMap(reqData *ReqFile, lockData *LockFile) (contract.PackageMap, error)`

BuildMap creates an efficient lookup map for composer packages.

### func [BuildMapFromBytes](/wrapper.go#L11)

`func BuildMapFromBytes(reqContent, lockContent []byte) (contract.PackageMap, error)`

### func [Diff](/main.go#L10)

`func Diff(previous, current *Input) (contract.DiffMap, error)`

## Types

### type [ComposerPackageWrapper](/wrapper_types.go#L5)

`type ComposerPackageWrapper struct { ... }`

#### func (*ComposerPackageWrapper) [GetLink](/wrapper_types.go#L27)

`func (w *ComposerPackageWrapper) GetLink() string`

#### func (*ComposerPackageWrapper) [GetName](/wrapper_types.go#L15)

`func (w *ComposerPackageWrapper) GetName() string`

#### func (*ComposerPackageWrapper) [GetVersion](/wrapper_types.go#L23)

`func (w *ComposerPackageWrapper) GetVersion() contract.PkgVersion`

#### func (*ComposerPackageWrapper) [IsAbandoned](/wrapper_types.go#L19)

`func (w *ComposerPackageWrapper) IsAbandoned() bool`

#### func (*ComposerPackageWrapper) [IsDevOnly](/wrapper_types.go#L31)

`func (w *ComposerPackageWrapper) IsDevOnly() bool`

#### func (*ComposerPackageWrapper) [IsRootDevRequirement](/wrapper_types.go#L39)

`func (w *ComposerPackageWrapper) IsRootDevRequirement() bool`

#### func (*ComposerPackageWrapper) [IsRootRequirement](/wrapper_types.go#L35)

`func (w *ComposerPackageWrapper) IsRootRequirement() bool`

### type [Input](/input.go#L3)

`type Input struct { ... }`

### type [InvalidFormatError](/parser.go#L18)

`type InvalidFormatError struct { ... }`

InvalidFormatError indicates the JSON is valid but doesn't match expected structure.

#### func (InvalidFormatError) [Error](/parser.go#L22)

`func (e InvalidFormatError) Error() string`

### type [InvalidJSONError](/parser.go#L9)

`type InvalidJSONError struct { ... }`

InvalidJSONError indicates the input JSON is malformed.

#### func (InvalidJSONError) [Error](/parser.go#L13)

`func (e InvalidJSONError) Error() string`

### type [LockFile](/types.go#L10)

`type LockFile struct { ... }`

LockFile represents the structure of a compdiff.lock file.

#### func [ParseLock](/parser.go#L27)

`func ParseLock(data []byte) (*LockFile, error)`

ParseLock parses a compdiff.lock file from JSON bytes.

### type [Package](/types.go#L16)

`type Package struct { ... }`

Package represents a single package entry in compdiff.lock.

### type [ReqFile](/types.go#L4)

`type ReqFile struct { ... }`

ReqFile represents the structure of a compdiff.json file (composer requirement).

#### func [ParseReq](/parser.go#L47)

`func ParseReq(data []byte) (*ReqFile, error)`

ParseReq parses a compdiff.json (composer requirement) file from JSON bytes.

### type [Support](/types.go#L32)

`type Support struct { ... }`

Support contains links to documentation and support.

### type [VersionReference](/types.go#L27)

`type VersionReference struct { ... }`

VersionReference contains the reference (commit hash or tag).

## Sub Packages

* [.tools](./.tools)

---
Readme created from Go doc with [goreadme](https://github.com/posener/goreadme)
