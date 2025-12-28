package dt

import (
	"os"
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

func (fp Filepath) Join(elems ...any) Filepath {
	return Filepath(EntryPath(fp).Join(elems...))
}

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
		err = NewErr(ErrIsADirectory)
	case IsDirEntry:
		destFP := FilepathJoin(dir, fp.Base())
		err = fp.CopyTo(destFP, opts)
	default:
		err = NewErr(
			ErrNotDirectory,
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
		err = NewErr(ErrIsADirectory)
		goto end
	}

	// Check if destination exists
	_, err = dest.Stat()
	destExists = err == nil

	// If dest exists and Overwrite is false, error
	if destExists && !opts.Overwrite {
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

func (fp Filepath) HasDotDotPrefix() bool {
	return EntryPath(fp).HasDotDotPrefix()
}

func (fp Filepath) Expand() (_ Filepath, err error) {
	var ep EntryPath
	ep, err = EntryPath(fp).Expand()
	return Filepath(ep), err
}

func (fp Filepath) ToTilde(opt TildeOption) TildeFilepath {
	return ToTilde[Filepath, TildeFilepath](fp, opt)
}

func (fp Filepath) TrimTilde() (tdp PathSegments) {
	return TrimTilde[Filepath](fp)
}

func (fp Filepath) ErrKV() ErrKV {
	return kv{k: "filepath", v: fp.ToTilde(OrFullPath)}
}

func (fp Filepath) Touch(mode os.FileMode) (err error) {
	err = fp.WriteFile([]byte{}, mode)
	if err != nil {
		err = NewErr(ErrFailedtoCreateFile, fp.ErrKV(), err)
	}
	return err
}
