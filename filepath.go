package dt

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type FilepathGetter interface {
	Filepath() Filepath
}

func ParseFilepath(s string) (fp Filepath, err error) {
	// TODO Add some validation here
	fp = Filepath(s)
	return fp, err
}

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

func (fp Filepath) Create() (*os.File, error) {
	return os.Create(string(fp))
}

func (fp Filepath) OpenFile(flag int, mode os.FileMode) (*os.File, error) {
	return os.OpenFile(string(fp), flag, mode)
}

func (fp Filepath) ReadFile(fileSys ...fs.FS) ([]byte, error) {
	if len(fileSys) == 0 {
		return os.ReadFile(string(fp))
	}
	return fs.ReadFile(fileSys[0], string(fp))
}

func (fp Filepath) WriteFile(data []byte, mode os.FileMode) error {
	return os.WriteFile(string(fp), data, mode)
}

func (fp Filepath) Rel(baseDir DirPath) (RelFilepath, error) {
	ps, err := filepath.Rel(string(baseDir), string(fp))
	return RelFilepath(ps), err
}

func (fp Filepath) Abs() (Filepath, error) {
	ps, err := filepath.Abs(string(fp))
	return Filepath(ps), err
}

func (fp Filepath) Remove() error {
	return os.Remove(string(fp))
}

func (fp Filepath) Rename(newFile Filepath) error {
	return os.Rename(string(fp), string(newFile))
}

func (fp Filepath) ValidPath() bool {
	return fs.ValidPath(string(fp))
}

func (fp Filepath) HasPrefix(prefix DirPath) bool {
	return strings.HasPrefix(string(fp), string(prefix))
}

func (fp Filepath) HasSuffix(suffix DirPath) bool {
	return strings.HasSuffix(string(fp), string(suffix))
}

func (fp Filepath) Open() (*os.File, error) {
	return os.Open(string(fp))
}

func (fp Filepath) Exists() (exists bool, err error) {
	var status EntryStatus
	status, err = fp.Status()
	if err != nil {
		goto end
	}
	exists = status == IsFileEntry
end:
	return exists, err
}

func (fp Filepath) EvalSymlinks() (_ Filepath, err error) {
	var path string
	path, err = filepath.EvalSymlinks(string(fp))
	return Filepath(path), err
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

func (fp Filepath) IsAbs() bool {
	return filepath.IsAbs(string(fp))
}

// ===[Enhancements]===
