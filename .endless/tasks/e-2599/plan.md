# Approach

Rationale is in the analysis. This is how.

## Files (go-dt root package)
- New `git_ref.go`: `GitRef`, `ParseGitRef`, `Validate`, `ToBranchName`, `ToTagName`, and the shared ref-format checker (unexported) that the other types reuse.
- New `git_branch_name.go`, `git_tag_name.go`, `git_remote_name.go`: each type, its `ParseGit…`, `Validate`, and `ToRef` where it has one.
- New `git_revision.go`: `GitRevision`, `ParseGitRevision`, `Validate`, and the `GitRevisionTypes` union constraint.
- New `git_object_id.go`: `GitObjectID`, `GitShortObjectID`, their Parse functions with `ParseGitObjectIDArgs` / `ParseGitShortObjectIDArgs` (`AllowUppercaseHex`), and `Validate`.
- `errors.go`: one sentinel per type (`ErrInvalidGitRef`, …) plus a sentinel per rule broken, wrapped with doterr `NewErr` in the identifier.go style (`goto end`).
- `README.md`: a git types section.

## Tests (go-dt `test/` module)
- Table tests per type: valid and invalid inputs, each invalid one naming the rule it breaks. The tables are cross-checked once by hand against `git check-ref-format` (and `--branch`) on git 2.52; the tests themselves never run git.
