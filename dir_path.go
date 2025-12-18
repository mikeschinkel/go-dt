package dt

import (
	"errors"
	"fmt"
	"io/fs"
	"iter"
	"os"
	"path/filepath"
	"strings"
)

func ParseDirPath(s string) (dp DirPath, err error) {
	if len(s) == 0 {
		err = ErrEmpty
		goto end
	}

	if s[0] != '~' {
		dp = DirPath(s)
		goto end
	}

	_, err = ParseTildeDirPath(s)
	if errors.Is(err, ErrNotTildePath) {
		dp = DirPath(s)
		err = nil
		goto end
	}
	if err != nil {
		goto end
	}

	dp, err = TildeDirPath(s).Expand()

end:
	return dp, err
}

func ParseDirPaths(dirs []string) (dps []DirPath, err error) {
	var errs []error
	var dp DirPath

	dps = make([]DirPath, 0, len(dirs))
	for _, dir := range dirs {
		dp, err = ParseDirPath(dir)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		dps = append(dps, dp)
	}
	return dps, CombineErrs(errs)
}

func DirPathRead(dp DirPath) (des []DirEntry, err error) {
	var entries []os.DirEntry
	entries, err = dp.ReadDir()
	if err != nil {
		goto end
	}
	des = make([]DirEntry, len(entries))
	for i, entry := range entries {
		des[i] = NewDirEntry(dp, entry)
	}
end:
	return des, err
}

var _ fmt.Stringer = (*DirPath)(nil)

// DirPath represents an absolute or relative filesystem directory path.
//
// It provides helper methods for working with directories while preserving
// type safety and semantic clarity over raw strings.  A zero value is valid
// but does not refer to any real directory.
type DirPath string

func (dp DirPath) String() string {
	return string(dp)
}

// EnsureExists verifies that the directory exists, creating it and any missing
// parent directories as needed.
//
// If the path already exists as a directory, EnsureExists is a no-op.
// If the path exists as a file, it returns ErrPathIsFile.
// Any other filesystem error is returned as-is.
func (dp DirPath) EnsureExists() (err error) {
	var info os.FileInfo
	info, err = os.Stat(string(dp))
	if errors.Is(err, os.ErrNotExist) {
		err = os.MkdirAll(string(dp), os.ModePerm)
		if err != nil {
			goto end
		}
		info, err = os.Stat(string(dp))
	}
	if err != nil {
		goto end
	}
	if !info.IsDir() {
		err = NewErr(ErrPathIsFile, err)
	}
end:
	return err
}

func (dp DirPath) DirFS() fs.FS {
	return os.DirFS(string(dp))
}

func (dp DirPath) Dir() DirPath {
	return DirPath(filepath.Dir(string(dp)))
}

func (dp DirPath) Base() PathSegment {
	return PathSegment(filepath.Base(string(dp)))
}

func (dp DirPath) Contains(part any) bool {
	return EntryPath(dp).Contains(part)
}

func (dp DirPath) Join(elems ...any) DirPath {
	return DirPath(EntryPath(dp).Join(elems...))
}

func (dp DirPath) Clean() DirPath {
	return DirPath(filepath.Clean(string(dp)))
}

func (dp DirPath) VolumeName() VolumeName {
	return VolumeName(filepath.VolumeName(string(dp)))
}

func (dp DirPath) MkdirAll(mode os.FileMode) error {
	return os.MkdirAll(string(dp), mode)
}

func (dp DirPath) RemoveAll() error {
	return os.RemoveAll(string(dp))
}

func (dp DirPath) Chmod(mode os.FileMode) error {
	return os.Chmod(string(dp), mode)
}

func (dp DirPath) Rel(baseDir DirPath) (PathSegments, error) {
	ps, err := filepath.Rel(string(baseDir), string(dp))
	return PathSegments(ps), err
}

func (dp DirPath) ReadDir() ([]os.DirEntry, error) {
	return os.ReadDir(string(dp))
}

func (dp DirPath) Stat(fileSys ...fs.FS) (os.FileInfo, error) {
	return EntryPath(dp).Stat(fileSys...)
}

func (dp DirPath) Lstat(fileSys ...fs.FS) (os.FileInfo, error) {
	return EntryPath(dp).Lstat(fileSys...)
}

func (dp DirPath) IsAbs() bool {
	return filepath.IsAbs(string(dp))
}

func (dp DirPath) Abs() (DirPath, error) {
	dir, err := filepath.Abs(string(dp))
	return DirPath(dir), err
}

func (dp DirPath) HasPrefix(prefix DirPath) bool {
	return strings.HasPrefix(string(dp), string(prefix))
}

func (dp DirPath) HasSuffix(suffix string) bool {
	return strings.HasSuffix(string(dp), suffix)
}

func (dp DirPath) TrimPrefix(prefix DirPath) DirPath {
	return DirPath(strings.TrimPrefix(string(dp), string(prefix)))
}

func (dp DirPath) TrimSuffix(TrimSuffix string) DirPath {
	return DirPath(strings.TrimSuffix(string(dp), TrimSuffix))
}
func (dp DirPath) ToSlash() DirPath {
	return DirPath(filepath.ToSlash(string(dp)))
}

func (dp DirPath) ToLower() DirPath {
	return DirPath(strings.ToLower(string(dp)))
}

func (dp DirPath) ToUpper() DirPath {
	return DirPath(strings.ToUpper(string(dp)))
}

func (dp DirPath) EvalSymlinks() (_ DirPath, err error) {
	var path string
	path, err = filepath.EvalSymlinks(string(dp))
	return DirPath(path), err
}

// ===[Enhancements]===

func (dp DirPath) Expand() (_ DirPath, err error) {
	var ep EntryPath
	ep, err = EntryPath(dp).Expand()
	return DirPath(ep), err
}

// Normalize expands a leading "~" to the current user's home directory (when
// it uses the correct OS path separator), then returns an absolute directory
// path.
func (dp DirPath) Normalize() (DirPath, error) {
	return TildeDirPath(dp).Expand()
}

func (dp DirPath) ToTilde() (tdp TildeDirPath, err error) {
	var home DirPath
	var rel PathSegments
	home, err = UserHomeDir()
	if err != nil {
		goto end
	}
	rel, err = dp.Rel(home)
	if err != nil {
		goto end
	}
	tdp = TildeDirPath(DirPathJoin(fmt.Sprintf("~%c", os.PathSeparator), rel))
end:
	return tdp, err
}

func (dp DirPath) Status(flags ...EntryStatusFlags) (status EntryStatus, err error) {
	return EntryPath(dp).Status(flags...)
}

func (dp DirPath) EnsureTrailSep() DirPath {
	return DirPath(EntryPath(dp).EnsureTrailSep())
}

func (dp DirPath) CanWrite() (bool, error) {
	return CanWrite(EntryPath(dp))
}

func (dp DirPath) Exists() (exists bool, err error) {
	var status EntryStatus
	status, err = dp.Status()
	if err != nil {
		goto end
	}
	exists = status == IsDirEntry
end:
	return exists, err
}

func (dp DirPath) HasDotDotPrefix() bool {
	return EntryPath(dp).HasDotDotPrefix()
}

// Walk walks the filesystem rooted at d using d.DirFS() and yields all entries
// as DirEntry values together with any per-entry errors encountered.
func (dp DirPath) Walk() iter.Seq2[DirEntry, error] {
	return dp.WalkFS(dp.DirFS())
}

// WalkFS walks the provided fsys (typically obtained from d.DirFS()) starting
// at "." and yields all entries as DirEntry values together with any per-entry
// errors encountered.
//
// This implementation is non-recursive and uses fs.ReadDir internally so that
// it can correctly honor the iterator contract: once the caller stops
// iteration (yield returns false), WalkFS will not call yield again.
func (dp DirPath) WalkFS(fsys fs.FS) iter.Seq2[DirEntry, error] {
	return func(yield func(DirEntry, error) bool) {
		// dirState tracks iteration state for a single directory.
		type dirState struct {
			dir     string
			entries []fs.DirEntry
			i       int
		}

		var skipDir bool

		// Start from the logical root "." inside fsys.
		stack := []dirState{{dir: "."}}

		for len(stack) > 0 {
			// Work on the directory at the top of the stack.
			s := &stack[len(stack)-1]

			// If we haven't read this directory yet, do so now.
			if s.entries == nil {
				ents, err := fs.ReadDir(fsys, s.dir)
				if err != nil {
					entry := NewDirEntryWithSkipDir(dp, RelPath(s.dir), &skipDir)
					skipDir = false
					if !yield(entry, err) {
						return
					}

					// On error, we cannot descend into this directory; pop and continue.
					stack = stack[:len(stack)-1]
					continue
				}

				s.entries = ents
				s.i = 0
			}

			// If we've exhausted this directory, pop it and continue with its parent.
			if s.i >= len(s.entries) {
				stack = stack[:len(stack)-1]
				continue
			}

			// Consume the next entry in this directory.
			de := s.entries[s.i]
			s.i++

			// Construct the relative path for this entry.
			var rel string
			if s.dir == "." {
				rel = de.Name()
			} else {
				rel = filepath.Join(s.dir, de.Name())
			}

			entry := NewDirEntryWithSkipDir(dp, RelPath(rel), &skipDir)
			entry.Entry = de

			skipDir = false

			if !yield(entry, nil) {
				goto end
			}

			if !de.IsDir() {
				continue
			}

			if skipDir {
				continue
			}

			// If this is a directory and the caller did not request SkipDir,
			// push it onto the stack to walk its children.
			stack = append(stack, dirState{dir: rel})
		}
	end:
		return
	}
}

// WalkFiles walks using d.DirFS() and yields only entries that represent
// regular files.
func (dp DirPath) WalkFiles() iter.Seq2[DirEntry, error] {
	return dp.WalkFilesFS(dp.DirFS())
}

// WalkFilesFS walks the provided fsys and yields only entries that represent
// regular files.
func (dp DirPath) WalkFilesFS(fsys fs.FS) iter.Seq2[DirEntry, error] {
	return func(yield func(DirEntry, error) bool) {
		for de, err := range dp.WalkFS(fsys) {
			if err != nil {
				if !yield(de, err) {
					return
				}
				continue
			}

			if de.Entry == nil {
				continue
			}

			if !de.IsFile() {
				continue
			}

			if !yield(de, nil) {
				return
			}
		}
	}
}

// WalkDirs walks using d.DirFS() and yields only entries that represent
// directories.
func (dp DirPath) WalkDirs() iter.Seq2[DirEntry, error] {
	return dp.WalkDirsFS(dp.DirFS())
}

// WalkDirsFS walks the provided fsys and yields only entries that represent
// directories.
func (dp DirPath) WalkDirsFS(fsys fs.FS) iter.Seq2[DirEntry, error] {
	return func(yield func(DirEntry, error) bool) {
		for de, err := range dp.WalkFS(fsys) {
			if err != nil {
				if !yield(de, err) {
					return
				}
				continue
			}

			if de.Entry == nil {
				continue
			}

			if !de.IsDir() {
				continue
			}

			if !yield(de, nil) {
				return
			}
		}
	}
}
