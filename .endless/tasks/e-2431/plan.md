## Approach

**Modules.** In dtx, appinfo and dtglob, drop the `replace github.com/mikeschinkel/go-dt => ../` line, then with `GOWORK=off` run `go get github.com/mikeschinkel/go-dt@v0.6.1` and `go mod tidy` so go.mod and go.sum name the published version. In each test and example module, keep only the replace pointing at the module under test (`dtx/test`: `dtx => ..`; `test` and `examples/basic_usage`: `go-dt => ..` / `../..`), drop any other, and tidy with `GOWORK=off`. dtx, appinfo and dtglob each depend only on root, so they have no order among themselves.

**Local development.** Convert the Makefile to a justfile recipe for recipe (help, test, test-unit, test-corpus, test-all, lint, fmt, vet, tidy, build, examples, clean, ci, ensure-valid), keeping `GOEXPERIMENT=jsonv2`. Add `work`, which runs `go work init` over every directory holding a go.mod (removing an existing go.work first). go.work and go.work.sum stay gitignored. Add a README section on local development covering `just work`, why go.work is not committed, and `GOWORK=off`.

**test.yml.** Set `GOWORK: off` at workflow level. Keep the root unit, fuzz-corpus and lint jobs; add a step per nested module that tests it and its test module; add a `go mod tidy -diff` check over every module.

**release.yml.** Replace the free-text `version` input with a `module` choice input (root, dtx, appinfo, dtglob) and a `version` input validated as `vX.Y.Z[-pre]`. Compute `TAG` and `DIR` from them. Before tagging: check `$DIR/go.mod` declares the module path for `$DIR` and contains no replace; run the module's tests (and its test module's) with `GOWORK=off`. Take `PREV_TAG` from `git tag -l "<prefix>v*" --sort=-v:refname` (root: tags without a slash). Build notes from `git log PREV_TAG..HEAD -- $DIR` (root: excluding the nested module directories). Title the release with the module name.

**Versions.** For each of dtx, appinfo and dtglob, compare its exported API at its last tag against HEAD (`go doc -all` or apidiff in a scratch copy). Per Go Stability at v0: a removal or incompatible change is a minor bump, otherwise patch. Give the user each tag, the reason and a short description, and wait for them to run the workflow.

**Gomion and Endless.** Add an open question to gomion E-2434, E-2212 and E-2347 and endless E-2433: whether release tooling should rely on go.work plus no replace in publishable modules, instead of removing and restoring replaces. Cite the trial results recorded in this task's analysis.
