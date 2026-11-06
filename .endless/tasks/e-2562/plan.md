# Plan: move go-dt's in-package tests into the test modules

Every test in the go-dt repo ends up in a `test/` module, testing the exported
API, per the go-pkgs testing rule (ED-1650). Two modules hold in-package tests
today: the root module (seven files) and `dtx` (four files). appinfo and
dtglob have none.

## Root module

Move each of `dir_path_test.go`, `entry_path_test.go`,
`path_segments_ext_test.go`, `segments_test.go`, `tilde_dir_path_test.go`,
`tilde_test.go` and `url_segments_test.go` to `test/` under its current name,
with `git mv` so the review sees renames. In each file: `package test`, import
`github.com/mikeschinkel/go-dt`, qualify every dt identifier with `dt.`. Six
files are already `package dt_test`, so only the package line and import
change. `tilde_test.go` is `package dt` but uses only exported identifiers
(`ToTilde`, `TrimTilde`, `EnsureUserHomeDir`, the `TildeOption` constants), so
it needs no wrapper. The existing `test/dirpath_test.go` stays as it is beside
the moved `dir_path_test.go`; no test function names collide.

## dtx

Rename the unexported `parsePathSegmentForOS(goos, s string)` in
`dtx/path_segments.go` to the exported `ParseOSPathSegmentWithOS(s, goos
string)`: the string first, the GOOS last, with a doc comment saying it is
`ParseOSPathSegment` for a given GOOS. `ParseOSPathSegment` calls it with
`runtime.GOOS`. This follows Mike's `With<Thing>` idiom: keep the base name,
append `WithOS`, pass the extra input last.

Then move `dir_path_scanner_test.go`, `path_segments_test.go`, `slices_test.go`
and `safe_map_test.go` to `dtx/test/` under their current names, as
`package test` with `dtx.`-qualified calls. `path_segments_test.go` calls
`dtx.ParseOSPathSegmentWithOS(tt.input, tt.goos)` with the arguments in the new
order. The other three use only exported identifiers. No test function names
collide with `dtx/test`'s existing files.

## Coverage

Remove `-coverprofile=... -covermode=atomic` from the root `go test` line in
the `justfile` (`test-unit`) and in `.github/workflows/test.yml` ("Run unit
tests"). Nothing consumes the profile and Mike does not use coverage yet; the
justfile gains the flags back when he does. The root `go test -v -race ./...`
steps stay so the root and dtx packages still compile under test.

## Verification

`GOWORK=off`: `go mod tidy -diff`, `go build ./...`, `go vet ./...` and
`golangci-lint run` in every module; `go test -race ./...` in root, `test/`,
`dtx`, `dtx/test`; `git ls-files '*_test.go'` shows no file outside a `test/`
directory; `go doc dtx ParseOSPathSegmentWithOS` shows the new signature; the
total number of `Test` functions across the repo is unchanged.
