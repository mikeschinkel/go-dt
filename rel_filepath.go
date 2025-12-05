package dt

import (
	"io/fs"
	"os"
	"path/filepath"
)

// RelFilepath is an relativate filepath with filename including extension if applicable
type RelFilepath string

func (fp RelFilepath) Dir() DirPath {
	return DirPath(filepath.Dir(string(fp)))
}

func (fp RelFilepath) Base() Filename {
	return Filename(filepath.Base(string(fp)))
}

func (fp RelFilepath) ValidPath() bool {
	return fs.ValidPath(string(fp))
}

func (fp RelFilepath) Stat(fileSys ...fs.FS) (os.FileInfo, error) {
	return EntryPath(fp).Stat(fileSys...)
}

func (fp RelFilepath) ReadFile(fileSys ...fs.FS) ([]byte, error) {
	if len(fileSys) == 0 {
		return os.ReadFile(string(fp))
	}
	return fs.ReadFile(fileSys[0], string(fp))
}

func (fp RelFilepath) Rel(baseDir DirPath) (RelFilepath, error) {
	ps, err := filepath.Rel(string(baseDir), string(fp))
	return RelFilepath(ps), err
}

func (fp RelFilepath) Exists() (exists bool, err error) {
	var status EntryStatus
	status, err = fp.Status()
	if err != nil {
		goto end
	}
	exists = status == IsFileEntry
end:
	return exists, err
}

func (fp RelFilepath) WriteFile(data []byte, mode os.FileMode) error {
	return os.WriteFile(string(fp), data, mode)
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
func (fp RelFilepath) Status(flags ...EntryStatusFlags) (status EntryStatus, err error) {
	return EntryPath(fp).Status(flags...)
}
