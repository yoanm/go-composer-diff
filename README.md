# compdiff<br/><sup><sub>Composer dependecies comparator library</sub></sup>

[![License](https://img.shields.io/github/license/yoanm/go-composer-diff.svg)](https://github.com/yoanm/go-composer-diff)
[![Code size](https://img.shields.io/github/languages/code-size/yoanm/go-composer-diff.svg)](https://github.com/yoanm/go-composer-diff)
[![Go Reference](https://pkg.go.dev/badge/github.com/yoanm/go-composer-diff.svg)](https://pkg.go.dev/github.com/yoanm/go-composer-diff)

![Dependabot Status](https://flat.badgen.net/github/dependabot/yoanm/go-composer-diff)
![Last commit](https://badgen.net/github/last-commit/yoanm/go-composer-diff)

[![Codacy Badge](https://app.codacy.com/project/badge/Grade/ebeacd3a91a74fef8a8ed4ea879ede72)](https://app.codacy.com/gh/yoanm/go-composer-diff/dashboard?utm_source=gh&utm_medium=referral&utm_content=&utm_campaign=Badge_grade)
[![Go Report Card](https://goreportcard.com/badge/github.com/yoanm/go-composer-diff?)](https://goreportcard.com/report/github.com/yoanm/go-composer-diff)

[![CI](https://github.com/yoanm/go-composer-diff/actions/workflows/CI.yml/badge.svg?branch=master)](https://github.com/yoanm/go-composer-diff/actions/workflows/CI.yml)
[![codecov](https://codecov.io/gh/yoanm/go-composer-diff/branch/master/graph/badge.svg?token=NHdwEBUFK5)](https://codecov.io/gh/yoanm/go-composer-diff)

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/yoanm/go-composer-diff)


A Go library for comparing two PHP Composer lock files and identifying package differences.

## Overview

`composer-diff` is a Go module that analyzes and compares Composer lock files to detect added, removed, and updated package dependencies. It parses `composer.lock` files and uses them in combination with `composer.json` files to provide comprehensive package analysis.

## Installation

```bash
go get github.com/yoanm/go-composer-diff
```

## Quick Start

The library provides two main entry points:

### Option 1: Compare from File Paths

```go
package main

import (
    "fmt"
    "log"

    compdiff "github.com/yoanm/go-composer-diff"
)

func main() {
    // Compare Composer lock files from disk
    changes, err := compdiff.FileDiff(
        &compdiff.FileInput{
            Lock:        "old/composer.lock",
            Requirement: "old/composer.json",
        },
        &compdiff.FileInput{
            Lock:        "new/composer.lock",
            Requirement: "new/composer.json",
        },
    )
    if err != nil {
        log.Fatal(err)
    }
    
    // Process results
    for pkgName, change := range changes {
        fmt.Printf("%s: %s\n", pkgName, change.Operation.Name)
    }
}
```

### Option 2: Compare from File Contents

```go
package main

import (
    "fmt"
    "log"
    "os"

    compdiff "github.com/yoanm/go-composer-diff"
)

func main() {
    // Read files
    prevReq, _ := os.ReadFile("old/composer.json")
    prevLock, _ := os.ReadFile("old/composer.lock")
    currReq, _ := os.ReadFile("new/composer.json")
    currLock, _ := os.ReadFile("new/composer.lock")
    
    // Compare from byte contents
    changes, err := compdiff.Diff(
        &compdiff.Input{
            Requirement: prevReq,
            Lock:        prevLock,
        },
        &compdiff.Input{
            Requirement: currReq,
            Lock:        currLock,
        },
    )
    if err != nil {
        log.Fatal(err)
    }
    
    // Process results
    for pkgName, change := range changes {
        fmt.Printf("%s: %s\n", pkgName, change.Operation.Name)
    }
}
```

## Features

- ✅ Detects added, removed, and updated packages, as well as unchanged packages
- ✅ Semantic version parsing (MAJOR/MINOR/PATCH/EXTRA)
- ✅ Update direction detection (DOWNGRADE/UPGRADE) for semver compatible versions (UNKNOWN_UPDATE otherwise)
- ✅ Support for both regular and dev dependencies
- ✅ Handles commit-based versions
- ✅ Identifies abandoned packages
- ✅ Tracks explicit root requirements vs transitive dependencies
- ✅ Efficient O(1) lookup-based comparison
- ✅ Comprehensive error handling with typed error responses

## How It Works

1. **Parse Files**: The library parses `composer.json` (requirements) and `composer.lock` (lock file)
2. **Build Package Map**: Creates an efficient lookup map containing package metadata
3. **Compare Maps**: Uses the dependency diff library to compare old and new package maps
4. **Return Results**: Provides detailed change information for each package

## API Reference

### Main Functions

#### `FileDiff(previous, current *FileInput) (contract.DiffMap, error)`

Compares two Composer projects from file paths.

**Parameters:**
- `previous`: FileInput with paths to old composer.lock and composer.json
- `current`: FileInput with paths to new composer.lock and composer.json

**Returns:**
- `contract.DiffMap`: Map of package names to their changes
- `error`: Non-nil if file reading or parsing fails

#### `Diff(previous, current *Input) (contract.DiffMap, error)`

Compares two Composer projects from byte contents.

**Parameters:**
- `previous`: Input with old composer.lock and composer.json contents
- `current`: Input with new composer.lock and composer.json contents

**Returns:**
- `contract.DiffMap`: Map of package names to their changes
- `error`: Non-nil if parsing or comparison fails

#### `BuildMapFromBytes(reqContent, lockContent []byte) (contract.PackageMap, error)`

Builds a package map from raw JSON content. Useful for custom processing.

**Parameters:**
- `reqContent`: Raw JSON bytes from composer.json
- `lockContent`: Raw JSON bytes from composer.lock

**Returns:**
- `contract.PackageMap`: Map of package names to package wrappers
- `error`: Non-nil if parsing fails

#### `BuildMap(reqData *ReqFile, lockData *LockFile) (contract.PackageMap, error)`

Builds a package map from parsed structures. Useful when you're already parsing files.

**Parameters:**
- `reqData`: Parsed ReqFile structure
- `lockData`: Parsed LockFile structure

**Returns:**
- `contract.PackageMap`: Map of package names to package wrappers
- `error`: Always nil (for API compatibility)

### Parser Functions

#### `ParseLock(data []byte) (*LockFile, error)`

Parses a composer.lock file from JSON bytes.

**Returns:**
- `error`: InvalidJSONError if JSON is malformed, InvalidFormatError if structure is invalid

#### `ParseReq(data []byte) (*ReqFile, error)`

Parses a composer.json file from JSON bytes.

**Returns:**
- `error`: InvalidJSONError if JSON is malformed, InvalidFormatError if structure is invalid

### Input Types

#### `Input`

Holds byte contents of Composer files for in-memory comparison.

```go
type Input struct {
    Lock []byte        // composer.lock contents
    Requirement []byte // composer.json contents
}
```

#### `FileInput`

Holds file paths for file-based comparison.

```go
type FileInput struct {
    Lock string        // path to composer.lock
    Requirement string // path to composer.json
}
```

### File Structure Types

#### `ReqFile`

Represents a parsed composer.json file.

```go
type ReqFile struct {
    Require    map[string]string // Regular dependencies
    RequireDev map[string]string // Development dependencies
}
```

#### `LockFile`

Represents a parsed composer.lock file.

```go
type LockFile struct {
    Packages    []Package // Installed packages
    PackagesDev []Package // Development packages
}
```

#### `Package`

Represents a package entry in the lock file.

```go
type Package struct {
    Name       string
    Version    string
    Source     *VersionReference
    Dist       *VersionReference
    Support    *Support
    Homepage   string
    Abandoned  interface{} // Can be bool or string
}
```

### Package Information Type

#### `ComposerPackageWrapper`

Provides access to package metadata. All objects in the PackageMap are ComposerPackageWrapper instances.

It implements `contract.PkgWrapper` interface from `github.com/yoanm/go-deps-diff/contract` package.

**Key Methods:**
- `GetName() string` - Package name
- `GetVersion() contract.PkgVersion` - Version information
- `IsAbandoned() bool` - Package abandoned status
- `GetLink() string` - Documentation link
- `IsDevOnly() bool` - Development-only dependency
- `IsRootRequirement() bool` - Explicitly required for prod
- `IsRootDevRequirement() bool` - Explicitly required for dev

### Error Types

#### `InvalidJSONError`

Indicates JSON parsing failed.

#### `InvalidFormatError`

Indicates JSON is valid but doesn't match Composer file structure (missing required fields).

## Examples

### Filter Changes by Operation Type

```go
changes, err := compdiff.FileDiff(prevInput, currInput)
if err != nil {
    log.Fatal(err)
}

// Find all upgraded packages
for pkgName, change := range changes {
    if change.Operation.Name == contract.UpgradeOperation {
        fmt.Printf("Upgraded: %s\n", pkgName)
    }
}

// Find all major version updates
for pkgName, change := range changes {
    if change.Operation.SemverType == contract.SemverMajorUpdate {
        fmt.Printf("Major update: %s\n", pkgName)
    }
}
```

### Analyze All Changes

```go
changes, err := compdiff.FileDiff(prevInput, currInput)
if err != nil {
    log.Fatal(err)
}

for pkgName, change := range changes {
    fmt.Printf("%s: %s (%s)\n", 
        pkgName,
        change.Operation.Name,
        change.Operation.SemverType)
    
    // Show version transition for updates
    if change.PreviousVersion.Raw != "" {
        fmt.Printf("  %s → %s\n",
            change.PreviousVersion.Label,
            change.Package.GetVersion().Label)
    }
    
    // Flag abandoned packages
    if change.Package.IsAbandoned() {
        fmt.Printf("  ⚠️  This package is abandoned\n")
    }
}
```

### Check for Security-Relevant Updates

```go
changes, err := compdiff.FileDiff(prevInput, currInput)
if err != nil {
    log.Fatal(err)
}

securityConcerns := 0
for pkgName, change := range changes {
    // Flag additions and removals
    if change.Operation.Name == contract.AdditionOperation ||
       change.Operation.Name == contract.RemovalOperation {
        fmt.Printf("⚠️  %s: %s\n", pkgName, change.Operation.Name)
        securityConcerns++
    }
    
    // Flag major/minor version updates
    switch change.Operation.SemverType {
    case contract.SemverMajorUpdate, contract.SemverMinorUpdate:
        fmt.Printf("📦 %s %s to %s\n", 
            pkgName, 
            change.Operation.Name,
            change.Package.GetVersion().Label)
    }
}

fmt.Printf("Total items to review: %d\n", securityConcerns)
```

### Process Dev Dependencies Separately

```go
changes, err := compdiff.FileDiff(prevInput, currInput)
if err != nil {
    log.Fatal(err)
}

prodChanges := 0
devChanges := 0

for pkgName, change := range changes {
    pkg := change.Package
    
    if pkg.IsDevOnly() || pkg.IsRootDevRequirement() {
        devChanges++
        fmt.Printf("[DEV] %s: %s\n", pkgName, change.Operation.Name)
    } else {
        prodChanges++
        fmt.Printf("[PROD] %s: %s\n", pkgName, change.Operation.Name)
    }
}

fmt.Printf("Production changes: %d, Dev changes: %d\n", prodChanges, devChanges)
```

## Testing

```bash
make test
```

## Related Projects

- [go-deps-diff](https://github.com/yoanm/go-deps-diff) - Generic dependency comparison library that powers this tool

## License

See LICENSE file for details.
