## Approach

**Go 1.27.** In every go-dt module (root, `dtx`, `appinfo`, `dtglob`, their `test/` modules, `test/`, `examples/basic_usage`), run `go mod edit -go=1.27.0` and then `GOWORK=off go mod tidy`. In `.github/workflows/test.yml` and `release.yml`, set `GO_VERSION: '1.27'` by hand. E-2471 later replaces both files with gomion's generated ones.

**justfile.** Remove the `export GOEXPERIMENT` line; go-dt doesn't use encoding/json/v2. Pin the linter at `golangci-lint@v2.15.0`, because v2.6.2 cannot read Go 1.27 export data.

**Lint.** Fix the five v2.15.0 findings without changing behaviour:
- `dtx/is_zero.go`, `dtx/nilable.go` and `dtx/test/nilable_test.go`: `reflect.Ptr` becomes `reflect.Pointer`.
- `dtx/ordered_map.go`: `WriteString(fmt.Sprintf(...))` becomes `fmt.Fprintf(...)`.

Then `just lint` reports 0 issues in every module.

**Already at the baseline.** The release notes leave out `.endless` commits, `/bin/` is ignored and untracked, and Endless's per-checkout files are ignored. The suite checks these; nothing changes them.

**Verify suite.** `.endless/tasks/e-2499/verify.toml`. The gotest checks run each test module with `dir`, followed by sh checks for the go directives, tidiness, the workflows, the justfile and lint, all with `GOWORK=off`.

**Release.** After the land, give Mike the four minor tags, each with a release description: `v0.7.0`, `dtx/v0.3.0`, `appinfo/v0.4.0`, `dtglob/v0.3.0`. He runs them with the current Release workflow.
