refactor: Unify path expansion in `EntryPath`

Centralize path expansion logic into `EntryPath.Expand()` to
eliminate code duplication and provide consistent tilde/dot/absolute
path handling across all path types. Previously, expansion logic was
duplicated in `TildeDirPath.Expand()` with subtle differences.

Changes:
- Add `EntryPath.Expand()` with comprehensive tilde, current dir,
  parent dir, and absolute path expansion support
- Add `EntryPath.Clean()` wrapper for `filepath.Clean()`
- Add `EvalSymlinks()` to `EntryPath`, `DirPath`, and `Filepath`
- Refactor `TildeDirPath.Expand()` to delegate to `EntryPath`
- Add `Expand()` to `DirPath` and `Filepath` delegating to
  `EntryPath`
- Add `entry_path_test.go` with platform-specific expansion tests
- Refactor `tilde_dir_path_test.go` to remove expansion tests (moved
  to `EntryPath`) and split remaining parse tests by platform

Benefits:
- Single source of truth for path expansion logic
- Consistent cross-platform behavior (Windows/Linux/macOS)
- Better test coverage with platform-specific test suites
- Reduced code duplication and maintenance burden
