package dt

// GitBranchName is the short name of a git branch, such as main or
// feature/login, as git-check-ref-format(1)'s --branch checks it.
//
// A name is valid exactly when refs/heads/<name> is a valid GitRef, the name
// does not begin with '-' and is not HEAD. The leading-dash rule stops git
// from reading the name as an option.
//
// Casting a GitBranchName to GitRevision can resolve to a different ref: git
// searches refs/tags/ before refs/heads/ (gitrevisions(7)), so with a tag
// also named main, git warns "refname 'main' is ambiguous" and uses the tag.
// Use ToRef to name the branch and nothing else.
type GitBranchName string

// ParseGitBranchName returns s as a GitBranchName, or an error wrapping
// ErrInvalidGitBranchName and the sentinel of the rule s breaks.
func ParseGitBranchName(s string) (name GitBranchName, err error) {
	err = GitBranchName(s).Validate()
	if err != nil {
		goto end
	}
	name = GitBranchName(s)
end:
	return name, err
}

// Validate applies ParseGitBranchName's rules to b.
func (b GitBranchName) Validate() error {
	return validateGitRefName(string(b), gitHeadsPrefix, "", ErrInvalidGitBranchName, true)
}

// ToRef returns b's full ref, refs/heads/<b>. It never fails.
func (b GitBranchName) ToRef() GitRef {
	return GitRef(gitHeadsPrefix + string(b))
}
