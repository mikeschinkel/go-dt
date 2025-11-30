package test

import (
	"testing"

	"github.com/mikeschinkel/go-dt"
	"github.com/mikeschinkel/go-dt/dtx"
)

// FuzzParseOSPathSegment tests ParseOSPathSegment with random inputs
func FuzzParseOSPathSegment(f *testing.F) {
	seeds := []string{
		"valid-segment",
		"path",
		"file.txt",
		"",
		"segment_with_underscore",
		"segment-with-dash",
		"123numeric",
		"MixedCase",
		"with spaces",
		"/invalid/slash",
		"back\\slash",
		"special!@#$",
	}

	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, input string) {
		// Should not panic
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("ParseOSPathSegment panicked with input: %q, panic: %v", input, r)
			}
		}()

		_, err := dtx.ParseOSPathSegment(input)
		if err == nil {
			// If it succeeds, verify the result is usable
			ps, _ := dtx.ParseOSPathSegment(input)
			_ = string(ps)
		}
	})
}

// FuzzIsZero tests IsZero with various types
func FuzzIsZero(f *testing.F) {
	seeds := []string{
		"",
		"non-empty",
		"0",
		"false",
		"nil",
	}

	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, input string) {
		// Should not panic
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("IsZero panicked with input: %q, panic: %v", input, r)
			}
		}()

		// Test with string
		_ = dtx.IsZero(input)

		// Test with converted types
		_ = dtx.IsZero(dt.Filename(input))
		_ = dtx.IsZero(dt.DirPath(input))
	})
}
