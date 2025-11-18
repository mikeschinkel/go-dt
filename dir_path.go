package dt

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

func ParseDirPath(s string) (dp DirPath, err error) {
	// TODO Add some validation here
	dp = DirPath(s)
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

func (dp DirPath) IsAbs() bool {
	return filepath.IsAbs(string(dp))
}

func (dp DirPath) Abs() (DirPath, error) {
	dir, err := filepath.Abs(string(dp))
	return DirPath(dir), err
}

// ===[Enhancements]===

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
