package dt

import (
	"path/filepath"
	"runtime"
)

type TildeDirPath string

func ParseTildeDirPath(s string) (tdp TildeDirPath, err error) {

	if len(s) == 0 {
		err = ErrEmpty
		goto end
	}

	if s[0] != '~' {
		err = ErrNotTildePath
		goto end
	}
	if len(s) == 1 {
		tdp = TildeDirPath(s)
		goto end
	}

	if runtime.GOOS == "windows" {
		if s[1] != '/' && s[1] != '\\' {
			err = ErrNotTildePath
			goto end
		}
		s = filepath.FromSlash(s)
	} else {
		if s[1] != '/' {
			err = ErrNotTildePath
			goto end
		}
	}

	// TODO: Add more validation here

	tdp = TildeDirPath(s)

end:
	return tdp, err
}

func (tdp TildeDirPath) Expand() (dp DirPath, err error) {
	var ep EntryPath
	ep, err = EntryPath(tdp).Expand()
	return DirPath(ep), err
}
