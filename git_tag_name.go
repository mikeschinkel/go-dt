package dt

// GitTagName is the short name of a git tag, such as v1.0.0.
//
// A name is valid exactly when refs/tags/<name> is a valid GitRef, the name
// does not begin with '-' and is not HEAD; git 2.52 refuses both as tag names.
//
// Casting a GitTagName to GitRevision can resolve to a different ref: git
// looks a short name up as $GIT_DIR/<name> and refs/<name> before
// refs/tags/<name> (gitrevisions(7)), so a tag named heads/main resolves to
// the branch refs/heads/main. Use ToRef to name the tag and nothing else.
type GitTagName string

// ParseGitTagName returns s as a GitTagName, or an error wrapping
// ErrInvalidGitTagName and the sentinel of the rule s breaks.
func ParseGitTagName(s string) (name GitTagName, err error) {
	err = GitTagName(s).Validate()
	if err != nil {
		goto end
	}
	name = GitTagName(s)
end:
	return name, err
}

// Validate applies ParseGitTagName's rules to t.
func (t GitTagName) Validate() error {
	return validateGitRefName(string(t), gitTagsPrefix, "", ErrInvalidGitTagName, true)
}

// ToRef returns t's full ref, refs/tags/<t>. It never fails.
func (t GitTagName) ToRef() GitRef {
	return GitRef(gitTagsPrefix + string(t))
}
