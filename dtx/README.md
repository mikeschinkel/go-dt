# dtx: Incubating experimental Domain Types and Helpers

`dtx` is the **experimental companion module** to [`dt`](https://github.com/mikeschinkel/go-dt). It contains domain types, helpers, and utilities that are being **evaluated for promotion** into the core `dt` module once they prove broadly useful and stable.

`dtx` is **usable in production**, but its APIs are **subject to change**. When an item is promoted into `dt`, an alias or shim will remain in `dtx` for at least **two major versions or two years** (whichever is longer) to provide a smooth migration path.

---

## Stability Policy

* ✅ **Production Use:** Allowed. Pin your dependency version and monitor the changelog.
* 🔁 **Promotion:** When promoted to `dt`, aliases/shims remain for ≥2 major versions or 2 years.
* ⚠️ **Change Policy:** APIs may evolve between versions to improve clarity or correctness.
* ❌ **Removal:** Rare and only for exceptional technical reasons (e.g., fundamental design flaws).

This approach encourages innovation while maintaining long-term stability for production users.

---

## Design Philosophy

The purpose of `dtx` is to incubate useful patterns and abstractions that align with the philosophy of the `dt` package:

* **Type Safety** over loose `interface{}` usage.
* **Zero-Value Validity** wherever practical.
* **ClearPath Code Style** — no early returns; explicit `end:` labels; readable, linear flow.
* **Cross-Platform Semantics** — code that works the same across Unix, macOS, and Windows.
* **Minimal Dependencies** — standard library only, no third-party packages.

---

## Example: `AssertType()`

`AssertType()` is a generic helper that safely performs type assertions and returns structured diagnostic errors:

```go
ro, err := dtx.AssertType[*RawOptions](opts)
if err != nil {
    fmt.Println(err)
}
```

Implementation excerpt:

```go
func AssertType[T any](value any) (t T, err error) {
    var ok bool

    if value == nil {
        err = dt.NewErr(de.ErrFailedTypeAssertion, "source", "nil")
        goto end
    }

    t, ok = value.(T)
    if !ok {
        err = dt.NewErr(de.ErrFailedTypeAssertion, "source", fmt.Sprintf("%T", value))
    }

end:
    if err != nil {
        err = dt.WithErr(err, "target", reflect.TypeOf((*T)(nil)).Elem().String())
    }
    return t, err
}
```

This pattern demonstrates how `dtx` functions complement `dt` by extending its semantics into type- and reflection-related helpers.

---

## Promotion Guidelines

1. **Adoption:** The helper or type is used in multiple packages or by external users.
2. **Stability:** No design changes for at least one minor release.
3. **Portability:** Works consistently across supported platforms.
4. **Zero Value:** Meaningful or explicitly invalid but documented.
5. **Error Semantics:** Uses `dt.NewErr()` and `dt.WithErr()` to preserve error metadata.

Once promoted, the symbol moves into `dt` with a type alias left behind in `dtx` for backwards compatibility.

---

## Versioning

* `dtx` follows [Semantic Versioning](https://semver.org/).
* Minor releases may include breaking changes.
* Major versions reset the two-version compatibility clock for promoted symbols.

---

## Contributing

If you have ideas for new domain types or helper functions, open an issue or pull request in the [`go-dt`](https://github.com/mikeschinkel/go-dt) repository. Please follow the **ClearPath** conventions and include rationale for inclusion or promotion.

---

## License

MIT License — see [LICENSE](../LICENSE) for details.

