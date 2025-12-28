package dt

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// EntryPath can be a Filepath or a DirPath
type EntryPath string

func (ep EntryPath) Dir() DirPath {
	return DirPath(filepath.Dir(string(ep)))
}

func (ep EntryPath) Clean() EntryPath {
	return EntryPath(filepath.Clean(string(ep)))
}

func (ep EntryPath) Base() PathSegment {
	return PathSegment(filepath.Base(string(ep)))
}

func (ep EntryPath) Stat(fileSys ...fs.FS) (fs.FileInfo, error) {
	if len(fileSys) == 0 {
		return os.Stat(string(ep))
	}
	return fs.Stat(fileSys[0], string(ep))
}

func (ep EntryPath) Lstat(fileSys ...fs.FS) (os.FileInfo, error) {
	if len(fileSys) == 0 {
		return os.Lstat(string(ep))
	}
	return fs.Lstat(fileSys[0], string(ep))
}

// Readlink returns the target of the symlink referred to by fp.
//
// It returns the resolved target as a EntryPath. If fp is not a symlink,
// it returns an empty EntryPath and a non-nil error from os.Readlink.
// On most systems the returned target is relative to the directory
// containing fp, not an absolute path.
func (ep EntryPath) Readlink() (target EntryPath, err error) {
	var linkTarget string
	linkTarget, err = os.Readlink(string(ep))
	if err != nil {
		goto end
	}
	target = EntryPath(linkTarget)
end:
	return target, err
}

func (ep EntryPath) HasSuffix(suffix DirPath) bool {
	return strings.HasSuffix(string(ep), string(suffix))
}

// Contains checks if ep contains the given substring.
// Accepts: string, DirPath, Filepath, EntryPath, PathSegment, or fmt.Stringer
// Panics on unsupported types.
func (ep EntryPath) Contains(substr any) bool {
	var s string

	switch v := substr.(type) {
	case string:
		s = v
	case DirPath:
		s = string(v)
	case Filepath:
		s = string(v)
	case EntryPath:
		s = string(v)
	case PathSegment:
		s = string(v)
	case interface{ String() string }:
		s = v.String()
	default:
		panic("EntryPath.Contains: unsupported type")
	}

	return strings.Contains(string(ep), s)
}

func (ep EntryPath) VolumeName() VolumeName {
	return VolumeName(filepath.VolumeName(string(ep)))
}

func (ep EntryPath) Abs() (EntryPath, error) {
	entry, err := filepath.Abs(string(ep))
	return EntryPath(entry), err
}

func (ep EntryPath) IsFile() bool {
	return filepath.IsAbs(string(ep))
}

func (ep EntryPath) IsAbs() bool {
	return filepath.IsAbs(string(ep))
}

func (ep EntryPath) EvalSymlinks() (_ EntryPath, err error) {
	var s string
	s, err = filepath.EvalSymlinks(string(ep))
	return EntryPath(s), err
}

func (ep EntryPath) Rel(baseDir EntryPath) (PathSegments, error) {
	ps, err := filepath.Rel(string(baseDir), string(ep))
	return PathSegments(ps), err
}
