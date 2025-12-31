package dt

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// EntryStatusFlags controls optional classification behavior.
// The zero value is safe and means "follow symlinks" (os.Stat).
type EntryStatusFlags uint32

const (
	// DontFollowSymlinks causes Status to inspect the entry itself
	// (os.Lstat) instead of following symlinks.
	DontFollowSymlinks EntryStatusFlags = 1 << iota

	// (Reserved for future flags)
	// TreatBrokenSymlinkAsMissing
	// ClassifyBlockVsCharDevice
	// ...
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
func (ep EntryPath) Status(flags ...EntryStatusFlags) (status EntryStatus, err error) {
	var info os.FileInfo
	var mode os.FileMode

	switch {
	case len(flags) == 0:
		info, err = os.Stat(string(ep))
	case flags[0]&DontFollowSymlinks != 0:
		info, err = os.Lstat(string(ep))
	default:
		info, err = os.Stat(string(ep))
	}

	if errors.Is(err, fs.ErrNotExist) {
		status = IsMissingEntry
		err = nil
		goto end
	}
	if err != nil {
		status = IsEntryError
		goto end
	}

	mode = info.Mode()
	switch {
	case mode.IsRegular():
		status = IsFileEntry
	case mode.IsDir():
		status = IsDirEntry
	case mode&fs.ModeSymlink != 0:
		status = IsSymlinkEntry
	case mode&fs.ModeSocket != 0:
		status = IsSocketEntry
	case mode&fs.ModeNamedPipe != 0:
		status = IsPipeEntry
	case mode&fs.ModeDevice != 0:
		status = IsDeviceEntry
	default:
		status = IsUnclassifiedEntryStatus
	}

end:
	return status, err
}

func (ep EntryPath) Join(elems ...any) EntryPath {
	ss := make([]string, 0, len(elems)+1)
	ss = append(ss, string(ep))
	for _, part := range elems {
		switch s := part.(type) {
		case string:
			ss = append(ss, s)
		case EntryPath:
			ss = append(ss, string(s))
		case DirPath:
			ss = append(ss, string(s))
		case Filepath:
			ss = append(ss, string(s))
		case RelFilepath:
			ss = append(ss, string(s))
		case PathSegment:
			ss = append(ss, string(s))
		case PathSegments:
			ss = append(ss, string(s))
		default:
			ss = append(ss, fmt.Sprintf("%s", part))
		}
	}
	return EntryPath(filepath.Join(ss...))
}

// EnsureTrailSep returns ep with exactly one trailing path separator
// when appropriate for the platform. It does not modify an empty string.
func (ep EntryPath) EnsureTrailSep() EntryPath {
	if ep == "" {
		goto end
	}

	// Already has a trailing native separator?
	if ep[len(ep)-1] == os.PathSeparator {
		goto end
	}

	// On Windows, also consider '/' as a valid existing trailing separator.
	if os.PathSeparator == '\\' && ep[len(ep)-1] == '/' {
		//goland:noinspection GoAssignmentToReceiver
		ep = ep[:len(ep)-1] + `\`
		goto end
	}

	ep += EntryPath(os.PathSeparator)
end:
	return ep
}

// HasDotDotPrefix reports whether p, interpreted as a relative path,
// starts with a ".." segment (e.g. ".." or "../foo" or "..\foo").
// It does NOT treat names like "..foo" as having a dot-dot prefix.
func (ep EntryPath) HasDotDotPrefix() bool {
	return ep == ".." || strings.HasPrefix(string(ep), ".."+string(os.PathSeparator))
}

func (ep EntryPath) Expand() (out EntryPath, err error) {
	var home DirPath
	s := string(ep)

	switch {
	case len(s) == 0:
		err = ErrEmpty
		goto end

	case s == ".":
		var dp DirPath
		// We are "current directory
		dp, err = Getwd()
		if err != nil {
			goto end
		}
		out = EntryPath(dp)
		goto end

	case s == "~":
		home, err = UserHomeDir()
		if err != nil {
			goto end
		}
		out = EntryPath(home)
		goto end

	case s[0] == '/':
		// We are an absolute path already, works on Windows, macOS and Linux.
		if runtime.GOOS == "windows" {
			s = filepath.FromSlash(s)
		}
		out, err = EntryPath(s).Clean().Abs()
		goto end

	case s[0] == '\\' && runtime.GOOS == "windows":
		// We are an absolute path on Windows already
		s = filepath.FromSlash(s)
		out, err = EntryPath(s).Clean().Abs()
		goto end

	case len(s) == 1:
		out, err = EntryPath(s).Clean().Abs()
		goto end

	case s[:2] == "..":
		// We are a parent path, just return it
		out, err = EntryPath(filepath.Dir(s)).Clean().Abs()
		goto end

	case s[:2] == "~/":
		// We start with ~/ so we are a tilde path; works on Windows or Linux/macOS
		if runtime.GOOS == "windows" {
			s = filepath.FromSlash(s)
		}
		// Go on to be handled by the tilde expansion

	case s[:2] == "~\\" && runtime.GOOS == "windows":
		// Go on to be handled by the tilde expansion

	default:
		// Not a special case, just a relative path
		out, err = EntryPath(s).Clean().Abs()
		goto end

	}

	home, err = UserHomeDir()
	if err != nil {
		goto end
	}

	if len(s) == 2 {
		out = EntryPath(home)
		goto end
	}

	out = EntryPathJoin(home, s[2:]).Clean()

end:
	if err != nil {
		err = WithErr(err, ErrFailedToExpandPath, ep.ErrKV())
	}
	return out, err
}

func (ep EntryPath) Exists() (exists bool, err error) {
	var status EntryStatus
	status, err = ep.Status()
	if err != nil {
		goto end
	}
	exists = status == IsDirEntry || status == IsFileEntry
end:
	return exists, err
}

func (ep EntryPath) ToTilde(opt TildeOption) (tep TildeEntryPath) {
	return ToTilde[EntryPath, TildeEntryPath](ep, opt)
}

func (ep EntryPath) TrimTilde() (tdp PathSegments) {
	return TrimTilde[EntryPath](ep)
}

func (ep EntryPath) IsFile() (isFile bool) {
	status, err := ep.Status()
	if err != nil {
		goto end
	}
	if status != IsFileEntry {
		goto end
	}
	isFile = true
end:
	return isFile
}

func (ep EntryPath) IsDir() (isDir bool) {
	status, err := ep.Status()
	if err != nil {
		goto end
	}
	if status != IsDirEntry {
		goto end
	}
	isDir = true
end:
	return isDir
}

func (ep EntryPath) ErrKV() ErrKV {
	return kv{k: "path", v: ep.ToTilde(OrFullPath)}
}

func (ep EntryPath) EnsureFilepath(defaultName Filename) (fp Filepath, err error) {
	var exists bool
	fp = Filepath(ep)
	if ep.IsDir() {
		fp = Filepath(ep.Join(defaultName))
	}
	exists, err = fp.Exists()
	if !exists {
		err = NewErr(ErrFileNotExists, fp.ErrKV(), err)
		goto end
	}
end:
	return fp, err
}

func EnsureFilepath(path string, defaultName Filename) (fp Filepath, err error) {
	var ep EntryPath
	ep, err = ParseEntryPath(path)
	if err != nil {
		err = NewErr(ErrInvalidFilepath, err)
		goto end
	}
	fp, err = ep.EnsureFilepath(defaultName)
end:
	return fp, err
}
