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

func NewDirEntry(root DirPath, skipDir *bool) DirEntry {
	return DirEntry{
		Root:    root,
		skipDir: skipDir,
	}
}

// SkipDir causes the current directory (if this DirEntry represents a
// directory during a Walk) to be skipped, analogous to returning fs.SkipDir
// from a WalkDir callback.
func (de DirEntry) SkipDir() {
	if de.skipDir == nil {
		goto end
	}
	*de.skipDir = true
end:
	return
}

// IsDir reports whether this DirEntry represents a directory.
func (de DirEntry) IsDir() (isDir bool) {
	if de.Entry == nil {
		goto end
	}
	isDir = de.Entry.IsDir()
end:
	return isDir
}

// IsFile reports whether this DirEntry represents a regular file (non-dir).
func (de DirEntry) IsFile() (isFile bool) {
	if de.Entry == nil {
		goto end
	}
	if de.Entry.IsDir() {
		goto end
	}
	isFile = de.Entry.Type().IsRegular()
end:
	return isFile
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

// DirPath returns the entry as a directory path. It is intended
// for directory entries and will panic if called on a non-directory.
func (de DirEntry) DirPath() DirPath {
	if !de.IsDir() {
		panic("dt.DirEntry.DirPath called on non-directory entry")
	}
	return DirPathJoin(de.Root, de.Rel)
}

// Filename returns the last path element as a Filename. It is intended for
// file entries and will panic if called on a non-file.
func (de DirEntry) Filename() Filename {
	if !de.IsFile() {
		panic("dt.DirEntry.Filename called on non-file entry")
	}
	return Filename(de.Base())
}
