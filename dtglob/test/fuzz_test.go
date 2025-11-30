package test

import (
	"testing"

	"github.com/mikeschinkel/go-dt/dtglob"
)

// FuzzGlobContains tests Glob.Contains with random inputs
func FuzzGlobContains(f *testing.F) {
	seeds := []struct {
		glob   string
		substr string
	}{
		{"*.go", "*"},
		{"**/*.txt", ".txt"},
		{"src/**/*.go", "/"},
		{"", "test"},
		{"test", ""},
		{"a/b/c", "/"},
	}

	for _, seed := range seeds {
		f.Add(seed.glob, seed.substr)
	}

	f.Fuzz(func(t *testing.T, glob, substr string) {
		// Should not panic
		g := dtglob.Glob(glob)
		_ = g.Contains(substr)
	})
}

// FuzzGlobSplit tests Glob.Split with random inputs
func FuzzGlobSplit(f *testing.F) {
	seeds := []struct {
		glob string
		sep  string
	}{
		{"src/test/*.go", "/"},
		{"**/*.txt", "/"},
		{"*.go", "/"},
		{"", "/"},
		{"test", ""},
		{"a/b/c", "/"},
		{"a:b:c", ":"},
	}

	for _, seed := range seeds {
		f.Add(seed.glob, seed.sep)
	}

	f.Fuzz(func(t *testing.T, glob, sep string) {
		// Should not panic
		g := dtglob.Glob(glob)
		_ = g.Split(sep)
	})
}
