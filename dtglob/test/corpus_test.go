package test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mikeschinkel/go-dt/dtglob"
)

// TestFuzzCorpus runs all fuzz corpus files as regression tests
func TestFuzzCorpus(t *testing.T) {
	corpusDir := "testdata/fuzz"

	if _, err := os.Stat(corpusDir); os.IsNotExist(err) {
		// t.Skip("No fuzz corpus found - run fuzzing locally to generate corpus")
		return
	}

	fuzzTests := []string{
		"FuzzGlobContains",
		"FuzzGlobSplit",
	}

	for _, fuzzTest := range fuzzTests {
		t.Run(fuzzTest, func(t *testing.T) {
			testDir := filepath.Join(corpusDir, fuzzTest)

			if _, err := os.Stat(testDir); os.IsNotExist(err) {
				// t.Skipf("No corpus for %s", fuzzTest)
				return
			}

			entries, err := os.ReadDir(testDir)
			if err != nil {
				t.Fatalf("Failed to read corpus directory: %v", err)
			}

			if len(entries) == 0 {
				// t.Skipf("No corpus files for %s", fuzzTest)
				return
			}

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

					switch fuzzTest {
					case "FuzzGlobContains":
						runGlobContainsCorpus(t, data)
					case "FuzzGlobSplit":
						runGlobSplitCorpus(t, data)
					}
				})
			}
		})
	}
}

func runGlobContainsCorpus(t *testing.T, data []byte) {
	parts := extractStringsFromCorpus(data, 2)

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Glob.Contains panicked with inputs: %q, panic: %v", parts, r)
		}
	}()

	if len(parts) >= 2 {
		g := dtglob.Glob(parts[0])
		_ = g.Contains(parts[1])
	}
}

func runGlobSplitCorpus(t *testing.T, data []byte) {
	parts := extractStringsFromCorpus(data, 2)

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Glob.Split panicked with inputs: %q, panic: %v", parts, r)
		}
	}()

	if len(parts) >= 2 {
		g := dtglob.Glob(parts[0])
		_ = g.Split(parts[1])
	}
}

func extractStringsFromCorpus(data []byte, count int) []string {
	result := make([]string, count)
	str := string(data)

	if count > 0 {
		result[0] = str
		if count > 1 {
			result[1] = str
		}
	}

	return result
}
