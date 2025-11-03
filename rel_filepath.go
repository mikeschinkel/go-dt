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
