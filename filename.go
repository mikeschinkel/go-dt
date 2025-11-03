package dt

import (
	"path/filepath"
)

func (fn Filename) Ext() FileExt {
	return FileExt(filepath.Ext(string(fn)))
}
