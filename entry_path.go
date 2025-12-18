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

// EntryPath can be a Filepath or a DirPath
type EntryPath string

func (ep EntryPath) Dir() DirPath {
	return DirPath(filepath.Dir(string(ep)))
}

func (ep EntryPath) Clean() EntryPath {
	return EntryPath(filepath.Clean(string(ep)))
}

func (ep EntryPath) Base() PathSegment {
	return PathSegment(filepath.Base(string(ep)))
}

func (ep EntryPath) Stat(fileSys ...fs.FS) (fs.FileInfo, error) {
	if len(fileSys) == 0 {
		return os.Stat(string(ep))
	}
	return fs.Stat(fileSys[0], string(ep))
}

func (ep EntryPath) Lstat(fileSys ...fs.FS) (os.FileInfo, error) {
	if len(fileSys) == 0 {
		return os.Lstat(string(ep))
	}
	return fs.Lstat(fileSys[0], string(ep))
}

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

// Readlink returns the target of the symlink referred to by fp.
//
// It returns the resolved target as a EntryPath. If fp is not a symlink,
// it returns an empty EntryPath and a non-nil error from os.Readlink.
// On most systems the returned target is relative to the directory
// containing fp, not an absolute path.
func (ep EntryPath) Readlink() (target EntryPath, err error) {
	var linkTarget string
	linkTarget, err = os.Readlink(string(ep))
	if err != nil {
		goto end
	}
	target = EntryPath(linkTarget)
end:
	return target, err
}

func (ep EntryPath) HasSuffix(suffix DirPath) bool {
	return strings.HasSuffix(string(ep), string(suffix))
}

// Contains checks if ep contains the given substring.
// Accepts: string, DirPath, Filepath, EntryPath, PathSegment, or fmt.Stringer
// Panics on unsupported types.
func (ep EntryPath) Contains(substr any) bool {
	var s string

	switch v := substr.(type) {
	case string:
		s = v
	case DirPath:
		s = string(v)
	case Filepath:
		s = string(v)
	case EntryPath:
		s = string(v)
	case PathSegment:
		s = string(v)
	case interface{ String() string }:
		s = v.String()
	default:
		panic("EntryPath.Contains: unsupported type")
	}

	return strings.Contains(string(ep), s)
}

func (ep EntryPath) VolumeName() VolumeName {
	return VolumeName(filepath.VolumeName(string(ep)))
}

func (ep EntryPath) Abs() (EntryPath, error) {
	entry, err := filepath.Abs(string(ep))
	return EntryPath(entry), err
}

func (ep EntryPath) IsFile() bool {
	return filepath.IsAbs(string(ep))
}

func (ep EntryPath) IsAbs() bool {
	return filepath.IsAbs(string(ep))
}

func (ep EntryPath) EvalSymlinks() (_ EntryPath, err error) {
	var s string
	s, err = filepath.EvalSymlinks(string(ep))
	return EntryPath(s), err
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

//====Extensions

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
		// Break out to be handled by the tilde expansion
		break

	case s[:2] == "~\\" && runtime.GOOS == "windows":
		// We start with ~\ so we are a tilde path on Windows
		// Break out to be handled by the tilde expansion
		break

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
	return out, err
}
