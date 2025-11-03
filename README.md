# Domain Types for Go (`dt`)

## Purpose

`dt` provides a **stable, dependency-free foundation** of common domain types for Go—types that are “good enough” for everyone to use, even if not perfect for every use-case.

Software ecosystems thrive when developers can build on shared assumptions instead of constantly reinventing them. `dt` aims to be that shared foundation: a **stake in the ground** for how common types like file paths, identifiers, and URLs that can be used in your packages with the knowledge that other packages can have access to those same types.

We’re not trying to design the ideal type for every situation. We’re trying to make it **easy to agree on something usable** so that code written by different teams and libraries can interoperate seamlessly.

When you write Go code using the standard library, you don’t have to ask whether `string` or `os.FileInfo` will be available to the next developer; they simply are. The goal of `dt` is to bring that same confidence and low-friction usability to domain types that the standard library never standardized.

--- 

## Status

This is **pre-alpha** and in development thus **subject to change**, although I am trying to bring to v1.0 as soon as I feel confident its architecture will not need to change. As of Novemeber 2025 I am actively working on it and using it in current projects.

If you find value in this project and want to use it, please start a discuss to let me know. If you discuver any issues with it, please open an issue or submit a pull request.

---

## Rationale

Many Go developers recognize that custom domain types can improve correctness and readability. Yet few actually use them, because in today’s ecosystem doing so requires too much effort.

The barrier is friction. Even simple types like Filename or `DirPath` often demand re-implementing helper methods, sacrificing interoperability with third-party libraries, and working around the standard library’s limited type flexibility.

`dt` envisions a future where that friction is greatly reduced.

---

## Why It Exists

> _"When two (2) different packages from two (2) different authors both use `dt.Filepath` or `dt.Identifier`, they can exchange values directly. No glue code, no type conversions, and no risk of mismatched assumptions."_<br>— The `dt` team

Today, developers who use custom types end up isolated. A package defining its own `FilePath` or `Identifier` can’t interoperate cleanly with others that do the same. Every ecosystem needs a shared vocabulary. 

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

Go's generics system _(as of Go 1.25)_ simply **cannot express the type relationships needed** to make package functions work ergonomically with domain types. Defaulting to package functions in `dt` would force constant type casting, defeating the impetutus behind `dt`, to provider greater type safety that `string` and other built-in types provide without having to constantly cast values from one derived type to another type dervied from the same base type:

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

For a detailed technical explanation of this decision, see **[ADR 2025-02-11: Methods Over Package Functions](adrs/2025-02-11-methods-over-package-functions.md)**.

---

## Companion Packages

| Package                    | Purpose                                                                                                                                                                                                                    |
| -------------------------- |----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| [`dt/de`](../de)           | **Domain Errors** — shared sentinel errors and helpers used throughout the ecosystem.                                                                                                                                      |
| [`dt/appinfo`](../appinfo) | **Application Information** — common interfaces for describing an app’s name, version, config paths, and metadata.                                                                                                         |
| [`dt/dtx`](../dtx)         | **Experimental Incubator** — new types and helpers under evaluation for future inclusion in `dt`. Usable in production with strong guarantees of continued existence, but subject to change when there is no other option. |

---

## Governance & Community

`dt` is **very** open to collaboration; we are actively seeking it. Our intention is to recruit enough active contributors that governance can eventually move to a dedicated GitHub organization. The aim is a community-led defacto-standard that remains practical, stable, and inclusive.

If you share this vision — whether as a library author, contributor, or just a developer who’s tired of having to use a `string` type instead of a bespoke domain type becausthe friction is just too great — start or join a discussion and/or submit a pull requestto help drive what `dt` can become.

---

## License

MIT License — see [LICENSE](LICENSE) for details.
