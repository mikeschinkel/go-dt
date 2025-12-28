package dtx

import (
	"testing"

	"github.com/mikeschinkel/go-dt"
)

// TestDirPathScanner_Determinism verifies that scanning the same directory
// multiple times always returns files in the same order.
func TestDirPathScanner_Determinism(t *testing.T) {
	const numRuns = 20 // Run many times to catch non-determinism

	// Use the go-dt package directory itself as test data
	// It has enough files to expose ordering issues if they exist
	testDir, err := dt.ParseDirPath("~/Projects/go-pkgs/go-dt")
	if err != nil {
		t.Fatalf("Failed to parse test directory: %v", err)
	}

	// First run - capture baseline
	scanner := NewDirPathScanner(testDir, DirPathScannerArgs{
		ContinueOnErr: false,
	})

	baseline, err := scanner.Scan()
	if err != nil {
		t.Fatalf("First scan failed: %v", err)
	}

	if len(baseline) == 0 {
		t.Fatal("Scan returned no files - test directory may not exist or be empty")
	}

	t.Logf("Baseline scan found %d files", len(baseline))

	// Run multiple times and compare to baseline
	for run := 1; run < numRuns; run++ {
		scanner := NewDirPathScanner(testDir, DirPathScannerArgs{
			ContinueOnErr: false,
		})

		results, err := scanner.Scan()
		if err != nil {
			t.Fatalf("Run %d: scan failed: %v", run, err)
		}

		// Check same number of files
		if len(results) != len(baseline) {
			t.Errorf("Run %d: got %d files, want %d", run, len(results), len(baseline))
			continue
		}

		// Check files are in same order
		for i, result := range results {
			if result != baseline[i] {
				t.Errorf("Run %d: file at index %d differs\n  got:  %s\n  want: %s",
					run, i, result, baseline[i])

				// Show first 5 differences to help debug
				if i < 5 {
					t.Logf("First few files from run %d:", run)
					for j := 0; j < 5 && j < len(results); j++ {
						t.Logf("  [%d] %s", j, results[j])
					}
					t.Logf("First few files from baseline:")
					for j := 0; j < 5 && j < len(baseline); j++ {
						t.Logf("  [%d] %s", j, baseline[j])
					}
				}
				break // Only report first difference per run
			}
		}
	}
}

// TestDirPathScanner_DeterminismWithSkips tests determinism with skip functions
func TestDirPathScanner_DeterminismWithSkips(t *testing.T) {
	const numRuns = 10

	testDir, err := dt.ParseDirPath("~/Projects/go-pkgs")
	if err != nil {
		t.Fatalf("Failed to parse test directory: %v", err)
	}

	skipPaths := []dt.PathSegment{
		".git",
		"vendor",
		"node_modules",
	}

	// First run
	scanner := NewDirPathScanner(testDir, DirPathScannerArgs{
		ContinueOnErr: false,
		SkipPaths:     skipPaths,
	})

	baseline, err := scanner.Scan()
	if err != nil {
		t.Fatalf("First scan failed: %v", err)
	}

	if len(baseline) == 0 {
		t.Fatal("Scan returned no files")
	}

	t.Logf("Baseline scan with skips found %d files", len(baseline))

	// Subsequent runs
	for run := 1; run < numRuns; run++ {
		scanner := NewDirPathScanner(testDir, DirPathScannerArgs{
			ContinueOnErr: false,
			SkipPaths:     skipPaths,
		})

		results, err := scanner.Scan()
		if err != nil {
			t.Fatalf("Run %d: scan failed: %v", run, err)
		}

		if len(results) != len(baseline) {
			t.Errorf("Run %d: got %d files, want %d", run, len(results), len(baseline))
			continue
		}

		for i, result := range results {
			if result != baseline[i] {
				t.Errorf("Run %d: file at index %d differs\n  got:  %s\n  want: %s",
					run, i, result, baseline[i])
				break
			}
		}
	}
}
