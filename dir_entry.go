package dt

import (
	"io/fs"
	"path/filepath"
)

// DirEntry represents a filesystem entry discovered while walking a DirPath.
// It wraps fs.DirEntry and uses EntryPath for the entry's path relative to the
// walked root.
type DirEntry struct {
	// Root is the logical root DirPath that was walked.
	// It is exactly what the caller passed to Walk / WalkFS and may be
	// absolute or relative.
	Root DirPath

	// Rel is the path of this entry relative to the fs.FS root / Root.
	// It is always relative.
	Rel EntryPath

	// Entry is the underlying fs.DirEntry. It may be nil for entries where
	// fs.WalkDir reported an error before obtaining a DirEntry.
	Entry fs.DirEntry

	skipDir *bool
}

// SkipDir causes the current directory (if this DirEntry represents a
// directory during a Walk) to be skipped, analogous to returning fs.SkipDir
// from a WalkDir callback.
func (de DirEntry) SkipDir() {
	if de.skipDir == nil {
		return
	}
	*de.skipDir = true
}

// IsDir reports whether this DirEntry represents a directory.
func (de DirEntry) IsDir() bool {
	if de.Entry == nil {
		return false
	}
	return de.Entry.IsDir()
}

// IsFile reports whether this DirEntry represents a regular file (non-dir).
func (de DirEntry) IsFile() bool {
	if de.Entry == nil {
		return false
	}
	if de.Entry.IsDir() {
		return false
	}
	mode := de.Entry.Type()
	return mode.IsRegular()
}

// Base returns the last path element of Rel as an EntryPath, similar to
// filepath.Base on the relative string. It does not panic and is valid for
// both files and directories.
func (de DirEntry) Base() EntryPath {
	return EntryPath(filepath.Base(string(de.Rel)))
}

// PathSegment returns the last path element as a PathSegment. It is intended
// for directory entries and will panic if called on a non-directory.
func (de DirEntry) PathSegment() PathSegment {
	if !de.IsDir() {
		panic("dt.DirEntry.PathSegment called on non-directory entry")
	}
	return PathSegment(de.Base())
}

// Filename returns the last path element as a Filename. It is intended for
// file entries and will panic if called on a non-file.
func (de DirEntry) Filename() Filename {
	if !de.IsFile() {
		panic("dt.DirEntry.Filename called on non-file entry")
	}
	return Filename(de.Base())
}
