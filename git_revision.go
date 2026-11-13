package dt

import (
	"strings"
)

// GitRevision is any expression git resolves to one object, as gitrevisions(7)
// uses the term: HEAD~1, main^{tree}, v1.0.0:path/to/file, @{upstream},
// :/fix bug, a ref, a branch name or an object ID.
//
// GitRevision checks only what keeps a value safe to hand to git: it is not
// empty, holds no NUL or other control character, has no leading or trailing
// space, and does not begin with '-', so git cannot read it as an option. It
// does not parse the gitrevisions(7) grammar; a malformed revision fails when
// git uses it, and git's error says why. git's '-' shorthand for the previous
// branch is therefore not a GitRevision.
//
// A revision range (A..B, A...B, ^A) names a set of commits rather than one
// object and has no type here.
type GitRevision string

// GitRevisionTypes is the union of the git types that name one object, so a
// generic function can take any of them and convert each as it needs; for
// example, turning a GitBranchName into its GitRef before handing it to git.
//
// Provisional: the name may change when a naming convention for union
// constraints is settled.
type GitRevisionTypes interface {
	GitRevision | GitObjectID | GitShortObjectID | GitRef | GitBranchName | GitTagName
}

// ParseGitRevision returns s as a GitRevision, or an error wrapping
// ErrInvalidGitRevision and the sentinel of the rule s breaks.
func ParseGitRevision(s string) (rev GitRevision, err error) {
	err = GitRevision(s).Validate()
	if err != nil {
		goto end
	}
	rev = GitRevision(s)
end:
	return rev, err
}

// Validate applies ParseGitRevision's rules to r.
func (r GitRevision) Validate() (err error) {
	var rule error
	var kvs []ErrKV
	s := string(r)

	switch {
	case s == "":
		rule = ErrEmpty
	case strings.IndexFunc(s, isControlByte) >= 0:
		rule = ErrControlCharacter
		kvs = []ErrKV{IntKV("position", strings.IndexFunc(s, isControlByte))}
	case s[0] == ' ' || s[len(s)-1] == ' ':
		rule = ErrSurroundingWhitespace
	case s[0] == '-':
		rule = ErrLeadingDash
	}
	if rule == nil {
		goto end
	}
	err = NewErr(ErrInvalidGitRevision, rule, StringKV("input", s), kvs)
end:
	return err
}

// isControlByte reports whether r is an ASCII control character: below 0x20,
// or DEL (0x7f).
func isControlByte(r rune) bool {
	return r < 0x20 || r == 0x7f
}
