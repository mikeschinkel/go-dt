package dt

import (
	"os"

	"github.com/mikeschinkel/go-dt/de"
)

func CanWrite(path EntryPath) (canWrite bool, err error) {
	// Resolve to a directory: if `path` is a dir, use it;
	// otherwise use its parent directory.
	var dir DirPath
	var tmpFile *os.File
	var status EntryStatus

	status, err = path.Status()
	if err != nil {
		goto end
	}

	switch status {
	case IsDirEntry:
		dir = DirPath(path)
	case IsMissingEntry, IsFileEntry:
		// If the path does not exist, treat it as a file path
		// and check its parent directory.
		dir = path.Dir()
	default:
		// Some other error (e.g., permission issue just to stat).
		//goland:noinspection GoDfaErrorMayBeNotNil
		err = WithErr(
			de.ErrNotFileOrDirectory,
			"entry_type", status.String(),
		)
		goto end
	}

	// Now dir should be a directory that *should* exist.
	// Try to create a temporary file there.
	tmpFile, err = CreateTemp(dir, ".canwrite-*")
	if err != nil {
		// Explicit permission error → definitely cannot write there.
		//if errors.Is(err, os.ErrPermission) {
		//	goto end
		//}
		// Other errors (e.g. directory doesn’t exist, read-only fs, etc.)
		// are surfaced so the caller can decide how to interpret them.
		goto end
	}
	if tmpFile == nil {
		goto end
	}
	defer func() {
		LogOnError(os.Remove(tmpFile.Name()))
		CloseOrLog(tmpFile)
	}()
	canWrite = true
end:
	return canWrite, err
}
