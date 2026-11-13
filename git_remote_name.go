package dt

// GitRemoteName is the name of a git remote, such as origin or upstream.
//
// A name is valid exactly when refs/remotes/<name>/test is a valid GitRef,
// the check git itself applies to a remote's name.
type GitRemoteName string

// ParseGitRemoteName returns s as a GitRemoteName, or an error wrapping
// ErrInvalidGitRemoteName and the sentinel of the rule s breaks.
func ParseGitRemoteName(s string) (name GitRemoteName, err error) {
	err = GitRemoteName(s).Validate()
	if err != nil {
		goto end
	}
	name = GitRemoteName(s)
end:
	return name, err
}

// Validate applies ParseGitRemoteName's rules to r.
func (r GitRemoteName) Validate() error {
	return validateGitRefName(string(r), gitRemotesPrefix, "/test", ErrInvalidGitRemoteName, false)
}
