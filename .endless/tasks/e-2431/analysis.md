## Why go.work instead of replace in publishable modules

A replace in a published module is ignored by consumers, so its require is the only thing a consumer sees, while locally the replace makes that require irrelevant. Nothing forces the number to be correct. Today's dtx shows it: built against its stated require (go-dt v0.5.0) it fails with `undefined: dt.ErrFileNotExist` and `undefined: dt.Logf`. go.work gives the same local behaviour (all modules built from disk) without touching go.mod, and since it is not committed, CI builds against the requires.

A replace is still correct in a module that is never published and exists to test the code in the same commit: without one, CI would test dtx/test against the published dtx. Keeping only the replace to the tested module means CI tests the commit's dtx, built against the go-dt version dtx itself requires.

## Trial (scratch copy of main, Go 1.26.2)

- All replaces removed, go.work over all nine modules: every module vets and tests. dtglob/test failed only because its own replace had been removed too and it requires dtglob at the zero pseudo-version; with its replace kept, it passes.
- Root and dtx changed together (dtx using a new root function): builds locally under go.work; with GOWORK=off and dtx requiring v0.6.1 it fails `undefined: dt.TrialHello`, so CI catches unpublished dependencies.
- GOWORK=off with requires bumped to v0.6.1 and go mod tidy: dtx, appinfo, dtglob, test and examples vet clean. Without the tidy, go.sum lacks go-dt's checksums (never recorded, because of the replace).
- dtx left out of go.work: go refuses loudly ("directory prefix . does not contain modules listed in go.work").
- dtx/test keeping `replace dtx => ..` alongside go.work: no conflict; GOWORK=off sees the commit's dtx.
- An outside consumer requiring dtx/v0.2.1 and go-dt v0.6.1 fails `undefined: dt.ErrFileDoesNotExist` (renamed in v0.6.0): what blocks go-cliutil.

## Rejected

- Committing go.work with GOWORK=off in every workflow: a step that forgets it silently builds against disk code.
- A committed go.local-dev.work copied by a recipe: must be edited when modules are added; generating from the go.mod files needs no upkeep. Endless E-1192 (go.work generation tool) can replace the recipe later.
- Free-text tag input: a typo fails only after starting the run; a module choice makes the tag and directory unable to disagree.
- CI running go mod tidy for you: CI must test what is committed; it checks tidiness and fails instead.
