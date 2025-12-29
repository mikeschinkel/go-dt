package dt

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func ParseDirPath(s string) (dp DirPath, err error) {
	ep, err := ParseEntryPath(s)
	return DirPath(ep), err
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

// IMPORTANT: Currently I explicitly do not want String() methods and want to see if casting to string will meet all needs.
//var _ fmt.Stringer = (*DirPath)(nil)

// DirPath represents an absolute or relative filesystem directory path.
//
// It provides helper methods for working with directories while preserving
// type safety and semantic clarity over raw strings.  A zero value is valid
// but does not refer to any real directory.
type DirPath string

// IMPORTANT: Currently I explicitly do not want String() methods and want to see if casting to string will meet all needs.
//func (dp DirPath) String() string {
//	return string(dp)
//}

// EnsureExists verifies that the directory exists, creating it and any missing
// parent directories as needed.
//
// If the path already exists as a directory, EnsureExists is a no-op.
// If the path exists as a file, it returns ErrPathIsFile.
// Any other filesystem error is returned as-is.
func (dp DirPath) EnsureExists(mod os.FileMode) (err error) {
	var info os.FileInfo
	info, err = os.Stat(string(dp))
	if errors.Is(err, os.ErrNotExist) {
		err = os.MkdirAll(string(dp), mod)
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
	if err != nil {
		err = WithErr(err, ErrFailedToEnsureDir, dp.ErrKV())
	}
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

func (dp DirPath) Remove() error {
	return os.Remove(string(dp))
}
