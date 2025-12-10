//go:build !windows

// Benchmark directory walking: stdlib filepath.WalkDir vs dt.DirPath.Walk.
//
// These benchmarks start from the current user's home directory,
// skip ".git" directories, and count files vs directories.
//
// NOTE: Scanning $HOME can be quite expensive on large systems.
// When running these benchmarks, you may want to constrain benchtime, e.g.:
//   go test -run=^$ -bench=WalkDir -benchtime=1x ./...

package test

import (
	"errors"
	"io/fs"
	"path/filepath"
	"syscall"
	"testing"

	"github.com/mikeschinkel/go-dt"
)

func BenchmarkFilepathWalkDir(b *testing.B) {
	b.ResetTimer()

	allFiles, allDirs := 0, 0

	for i := 0; i < b.N; i++ {
		files, dirs, err := filepathWalk(benchmarkDir(b))
		if err != nil {
			b.Fatalf("FilepathWalk: %v", err)
		}

		allFiles += files
		allDirs += dirs
	}
	b.Logf("FilepathWalk:")
	b.Logf("- Files: %d", allFiles)
	b.Logf("- Dirs: %d", allDirs)
}

func BenchmarkDirPathWalk(b *testing.B) {

	b.ResetTimer()

	allFiles, allDirs := 0, 0

	for i := 0; i < b.N; i++ {
		files, dirs, err := dirPathWalk(benchmarkDir(b))
		if err != nil {
			b.Fatalf("DirPathWalk: %v", err)
		}

		allFiles += files
		allDirs += dirs
	}
	b.Logf("DirPathWalk:")
	b.Logf("- Files: %d", allFiles)
	b.Logf("- Dirs: %d", allDirs)

}
func BenchmarkDTWalkDir(b *testing.B) {

	b.ResetTimer()

	allFiles, allDirs := 0, 0

	for i := 0; i < b.N; i++ {
		files, dirs, err := dtWalkDir(benchmarkDir(b))
		if err != nil {
			b.Fatalf("DTWalkDir: %v", err)
		}

		allFiles += files
		allDirs += dirs
	}
	b.Logf("DTWalkDir:")
	b.Logf("- Files: %d", allFiles)
	b.Logf("- Dirs: %d", allDirs)

}

//goland:noinspection GoDirectComparisonOfErrors
func isPermErr(err error) (isErr bool) {
	var perr *fs.PathError
	switch {
	case err == nil:
		goto end
	case !errors.As(err, &perr):
		goto end
	case perr.Err == syscall.EPERM:
		isErr = true
		goto end
	case perr.Err == syscall.EACCES:
		isErr = true
		goto end
	}
end:
	return isErr
}

// filepathWalk is a helper that uses filepath.WalkDir, skips ".git"
// directories, and returns counts of files and dirs.
//
//goland:noinspection GoDirectComparisonOfErrors
func filepathWalk(root dt.DirPath) (files, dirs int, err error) {
	var first bool
	files = 0
	dirs = 0

	err = filepath.WalkDir(string(root), func(path string, d fs.DirEntry, err error) error {
		switch {
		case !first:
			// Don't count the root entry
			first = d.IsDir()
		case isPermErr(err):
			// Ignore permission errors
			err = nil
		case err != nil:
			// Return an error other than permission
		case d.IsDir():
			if d.Name() == ".git" {
				err = fs.SkipDir
				break
			}
			if path == string(root) {
				// don’t count the root itself
				break
			}
			dirs++
		case d.Type().IsRegular():
			files++
		}
		return err
	})
	return files, dirs, err
}

// dirPathWalk is a helper that uses dt.DirPath.Walk, skips ".git"
// directories, and returns counts of files and dirs.
func dirPathWalk(root dt.DirPath) (files, dirs int, err error) {
	files = 0
	dirs = 0
	var de dt.DirEntry
	for de, err = range root.Walk() {
		switch {
		case isPermErr(err):
			// Ignore permission errors
			err = nil
		case err != nil:
			// Return an error other than permission
		case de.Entry == nil:
			// Skip root
			continue
		case de.IsDir():
			if de.Entry.Name() == ".git" {
				de.SkipDir()
				break
			}
			dirs++
		case de.IsFile():
			files++
		}
	}
	return files, dirs, err
}

// dtWalkDir is a helper that uses dt.WalkDir, skips ".git"
// directories, and returns counts of files and dirs.
func dtWalkDir(root dt.DirPath) (files, dirs int, err error) {
	var first bool
	files = 0
	dirs = 0
	var de dt.DirEntry
	for de, err = range dt.WalkDir(root) {
		switch {
		case !first:
			// Don't count the root entry
			first = de.IsDir()
		case isPermErr(err):
			// Ignore permission errors
			err = nil
		case err != nil:
			// Return an error other than permission
		case de.IsDir():
			if de.Entry.Name() == ".git" {
				de.SkipDir()
				break
			}
			dirs++
		case de.IsFile():
			files++
		}
	}
	return files, dirs, err
}

func benchmarkDir(b *testing.B) dt.DirPath {
	home, err := dt.UserHomeDir()
	if err != nil {
		b.Fatalf("os.UserHomeDir: %v", err)
	}
	return dt.DirPathJoin(home, "Projects")
}
