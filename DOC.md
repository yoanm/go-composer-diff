# compdiff

## Functions

### func [BuildMap](/package_map.go#L24)

`func BuildMap(reqData *ReqFile, lockData *LockFile) (contract.PackageMap, error)`

BuildMap creates an efficient lookup map for composer packages.

### func [BuildMapFromBytes](/package_map.go#L9)

`func BuildMapFromBytes(reqContent, lockContent []byte) (contract.PackageMap, error)`

### func [Diff](/main.go#L32)

`func Diff(previous, current *Input) (contract.DiffMap, error)`

### func [FileDiff](/main.go#L11)

`func FileDiff(previous, current *FileInput) (contract.DiffMap, error)`

## Types

### type [ComposerPackageWrapper](/wrapper_types.go#L5)

`type ComposerPackageWrapper struct { ... }`

#### func (*ComposerPackageWrapper) [GetLink](/wrapper_types.go#L30)

`func (w *ComposerPackageWrapper) GetLink() string`

#### func (*ComposerPackageWrapper) [GetName](/wrapper_types.go#L18)

`func (w *ComposerPackageWrapper) GetName() string`

#### func (*ComposerPackageWrapper) [GetVersion](/wrapper_types.go#L26)

`func (w *ComposerPackageWrapper) GetVersion() contract.PkgVersion`

#### func (*ComposerPackageWrapper) [IsAbandoned](/wrapper_types.go#L22)

`func (w *ComposerPackageWrapper) IsAbandoned() bool`

#### func (*ComposerPackageWrapper) [IsDevOnly](/wrapper_types.go#L34)

`func (w *ComposerPackageWrapper) IsDevOnly() bool`

#### func (*ComposerPackageWrapper) [IsRootDevRequirement](/wrapper_types.go#L42)

`func (w *ComposerPackageWrapper) IsRootDevRequirement() bool`

#### func (*ComposerPackageWrapper) [IsRootRequirement](/wrapper_types.go#L38)

`func (w *ComposerPackageWrapper) IsRootRequirement() bool`

### type [FileInput](/input.go#L11)

`type FileInput struct { ... }`

### type [Input](/input.go#L3)

`type Input struct { ... }`

### type [InvalidFormatError](/parser.go#L57)

`type InvalidFormatError struct { ... }`

InvalidFormatError indicates the JSON is valid but doesn't match expected structure.

#### func (InvalidFormatError) [Error](/parser.go#L61)

`func (e InvalidFormatError) Error() string`

### type [InvalidJSONError](/parser.go#L48)

`type InvalidJSONError struct { ... }`

InvalidJSONError indicates the input JSON is malformed.

#### func (InvalidJSONError) [Error](/parser.go#L52)

`func (e InvalidJSONError) Error() string`

### type [LockFile](/types.go#L10)

`type LockFile struct { ... }`

LockFile represents the structure of a composer.lock file.

#### func [ParseLock](/parser.go#L8)

`func ParseLock(data []byte) (*LockFile, error)`

ParseLock parses a composer.lock file from JSON bytes.

### type [Package](/types.go#L16)

`type Package struct { ... }`

Package represents a single package entry in composer.lock.

### type [ReqFile](/types.go#L4)

`type ReqFile struct { ... }`

ReqFile represents the structure of a composer.json file (composer requirement).

#### func [ParseReq](/parser.go#L28)

`func ParseReq(data []byte) (*ReqFile, error)`

ParseReq parses a composer.json file from JSON bytes.

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
