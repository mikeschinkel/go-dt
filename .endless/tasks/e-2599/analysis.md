# Decisions (planned with Mike, 2026-11-12)

## The family and where it lives
- go-dt's root package, every name prefixed with Git (Mike). A git-only subpackage was rejected: developers would have to discover a second domain-types package and choose between it and go-dt.
- `Name` is dropped only from `GitRef`: a ref already is a name. A branch (tip, upstream, reflog), an annotated tag (an object with tagger and message) and a remote (URLs, fetch rules) are each more than their name, so `GitBranchName`, `GitTagName`, `GitRemoteName` keep it.
- Short object IDs get their own type, `GitShortObjectID` (Mike: C). git spells the concept `--short` (`git rev-parse --short`); "Abbrev" was rejected because Mike avoids abbreviations without a strong reason.

## GitRef holds only full refnames (Mike: Q4 A)
- `refs/...` only, matching `git check-ref-format` with no options. `HEAD`, `FETCH_HEAD`, `ORIG_HEAD`, `MERGE_HEAD` are not GitRefs; they are GitRevisions. Every GitRef therefore splits into a namespace and a name.
- No separate full-refname type: GitRef already holds exactly those values (Mike: ok).

## GitRevision validates only the basics (Mike: Q1 B)
- Not empty, no NUL or control characters, no leading or trailing whitespace, no leading `-`. The leading-dash rule stops a revision from being read by git as an option.
- A malformed but basic-valid revision fails when git uses it (for example `fatal: ambiguous argument 'main^{tre}'`), and go-gitutil returns that error. Mike: better than maintaining a parallel copy of the gitrevisions(7) grammar, which grows as git does.
- Rejected cost accepted: git's `-` shorthand for the previous branch is not a GitRevision.

## Object IDs
- `GitObjectID` holds a full ID only: 40 hex digits (SHA-1) or 64 (SHA-256). It always names exactly one object without asking a repository.
- `GitShortObjectID` holds 4 to 64 hex digits, so every GitObjectID is also a valid GitShortObjectID (plain cast upward).
- Uppercase hex (Mike: D): git accepts it (`git rev-parse --verify` resolves an uppercased ID, verified on git 2.52) but always prints lowercase. `ParseGitObjectID` and `ParseGitShortObjectID` reject uppercase by default; `AllowUppercaseHex` in their args accepts it and stores lowercase, so `==` stays reliable and the arg-less `Validate()` never rejects a value Parse accepted. Args shape follows E-2567's `ParseEmailAddressArgs` (a pointer to an args struct; nil means defaults).

## Conversions (Mike: Q5 B, then ToRef)
- Method names drop the Git prefix, since the receiver already is a git type: `b.ToRef()`, `t.ToRef()`, `r.ToBranchName()`, `r.ToTagName()`.
- Mike chose `ToRef()` over `Ref()` though it never fails: every conversion starts with To, and the signature shows whether it can fail. This departs from E-2567, where every To* returns an error.
- Upward conversions are plain casts: GitObjectID → GitShortObjectID → GitRevision, GitRef → GitRevision.
- Casting a short branch name to GitRevision can resolve to the wrong ref: with a tag `main`, git warns "refname 'main' is ambiguous" and the tag wins (gitrevisions(7) search order). Doc comments say so and point to `ToRef()`.

## A shared union constraint (Mike: Q6 A)
- go-gitutil will take any revision-like git type through Go 1.27 generic methods. A generic call can turn a GitBranchName into its full ref before reaching git, which a plain cast to GitRevision cannot.
- go-dt exports the union so every package shares one list: provisionally `GitRevisionTypes`. No naming convention for union constraints exists yet (E-2567 writes its unions inline); Mike wants one established in gomion, filed separately. The name may change when that lands; go-dt is before v1.0.

## Validation rules (git 2.52)
- Ref rules, from `git check-ref-format`: no component begins with `.` or ends with `.lock`; no `..`; no control characters, space, `~`, `^`, `:`, `?`, `*`, `[` or `\`; does not begin or end with `/` and has no `//`; does not end with `.`; no `@{`; is not `@`.
- GitBranchName: `refs/heads/<name>` passes the ref rules, `<name>` does not begin with `-` and is not `HEAD` (as `check-ref-format --branch`).
- GitTagName: `refs/tags/<name>` passes the ref rules, `<name>` does not begin with `-` and is not `HEAD` (git 2.52 refuses both as tag names).
- GitRemoteName: git's remote-name check, `refs/remotes/<name>/test` passes the ref rules.
- All written in Go; go-dt never runs git and stays dependency-free.

## Not here
- Revision ranges (`A..B`, `A...B`, `^A`) name a set of commits, not one object; no type.
- go-gitutil's move from its own GitRef (a revision) and RemoteName to these types is go-gitutil's work (E-2600).
