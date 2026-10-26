## Approach

Three recipes in go-dt's justfile, each a thin call to a script under `scripts/`, so the logic can be shellchecked and read on its own.

**`just propose <module>`** (`scripts/propose.sh`). `module` is one of root, dtx, appinfo, dtglob. Find the previous tag the way release.yml does: `git tag -l '<dir>/v*' --sort=-v:refname | head -n 1` for a nested module, `git tag -l 'v*' …` for root. With no previous tag, say so and stop without a verdict. Extract the whole repo at that tag (`git archive <tag>`) to a temp directory, so the old module builds with whatever go.mod and in-repo replace it had at the time, and remove the directory on exit. Run `go run golang.org/x/exp/cmd/apidiff@<pinned>` with `GOWORK=off`: `-w` on the old tree, `-w` on HEAD, then the comparison. Print the changes, then the verdict under 0.X.Y: any incompatible change gives minor (`v0.(Y+1).0`), otherwise patch (`v0.Y.(Z+1)`), with the resulting tag. Separately, compare the go-dt require at the tag with HEAD's. When the X of 0.X.Y differs (v0.5.0 to v0.6.1), print a warning that upgrading this module moves consumers across an incompatible go-dt minor; the verdict is unchanged.

**`just floors [dir]`** (`scripts/floors.sh`). `dir` defaults to `..`. Find every go.mod under it (pruning `.endless/worktrees`, `vendor` and `.git`) that requires github.com/mikeschinkel/go-dt or one of its nested modules, and print module, path required and version. Label each required version:
- **published:** `go list -m -versions <path>` lists it;
- **tagged, not resolvable:** the proxy doesn't list it, but `git ls-remote --tags` on the go-dt remote has the tag (`<dir>/` prefix for nested modules);
- **never tagged:** neither has it.
Exit non-zero when any version is not published.

**`just release <module> <version> [description]`** (`scripts/release.sh`). Fetch origin; refuse when the worktree is dirty or local HEAD differs from `origin/main`, naming which. Print the tag the workflow will make and the commit it will tag, run `gh workflow run release.yml -f module=… -f version=… [-f description=…]`, and print the run URL. The workflow keeps its own checks (tag exists, module path, replace, tests).

**Docs.** Add a "Releasing" part to README's Local Development section covering the three recipes; mention `just propose` and `just floors` in AGENTS.md.

**Verify suite.** `.endless/tasks/e-2441/verify.sh` sourcing the harness: shellcheck the scripts; run `just propose` on each nested module against its real last tag and check the verdict and tag lines; run `just floors` on a scratch tree with go.mod files requiring a published version, a tag that does not exist (v0.7.0) and check the labels and exit code; check `just release` refuses on a dirty tree and on HEAD ≠ origin/main without calling gh (a stub `gh` on PATH that records any call).
