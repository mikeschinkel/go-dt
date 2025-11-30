package test

import (
	"testing"

	"github.com/mikeschinkel/go-dt"
)

// FuzzParseDirPath tests ParseDirPath with random inputs
func FuzzParseDirPath(f *testing.F) {
	// Seed corpus with valid and edge cases
	seeds := []string{
		"/",
		"/usr/local/bin",
		"relative/path",
		".",
		"..",
		"",
		"/path/with/trailing/slash/",
		"//double//slashes",
		"/path/with/../parent",
		"/path/with/./current",
		"C:\\Windows\\Path", // Windows-style path
		"/path with spaces",
		"/path-with-dashes",
		"/path_with_underscores",
		"/path.with.dots",
		"/very/long/path/that/goes/on/and/on/for/many/levels/deep",
	}

	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, input string) {
		// Should not panic
		_, err := dt.ParseDirPath(input)

		// If it returns without error, the result should be usable
		if err == nil {
			// Test that we can convert it back to string
			dp, _ := dt.ParseDirPath(input)
			_ = string(dp)
		}
	})
}

// FuzzParseDirPaths tests ParseDirPaths with random inputs
func FuzzParseDirPaths(f *testing.F) {
	// Seed corpus with various path combinations (using single strings)
	seeds := []string{
		"/usr/local,/var/log,/tmp",
		"/,/home,/etc",
		"relative,paths,here",
		"",
		"/single",
		",,,",
		"/path1,,/path2",
	}

	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, input string) {
		// Split by comma for testing
		var paths []string
		if input != "" {
			// Simple split - in real use, paths wouldn't have commas
			paths = []string{input}
		}

		// Should not panic
		_, _ = dt.ParseDirPaths(paths)
	})
}

// FuzzDirPathJoin tests DirPathJoin with random inputs
func FuzzDirPathJoin(f *testing.F) {
	// Seed corpus with various join scenarios
	seeds := []struct {
		part1 string
		part2 string
	}{
		{"/usr/local", "bin"},
		{"relative", "path"},
		{"", "empty"},
		{"/", "root"},
		{".", "current"},
		{"..", "parent"},
		{"/trailing/", "slash"},
		{"multiple/parts", "here"},
	}

	for _, seed := range seeds {
		f.Add(seed.part1, seed.part2)
	}

	f.Fuzz(func(t *testing.T, part1, part2 string) {
		// Should not panic
		_ = dt.DirPathJoin(part1, part2)

		// Also test with DirPath.Join method
		dp := dt.DirPath(part1)
		_ = dp.Join(part2)
	})
}

// FuzzFilepath tests Filepath operations with random inputs
func FuzzFilepath(f *testing.F) {
	seeds := []string{
		"/etc/config.json",
		"/usr/local/bin/app",
		"relative/file.txt",
		"noextension",
		".hidden",
		"file.with.multiple.dots.txt",
		"/path/to/file",
		"",
	}

	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, input string) {
		// Create a Filepath - should not panic
		fp := dt.Filepath(input)

		// Basic operations should not panic
		_ = string(fp)
	})
}

// FuzzFilename tests Filename operations with random inputs
func FuzzFilename(f *testing.F) {
	seeds := []string{
		"config.json",
		"app",
		"file.txt",
		"",
		".hidden",
		"file.with.multiple.dots.txt",
		"very-long-filename-with-many-characters-in-it.extension",
		"special!@#$%chars",
	}

	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, input string) {
		// Create a Filename - should not panic
		fn := dt.Filename(input)

		// Basic operations should not panic
		_ = string(fn)
	})
}
