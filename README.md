# Domain Types for Go (`dt`)

## Purpose

`dt` provides a **stable, dependency-free foundation** of common domain types for Go—types that are “good enough” for everyone to use, even if not perfect for every use-case.

Software ecosystems thrive when developers can build on shared assumptions instead of constantly reinventing them. `dt` aims to be that shared foundation: a **stake in the ground** for how common types like file paths, identifiers, and URLs that can be used in your packages with the knowledge that other packages can have access to those same types.

We’re not trying to design the ideal type for every situation. We’re trying to make it **easy to agree on something usable** so that code written by different teams and libraries can interoperate seamlessly.

When you write Go code using the standard library, you don’t have to ask whether `string` or `os.FileInfo` will be available to the next developer; they simply are. The goal of `dt` is to bring that same confidence and low-friction usability to domain types that the standard library never standardized.

--- 

## Status

This is **pre-alpha** and in development, thus **subject to change**—though I'm working toward v1.0 with confidence that the architecture is approaching stability. As of December 2025, I'm actively developing and using it in current projects.

If you find value in this project and want to use it, please open a discussion to let me know. If you discover any issues, please open an issue or submit a pull request.

---

## Stability & Compatibility Policy

dt uses explicit stability levels for exported symbols:

- **Stable**: no breaking changes to symbol identity or documented behavior
- **Provisional**: likely safe for production, but may still change as naming/shape matures
- **Experimental**: may change or disappear without notice
- **Deprecated**: scheduled for removal with a published removal floor
- **Obsolete**: retained for compatibility but should not be used
- **Internal**: exported for technical reasons but not part of the public API

Alignment note: OpenTelemetry provides useful prior art for maturity ladders and deprecation floors; this repo borrows vocabulary and minimum floors, but retains symbol-level enforcement goals that OTel does not target for Go API identifiers.

## Deprecation Policy

For Deprecated symbols:

- Minimum removal floor: **2 minor releases or 6 months, whichever is later**
- Required metadata: `RemoveAfter` and `Replacement`
- We meet OTel's minimum and may exceed it; the symbol-level contract is authoritative

## How to Read Contract Blocks

Contract blocks are machine-readable doc comments attached to exported identifiers.

```go
// Filepath represents a filesystem path including a filename.
//
// Contract:
// - Stability: Stable
// - Since: v0.6.0
// - Note: Methods are stable and will not change without a major version.
type Filepath string
```

Deprecated items also include:

```go
// - Deprecated: v0.8.0
// - RemoveAfter: v1.0.0 or 2026-01-01
// - Replacement: dt.NewFilepath
```

## Tooling / CI (planned)

Planned tooling makes contracts enforceable:

- Extract `Contract:` blocks into a generated JSON index
- Validate that Deprecated items include `RemoveAfter` and `Replacement`
- Use `golang.org/x/exp/apidiff` to detect breaking changes
- Fail CI when Stable contracts break without a major version bump

## Rationale

Many Go developers recognize that custom domain types can improve correctness and readability. Yet few actually use them, because in today’s ecosystem doing so requires too much effort.

The barrier is friction. Even simple types like Filename or `DirPath` often demand re-implementing helper methods, sacrificing interoperability with third-party libraries, and working around the standard library’s limited type flexibility.

`dt` envisions a future where that friction is greatly reduced.

---

## Why `dt` Exists

> _"When different packages from different authors both use `dt.Filepath` or `dt.DirPath`, they can exchange values directly. No glue code, no type conversions, and greatly reduced risk of mismatched assumptions."_<br>— The `dt` team

Today, developers who use custom types end up isolated. A package defining its own `FilePath` or `DirPath` cannot interoperate cleanly with others that do the same. Every ecosystem needs a shared vocabulary. 

For Go, we hope that `dt` can be that vocabulary.


---

## The Problem We’re Solving

> Using domain types in Go shouldn’t feel like swimming upstream.

Because the Go standard library works almost exclusively with built-in types (`string`, `int`, etc.), developers are discouraged from introducing semantic wrappers—even when they make code safer and more expressive. `dt` solves this by doing the hard part once, in one place, and promising to keep it stable.

We want to make it **easy and obvious** to choose `dt.Filename` instead of `string`, and to know that everyone else is doing the same.

---

## Our Goal

A community-driven, well-designed set of types that:

* 🧱 Provide **semantic clarity** for common domains _(files, identifiers, URLs, etc.)._
* ⚙️ Work seamlessly with the Go standard library and third-party packages.
* 🔄 Enable **interoperability** across projects and teams.
* 🧭 Are **stable and intuitive**, with minimal learning curve.
* 🌍 Encourage a **shared ecosystem vocabulary**.

Our _reasonable goal_ is broad adoption: if enough developers and library authors agree to depend on `dt`, interoperability naturally follows.

Our **ultimate goal** — _said with a wink_ — is to see the Go team adopt something like `dt` officially, perhaps as `golang.org/x/dt`. Whether or not that ever happens, the mission remains the same: to make domain types first-class citizens in Go.

---

## The Gofmt Analogy

As [Rob Pike](https://twitter.com/rob_pike) noted:

> _"Gofmt's style is no one's favorite, yet gofmt is everyone's favorite."_

We believe the same principle could _(should?)_ apply to this Domain Types package:

> _"Types in Domain Types are no one's favorite, yet the Types in Domain Types are everyone's favorite."_

Just as `gofmt` succeeds by being **good enough** rather than perfect, `dt` succeeds by establishing a **shared vocabulary** that everyone can build on. It's not about designing the ideal type for every situation—it's about making it easy to agree on something usable, so developers and libraries can interoperate seamlessly without constant friction.

---

## Our Promise

* 🧩 **Simplicity:** Each type does one thing and does it well.
* 🛠️ **Minimal Logic:** No complex behavior—just useful helpers like `ReadFile()` or `EnsureExists()` that feel native to Go.
* 🧭 **Stability:** After `v1.0.0`, no breaking changes without a major bump.
* 🕰️ **Longevity:** Designed to remain relevant for many years without redesign.
* 🧑‍🤝‍🧑 **Interoperability:** Always safe to use in your own packages and libraries.

---

## Design Choice: Methods Over Package Functions

If you're coming from the Go standard library, you might notice that `dt` frequently uses **methods on types** rather than package-level functions:

```go
// Standard library typically uses package functions
var fp string
...
content, err := os.ReadFile(fp)
dir := filepath.Dir(dp)

// dt typically uses methods
var fp dt.Filepath
...
content, err := fp.ReadFile()
dir := fp.Dir()
```

**This is not a stylistic preference.** It's a pragmatic response to a technical limitation.

Go's generics system _(as of Go 1.25)_ simply **cannot express the type relationships needed** to make package functions work ergonomically with domain types. Defaulting to package functions in `dt` would force constant type casting, defeating the impetus behind `dt`, to provide greater type safety that `string` and other built-in types provide without having to constantly cast values from one derived type to another type derived from the same base type:

```go
// A package function approach would require verbose casting
var myPath dt.Filepath
var myDir dt.DirPath
...
myDir = dt.DirPath(filepath.Dir(string(myPath))

// Whereas a method approach is clean and type-safe
var myPath dt.Filepath
var myDir dt.DirPath
...
myDir = myPath.Dir()
```

The method-based approach provides:
- **Type-casting rarely required** – operations are attached to the type for better ergonomics
- **Better discoverability** – IDE autocompletes will shows all available operations
- **Cleaner composition** – Chaining for those who prefer it, e.g.: `dir.Join(file).ReadFile()`
- _And best of all:_ **Type safety** – The Go compiler prevents invalid type compositions

For a detailed technical explanation of this decision, see **[ADR 2025-02-11: Methods Over Package Functions](adrs/adr-2025-11-02-methods-over-package-functions.md)**.

---

## Type Classification

`dt` provides two categories of types:

### Core Types

Domain types that model real-world concepts and are intended to be used as replacements for built-in types like `string`. These represent semantic entities: file paths, identifiers, URLs, versions, and the like.

**Examples:** `Filename`, `Filepath`, `DirPath`, `URL`, `Identifier`, `Version`

### Supplemental Types

Types that provide essential functionality for core types—either as containers or enumerations. These are used alongside core types to enable safe operations and type-safe classification.

**Examples:** `DirEntry` (container), `EntryStatus` (enumeration)

---

## Type Hierarchy

The following diagram shows how `dt` types relate to one another, from broad to narrow:

```
Types
├── Supplemental Types
│   ├── DirEntry (used during directory walking)
│   └── EntryStatus (enumeration: FileEntry, DirEntry, SymlinkEntry, etc.)
└── Core Types (all extend string)
    ├── Identifier (validated identifiers)
    ├── TimeFormat (time layout strings)
    ├── Version (software version strings)
    ├── EntryPath (generic file or directory path)
    │   ├── DirPath (directory path)
    │   │   ├── TildeDirPath (tilde-prefixed directory path)
    │   │   ├── VolumeName (Windows volume identifier)
    │   │   └── PathSegments
    │   │       └── PathSegment (single path component)
    │   └── Filepath (file path with filename)
    │       ├── RelFilepath (relative file path with traversal protection)
    │       │   └── Filename (filename without path)
    │       │       └── FileExt (file extension with leading period)
    │       └── Filename (filename without path)
    │           └── FileExt (file extension with leading period)
    ├── RelPath (generic relative file or directory path)
    │   ├── PathSegments
    │   │   └── PathSegment (single path component)
    │   └── RelFilepath (relative file path with traversal protection)
    │       └── Filename (filename without path)
    │           └── FileExt (file extension with leading period)
    ├── InternetDomain (internet domain names)
    └── URL (syntactically valid URLs)
        ├── URLSegments (URL path components)
        │   └── URLSegment (single URL path component)
        └── Filename (when used in URL contexts)
            └── FileExt
```

---

## Parse Functions

Many types in `dt` have associated `Parse<Type>()` functions:

```go
filename, err := dt.ParseFilename("config.json")
url, err := dt.ParseURL("https://example.com")
version, err := dt.ParseVersion("1.2.3")
```

**Current Status:** These functions are currently lightweight casting functions. Over time, they will evolve to include robust validation similar to `url.Parse()` from the standard library.

**Future Intent:** As validation is implemented progressively, the type hierarchy understanding will evolve to reflect validated constraints. This design ensures `dt` can add validation without breaking existing code.

---

## File Type Semantics

Throughout the `dt` API, the term **"file"** specifically means **regular file**, as determined by `fs.FileMode.IsRegular()`. This aligns with Go standard library conventions.

### What is a "Regular File"?

A regular file is a normal data file—the most common file system entry type. It excludes:
- Directories
- Symbolic links
- Named pipes (FIFOs)
- Unix sockets
- Device files (block or character)
- Other irregular entries

### API Implications

Methods and functions with "File" in their name operate on regular files only:

- `IsFile()` — Returns `true` only for regular files
- `WalkFiles()` — Yields only regular files
- `Filename()` — Available for regular files

To process all entry types including symlinks and special files, use:
- `Walk()` — Yields all entry types
- `Status()` — Classify entry type
- Pattern matching on `EntryStatus` values

### Example: Handling Different Entry Types

```go
func processAllEntries(root dt.DirPath) (err error) {
    for entry := range root.Walk() {
        var status dt.EntryStatus

        status = entry.Status()

        switch status {
        case dt.IsFileEntry:
            // Regular file
            fmt.Println("File:", entry.Filename())

        case dt.IsDirEntry:
            // Directory
            fmt.Println("Dir:", entry.DirPath())

        case dt.IsSymlinkEntry:
            // Symbolic link
            var target dt.EntryPath
            var path dt.EntryPath

            path = entry.EntryPath()
            target, err = path.Readlink()
            if err != nil {
                goto end
            }
            fmt.Println("Symlink:", path, "->", target)

        case dt.IsSocketEntry, dt.IsPipeEntry, dt.IsDeviceEntry:
            // Special file types
            fmt.Println("Special:", entry.EntryPath())

        default:
            // Unknown or error
            fmt.Println("Other:", entry.EntryPath())
        }
    }

end:
    return err
}
```

This distinction ensures type safety and predictable behavior across the API.

---

## Types and APIs

### Filesystem Path Types

Filesystem paths are the core focus of `dt`. They represent locations on the filesystem and provide type-safe operations for reading, writing, and navigation.

#### DirPath

Represents a filesystem directory path (absolute or relative).

**Key Methods:**
- `EnsureExists()` — Create directory and parents if needed; error if path exists as file
- `ReadDir()` — List directory contents
- `Walk()` — Iterate through directory tree with `SkipDir()` support
- `WalkFiles()` — Iterate through regular files only
- `WalkDirs()` — Iterate through directories only
- `Join(...any)` — Join path components
- `Dir()` — Parent directory
- `Base()` — Directory name as `PathSegment`
- `Clean()` — Normalize path
- `Stat()` — Get file info (follows symlinks)
- `Lstat()` — Get file info without following symlinks
- `Exists()` — Check existence
- `DirFS()` — Convert to `fs.FS`

**Comprehensive Example:**
```go
func processDirPath() (err error) {
    var dir dt.DirPath
    var subDir dt.DirPath
    var configFile dt.Filepath

    dir = dt.DirPath("/home/user/projects")
    err = dir.EnsureExists()
    if err != nil {
        goto end
    }

    // Walk all regular files recursively
    for entry := range dir.WalkFiles() {
        fmt.Println("File:", entry.Filename())
    }

    // Create and navigate subdirectory
    subDir = dir.Join("src", "main")
    err = subDir.EnsureExists()
    if err != nil {
        goto end
    }

    // Work with files in subdirectory
    configFile = dt.FilepathJoin(subDir, "config.json")
    err = configFile.WriteFile([]byte("{}"), 0o644)
    if err != nil {
        goto end
    }

end:
    return err
}
```

#### TildeDirPath

Directory path with tilde (`~`) prefix for user home directory expansion.

**Key Methods:**
- `Expand()` — Expand tilde to full path

**Package Function:**
- `ParseTildeDirPath(s string)` — Parse string as tilde directory path

**Example:**
```go
func processHomePath() (dir dt.DirPath, err error) {
    var tildePath dt.TildeDirPath

    tildePath, err = dt.ParseTildeDirPath("~/projects/go-dt")
    if err != nil {
        goto end
    }

    dir, err = tildePath.Expand()
    if err != nil {
        goto end
    }

    fmt.Println("Expanded path:", dir)

end:
    return dir, err
}
```

**Note:** `Expand()` is not strict about tilde path validity. It will expand any path, including non-tilde paths like `"."`, by resolving them relative to the user's home directory. It only returns an error if the underlying `os.UserHomeDir()` call fails.

#### Filepath

Represents a complete file path including filename and extension.

**Key Methods:**
- `ReadFile()` — Read file contents
- `WriteFile(data, mode)` — Write file
- `Create()` — Create file
- `OpenFile(flag, mode)` — Open with flags
- `Dir()` — Parent directory as `DirPath`
- `Base()` — Filename as `Filename`
- `Ext()` — File extension as `FileExt`
- `Stat()`, `Lstat()` — File info
- `Exists()` — Check existence
- `CopyTo(dest, opts)` — Copy file to destination with optional settings
- `CopyToDir(dest, opts)` — Copy file to destination directory
- `Remove()` — Delete file

**Comprehensive Example:**
```go
func processFile() (data []byte, err error) {
    var file dt.Filepath
    var dir dt.DirPath

    file = dt.Filepath("/home/user/config.json")

    // Check and read file
    if ok, statErr := file.Exists(); !ok {
        err = statErr
        goto end
    }

    data, err = file.ReadFile()
    if err != nil {
        goto end
    }

    // Get parent directory
    dir = file.Dir()
    fmt.Println("Config location:", dir)

    // Copy to backup
    backup := dt.FilepathJoin(dir, "config.json.bak")
    err = file.CopyTo(backup, nil)
    if err != nil {
        goto end
    }

end:
    return data, err
}
```

#### RelFilepath

Represents a relative file path with protections against directory traversal attacks. Validates that the path does not attempt to escape the intended directory using `../` sequences.

**Key Methods:**
- `Dir()` — Parent directory
- `Base()` — Filename
- `ValidPath()` — Check path validity
- `Stat(fileSys ...fs.FS)` — Get file info
- `ReadFile(fileSys ...fs.FS)` — Read file contents
- `WriteFile(data, mode)` — Write file
- `Rel(baseDir)` — Get relative path
- `Exists()` — Check existence
- `Status(flags)` — Get entry status

**Note:** RelFilepath currently has fewer methods than Filepath. This is not a fundamental design constraint—additional methods will be added as use cases arise. The limited method set reflects practical usage patterns rather than architectural limitations.

**Example:**
```go
func processUserFile(userPath string) (data []byte, err error) {
    var relPath dt.RelFilepath
    var valid bool

    relPath = dt.RelFilepath(userPath)
    valid = relPath.ValidPath()
    if !valid {
        err = errors.New("invalid relative path")
        goto end
    }

    data, err = relPath.ReadFile()
    if err != nil {
        goto end
    }

end:
    return data, err
}
```

#### RelPath

Generic relative path that can represent either a file or directory. Similar to `EntryPath` but constrained to relative paths.

**Key Methods:**
- `Dir()` — Parent directory
- `Base()` — Base name as `PathSegment`
- `Stat(fileSys ...fs.FS)` — Get file info
- `Lstat()` — Get file info without following symlinks
- `Status(flags)` — Get entry status
- `Readlink()` — Resolve symlink target
- `VolumeName()` — Get volume name
- `Abs()` — Convert to absolute path
- `Join(elems ...any)` — Join path components
- `HasSuffix(suffix)`, `Contains(substr)` — String operations

**Example:**
```go
func processRelativePath(rel dt.RelPath) (err error) {
    var abs dt.RelPath
    var status dt.EntryStatus

    abs, err = rel.Abs()
    if err != nil {
        goto end
    }

    status, err = abs.Status()
    if err != nil {
        goto end
    }

    fmt.Printf("Path: %s, Type: %s\n", abs, status)

end:
    return err
}
```

#### EntryPath

Generic filesystem entry path that can represent either a file or directory. Use when the type is unknown until runtime.

**Key Methods:**
- `Dir()` — Parent directory
- `Base()` — Base name as `PathSegment`
- `Stat(fileSys ...fs.FS)` — Get file info
- `Lstat()` — Get file info without following symlinks
- `Status(flags ...EntryStatusFlags)` — Get entry type classification
- `Readlink()` — Resolve symlink target
- `VolumeName()` — Get volume name (Windows)
- `Abs()` — Get absolute path
- `Join(elems ...any)` — Join path components
- `HasSuffix(suffix)`, `Contains(substr)` — String operations
- `EnsureTrailSep()` — Ensure trailing separator
- `HasDotDotPrefix()` — Check for `..` prefix

**Type Checking Pattern:**
```go
func processEntry(ep dt.EntryPath) (err error) {
    var status dt.EntryStatus

    status, err = ep.Status()
    if err != nil {
        goto end
    }

    switch status {
    case dt.IsFileEntry:
        fp := dt.Filepath(ep)
        // Handle regular file
    case dt.IsDirEntry:
        dp := dt.DirPath(ep)
        // Handle directory
    case dt.IsSymlinkEntry:
        var target dt.EntryPath

        target, err = ep.Readlink()
        if err != nil {
            goto end
        }
        // Handle symlink
    case dt.IsSocketEntry, dt.IsPipeEntry, dt.IsDeviceEntry:
        // Handle or skip special types
    }

end:
    return err
}
```

**Note:** EntryPath does not have `IsFile()` or `IsDir()` methods. Use `Status()` for type classification and cast to the appropriate type.

#### Filename

Filename without any path component.

**Key Methods:**
- `Ext()` — Extension as `FileExt`
- `String()` — Get filename as string

**Example:**
```go
fn := dt.Filename("document.txt")
ext := fn.Ext()        // FileExt(".txt")
name := fn.String()    // "document.txt"
```

#### FileExt

File extension including the leading period.

**Type:** `FileExt` is a bare `string` type with no public methods.

**Usage:**
```go
ext := dt.FileExt(".txt")
name := string(ext)  // Convert to string when needed
```

**Note:** FileExt currently has no methods. Extension validation and manipulation can be done through standard string operations after casting.

### Path Components and Segments

Working with components of filesystem paths.

#### PathSegments

Represents a filesystem path as a string with segment operations.

**Key Methods:**
- `Split()` — Split into `[]PathSegment` using OS separator
- `Segment(index)` — Get segment at index
- `Slice(start, end)` — Get segment slice (supports `end == -1` for "to last")
- `SliceScalar(start, end)` — Get joined substring of segments without intermediate allocations
- `LastIndex(sep)` — Find last occurrence of separator

**Example:**
```go
path := dt.PathSegments("home/user/projects/file.go")
segments := path.Split()                // []PathSegment{"home", "user", ...}
segment := path.Segment(1)              // "user"
slice := path.Slice(0, 3)               // first 3 segments
scalar := path.SliceScalar(0, 3)        // "home/user/projects"
```

#### PathSegment

A single filesystem path component.

**Key Methods:**
- `HasPrefix(prefix)`, `HasSuffix(suffix)` — String prefix/suffix checks
- `TrimPrefix(prefix)`, `TrimSuffix(suffix)` — Remove prefix/suffix
- `Contains(substr)` — Check substring presence

### URL Types

Working with URLs and URL components.

#### URL

Represents a syntactically valid Uniform Resource Locator.

**Key Methods:**
- `Parse()` — Parse into `*url.URL` for detailed access
- `GET(client)` — Perform HTTP GET request
- `HTTPGet(client)` — Perform HTTP GET request (alias)

**Comprehensive Example:**
```go
func fetchData(endpoint string) (body []byte, err error) {
    var apiURL dt.URL
    var resp *http.Response

    apiURL, err = dt.ParseURL(endpoint)
    if err != nil {
        goto end
    }

    resp, err = apiURL.GET(http.DefaultClient)
    if err != nil {
        goto end
    }
    defer resp.Body.Close()

    body, err = io.ReadAll(resp.Body)
    if err != nil {
        goto end
    }

end:
    return body, err
}
```

#### URLSegments

Represents URL path segments (parts separated by `/`).

**Key Methods:**
- `Split()` — Split into `[]URLSegment`
- `Segment(index)` — Get segment at index
- `Slice(start, end)` — Get segment slice
- `SliceScalar(start, end, sep)` — Get joined scalar with custom separator
- `LastIndex(sep)` — Find last occurrence of separator
- `Base()` — Last segment as `URLSegment`

**Example:**
```go
segments := dt.URLSegments("api/v1/users/123")
all := segments.Split()                     // []URLSegment{"api", "v1", "users", "123"}
middle := segments.Slice(1, 3)              // []URLSegment{"v1", "users"}
scalar := segments.SliceScalar(1, 3, "/")   // "v1/users"
last := segments.Base()                     // "123"
```

#### URLSegment

A single URL path component, semantically different from a `PathSegment`.

### Helper Types

#### Identifier

A string type representing validated identifiers suitable for Git references, semantic version components, and similar uses.

**Example:**
```go
ref := dt.Identifier("main")
tag := dt.Identifier("v1.2.3")
```

#### Version

Software version string following semantic versioning conventions.

**Key Methods:**
- `Major()`, `Minor()`, `Patch()` — Extract version components
- `Valid()` — Check if version is valid

#### InternetDomain

Internet domain name (e.g., `example.com`).

**Key Methods:**
- `Valid()` — Validate domain format
- `TLD()` — Extract top-level domain
- `String()` — Get as string

#### TimeFormat

Time layout string for use with `time.Parse()` and `time.Format()`.

#### VolumeName

Mounted volume name, primarily for Windows support (e.g., `C:`).

### Directory Entry and Status

#### DirEntry

Represents a filesystem entry encountered during directory walking.

**Key Methods:**
- `Name()` — Entry name as `PathSegment`
- `Filename()` — Filename as `Filename` (if file)
- `DirPath()` — Full path as `DirPath`
- `Filepath()` — Full path as `Filepath` (if file)
- `Status()` — Entry type as `EntryStatus`
- `IsFile()`, `IsDir()`, `IsSymlink()` — Type checks
- `Info()` — Get underlying `fs.FileInfo`

#### EntryStatus

Enumeration classifying filesystem entries:
- `FileEntry` — Regular file
- `DirEntry` — Directory
- `SymlinkEntry` — Symbolic link
- `UnknownEntry` — Other or unknown type

---

## Generic Segment Operations

These functions work with any string type and separator, providing zero-cost abstractions for delimited string manipulation.

### SplitSegments

```go
func SplitSegments[S ~string](s, sep string) []S
```

Splits a string by separator into typed segments. Pre-counts separators for optimal memory allocation.

**Example:**
```go
func example() {
    segments := dt.SplitSegments[dt.PathSegment]("home/user/projects", "/")
    // segments = []PathSegment{"home", "user", "projects"}
}
```

### IndexSegments

```go
func IndexSegments[S ~string](s, sep string, index int) S
```

Returns the segment at the given index. Returns empty string if index is out of bounds or negative.

**Example:**
```go
func example() {
    segment := dt.IndexSegments[dt.PathSegment]("home/user/projects", "/", 1)
    // segment = "user"
}
```

### SliceSegments

```go
func SliceSegments[S ~string](s, sep string, start, end int) []S
```

Returns a slice of segments from `start` (inclusive) to `end` (exclusive). Supports `end == -1` to mean "to the last segment". Returns empty slice if indices are invalid.

**Example:**
```go
func example() {
    segments := dt.SliceSegments[dt.PathSegment]("home/user/projects/file.go", "/", 1, 3)
    // segments = []PathSegment{"user", "projects"}
}
```

### SliceSegmentsScalar

```go
func SliceSegmentsScalar[S ~string](s, sep string, start, end int) S
```

Like `SliceSegments`, but returns a joined scalar string instead of a slice. **Zero heap allocations** — uses single-pass byte position tracking, ideal for extracting contiguous segment ranges.

**Performance Note:** This function is optimized for memory efficiency with exactly one pass through the input string without creating intermediate data structures.

**Example:**
```go
func example() {
    result := dt.SliceSegmentsScalar[dt.PathSegment]("home/user/projects/file.go", "/", 1, 3)
    // result = "user/projects" (no intermediate allocations)
}
```

### JoinSegments

```go
func JoinSegments[S ~string](ss []S, sep string) S
```

Joins a slice of segments with a separator. Pre-calculates required capacity for minimal allocations.

**Example:**
```go
func example() {
    segments := []dt.PathSegment{"home", "user", "projects"}
    result := dt.JoinSegments(segments, "/")
    // result = "home/user/projects"
}
```

---

## Type-Safe Join Functions

`dt` provides generic join functions for safely combining path and URL components with type preservation. These functions use generics to accept any string-like types.

### Filesystem Path Join Functions

```go
// Generic join functions for any string-like types
func DirPathJoin[T1, T2 ~string](a T1, b T2) DirPath
func DirPathJoin3[T1, T2, T3 ~string](a T1, b T2, c T3) DirPath
func DirPathJoin4[T1, T2, T3, T4 ~string](a T1, b T2, c T3, d T4) DirPath
func DirPathJoin5[T1, T2, T3, T4, T5 ~string](a T1, b T2, c T3, d T4, e T5) DirPath

func FilepathJoin[T1, T2 ~string](a T1, b T2) Filepath
func FilepathJoin3[T1, T2, T3 ~string](a T1, b T2, c T3) Filepath
func FilepathJoin4[T1, T2, T3, T4 ~string](a T1, b T2, c T3, d T4) Filepath
func FilepathJoin5[T1, T2, T3, T4, T5 ~string](a T1, b T2, c T3, d T4, e T5) Filepath

func RelFilepathJoin[T1, T2 ~string](a T1, b T2) RelFilepath
func RelFilepathJoin3[T1, T2, T3 ~string](a T1, b T2, c T3) RelFilepath
func RelFilepathJoin4[T1, T2, T3, T4 ~string](a T1, b T2, c T3, d T4) RelFilepath
func RelFilepathJoin5[T1, T2, T3, T4, T5 ~string](a T1, b T2, c T3, d T4, e T5) RelFilepath

func EntryPathJoin[T1, T2 ~string](a T1, b T2) EntryPath
func EntryPathJoin3[T1, T2, T3 ~string](a T1, b T2, c T3) EntryPath
func EntryPathJoin4[T1, T2, T3, T4 ~string](a T1, b T2, c T3, d T4) EntryPath
func EntryPathJoin5[T1, T2, T3, T4, T5 ~string](a T1, b T2, c T3, d T4, e T5) EntryPath
```

### Path Segment Join Functions

```go
func PathSegmentsJoin[T1, T2 ~string](a T1, b T2) PathSegments
func PathSegmentsJoin3[T1, T2, T3 ~string](a T1, b T2, c T3) PathSegments
```

### URL Join Functions

```go
func URLJoin[T1, T2 ~string](a T1, b T2) URL
func URLJoin3[T1, T2, T3 ~string](a T1, b T2, c T3) URL
func URLJoin4[T1, T2, T3, T4 ~string](a T1, b T2, c T3, d T4) URL
func URLJoin5[T1, T2, T3, T4, T5 ~string](a T1, b T2, c T3, d T4, e T5) URL

func URLSegmentsJoin[T1, T2 ~string](a T1, b T2) URLSegments
func URLSegmentsJoin3[T1, T2, T3 ~string](a T1, b T2, c T3) URLSegments
```

**Example:**
```go
func buildPaths() {
    baseDir := dt.DirPath("/home/user")

    // Join accepts any string-like types
    subDir := dt.DirPathJoin(baseDir, "projects")
    configFile := dt.FilepathJoin(baseDir, "config.json")

    // Fixed-arity versions for multiple components
    logFile := dt.FilepathJoin3(baseDir, "logs", "app.log")
}
```

---

## OS Wrapper Functions

`dt` wraps common `os` package functions with type safety.

### Directory Operations

```go
func MkdirAll(path dt.DirPath, perm os.FileMode) error
func MkdirTemp(dir dt.DirPath, pattern string) (dt.DirPath, error)
func RemoveAll(path dt.DirPath) error
```

### File Operations

```go
func CreateFile(path dt.Filepath) (*os.File, error)
func CreateTemp(dir dt.DirPath, pattern string) (*os.File, error)
func ReadFile(path dt.Filepath) ([]byte, error)
func WriteFile(path dt.Filepath, data []byte, perm os.FileMode) error
```

### Working Directory and Home

```go
func Getwd() (dt.DirPath, error)
func TempDir() dt.DirPath
func UserHomeDir() (dt.DirPath, error)
func UserConfigDir() (dt.DirPath, error)
func UserCacheDir() (dt.DirPath, error)
```

### Filesystem Introspection

```go
func DirFS(path dt.DirPath) fs.FS
```

**Example:**
```go
func setupApplicationDirs() (err error) {
    var homeDir dt.DirPath
    var configDir dt.DirPath

    homeDir, err = dt.UserHomeDir()
    if err != nil {
        goto end
    }

    configDir = dt.DirPathJoin(homeDir, ".config", "myapp")
    err = dt.MkdirAll(configDir, 0o755)
    if err != nil {
        goto end
    }

    tempDir, err := dt.MkdirTemp(dt.TempDir(), "myapp-*")
    if err != nil {
        goto end
    }
    defer dt.RemoveAll(tempDir)

end:
    return err
}
```

---

## Error Handling

`dt` uses the structured error system from the `doterr` package internally. We recommend you use `doterr` as well for consistent error handling throughout your applications.

**Reference:** See [`go-doterr`](https://github.com/mikeschinkel/go-doterr) for complete documentation on structured error metadata and chaining.

---

## Utility Functions

### File Status Checking

```go
func CanWrite(path dt.Filepath) (bool, error)
func Stat(path dt.EntryPath) (os.FileInfo, error)
func StatFile(path dt.Filepath) (os.FileInfo, error)
func StatDir(path dt.DirPath) (os.FileInfo, error)
```

### File Time Operations

```go
func Chtimes(path dt.EntryPath, atime, mtime time.Time) error
func ChangeFileTimes(path dt.Filepath, atime, mtime time.Time) error
func ChangeDirTimes(path dt.DirPath, atime, mtime time.Time) error
```

### Time Formatting Utilities

```go
func ParseTimeDurationEx(duration string) (time.Duration, error)
```

### Logging

```go
func Logger() *slog.Logger
func SetLogger(logger *slog.Logger)
func EnsureLogger()
func LogOnError(err error)
func CloseOrLog(closer io.Closer)
```

**Example:**
```go
func checkFileAccess(filePath dt.Filepath) (err error) {
    var info os.FileInfo
    var canWrite bool

    canWrite, err = dt.CanWrite(filePath)
    if err != nil {
        goto end
    }

    info, err = dt.StatFile(filePath)
    if err != nil {
        goto end
    }

    fmt.Printf("File: %s, Writable: %v, Size: %d\n",
        filePath, canWrite, info.Size())

end:
    return err
}
```

---

## Companion Packages

### dtx (Experimental Extensions)

Experimental types and utilities under evaluation for potential inclusion in the main `dt` package. Safe for production use with strong compatibility guarantees, though subject to evolution when necessary.

**Key Features:**
- `GetWorkingDir()` — OS-aware working directory detection with hint support
- `IsZero()`, `IsNil()`, `IsNilable()`, `IsNilableKind()` — Type introspection helpers
- `Must()` — Panic on error helper for fail-fast patterns
- `Panicf()` — Formatted panic function
- `AssertType()` — Safe type assertion with panic fallback
- `TempTestDir()`, `SetTestEnv()` — Testing and environment helpers
- OS-specific path segment parsers (Windows, Darwin, Linux)
- `EntryStatusError()` — Convert `EntryStatus` to error types

**Package:** [`go-dt/dtx`](dtx)

### dtglob (Glob Patterns)

Pattern-based file operations with glob-style matching and bulk copy operations.

**Key Features:**
- `Glob` — Type-safe glob pattern representation
- `GlobRule` — Single file copy operation specification
- `GlobRules` — Container for multiple rules with batch `CopyTo()` operation

**Package:** `go-dt/dtglob` _(if available in your installation)_

### appinfo (Application Metadata)

Standard interface for describing application metadata across the ecosystem.

**Key Features:**
- `AppInfo` interface — Contract for application metadata including name, version, config paths
- `New(Args)` — Create concrete `AppInfo` implementations
- Test helpers for verifying `AppInfo` implementations

**Package:** `go-dt/appinfo`

---

## Local Development

This repository holds several Go modules: the root `go-dt`, the publishable nested modules `dtx`, `appinfo` and `dtglob`, and the test and example modules beside them. The publishable modules require a **published** go-dt version and carry no `replace` directive, because a consumer never sees a replace — only the require — and a replace would let a wrong require go unnoticed.

To build every module from the code on disk, create a Go workspace:

```sh
just work
```

`just work` writes a `go.work` listing every directory that holds a `go.mod`, so it needs no edit when a module is added. Re-run it after adding one.

`go.work` is not committed (it is in `.gitignore`). CI builds and tests with `GOWORK=off`, against the versions each `go.mod` requires, so a change in a nested module that depends on an unpublished root change fails in CI until that root version is tagged and required. To see what CI sees locally, prefix a command with `GOWORK=off`:

```sh
cd dtx && GOWORK=off go test ./...
```

Each test or example module keeps exactly one `replace`, pointing at the module it tests, so CI tests the commit's code rather than the last published tag. `just tidy` runs `go mod tidy` in every module with `GOWORK=off`, and CI fails when any `go.mod` is not tidy.

Run `just` to list the other recipes (`test`, `test-all`, `lint`, `vet`, `examples`, `ci`, …).

---

## Objections

### Why do we need Domain Types at all?!?

| |                                                                                                                                                                                                                                                                         |
| --- |-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| **Short answer** | Shared types reduce friction and make codebases interoperate.                                                                                                                                                                                                           |
| **Longer answer** | When two packages both use a domain type like `dt.DirPath` they exchange values without constant type casting.<br> Using a 3rd package eliminates _dreaded_ import cycles. <br/>Using `dt` saves time, reduces bugs, and makes APIs easier to compose across teams. |
| **<nobr>Minimal adoption</nobr>** | Use a single type (for example `dt.Filepath`) at package boundaries.                                                                                                                                                                                                    |
| **<nobr>Compatibility angle</nobr>** | Stable shared types are a stronger guarantee than ad-hoc local types that can drift or be renamed.                                                                                                                                                                      |

### Why not define my own types per package?

| |                                                                                                                                                                                                                                                                                                                                                          |
| --- |----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| **Short answer** | You lose the **key benefit** of domains types: _interoperability._                                                                                                                                                                                                                                                                                       |
| **Longer answer** | Local types fracture the ecosystem because each package invents its own names and helpers.<br>Local types require constant type casting across packages. <br>`dt` aims to be the **stable**, no _(other)_ dependency. shared vocabulary. <br>You often **must** create your own 3rd package **to avoid import cycles** anyway, why not use `dt` instead? |
| **<nobr>Minimal adoption</nobr>** | Use `dt` types at package boundaries, keep internals local.                                                                                                                                                                                                                                                                                              |
| **<nobr>Compatibility angle</nobr>** | Shared stable types reduce breakage at integration points.                                                                                                                                                                                                                                                                                               |

### This is overkill for small libraries.

| | |
| --- | --- |
| **Short answer** | Small libraries benefit most from stable compatibility signals. |
| **Longer answer** | Small teams often ship faster and have fewer release resources, so breaking changes are more costly.<br>A small set of stable types reduces downstream churn. |
| **<nobr>Minimal adoption</nobr>** | Use only 1-2 core types and ignore the rest. |
| **<nobr>Compatibility angle</nobr>** | A tiny stable surface is easier to preserve than a large, implicit one. |

### Why not just use plain strings?

| |                                                                                                                                                                                                                                                                                           |
| --- |-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| **Short answer** | Types encode intent and reduce mistakes.                                                                                                                                                                                                                                                  |
| **Longer answer** | `string` does not tell readers or tooling whether the value is a filename, a URL, or an identifier.<br>`dt` types make intent explicit and enable safe helper methods.<br>Validating `Parse<type>()(<type>,error)` funcs for each type make baked-in assumptions explicit and actionable. |
| **<nobr>Minimal adoption</nobr>** | Replace only the most error-prone `string` fields in public APIs.                                                                                                                                                                                                                         |
| **<nobr>Compatibility angle</nobr>** | Stable types keep public contracts clear even as internals change.                                                                                                                                                                                                                        |

### I do not want to teach new developers a new way.

| |                                                                                                                                                              |
| --- |--------------------------------------------------------------------------------------------------------------------------------------------------------------|
| **Short answer** | Use a small set and let [tooling](https://github.com/mikeschinkel/gomion) teach the rest.                                                                    |
| **Longer answer** | A handful of domain types is easy and obvious to learn. <br>Method-based API are discoverable in IDEs.<br>Linters can enforce conventions at the boundaries. |
| **<nobr>Minimal adoption</nobr>** | Use only the types that already map to your domain language.<br>Start with `dt.DirPath` and `dt.Filepath`, for example.                                      |
| **<nobr>Compatibility angle</nobr>** | Clear, stable types reduce downstream friction and support burden.                                                                                           |

## Governance & Community

`dt` is **very** open to collaboration; we are actively seeking it. Our intention is to recruit enough active contributors that governance can eventually move to a dedicated GitHub organization. The aim is a community-led defacto-standard that remains practical, stable, and inclusive.

If you share this vision — whether as a library author, contributor, or just a developer who’s tired of having to use a `string` type instead of a bespoke domain type becausthe friction is just too great — start or join a discussion and/or submit a pull requestto help drive what `dt` can become.

---

## License

MIT License — see [LICENSE](LICENSE) for details.
