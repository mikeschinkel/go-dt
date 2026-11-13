package dt

import (
	"strings"
)

// GitObjectID is a full git object name: 40 hex digits for a SHA-1
// repository or 64 for a SHA-256 one. It names exactly one object without
// asking a repository.
//
// git accepts uppercase hex but always prints lowercase, and a GitObjectID
// holds lowercase only, so == compares two IDs reliably. ParseGitObjectID
// rejects uppercase unless ParseGitObjectIDArgs.AllowUppercaseHex is set, and
// then stores the value in lowercase.
//
// Every GitObjectID is also a valid GitShortObjectID and GitRevision, so a
// plain cast converts upward.
type GitObjectID string

// GitShortObjectID is a git object name of 4 to 64 hex digits, as git
// rev-parse --short prints one. It may be a full GitObjectID, or a prefix that
// a repository resolves to one object only while no other object shares it.
//
// Like GitObjectID it holds lowercase hex only.
//
// Every GitShortObjectID is also a valid GitRevision, so a plain cast
// converts upward.
type GitShortObjectID string

const (
	gitSHA1HexLen       = 40
	gitSHA256HexLen     = 64
	gitShortObjectIDMin = 4
	gitShortObjectIDMax = gitSHA256HexLen
)

// ParseGitObjectIDArgs holds ParseGitObjectID's options; nil means the
// defaults.
type ParseGitObjectIDArgs struct {
	// AllowUppercaseHex accepts uppercase hex digits and stores them in
	// lowercase.
	AllowUppercaseHex bool
}

// ParseGitShortObjectIDArgs holds ParseGitShortObjectID's options; nil means
// the defaults.
type ParseGitShortObjectIDArgs struct {
	// AllowUppercaseHex accepts uppercase hex digits and stores them in
	// lowercase.
	AllowUppercaseHex bool
}

// ParseGitObjectID returns s as a GitObjectID, or an error wrapping
// ErrInvalidGitObjectID and the sentinel of the rule s breaks. args may be
// nil.
func ParseGitObjectID(s string, args *ParseGitObjectIDArgs) (id GitObjectID, err error) {
	var rule error
	var kvs []ErrKV
	var allowUpper bool

	if args != nil {
		allowUpper = args.AllowUppercaseHex
	}
	kvs, rule = checkGitHex(s, allowUpper)
	if rule == nil && len(s) != gitSHA1HexLen && len(s) != gitSHA256HexLen {
		rule = ErrObjectIDLength
		kvs = []ErrKV{IntKV("length", len(s))}
	}
	if rule != nil {
		err = NewErr(ErrInvalidGitObjectID, rule, StringKV("input", s), kvs)
		goto end
	}
	id = GitObjectID(strings.ToLower(s))
end:
	return id, err
}

// Validate applies ParseGitObjectID's rules to id, with the default args.
func (id GitObjectID) Validate() (err error) {
	_, err = ParseGitObjectID(string(id), nil)
	return err
}

// ParseGitShortObjectID returns s as a GitShortObjectID, or an error wrapping
// ErrInvalidGitShortObjectID and the sentinel of the rule s breaks. args may
// be nil.
func ParseGitShortObjectID(s string, args *ParseGitShortObjectIDArgs) (id GitShortObjectID, err error) {
	var rule error
	var kvs []ErrKV
	var allowUpper bool

	if args != nil {
		allowUpper = args.AllowUppercaseHex
	}
	kvs, rule = checkGitHex(s, allowUpper)
	switch {
	case rule != nil:
	case len(s) < gitShortObjectIDMin:
		rule = ErrTooShort
		kvs = []ErrKV{IntKV("length", len(s))}
	case len(s) > gitShortObjectIDMax:
		rule = ErrTooLong
		kvs = []ErrKV{IntKV("length", len(s))}
	}
	if rule != nil {
		err = NewErr(ErrInvalidGitShortObjectID, rule, StringKV("input", s), kvs)
		goto end
	}
	id = GitShortObjectID(strings.ToLower(s))
end:
	return id, err
}

// Validate applies ParseGitShortObjectID's rules to id, with the default args.
func (id GitShortObjectID) Validate() (err error) {
	_, err = ParseGitShortObjectID(string(id), nil)
	return err
}

// checkGitHex returns the sentinel of the first hex rule s breaks, with
// metadata naming the position, or a nil rule when s is non-empty hex.
func checkGitHex(s string, allowUpper bool) (kvs []ErrKV, rule error) {
	if s == "" {
		rule = ErrEmpty
		goto end
	}
	for i := 0; i < len(s); i++ {
		ch := s[i]
		switch {
		case ch >= '0' && ch <= '9', ch >= 'a' && ch <= 'f':
			continue
		case ch >= 'A' && ch <= 'F':
			if allowUpper {
				continue
			}
			rule = ErrUppercaseHex
		default:
			rule = ErrNonHexDigit
		}
		kvs = []ErrKV{IntKV("position", i)}
		goto end
	}
end:
	return kvs, rule
}
