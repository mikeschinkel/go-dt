package dt

import (
	"strings"
)

// GitRef is a git refname in its full form, such as refs/heads/main,
// refs/tags/v1.0.0 or refs/remotes/origin/main, as gitrevisions(7) and
// git-check-ref-format(1) use the term.
//
// A GitRef always begins with refs/ and passes git check-ref-format with no
// options, so it always splits into a namespace and a name. HEAD, FETCH_HEAD,
// ORIG_HEAD and MERGE_HEAD are not GitRefs; they are GitRevisions.
//
// It rejects, as git does:
//   - a path component that is empty (a leading, trailing or doubled '/'),
//     begins with '.' or ends with '.lock'
//   - '..' and '@{' anywhere
//   - control characters, space, '~', '^', ':', '?', '*', '[' and '\'
//   - a trailing '.'
//
// git also rejects the refname '@', which cannot begin with refs/ and so is
// already rejected as ErrNotUnderRefs.
//
// Every GitRef is also a valid GitRevision, so a plain cast converts upward.
type GitRef string

const (
	gitRefsPrefix    = "refs/"
	gitHeadsPrefix   = "refs/heads/"
	gitTagsPrefix    = "refs/tags/"
	gitRemotesPrefix = "refs/remotes/"
)

// ParseGitRef returns s as a GitRef, or an error wrapping ErrInvalidGitRef and
// the sentinel of the rule s breaks.
func ParseGitRef(s string) (ref GitRef, err error) {
	err = GitRef(s).Validate()
	if err != nil {
		goto end
	}
	ref = GitRef(s)
end:
	return ref, err
}

// Validate applies ParseGitRef's rules to r.
func (r GitRef) Validate() (err error) {
	var rule error
	var kvs []ErrKV

	switch {
	case r == "":
		rule = ErrEmpty
	case !strings.HasPrefix(string(r), gitRefsPrefix):
		rule = ErrNotUnderRefs
	default:
		kvs, rule = checkGitRefFormat(string(r))
	}
	if rule == nil {
		goto end
	}
	err = NewErr(ErrInvalidGitRef, rule, StringKV("input", string(r)), kvs)
end:
	return err
}

// ToBranchName returns the branch name of a ref under refs/heads/, so
// refs/heads/main returns main. It returns an error wrapping ErrInvalidGitRef
// and ErrNotBranchRef when r is not under refs/heads/, or the
// ParseGitBranchName error when the name after refs/heads/ is not a valid
// GitBranchName (refs/heads/HEAD, for example).
func (r GitRef) ToBranchName() (name GitBranchName, err error) {
	s, ok := strings.CutPrefix(string(r), gitHeadsPrefix)
	if !ok {
		err = NewErr(ErrInvalidGitRef, ErrNotBranchRef, StringKV("input", string(r)))
		goto end
	}
	name, err = ParseGitBranchName(s)
end:
	return name, err
}

// ToTagName returns the tag name of a ref under refs/tags/, so
// refs/tags/v1.0.0 returns v1.0.0. It returns an error wrapping
// ErrInvalidGitRef and ErrNotTagRef when r is not under refs/tags/, or the
// ParseGitTagName error when the name after refs/tags/ is not a valid
// GitTagName (refs/tags/HEAD, for example).
func (r GitRef) ToTagName() (name GitTagName, err error) {
	s, ok := strings.CutPrefix(string(r), gitTagsPrefix)
	if !ok {
		err = NewErr(ErrInvalidGitRef, ErrNotTagRef, StringKV("input", string(r)))
		goto end
	}
	name, err = ParseGitTagName(s)
end:
	return name, err
}

// validateGitRefName validates a short name that git checks by placing it in a
// full refname, prefix + name + suffix, as git does for branch, tag and remote
// names. With headRules it also rejects a name that begins with '-' or is HEAD,
// as git-check-ref-format(1)'s --branch does.
func validateGitRefName(name, prefix, suffix string, typeErr error, headRules bool) (err error) {
	var rule error
	var kvs []ErrKV
	var ref string

	switch {
	case name == "":
		rule = ErrEmpty
	case headRules && strings.HasPrefix(name, "-"):
		rule = ErrLeadingDash
	case headRules && name == "HEAD":
		rule = ErrIsHEAD
	default:
		ref = prefix + name + suffix
		kvs, rule = checkGitRefFormat(ref)
		if rule != nil {
			kvs = append([]ErrKV{StringKV("ref", ref)}, kvs...)
		}
	}
	if rule == nil {
		goto end
	}
	err = NewErr(typeErr, rule, StringKV("input", name), kvs)
end:
	return err
}

// checkGitRefFormat returns the sentinel of the first git-check-ref-format(1)
// rule that ref breaks, with metadata naming the component or character, or a
// nil rule when ref passes them all. It follows git's own order: each
// '/'-separated component left to right, then a trailing '.'. It does not
// check that ref has two components; the callers' prefixes guarantee that.
func checkGitRefFormat(ref string) (kvs []ErrKV, rule error) {
	// git checks for a refname of exactly '@' first (rule 9). The callers'
	// prefixes mean ref is never '@', so the check and its sentinel are
	// commented out; see ErrLoneAt in errors.go before uncommenting.
	//
	// if ref == "@" {
	// 	rule = ErrLoneAt
	// 	goto end
	// }
	for component := range strings.SplitSeq(ref, "/") {
		kvs, rule = checkGitRefComponent(component)
		if rule != nil {
			goto end
		}
	}
	if strings.HasSuffix(ref, ".") {
		rule = ErrTrailingPeriod
	}
end:
	return kvs, rule
}

// checkGitRefComponent applies git's rules for one '/'-separated component of
// a refname.
func checkGitRefComponent(c string) (kvs []ErrKV, rule error) {
	var last byte

	for i := 0; i < len(c); i++ {
		ch := c[i]
		switch {
		case ch < 0x20 || ch == 0x7f:
			rule = ErrControlCharacter
		case strings.IndexByte(" ~^:?*[\\", ch) >= 0:
			rule = ErrInvalidCharacter
		case ch == '.' && last == '.':
			rule = ErrDoubleDot
		case ch == '{' && last == '@':
			rule = ErrAtBrace
		}
		if rule != nil {
			kvs = []ErrKV{StringKV("component", c), StringKV("char", string(ch))}
			goto end
		}
		last = ch
	}
	switch {
	case c == "":
		rule = ErrEmptyRefComponent
	case c[0] == '.':
		rule = ErrRefComponentLeadingDot
	case strings.HasSuffix(c, ".lock"):
		rule = ErrRefComponentLockSuffix
	}
	if rule != nil {
		kvs = []ErrKV{StringKV("component", c)}
	}
end:
	return kvs, rule
}
