go-dt's root module still carries seven `*_test.go` files beside the code:
`dir_path_test.go`, `entry_path_test.go`, `path_segments_ext_test.go`,
`segments_test.go`, `tilde_dir_path_test.go`, `tilde_test.go` and
`url_segments_test.go`. They predate the go-pkgs testing rule (ED-1650, in the
gomion project) that a package's tests live in its `test/` module against the
exported API, and that an unexported function needing a test gets an exported
wrapper in a build-tagged file rather than an in-package test.

E-2549 followed the rule for its own new test (`test/stdio_test.go`) and left
these alone as a separate, larger change. Noticed on 2026-11-03 while writing
that task's verify suite; Mike asked for this task on 2026-11-05.

Six of the seven already use the external package `dt_test`, so they only
reach the exported API and can move with an import change. `tilde_test.go` is
`package dt`, but every identifier it uses (`ToTilde`, `TrimTilde`,
`EnsureUserHomeDir`, the `TildeOption` constants) is exported, so it needs no
wrapper either.

The root `justfile` and `.github/workflows/test.yml` collect coverage with
`go test -coverprofile ./...` in the root module; once no test lives there,
that coverage comes from `test/` instead (`go test -coverpkg` from the test
module, or dropping the root profile), which the task has to settle so CI
keeps reporting coverage.
