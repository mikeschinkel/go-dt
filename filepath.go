package dt

import (
	"io/fs"
	"os"
	"path/filepath"
)

const MissingFile = IsMissingEntry

// Filepath is an absolute or relativate filepath with filename including extension if applicable
type Filepath string

func (Filepath) StringLike() {}

func (fp Filepath) Dir() DirPath {
	return DirPath(filepath.Dir(string(fp)))
}

func (fp Filepath) Base() Filename {
	return Filename(filepath.Base(string(fp)))
}

func (fp Filepath) Ext() FileExt {
	return FileExt(filepath.Ext(string(fp)))
}

func (fp Filepath) Stat(fileSys ...fs.FS) (os.FileInfo, error) {
	return EntryPath(fp).Stat(fileSys...)
}

func (fp Filepath) Lstat() (os.FileInfo, error) {
	return os.Lstat(string(fp))
}

func (fp Filepath) ReadFile(fileSys ...fs.FS) ([]byte, error) {
	if len(fileSys) == 0 {
		return os.ReadFile(string(fp))
	}
	return fs.ReadFile(fileSys[0], string(fp))
}

func (fp Filepath) Rel(baseDir DirPath) (PathSegments, error) {
	ps, err := filepath.Rel(string(baseDir), string(fp))
	return PathSegments(ps), err
}

func (fp Filepath) Remove() error {
	return os.Remove(string(fp))
}

func (fp Filepath) ValidPath() bool {
	return fs.ValidPath(string(fp))
}

// EntryStatusFlags controls optional classification behavior.
// The zero value is safe and means "follow symlinks" (os.Stat).
type EntryStatusFlags uint32

const (
	// DontFollowSymlinks causes Status to inspect the entry itself
	// (os.Lstat) instead of following symlinks.
	DontFollowSymlinks EntryStatusFlags = 1 << iota

	// (Reserved for future flags)
	// TreatBrokenSymlinkAsMissing
	// ClassifyBlockVsCharDevice
	// ...
)

// Status classifies the filesystem entry referred to by fp.
//
// It returns IsMissingEntry when the entry does not exist (err == nil).
// It returns IsEntryError for all other filesystem errors (err != nil).
// By default it follows symlinks (like os.Stat). To inspect the entry
// itself, pass FlagDontFollowSymlinks.
//
// On platforms that don't support certain kinds (e.g., sockets/devices on
// Windows), those statuses will never be returned.
func (fp Filepath) Status(flags ...EntryStatusFlags) (status EntryStatus, err error) {
	return EntryPath(fp).Status(flags...)
}

// Readlink returns the target of the symlink referred to by fp.
//
// It returns the resolved target as a Filepath. If fp is not a symlink,
// it returns an empty Filepath and a non-nil error from os.Readlink.
// On most systems the returned target is relative to the directory
// containing fp, not an absolute path.
func (fp Filepath) Readlink() (target Filepath, err error) {
	var ep EntryPath
	ep, err = EntryPath(fp).Readlink()
	return Filepath(ep), err
}

//

func ParseFilepath(s string) (fp Filepath, err error) {
	// TODO Add some validation here
	fp = Filepath(s)
	return fp, err
}
