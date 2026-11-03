# Plan: go-dt StdIO

Mike's edit, moved from the main checkout onto this branch unchanged (commit 477e579): `stdio.go` adds the `StdIO` interface with `Output()`, `ErrOut()` and `Input()`, the unexported `stdIO` struct, `NewStdIO(out, err, in...)` and `NewDefaultStdIO()`. Nothing else changes. go-cliutil E-2550, go-cfgstore E-2551 and gomion E-2540 consume it in that order.

## After landing (Mike)

From the main checkout, land and push main:

```
endless worktree land E-2549 && git push origin main
```

Also from the main checkout, release the root module as v0.7.1 (previous root tag v0.7.0). The empty assignments make `gh` ignore the token in the environment and use the keyring PAT:

```
GITHUB_TOKEN= GH_TOKEN= gh workflow run release.yml -f module=root -f version=v0.7.1 -f description='Adds `dt.StdIO`, an interface exposing `Output()`, `ErrOut()` and `Input()` streams, with `dt.NewStdIO` and `dt.NewDefaultStdIO` constructors, so CLI packages take their standard streams by injection.'
```

E-2550 waits for this tag to resolve from the module proxy.
