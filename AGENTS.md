# Repository Guidelines

## Project Structure & Module Organization
- Root module `github.com/mikeschinkel/go-dt` holds core domain types (paths, URLs, identifiers) alongside helpers; see `*.go` in the repo root.
- `test/` exercises end-to-end behavior and houses fuzz corpus tests; each submodule with a `/test` folder mirrors this pattern.
- Optional submodules: `appinfo/`, `dtglob/`, and `dtx/` contain auxiliary packages; each may have its own `go.mod` and tests.
- `examples/` contains small binaries; `bin/` (created by builds) is disposable.
- `adrs/` documents architectural decisions; consult when proposing API changes.

## Build, Test, and Development Commands
- `GOEXPERIMENT=jsonv2 make test` — run race-enabled unit tests for root and all submodules; writes `coverage.txt`.
- `GOEXPERIMENT=jsonv2 make test-corpus` — execute fuzz corpus regressions via `TestFuzzCorpus` in each module.
- `GOEXPERIMENT=jsonv2 make lint` — run golangci-lint (v2.6.x) across modules.
- `make fmt` / `make vet` — gofmt all files; go vet with the jsonv2 experiment enabled.
- `GOEXPERIMENT=jsonv2 make build` — compile packages; `make examples` builds each example into `./bin/`.
- `make ci` — convenience target that runs fmt, vet, lint, and full test suite locally.

## Coding Style & Naming Conventions
- Go 1.25.x with `GOEXPERIMENT=jsonv2` expected; keep modules tidy with `make tidy`.
- Always `gofmt`; prefer small, composable methods on domain types (DirPath, Filepath, URL) over package functions.
- Exported types and methods use PascalCase; keep package-level names minimal and consistent with existing DT vocabulary.
- Use `context.Context` as the first arg when applicable; avoid stutter (`dt.DirPath`, not `dt.PathDirPath`).

## Testing Guidelines
- Add `*_test.go` next to the code under test; integration-style tests belong in `test/` (or `<module>/test/`).
- Include race detection and coverage locally (`make test`); keep fuzz corpora current when adding fuzzers.
- Prefer table-driven tests; use temporary dirs via helpers in `dtx` when available.

## Commit & Pull Request Guidelines
- Follow the existing Conventional Commit style seen in history (`feat:`, `fix:`, `chore:`); scope names optional but encouraged for modules.
- PRs should explain motivation, summarize behavioral changes, and list validation (`make test`, `make lint`, etc.).
- Link related issues/ADRs and include repro steps or logs when fixing bugs; screenshots only when UI output is relevant (rare here).
- Keep changes focused; cross-module edits should note downstream impact on domain types and semantics.
