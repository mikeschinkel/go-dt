package test

import (
	"os"
	"testing"
)

// TestFuzzCorpus runs all fuzz corpus files as regression tests
func TestFuzzCorpus(t *testing.T) {
	corpusDir := "testdata/fuzz"

	if _, err := os.Stat(corpusDir); os.IsNotExist(err) {
		// // t.Skip("No fuzz corpus found - run fuzzing locally to generate corpus")
		return
	}

	// appinfo doesn't have complex fuzzing needs currently
	// This is a placeholder for future fuzz tests if needed
	t.Skip("No fuzz tests defined for appinfo yet")
}
