# Plan: Update go-dt README.md

## Overview
Update the README.md to fully document the current state of the go-dt package. The documentation should be comprehensive but not overly detailed - focusing on types, methods, and usage patterns rather than per-function reference material.

## File to Modify
- `/Users/mikeschinkel/Projects/go-pkgs/go-dt/README.md`

## Key Philosophy to Convey
- Modeling the real world is hard; abstractions are leaky. But `dt` is an attempt to model better than ad-hoc per-package approaches.
- Goal: Like gofmt standardizes style, dt aims to standardize domain types
  - _"Gofmt's style is no one's favorite, yet gofmt is everyone's favorite."_
  - _"Types in Domain Types are no one's favorite, yet the Types in Domain Types are everyone's favorite."_

## Type Classification

### Core Types
Domain types to be used in place of scalars like `string`. These model real-world concepts.

### Supplemental Types
Types needed to provide useful functionality for core types (e.g., DirEntry, EntryStatus).

## Type Hierarchy (broad to narrow)
```
Types
├── Supplemental Types
│        ├── DirEntry
│        └── Enum Types
│            └── EntryStatus
└── Core Types
    └── string
            ├── Identifier
            │        ├── PathSegment
            │        └── URLSegment
            ├── TimeFormat
            ├── Version
            ├── URL
            │   ├── InternetDomain
            │   ├── URLSegment
            │   │   └── Filename
            │   │       └── FileExt
            │   ├── URLSegments
            │   │   └── URLSegment
            │   │       └── Filename
            │   │           └── FileExt
            │   └── Filename
            │       └── FileExt
            └── EntryPath
                ├── VolumeName
                ├── DirPath
                │   ├── VolumeName
                │   └── PathSegments
                │       └── PathSegment
                └── Filepath
                    ├── RelFilepath
                    │   └── Filename
                    │       └── FileExt
                    └── Filename
                        └── FileExt
```

## Parse Functions
- `Parse<Type>(baseType) (<Type>, error)` functions exist for many types
- Currently little more than casting functions (placeholders)
- Intent: Future robust validation like `url.Parse(string)` for URLs
- As validation is implemented, type hierarchy understanding will evolve

## Code Style for Examples
Use "clear path" pattern with `goto end`, NOT inline error handling:

```go
// CORRECT - clear path style
func foo() (err error) {
    var dir dt.DirPath

    dir = dt.DirPath("/home/user/projects")
    err = dir.EnsureExists()
    if err != nil {
        goto end
    }
    // ... more operations
end:
    return err
}

// AVOID - inline error handling
if err := dir.EnsureExists(); err != nil {
    // handle error
}
```

## Current State Analysis
The existing README covers:
- Purpose, rationale, and philosophy (keep these)
- Methods over package functions design choice (keep)
- Basic type documentation (outdated/incomplete)
- Companion packages section (needs updating)

## Documentation to Add/Update

### Main Package Types
1. **DirPath** - Directory path with walking, existence checks, file system operations
2. **Filepath** - File path with read/write, copy, status operations
3. **RelFilepath** - Relative file path with validation
4. **EntryPath** - Generic entry (file or dir) path
5. **Filename** - Filename with extension, no path
6. **FileExt** - File extension with leading period
7. **PathSegment** - Single filesystem path component
8. **PathSegments** - Multiple path segments with Split/Slice operations
9. **URL** - Syntactically valid URL with HTTP methods
10. **URLSegment** - Single URL path component
11. **URLSegments** - Multiple URL segments with Split/Slice operations
12. **Identifier** - Validated identifier string
13. **Version** - Software version string
14. **VolumeName** - Mounted volume name (Windows)
15. **InternetDomain** - Internet domain
16. **TimeFormat** - Time layout string
17. **DirEntry** - Directory entry during walking
18. **EntryStatus** - Filesystem entry classification

### Generic Functions
- **SplitSegments** - Split string by separator into typed segments
- **IndexSegments** - Get segment at index
- **SliceSegments** - Get slice of segments
- **SliceSegmentsScalar** - Get joined scalar from segment range (zero allocation)
- **JoinSegments** - Join segments with separator

### Join Functions (type-safe generic)
- DirPathJoin, DirPathJoin3-5
- FilepathJoin, FilepathJoin3-5
- RelFilepathJoin, RelFilepathJoin3-5
- EntryPathJoin, EntryPathJoin3-5
- PathSegmentsJoin, PathSegmentsJoin3-5
- URLJoin, URLJoin3-5
- URLSegmentsJoin, URLSegmentsJoin3-5

### OS Wrapper Functions
- DirFS, MkdirAll, MkdirTemp, RemoveAll
- UserHomeDir, UserConfigDir, UserCacheDir
- CreateFile, CreateTemp, ReadFile, WriteFile
- Getwd, TempDir

### Error Handling (Internal - minimal mention)
- Single sentence: "dt uses doterr internally; we recommend you do so as well"
- Link to go-doterr README
- Do NOT document doterr functions as part of dt's API

### Utility Functions
- CanWrite, Stat, StatFile, StatDir
- Chtimes, ChangeFileTimes, ChangeDirTimes
- ParseTimeDurationEx
- Logger, SetLogger, EnsureLogger
- LogOnError, CloseOrLog

### Companion Packages (Remove dt/de reference)

#### dtx (Experimental)
- GetWorkingDir - OS-aware working directory with hints
- IsZero, IsNil, IsNilable, IsNilableKind
- Must - Panic on error helper
- Panicf - Formatted panic
- AssertType - Safe type assertion
- TempTestDir, SetTestEnv - Testing helpers
- OS-specific path segment parsers (Windows, Darwin, Linux)
- EntryStatusError - Convert status to error

#### dtglob (Glob patterns)
- Glob - Glob pattern type
- GlobRule - Single copy operation
- GlobRules - Container for rules with CopyTo

#### appinfo (Application metadata)
- AppInfo interface - Application metadata contract
- New(Args) - Create AppInfo implementation
- Test helper for testing appinfo

## Documentation Structure

### Sections to Update/Add

1. **Status** - Keep but verify current
2. **Purpose/Rationale/Why** - Keep as-is, enhance with gofmt-style analogy
3. **Design Choice: Methods Over Functions** - Keep as-is
4. **Type Classification** - NEW: Core vs Supplemental types explanation
5. **Type Hierarchy** - NEW: Include the broad-to-narrow hierarchy tree
6. **Parse Functions** - NEW: Explain as validation placeholders
7. **Types and APIs** - Major rewrite, organized by hierarchy
   - File System Path Types (DirPath, Filepath, RelFilepath, EntryPath)
   - Path Components (Filename, FileExt, PathSegment, PathSegments)
   - URL Types (URL, URLSegment, URLSegments)
   - Helper Types (Identifier, Version, etc.)
   - Entry Status and DirEntry
8. **Generic Segment Operations** - Add/expand
9. **Type-Safe Join Functions** - New section
10. **OS Wrapper Functions** - New section
11. **Error Handling** - Minimal doterr mention with link
12. **Utility Functions** - New section
13. **Companion Packages** - Update with dtx, dtglob, appinfo details
14. **Examples** - Comprehensive per-type examples
15. **Governance & Community** - Keep
16. **License** - Keep

## Example Structure (per major type)
Each major type should have a comprehensive example showing:
- Type creation/parsing
- Key methods in realistic context
- Common patterns
- Use clear path style with `goto end`

Example format:
```go
// DirPath Example - using clear path style
func processDirPath() (err error) {
    var dir dt.DirPath
    var subDir dt.DirPath
    var configFile dt.Filepath

    dir = dt.DirPath("/home/user/projects")
    err = dir.EnsureExists()
    if err != nil {
        goto end
    }

    // Walk files in directory
    for entry, walkErr := range dir.WalkFiles() {
        if walkErr != nil {
            continue
        }
        fmt.Println(entry.Filename())
    }

    // Join paths type-safely
    subDir = dir.Join("src", "main")
    configFile = dt.FilepathJoin(dir, "config.json")

    fmt.Println(subDir, configFile)

end:
    return err
}
```

## Implementation Steps

1. Read current README structure
2. Preserve Purpose, Rationale sections - enhance with gofmt-style analogy
3. Keep Design Choice: Methods Over Functions section
4. ADD new "Type Classification" section (Core vs Supplemental)
5. ADD new "Type Hierarchy" section with the tree diagram
6. ADD new "Parse Functions" section explaining validation placeholders
7. Rewrite Types and APIs section organized by hierarchy
8. Add comprehensive examples using clear path style (`goto end`)
9. Add Generic Segment Operations section
10. Add Type-Safe Join Functions section
11. Add OS Wrapper Functions section
12. Add minimal doterr mention (single sentence + link to go-doterr README)
13. Add Utility Functions section
14. Update Companion Packages (dtx, dtglob, appinfo) - REMOVE dt/de reference
15. Keep ADR reference (adrs/adr-2025-11-02-methods-over-package-functions.md)
16. Verify all method/function names match actual code
17. Update Status if needed
