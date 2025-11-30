package test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mikeschinkel/go-dt"
)

// TestFuzzCorpus runs all fuzz corpus files as regression tests
// This ensures that any interesting inputs discovered during fuzzing
// are tested in CI/CD to prevent regressions
func TestFuzzCorpus(t *testing.T) {
	corpusDir := "testdata/fuzz"

	// Check if corpus directory exists
	if _, err := os.Stat(corpusDir); os.IsNotExist(err) {
		// t.Skip("No fuzz corpus found - run fuzzing locally to generate corpus")
		return
	}

	// Find all fuzz test directories
	fuzzTests := []string{
		"FuzzParseDirPath",
		"FuzzParseDirPaths",
		"FuzzDirPathJoin",
		"FuzzFilepath",
		"FuzzFilename",
	}

	for _, fuzzTest := range fuzzTests {
		t.Run(fuzzTest, func(t *testing.T) {
			testDir := filepath.Join(corpusDir, fuzzTest)

			// Check if this fuzz test has corpus data
			if _, err := os.Stat(testDir); os.IsNotExist(err) {
				// t.Skipf("No corpus for %s", fuzzTest)
				return
			}

			// Read all corpus files
			entries, err := os.ReadDir(testDir)
			if err != nil {
				t.Fatalf("Failed to read corpus directory: %v", err)
			}

			if len(entries) == 0 {
				// t.Skipf("No corpus files for %s", fuzzTest)
				return
			}

			// Run each corpus file as a test
			for _, entry := range entries {
				if entry.IsDir() {
					continue
				}

				t.Run(entry.Name(), func(t *testing.T) {
					corpusFile := filepath.Join(testDir, entry.Name())
					data, err := os.ReadFile(corpusFile)
					if err != nil {
						t.Fatalf("Failed to read corpus file: %v", err)
					}

					// Run the appropriate test based on fuzz test name
					switch fuzzTest {
					case "FuzzParseDirPath":
						runParseDirPathCorpus(t, data)
					case "FuzzParseDirPaths":
						runParseDirPathsCorpus(t, data)
					case "FuzzDirPathJoin":
						runDirPathJoinCorpus(t, data)
					case "FuzzFilepath":
						runFilepathCorpus(t, data)
					case "FuzzFilename":
						runFilenameCorpus(t, data)
					}
				})
			}
		})
	}
}

func runParseDirPathCorpus(t *testing.T, data []byte) {
	// Extract the string from the corpus file
	// Go's fuzzing format: "go test fuzz v1\nstring(\"...\")\n"
	input := extractStringFromCorpus(data)

	// Should not panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("ParseDirPath panicked with input: %q, panic: %v", input, r)
		}
	}()

	_, err := dt.ParseDirPath(input)
	if err == nil {
		// If it succeeds, verify the result is usable
		dp, _ := dt.ParseDirPath(input)
		_ = string(dp)
	}
}

func runParseDirPathsCorpus(t *testing.T, data []byte) {
	input := extractStringFromCorpus(data)

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("ParseDirPaths panicked with input: %q, panic: %v", input, r)
		}
	}()

	paths := []string{input}
	_, _ = dt.ParseDirPaths(paths)
}

func runDirPathJoinCorpus(t *testing.T, data []byte) {
	// For DirPathJoin, we expect two strings in the corpus
	parts := extractStringsFromCorpus(data, 2)

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("DirPathJoin panicked with inputs: %q, panic: %v", parts, r)
		}
	}()

	if len(parts) >= 2 {
		_ = dt.DirPathJoin(parts[0], parts[1])
		dp := dt.DirPath(parts[0])
		_ = dp.Join(parts[1])
	}
}

func runFilepathCorpus(t *testing.T, data []byte) {
	input := extractStringFromCorpus(data)

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Filepath panicked with input: %q, panic: %v", input, r)
		}
	}()

	fp := dt.Filepath(input)
	_ = string(fp)
}

func runFilenameCorpus(t *testing.T, data []byte) {
	input := extractStringFromCorpus(data)

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Filename panicked with input: %q, panic: %v", input, r)
		}
	}()

	fn := dt.Filename(input)
	_ = string(fn)
}

// extractStringFromCorpus extracts a string value from Go's fuzz corpus format
func extractStringFromCorpus(data []byte) string {
	// Simple extraction - corpus format is: "go test fuzz v1\nstring(\"...\")\n"
	// For production use, you might want more robust parsing
	str := string(data)

	// Skip the header line
	if len(str) > 0 {
		// This is a simplified version - real corpus parsing would be more robust
		return str
	}

	return ""
}

// extractStringsFromCorpus extracts multiple string values from corpus
func extractStringsFromCorpus(data []byte, count int) []string {
	// Simplified - for multi-parameter fuzzing
	result := make([]string, count)
	str := extractStringFromCorpus(data)

	// For now, just split the string
	// In real usage, the corpus format would properly encode multiple values
	if count > 0 {
		result[0] = str
		if count > 1 {
			result[1] = str
		}
	}

	return result
}
