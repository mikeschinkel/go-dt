# ADR 2025-11-02: Methods Over Package Functions

## Status

Accepted

## Context

The Go standard library predominantly uses **package-level functions** for operations on built-in types:

```go
// Standard library pattern
content, err := os.ReadFile(filepath)
joined := filepath.Join(dir, file)
cleaned := filepath.Clean(path)
```

When designing `dt`, we faced a fundamental question: should domain types like `Filepath`, `DirPath`, and `Filename` follow this same pattern, or should they expose operations as **methods on the types themselves**?

```go
// Package function approach (stdlib-like)
content, err := dt.ReadFile(filepath)
joined := dt.Join(dir, file)

// Method approach (chosen by dt)
content, err := filepath.ReadFile()
joined := dir.Join(file)
```

At first glance, the package function approach appears more idiomatic—it mirrors familiar Go patterns. However, this surface-level assessment misses a critical technical constraint.

## Problem

Go's generic system, while powerful for many use cases, has **fundamental limitations** when it comes to expressing relationships between domain types and their operations.

### The Casting Problem

Consider a simple scenario: joining a directory path with a filename to produce a filepath.

With package functions, the signature must be:

```go
func Join(dir DirPath, file Filename) Filepath {
    return Filepath(filepath.Join(string(dir), string(file)))
}
```

This works fine for the simple case. But what about joining multiple segments? Or accepting either absolute or relative paths? The combinatorial explosion of type combinations makes package functions impractical:

```go
// Every combination requires a separate function
func JoinDirFile(dir DirPath, file Filename) Filepath
func JoinDirSegments(dir DirPath, segments PathSegments) Filepath
func JoinRelativeFile(rel RelFilepath, segments PathSegments) Filepath
// ... dozens more variations
```

Even with generics, Go cannot express "a function that accepts any path-like type and returns the appropriate result type" without forcing the caller to cast:

```go
// Attempting to use generics
func Join[T PathType](parts ...string) T {
    return T(filepath.Join(parts...)) // ❌ Cannot convert string to T
}

// Caller must cast anyway
fp := dt.Join[Filepath](string(dir), string(file)) // verbose and defeats the purpose
```

### The Interoperability Problem

Package functions also create friction when types need to interoperate:

```go
// With package functions - requires type assertions and casts
var path any = getPathFromSomewhere()
if fp, ok := path.(dt.Filepath); ok {
    content, err := dt.ReadFile(fp) // works, but verbose
}

// With methods - type assertion is sufficient
if fp, ok := path.(dt.Filepath); ok {
    content, err := fp.ReadFile() // the method comes with the type
}
```

### The Discoverability Problem

When using package functions, developers must:
1. Import the `dt` package
2. Remember which function to use for which operation
3. Often reference documentation to find the right function name

With methods:
1. Type `filepath.` and let your IDE show available operations
2. The method is directly attached to the type—no guessing needed

## Decision

**We use methods on domain types rather than package-level functions.**

This decision is **not** based on stylistic preference. It is a pragmatic response to Go's generic system limitations. Methods provide:

1. **Zero-casting ergonomics** – The type carries its operations, eliminating the need for type conversions
2. **Better discoverability** – IDE autocomplete shows available operations immediately
3. **Cleaner composition** – Chaining operations feels natural: `dir.Join(file).ReadFile()`
4. **Simpler API surface** – Fewer function variations needed because methods dispatch on the receiver type

## Constraints

- Go generics (as of Go 1.25) cannot express "return the same type as the receiver but transformed" without casting
- Type switches and assertions are necessary for polymorphic path operations with package functions
- The standard library's approach works for built-in types (`string`, `int`) but breaks down for custom domain types

## Consequences

### Positive

- **Minimal casting** – Users rarely need to cast between types
- **Intuitive API** – `filepath.ReadFile()` reads naturally and matches developer expectations
- **Extensibility** – New methods can be added without polluting the package namespace
- **Type safety** – The compiler prevents nonsensical operations (e.g., `filename.MkdirAll()` doesn't exist)

### Negative

- **Feels different from stdlib** – Experienced Go developers may initially find the method-heavy API unfamiliar
- **Not "idiomatic"** – Deviates from the Go proverb "accept interfaces, return structs" (though we do both)
- **Potential confusion** – Developers may wonder why we didn't follow the standard library pattern

### Mitigation

We address the "feels different" concern through:

1. **Clear documentation** – This ADR and README sections explain the rationale
2. **Consistent naming** – Methods mirror stdlib function names where applicable (`ReadFile`, `MkdirAll`, etc.)
3. **Interoperability** – Types implement standard interfaces (`fmt.Stringer`, `fs.FS` compatible types, etc.)

## Alternatives Considered

### 1. Package Functions with Extensive Overloading

Create separate functions for every type combination:

```go
func JoinDirFile(dir DirPath, file Filename) Filepath
func JoinDirSegments(dir DirPath, seg PathSegments) DirPath
func JoinFilepathSegments(fp Filepath, seg PathSegments) Filepath
```

**Rejected because:**
- Unmaintainable API surface (50+ functions)
- Poor discoverability
- Still requires casting in generic code

### 2. Generic Package Functions with Manual Casting

Use generics but require callers to cast:

```go
func ReadFile[T ~string](path T) ([]byte, error)

// Usage
content, err := dt.ReadFile(dt.Filepath(myPath)) // verbose
```

**Rejected because:**
- Defeats the purpose of type safety
- More verbose than methods
- No better than using `os.ReadFile(string(myPath))`

### 3. Hybrid Approach (Package Functions + Methods)

Provide both package functions and methods:

```go
// Both available
content1, _ := dt.ReadFile(filepath)
content2, _ := filepath.ReadFile()
```

**Rejected because:**
- Confusing API with two ways to do everything
- Maintenance burden
- Doesn't solve the casting problem for package functions

## References

- [Go Generics Limitations](https://go.dev/blog/type-parameters) – Official Go blog on generic constraints
- [Golang spec on method sets](https://go.dev/ref/spec#Method_sets)
- [Discussion: Methods vs Functions in Go](https://github.com/golang/go/wiki/CodeReviewComments#receiver-type)

## Conclusion

The decision to use methods instead of package functions is **technically motivated, not stylistic**. Given Go's current generic system, methods provide the best developer experience with minimal casting and maximum type safety.

If Go's generics evolve to support the type relationships we need (e.g., associated types, higher-kinded types), we may revisit this decision in a future major version. Until then, methods are the pragmatic choice.

---

**TL;DR:** We use methods because Go generics can't eliminate casting for package functions operating on domain types. It's about ergonomics, not preference.
