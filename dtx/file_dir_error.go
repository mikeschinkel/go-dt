package dtx

import (
	"errors"
	"io/fs"
	"syscall"
)

func NoFileOrDirErr(err error) (is bool) {
	var pathError *fs.PathError
	var errNo syscall.Errno
	if err == nil {
		goto end
	}
	if !errors.As(err, &pathError) {
		goto end
	}
	if pathError.Op != "open" {
		goto end
	}
	if !errors.As(pathError.Err, &errNo) {
		goto end
	}
	//goland:noinspection GoDirectComparisonOfErrors
	if errNo != syscall.ENOENT {
		goto end
	}
	is = true
end:
	return is
}
