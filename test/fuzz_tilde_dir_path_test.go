package test

import (
	"errors"
	"testing"

	"github.com/mikeschinkel/go-dt"
)

func FuzzParseTildeDirPath(f *testing.F) {
	seeds := []string{
		"~",
		"~/",
		"~/sub",
		"~//deep//path",
		"~\\windows",
		"~noslash",
		"",
		"/not/tilde",
	}

	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, input string) {
		home := t.TempDir()
		t.Setenv("HOME", home)
		t.Setenv("USERPROFILE", home)

		tdp, err := dt.ParseTildeDirPath(input)
		if err != nil {
			return
		}

		if _, err = tdp.Expand(); err != nil {
			t.Fatalf("Expand() failed after successful parse: %v", err)
		}
	})
}

func FuzzTildeDirPathExpand(f *testing.F) {
	seeds := []string{
		"~",
		"~/",
		"~/sub/dir",
		"~//deep//path",
		"~\\windows",
		"~noslash",
		"",
		"/not/tilde",
	}

	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, input string) {
		home := t.TempDir()
		t.Setenv("HOME", home)
		t.Setenv("USERPROFILE", home)

		tdp := dt.TildeDirPath(input)

		parsed, parseErr := dt.ParseTildeDirPath(string(tdp))
		_, expandErr := tdp.Expand()

		if parseErr == nil && expandErr != nil {
			t.Fatalf("Expand() returned error %v after Parse success for %q", expandErr, parsed)
		}

		if parseErr != nil && expandErr != nil && !errors.Is(expandErr, parseErr) {
			t.Fatalf("Expand() error mismatch: parse=%v expand=%v", parseErr, expandErr)
		}
	})
}
