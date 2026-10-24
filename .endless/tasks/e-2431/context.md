Found 2026-10-24: root v0.6.0 lacks PathSegments.IsRoot(); dtx/v0.2.1 (2025-11-18) does not compile against go-dt v0.6.0 (undefined: dt.ErrFileDoesNotExist), with 25 files changed since; appinfo/v0.2.1 and dtglob/v0.2.1 have 7 and 9 files changed since. All three submodules require stale go-dt versions (v0.5.0, v0.3.2) and build only through their in-repo replace. Tagging, for the short term, is done by the user: this task readies each module (requires on published tags, no cross-project replaces, builds and tests pass), then tells the user which tag to make and why.

What go-cliutil's release needs (from its release task): go-cliutil requires go-dt root, dtx and appinfo, not dtglob. It builds against go-dt v0.6.0 and needs nothing newer from the root (it does not use IsRoot). It needs published dtx and appinfo tags whose go.mod require go-dt >= v0.6.0 with no replaces, and whose code compiles against v0.6.0.

go-dt's .github/workflows/release.yml (workflow_dispatch, version input) validates the version against ^v[0-9]+\.[0-9]+\.[0-9]+(-...)?$, so it cannot make directory-prefixed tags such as dtx/v0.3.0. It takes its previous tag from `git tag -l --sort=-v:refname | head -n 1` over all tags, which would mix prefixed and root tags. go-cliutil's release.yml has the same prefix restriction.



Decided by Mike on 2026-10-24 (in E-2374's session): extend go-dt's release.yml so it can
make directory-prefixed tags (dtx/vX.Y.Z, appinfo/vX.Y.Z), rather than tagging dtx and
appinfo by hand. This includes taking the previous tag from the same module's prefix only.
Mike still runs the workflow himself. go-dt v0.6.1 is published (8271ac5), so the dtx and
appinfo tags will require go-dt v0.6.1.
