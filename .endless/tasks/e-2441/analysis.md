## Decisions (Mike, 2026-10-26, in E-2431's session)

- **floors scans a directory, default `..`.** Inside go-dt, CI already builds with GOWORK=off, so an unpublished require there already fails; the value is in sibling repos (go-tealeaves' unpublished v0.7.0, E-2440). The default fits Mike's layout, where go-dt sits beside the other go-pkgs repos; anyone else passes a directory.
- **The Go module proxy decides "published".** It is what a consumer's build resolves. `git ls-remote` is used only to label a miss: "tagged, not resolvable" (proxy lag after tagging, or a broken go.mod at the tag) versus "never tagged". Counting a GitHub tag as published would pass a tag no consumer can resolve.
- **A raised go-dt floor is a warning, not a verdict change.** The verdict follows the module's own API under 0.X.Y, as applied to dtx/v0.2.2 in E-2431. Whether a floor raised across an incompatible go-dt minor should itself force a minor is open in EQ-126 to EQ-128 on the Gomion and Endless brainstorms; this task does not settle it.
- **release refuses a dirty tree or HEAD ≠ origin/main.** The workflow tags origin's main, so this makes the tagged commit the one tested locally.

## Decided without asking (Mike could object; none did)

- apidiff runs via `go run …@<pinned>`, like golangci-lint in the same justfile, so root's go.mod (dependency-free) gains nothing.
- The old tag is built from the whole repo at that tag, extracted to a temp dir.
- Verdict: incompatible → minor, otherwise patch, under 0.X.Y.
- `propose root` compares against the previous unprefixed tag.
- Script logic lives in `scripts/*.sh`.
- No release-description writing; agents do not run `just release`.

## Checked against Gomion (survey 2026-10-26)

No conflicts. Gomion's `goutils.SuggestNextVersions` uses the same v0 rule. It builds the old version from a local clone checked out at the tag, loaded with GOWORK=off, which is the same idea as extracting the tag. It imports apidiff as a library rather than `go run`, but only from a release engine that is commented out. Its published check uses git tags only; floors asks a different question (can a consumer resolve it), which only the proxy answers. Gomion's remove-and-restore-replaces flow is not used by anything here. Findings filed or noted: gomion E-2442 (prefixed tags), notes on E-2212 and E-2434.

## Considered

An Endless Go command that the justfile calls (Mike's suggestion) would allow library use of apidiff and x/mod. Not needed for these three recipes; worth revisiting if E-2433 decides Endless itself owns release planning, at which point these recipes become thin callers of it.
