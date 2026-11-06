# Analysis

**Why dtx is in scope (Mike, 2026-11-05, Q1 → B).** One pass brings the whole
repo under ED-1650. The only design point dtx raised was its one test of an
unexported function, and Mike settled it: expose the function rather than
wrap it, since `ParseOSPathSegment` already exists and a `WithOS` variant
taking the GOOS explicitly is a standard shape in his packages.

**Why `ParseOSPathSegmentWithOS(s, goos string)`.** Mike's words: `With<Thing>`
is the more common idiom than `For<Thing>`; the thing is passed after the
other parameters; and since `ParseOSPathSegment()` is the one without a GOOS,
adding `WithOS` to its name is the pattern he uses. Recorded as a lesson in
`.endless/LESSONS.md`.

**Why the coverage flags go (Mike, 2026-11-05, Q2 → A).** The profile the
root `go test -coverprofile` writes has no consumer: no upload step, no badge,
no README mention. After the move the root module has no tests, so the flag
would write an empty profile. Mike does not use coverage today; the justfile
is updated when he starts.

**Why a pure move, not a merge.** `test/dirpath_test.go` (two parse tests) and
the root `dir_path_test.go` (34 method tests) cover the same type. Keeping both
names makes every file a `git mv` the review can see as a rename; merging
would hide 34 tests inside a content diff for no behavioral gain.

**What makes the move safe.** Six root files are already `package dt_test`,
so they can only reach the exported API now. `tilde_test.go` and the three
in-package dtx files other than `path_segments_test.go` reference no
unexported identifier (checked by grepping each lowercase identifier against
the non-test sources). No root test uses `testdata`, go-fsfix or go-testutil,
so no module gains or loses a dependency. `test/go.mod` already replaces
`go-dt => ..` and `dtx/test/go.mod` replaces `dtx => ..`, so the moved tests
run against the code on disk.
