package dtx

import (
	"runtime"

	"github.com/mikeschinkel/go-dt"
)

func GetWorkingDir() (wd dt.DirPath, err error) {
	wd, err = dt.Getwd()
	if err != nil {
		var hint string

		switch runtime.GOOS {
		case "windows":
			hint = "it may have been deleted, moved, or become inaccessible; " +
				"try running `cd` and then `echo %CD%` (in cmd.exe) or `Get-Location` (in PowerShell) " +
				"in the same terminal session"
		case "darwin", "linux":
			hint = "it may have been deleted, moved, or become inaccessible; " +
				"try running `pwd` and `ls -ld .` in your shell"
		default:
			hint = "it may have been deleted, moved, or become inaccessible; " +
				"check your shell's equivalent of `pwd` and directory listing commands"
		}

		err = dt.NewErr(
			dt.ErrCannotDetermineWorkingDirectory,
			"hint", hint,
			err,
		)
	}

	return wd, err
}
