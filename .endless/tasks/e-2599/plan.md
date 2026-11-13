# Approach

Rationale is in the analysis. This is how.

## Files (go-dt root package)
- New `git_ref.go`: `GitRef`, `ParseGitRef`, `Validate`, `ToBranchName`, `ToTagName`, and the shared ref-format checker (unexported) that the other types reuse.
- New `git_branch_name.go`, `git_tag_name.go`, `git_remote_name.go`: each type, its `ParseGit…`, `Validate`, and `ToRef` where it has one.
- New `git_revision.go`: `GitRevision`, `ParseGitRevision`, `Validate`, and the `GitRevisionTypes` union constraint.
- New `git_object_id.go`: `GitObjectID`, `GitShortObjectID`, their Parse functions with `ParseGitObjectIDArgs` / `ParseGitShortObjectIDArgs` (`AllowUppercaseHex`), and `Validate`.
- `errors.go`: the sentinels below, in their own commented block like E-2567's. Each parse error is `NewErr(ErrInvalidGit<Type>, <rule sentinel>, "input", s, ...)` in the identifier.go style (`goto end`), adding `"char"` or `"component"` where a rule names one.

## Sentinel errors

Type sentinels (new), one per type:

    ErrInvalidGitRef           = errors.New("invalid git ref")
    ErrInvalidGitBranchName    = errors.New("invalid git branch name")
    ErrInvalidGitTagName       = errors.New("invalid git tag name")
    ErrInvalidGitRemoteName    = errors.New("invalid git remote name")
    ErrInvalidGitRevision      = errors.New("invalid git revision")
    ErrInvalidGitObjectID      = errors.New("invalid git object ID")
    ErrInvalidGitShortObjectID = errors.New("invalid git short object ID")

Rule sentinels (new):

    ErrNotUnderRefs           = errors.New("not under refs/")
    ErrNotBranchRef           = errors.New("not under refs/heads/")          // GitRef.ToBranchName
    ErrNotTagRef              = errors.New("not under refs/tags/")           // GitRef.ToTagName
    ErrEmptyRefComponent      = errors.New("empty path component: leading, trailing or doubled '/'")
    ErrRefComponentLeadingDot = errors.New("path component begins with '.'")
    ErrRefComponentLockSuffix = errors.New("path component ends with '.lock'")
    ErrDoubleDot              = errors.New("contains '..'")
    ErrAtBrace                = errors.New("contains '@{'")
    ErrLoneAt                 = errors.New("is '@'")
    ErrLeadingDash            = errors.New("begins with '-'")
    ErrIsHEAD                 = errors.New("is 'HEAD'")
    ErrSurroundingWhitespace  = errors.New("leading or trailing whitespace")
    ErrNonHexDigit            = errors.New("contains a non-hex digit")
    ErrUppercaseHex           = errors.New("uppercase hex digit; set AllowUppercaseHex to accept it")
    ErrObjectIDLength         = errors.New("not 40 or 64 hex digits")

Rule sentinels go-dt already has, reused:

- `ErrEmpty`: empty value (every type).
- `ErrControlCharacter`: NUL, other bytes below 0x20, or 0x7f (refs and revisions).
- `ErrInvalidCharacter`: space, `~`, `^`, `:`, `?`, `*`, `[` or `\` in a ref, with `"char"` metadata.
- `ErrTrailingPeriod`: a ref ends with `.`.
- `ErrTooShort`, `ErrTooLong`: a short object ID outside 4 to 64 hex digits.

Which type uses which rule:

| Type | Rule sentinels |
|---|---|
| GitRef | ErrEmpty, ErrNotUnderRefs, ErrEmptyRefComponent, ErrRefComponentLeadingDot, ErrRefComponentLockSuffix, ErrDoubleDot, ErrControlCharacter, ErrInvalidCharacter, ErrTrailingPeriod, ErrAtBrace, ErrLoneAt |
| GitBranchName, GitTagName | the GitRef rules on refs/heads/ or refs/tags/ + name (except ErrNotUnderRefs), plus ErrLeadingDash, ErrIsHEAD |
| GitRemoteName | the GitRef rules on refs/remotes/<name>/test (except ErrNotUnderRefs) |
| GitRevision | ErrEmpty, ErrControlCharacter, ErrSurroundingWhitespace, ErrLeadingDash |
| GitObjectID | ErrEmpty, ErrNonHexDigit, ErrUppercaseHex, ErrObjectIDLength |
| GitShortObjectID | ErrEmpty, ErrNonHexDigit, ErrUppercaseHex, ErrTooShort, ErrTooLong |

ToBranchName and ToTagName return ErrInvalidGitRef wrapping ErrNotBranchRef or ErrNotTagRef.
- `README.md`: a git types section.

## Tests (go-dt `test/` module)
- Table tests per type: valid and invalid inputs, each invalid one naming the rule it breaks. The tables are cross-checked once by hand against `git check-ref-format` (and `--branch`) on git 2.52; the tests themselves never run git.
