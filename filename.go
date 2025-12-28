package dt

import (
	"io/fs"
	"os"
	"path/filepath"
)

type Filename string

func (fn Filename) Ext() FileExt {
	return FileExt(filepath.Ext(string(fn)))
}

func (fn Filename) ReadFile(fileSys ...fs.FS) ([]byte, error) {
	if len(fileSys) == 0 {
		return os.ReadFile(string(fn))
	}
	return fs.ReadFile(fileSys[0], string(fn))
}

func (fn Filename) OpenFile(flag int, mode os.FileMode) (*os.File, error) {
	return os.OpenFile(string(fn), flag, mode)
}
