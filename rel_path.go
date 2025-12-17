package dt

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// RelPath can be a RelFilepath or a PathSegments
type RelPath string

func (rp RelPath) Dir() DirPath {
	return DirPath(filepath.Dir(string(rp)))
}

func (rp RelPath) Base() PathSegment {
	return PathSegment(filepath.Base(string(rp)))
}

func (rp RelPath) Stat(fileSys ...fs.FS) (fs.FileInfo, error) {
	if len(fileSys) == 0 {
		return os.Stat(string(rp))
	}
	return fs.Stat(fileSys[0], string(rp))
}

func (rp RelPath) Lstat() (os.FileInfo, error) {
	return os.Lstat(string(rp))
}

// Status classifies the filesystem entry referred to by fp.
//
// It returns IsMissingEntry when the entry does not exist (err == nil).
// It returns IsEntryError for all other filesystem errors (err != nil).
// By default it follows symlinks (like os.Stat). To inspect the entry
// itself, pass FlagDontFollowSymlinks.
//
// On platforms that don't support certain kinds (e.g., sockets/devices on
// Windows), those statuses will never be returned.
func (rp RelPath) Status(flags ...EntryStatusFlags) (status EntryStatus, err error) {
	return EntryPath(rp).Status(flags...)
}

// Readlink returns the target of the symlink referred to by fp.
//
// It returns the resolved target as a RelPath. If fp is not a symlink,
// it returns an empty RelPath and a non-nil error from os.Readlink.
// On most systems the returned target is relative to the directory
// containing fp, not an absolute path.
func (rp RelPath) Readlink() (target RelPath, err error) {
	var ep EntryPath
	ep, err = EntryPath(rp).Readlink()
	return RelPath(ep), err
}

func (rp RelPath) HasSuffix(suffix DirPath) bool {
	return strings.HasSuffix(string(rp), string(suffix))
}

// Contains checks if ep contains the given substring.
// Accepts: string, DirPath, Filepath, RelPath, PathSegment, or fmt.Stringer
// Panics on unsupported types.
func (rp RelPath) Contains(substr any) bool {
	return EntryPath(rp).Contains(substr)
}

func (rp RelPath) VolumeName() VolumeName {
	return VolumeName(filepath.VolumeName(string(rp)))
}

func (rp RelPath) Abs() (RelPath, error) {
	entry, err := filepath.Abs(string(rp))
	return RelPath(entry), err
}

func (rp RelPath) Join(elems ...any) RelPath {
	return RelPath(EntryPath(rp).Join(elems...))
}
