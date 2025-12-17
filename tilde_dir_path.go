package dt

import (
	"path/filepath"
	"runtime"
	"strings"
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
	var home DirPath
	var remainder string
	s := string(tdp)

	if len(s) == 0 {
		err = ErrEmpty
		goto end
	}
	if s[0] != '~' {
		dp, err = DirPath(s).Abs()
		goto end
	}

	if len(s) == 1 {
		dp, err = UserHomeDir()
		goto end
	}

	if runtime.GOOS == "windows" {
		if s[1] != '/' && s[1] != '\\' {
			dp, err = DirPath(s).Abs()
			goto end
		}
		s = filepath.FromSlash(s)
	} else {
		if s[1] != '/' {
			dp, err = DirPath(s).Abs()
			goto end
		}
	}

	home, err = UserHomeDir()
	if err != nil {
		goto end
	}

	remainder = ""
	if len(s) > 2 {
		remainder = strings.TrimLeft(s[2:], "/\\")
	}
	if len(remainder) == 0 {
		dp = home
		goto end
	}

	dp = DirPathJoin(home, remainder).Clean()

end:
	return dp, err
}
