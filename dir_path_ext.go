package dt

import (
	"io/fs"
	"iter"
	"os"
	"path/filepath"
)

func (dp DirPath) Expand() (_ DirPath, err error) {
	var ep EntryPath
	ep, err = EntryPath(dp).Expand()
	return DirPath(ep), err
}

// Normalize expands a leading "~" to the current user's home directory (when
// it uses the correct OS path separator), then returns an absolute directory
// path.
func (dp DirPath) Normalize() (DirPath, error) {
	return dp.Expand()
}

func (dp DirPath) ToTilde(opt TildeOption) (tdp TildeDirPath) {
	return ToTilde[DirPath, TildeDirPath](dp, opt)
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

func (dp DirPath) IsTidlePath() bool {
	_, err := ParseTildeDirPath(string(dp))
	return err == nil
}

func (dp DirPath) ErrKV() ErrKV {
	return kv{k: "dir_path", v: dp.ToTilde(OrFullPath)}
}

func (dp DirPath) TrimTilde() (tdp PathSegments) {
	return TrimTilde[DirPath](dp)
}

func (dp DirPath) MkSubdirs(subdirs []PathSegments, mode os.FileMode) (err error) {
	var errs []error
	for _, dir := range subdirs {
		errs = AppendErr(errs, DirPathJoin(dp, dir).MkdirAll(mode))
	}
	return CombineErrs(errs)
}

func (dp DirPath) TouchFiles(files []RelFilepath, mode os.FileMode) (err error) {
	var errs []error
	for _, file := range files {
		errs = AppendErr(errs, FilepathJoin(dp, file).Touch(mode))
	}
	return CombineErrs(errs)
}
