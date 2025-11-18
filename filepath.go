package dt

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/mikeschinkel/go-dt/de"
)

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

func (fp Filepath) Join(elems ...any) Filepath {
	return Filepath(EntryPath(fp).Join(elems...))
}

//

// CopyToDir copies the file to the destination directory path with optional
// permission control
func (fp Filepath) CopyToDir(dest DirPath, opts *CopyOptions) (err error) {
	var status EntryStatus

	dir := dest
	status, err = fp.Status()
	if err != nil {
		goto end
	}
	switch status {
	case IsFileEntry:
		err = NewErr(dt.ErrIsADirectory)
	case IsDirEntry:
		destFP := FilepathJoin(dir, fp.Base())
		err = fp.CopyTo(destFP, opts)
	default:
		err = NewErr(
			dt.ErrNotADirectory,
			"entry_status", status,
		)
	}
end:
	return err
}

// CopyTo copies the file to the destination filepath with optional permission
// control
func (fp Filepath) CopyTo(dest Filepath, opts *CopyOptions) (err error) {
	var srcFile *os.File
	var destFile *os.File
	var srcInfo os.FileInfo
	var destMode os.FileMode
	var destExists bool

	// Normalize opts
	if opts == nil {
		opts = new(CopyOptions)
	}

	// Read source file info
	srcInfo, err = fp.Stat()
	if err != nil {
		goto end
	}

	if srcInfo.IsDir() {
		err = NewErr(dt.ErrIsADirectory)
		goto end
	}

	// Check if destination exists
	_, err = dest.Stat()
	destExists = err == nil

	// If dest exists and Force is false, error
	if destExists && !opts.Force {
		err = os.ErrExist
		goto end
	}

	// Determine destination permissions
	if opts.DestModeFunc != nil {
		destMode = opts.DestModeFunc(EntryPath(dest))
		if destMode == 0 {
			// 0 means preserve source permissions
			destMode = srcInfo.Mode()
		}
	} else {
		// No callback, preserve source permissions
		destMode = srcInfo.Mode()
	}

	// Create parent directory if needed
	err = dest.Dir().MkdirAll(0755)
	if err != nil {
		goto end
	}

	// Open source file
	srcFile, err = fp.Open()
	if err != nil {
		goto end
	}
	defer CloseOrLog(srcFile)

	// Create destination file
	destFile, err = dest.OpenFile(os.O_WRONLY|os.O_CREATE|os.O_TRUNC, destMode)
	if err != nil {
		goto end
	}
	defer CloseOrLog(destFile)

	// Copy contents
	_, err = srcFile.WriteTo(destFile)
	if err != nil {
		goto end
	}

end:
	return err
}
