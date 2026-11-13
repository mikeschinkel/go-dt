package test

import (
	"errors"
	"strings"
	"testing"

	"github.com/mikeschinkel/go-dt"
)

// gitCase is one row of a git type's table. A nil rule means the input is
// valid. The ref tables were cross-checked by hand against git
// check-ref-format (with and without --branch) on git 2.52; the tests never
// run git.
type gitCase struct {
	input string
	rule  error
}

// checkGitCase asserts that err matches c: nil for a valid row, and otherwise
// an error that wraps both typeErr and c.rule.
func checkGitCase(t *testing.T, c gitCase, typeErr, err error) {
	t.Helper()
	switch {
	case c.rule == nil && err != nil:
		t.Errorf("%q: unexpected error: %v", c.input, err)
	case c.rule == nil:
	case err == nil:
		t.Errorf("%q: expected an error wrapping %v, got nil", c.input, c.rule)
	case !errors.Is(err, typeErr):
		t.Errorf("%q: error %v does not wrap %v", c.input, err, typeErr)
	case !errors.Is(err, c.rule):
		t.Errorf("%q: error %v does not wrap %v", c.input, err, c.rule)
	}
}

func TestParseGitRef(t *testing.T) {
	cases := []gitCase{
		{"refs/heads/main", nil},
		{"refs/heads/feature/login", nil},
		{"refs/tags/v1.0.0", nil},
		{"refs/remotes/origin/main", nil},
		{"refs/heads/@", nil},
		{"refs/heads/a@b", nil},
		{"refs/heads/{", nil},
		{"refs/heads/é", nil},
		{"refs/heads/-x", nil},
		{"refs/heads/HEAD", nil},
		{"", dt.ErrEmpty},
		{"HEAD", dt.ErrNotUnderRefs},
		{"FETCH_HEAD", dt.ErrNotUnderRefs},
		{"ORIG_HEAD", dt.ErrNotUnderRefs},
		{"MERGE_HEAD", dt.ErrNotUnderRefs},
		{"heads/main", dt.ErrNotUnderRefs},
		{"refs", dt.ErrNotUnderRefs},
		{"refs/", dt.ErrEmptyRefComponent},
		{"refs/heads/", dt.ErrEmptyRefComponent},
		{"refs//heads", dt.ErrEmptyRefComponent},
		{"refs/heads/.a", dt.ErrRefComponentLeadingDot},
		{"refs/heads/a/.b", dt.ErrRefComponentLeadingDot},
		{"refs/heads/a.lock", dt.ErrRefComponentLockSuffix},
		{"refs/heads/a.lock/b", dt.ErrRefComponentLockSuffix},
		{"refs/heads/a..b", dt.ErrDoubleDot},
		{"refs/heads/a\x00b", dt.ErrControlCharacter},
		{"refs/heads/a\tb", dt.ErrControlCharacter},
		{"refs/heads/a\x7fb", dt.ErrControlCharacter},
		{"refs/heads/a b", dt.ErrInvalidCharacter},
		{"refs/heads/a~1", dt.ErrInvalidCharacter},
		{"refs/heads/a^", dt.ErrInvalidCharacter},
		{"refs/heads/a:b", dt.ErrInvalidCharacter},
		{"refs/heads/a?", dt.ErrInvalidCharacter},
		{"refs/heads/a*", dt.ErrInvalidCharacter},
		{"refs/heads/a[b", dt.ErrInvalidCharacter},
		{"refs/heads/a\\b", dt.ErrInvalidCharacter},
		{"refs/heads/a.", dt.ErrTrailingPeriod},
		{"refs/heads/a@{b", dt.ErrAtBrace},
	}
	for _, c := range cases {
		ref, err := dt.ParseGitRef(c.input)
		checkGitCase(t, c, dt.ErrInvalidGitRef, err)
		checkGitCase(t, c, dt.ErrInvalidGitRef, dt.GitRef(c.input).Validate())
		if c.rule == nil && string(ref) != c.input {
			t.Errorf("%q: parsed as %q", c.input, ref)
		}
	}
}

// branchAndTagCases holds the rows GitBranchName and GitTagName share: both
// check refs/<namespace>/<name> plus the leading-dash and HEAD rules.
var branchAndTagCases = []gitCase{
	{"main", nil},
	{"feature/login", nil},
	{"v1.0.0", nil},
	{"@", nil},
	{"a@b", nil},
	{"{", nil},
	{"é", nil},
	{"x-y", nil},
	{"", dt.ErrEmpty},
	{"-x", dt.ErrLeadingDash},
	{"HEAD", dt.ErrIsHEAD},
	{"/a", dt.ErrEmptyRefComponent},
	{"a/", dt.ErrEmptyRefComponent},
	{"a//b", dt.ErrEmptyRefComponent},
	{".a", dt.ErrRefComponentLeadingDot},
	{"a/.b", dt.ErrRefComponentLeadingDot},
	{"a.lock", dt.ErrRefComponentLockSuffix},
	{"foo.lock/x", dt.ErrRefComponentLockSuffix},
	{"a..b", dt.ErrDoubleDot},
	{"a\x00b", dt.ErrControlCharacter},
	{"a b", dt.ErrInvalidCharacter},
	{"x*", dt.ErrInvalidCharacter},
	{"a~1", dt.ErrInvalidCharacter},
	{"a.", dt.ErrTrailingPeriod},
	{"a@{b", dt.ErrAtBrace},
	{"@{-1}", dt.ErrAtBrace},
}

func TestParseGitBranchName(t *testing.T) {
	for _, c := range branchAndTagCases {
		name, err := dt.ParseGitBranchName(c.input)
		checkGitCase(t, c, dt.ErrInvalidGitBranchName, err)
		checkGitCase(t, c, dt.ErrInvalidGitBranchName, dt.GitBranchName(c.input).Validate())
		if c.rule == nil && string(name) != c.input {
			t.Errorf("%q: parsed as %q", c.input, name)
		}
	}
}

func TestParseGitTagName(t *testing.T) {
	for _, c := range branchAndTagCases {
		name, err := dt.ParseGitTagName(c.input)
		checkGitCase(t, c, dt.ErrInvalidGitTagName, err)
		checkGitCase(t, c, dt.ErrInvalidGitTagName, dt.GitTagName(c.input).Validate())
		if c.rule == nil && string(name) != c.input {
			t.Errorf("%q: parsed as %q", c.input, name)
		}
	}
}

func TestParseGitRemoteName(t *testing.T) {
	cases := []gitCase{
		{"origin", nil},
		{"upstream", nil},
		{"a.b", nil},
		{"a/b", nil},
		{"HEAD", nil},
		{"", dt.ErrEmpty},
		{"/a", dt.ErrEmptyRefComponent},
		{"a/", dt.ErrEmptyRefComponent},
		{".a", dt.ErrRefComponentLeadingDot},
		{"a.lock", dt.ErrRefComponentLockSuffix},
		{"a..b", dt.ErrDoubleDot},
		{"a\nb", dt.ErrControlCharacter},
		{"a b", dt.ErrInvalidCharacter},
		{"x*", dt.ErrInvalidCharacter},
		{"a@{b", dt.ErrAtBrace},
	}
	for _, c := range cases {
		name, err := dt.ParseGitRemoteName(c.input)
		checkGitCase(t, c, dt.ErrInvalidGitRemoteName, err)
		checkGitCase(t, c, dt.ErrInvalidGitRemoteName, dt.GitRemoteName(c.input).Validate())
		if c.rule == nil && string(name) != c.input {
			t.Errorf("%q: parsed as %q", c.input, name)
		}
	}
}

func TestParseGitRevision(t *testing.T) {
	cases := []gitCase{
		{"HEAD", nil},
		{"HEAD~1", nil},
		{"main^{tree}", nil},
		{"v1:path", nil},
		{"@{upstream}", nil},
		{":/fix bug", nil},
		{"refs/heads/main", nil},
		{"0123abcd", nil},
		{"main..dev", nil},
		{"", dt.ErrEmpty},
		{"a\x00b", dt.ErrControlCharacter},
		{"a\tb", dt.ErrControlCharacter},
		{"a\x7f", dt.ErrControlCharacter},
		{" HEAD", dt.ErrSurroundingWhitespace},
		{"HEAD ", dt.ErrSurroundingWhitespace},
		{"-x", dt.ErrLeadingDash},
		{"-", dt.ErrLeadingDash},
	}
	for _, c := range cases {
		rev, err := dt.ParseGitRevision(c.input)
		checkGitCase(t, c, dt.ErrInvalidGitRevision, err)
		checkGitCase(t, c, dt.ErrInvalidGitRevision, dt.GitRevision(c.input).Validate())
		if c.rule == nil && string(rev) != c.input {
			t.Errorf("%q: parsed as %q", c.input, rev)
		}
	}
}

const (
	sha1ID   = "2603b9065734ed8d35b946422a72d72b95af6249"
	sha256ID = "2603b9065734ed8d35b946422a72d72b95af62492603b9065734ed8d35b94642"
)

func TestParseGitObjectID(t *testing.T) {
	cases := []gitCase{
		{sha1ID, nil},
		{sha256ID, nil},
		{"", dt.ErrEmpty},
		{strings.Replace(sha1ID, "2", "g", 1), dt.ErrNonHexDigit},
		{strings.ToUpper(sha1ID), dt.ErrUppercaseHex},
		{sha1ID[:39], dt.ErrObjectIDLength},
		{sha1ID + "0", dt.ErrObjectIDLength},
		{sha256ID[:63], dt.ErrObjectIDLength},
		{sha256ID + "0", dt.ErrObjectIDLength},
		{"abcd", dt.ErrObjectIDLength},
	}
	for _, c := range cases {
		id, err := dt.ParseGitObjectID(c.input, nil)
		checkGitCase(t, c, dt.ErrInvalidGitObjectID, err)
		checkGitCase(t, c, dt.ErrInvalidGitObjectID, dt.GitObjectID(c.input).Validate())
		if c.rule == nil && string(id) != c.input {
			t.Errorf("%q: parsed as %q", c.input, id)
		}
	}
}

func TestParseGitShortObjectID(t *testing.T) {
	cases := []gitCase{
		{"abcd", nil},
		{"0123456", nil},
		{sha1ID, nil},
		{sha256ID, nil},
		{"", dt.ErrEmpty},
		{"abcx", dt.ErrNonHexDigit},
		{"ABCD", dt.ErrUppercaseHex},
		{"abc", dt.ErrTooShort},
		{sha256ID + "0", dt.ErrTooLong},
	}
	for _, c := range cases {
		id, err := dt.ParseGitShortObjectID(c.input, nil)
		checkGitCase(t, c, dt.ErrInvalidGitShortObjectID, err)
		checkGitCase(t, c, dt.ErrInvalidGitShortObjectID, dt.GitShortObjectID(c.input).Validate())
		if c.rule == nil && string(id) != c.input {
			t.Errorf("%q: parsed as %q", c.input, id)
		}
	}
}

func TestGitObjectIDAllowUppercaseHex(t *testing.T) {
	upper := strings.ToUpper(sha1ID)

	id, err := dt.ParseGitObjectID(upper, &dt.ParseGitObjectIDArgs{AllowUppercaseHex: true})
	if err != nil {
		t.Fatalf("ParseGitObjectID(%q, AllowUppercaseHex): %v", upper, err)
	}
	if id != dt.GitObjectID(sha1ID) {
		t.Errorf("got %q, want the lowercase %q", id, sha1ID)
	}
	if err = id.Validate(); err != nil {
		t.Errorf("Validate() on a parsed uppercase ID: %v", err)
	}

	short, err := dt.ParseGitShortObjectID("ABCD", &dt.ParseGitShortObjectIDArgs{AllowUppercaseHex: true})
	if err != nil {
		t.Fatalf("ParseGitShortObjectID(ABCD, AllowUppercaseHex): %v", err)
	}
	if short != "abcd" {
		t.Errorf("got %q, want abcd", short)
	}
	if err = short.Validate(); err != nil {
		t.Errorf("Validate() on a parsed uppercase short ID: %v", err)
	}

	_, err = dt.ParseGitObjectID(upper, &dt.ParseGitObjectIDArgs{})
	if !errors.Is(err, dt.ErrUppercaseHex) {
		t.Errorf("zero args: got %v, want ErrUppercaseHex", err)
	}
}

func TestGitToRef(t *testing.T) {
	if got := dt.GitBranchName("feature/login").ToRef(); got != "refs/heads/feature/login" {
		t.Errorf("GitBranchName.ToRef() = %q", got)
	}
	if got := dt.GitTagName("v1.0.0").ToRef(); got != "refs/tags/v1.0.0" {
		t.Errorf("GitTagName.ToRef() = %q", got)
	}
}

func TestGitRefToBranchName(t *testing.T) {
	name, err := dt.GitRef("refs/heads/feature/login").ToBranchName()
	if err != nil || name != "feature/login" {
		t.Errorf("ToBranchName() = %q, %v", name, err)
	}
	cases := []struct {
		ref     dt.GitRef
		typeErr error
		ruleErr error
	}{
		{"refs/tags/v1", dt.ErrInvalidGitRef, dt.ErrNotBranchRef},
		{"HEAD", dt.ErrInvalidGitRef, dt.ErrNotBranchRef},
		{"refs/heads/HEAD", dt.ErrInvalidGitBranchName, dt.ErrIsHEAD},
		{"refs/heads/-x", dt.ErrInvalidGitBranchName, dt.ErrLeadingDash},
	}
	for _, c := range cases {
		_, err = c.ref.ToBranchName()
		if !errors.Is(err, c.typeErr) || !errors.Is(err, c.ruleErr) {
			t.Errorf("%q.ToBranchName(): got %v, want %v and %v", c.ref, err, c.typeErr, c.ruleErr)
		}
	}
}

func TestGitRefToTagName(t *testing.T) {
	name, err := dt.GitRef("refs/tags/v1.0.0").ToTagName()
	if err != nil || name != "v1.0.0" {
		t.Errorf("ToTagName() = %q, %v", name, err)
	}
	cases := []struct {
		ref     dt.GitRef
		typeErr error
		ruleErr error
	}{
		{"refs/heads/main", dt.ErrInvalidGitRef, dt.ErrNotTagRef},
		{"refs/tagsx/v1", dt.ErrInvalidGitRef, dt.ErrNotTagRef},
		{"refs/tags/HEAD", dt.ErrInvalidGitTagName, dt.ErrIsHEAD},
	}
	for _, c := range cases {
		_, err = c.ref.ToTagName()
		if !errors.Is(err, c.typeErr) || !errors.Is(err, c.ruleErr) {
			t.Errorf("%q.ToTagName(): got %v, want %v and %v", c.ref, err, c.typeErr, c.ruleErr)
		}
	}
}

// firstGitRevision is generic over GitRevisionTypes, so it compiles only
// while every listed type is in the union.
func firstGitRevision[T dt.GitRevisionTypes](v T) dt.GitRevision {
	return dt.GitRevision(v)
}

func TestGitRevisionTypesUnion(t *testing.T) {
	revs := []dt.GitRevision{
		firstGitRevision(dt.GitRevision("HEAD~1")),
		firstGitRevision(dt.GitObjectID(sha1ID)),
		firstGitRevision(dt.GitShortObjectID("abcd")),
		firstGitRevision(dt.GitRef("refs/heads/main")),
		firstGitRevision(dt.GitBranchName("main")),
		firstGitRevision(dt.GitTagName("v1.0.0")),
	}
	for _, r := range revs {
		if err := r.Validate(); err != nil {
			t.Errorf("upward cast %q is not a valid GitRevision: %v", r, err)
		}
	}
}
