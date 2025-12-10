package dt

import (
	"os"
	"strings"
)

type PathSegment string

func (ps PathSegment) Contains(part any) bool {
	return EntryPath(ps).Contains(part)
}

func (ps PathSegment) HasDotDotPrefix() bool {
	return EntryPath(ps).HasDotDotPrefix()
}

func (ps PathSegment) Exists() (exists bool, err error) {
	var status EntryStatus
	status, err = ps.Status()
	if err != nil {
		goto end
	}
	exists = status == IsFileEntry
end:
	return exists, err
}

func (ps PathSegment) WriteFile(data []byte, mode os.FileMode) error {
	return os.WriteFile(string(ps), data, mode)
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
func (ps PathSegment) Status(flags ...EntryStatusFlags) (status EntryStatus, err error) {
	return EntryPath(ps).Status(flags...)
}

func (ps PathSegment) MkdirAll(mode os.FileMode) error {
	return os.MkdirAll(string(ps), mode)
}

func (ps PathSegment) TrimPrefix(prefix DirPath) PathSegment {
	return PathSegment(strings.TrimPrefix(string(ps), string(prefix)))
}

func (ps PathSegment) TrimSuffix(TrimSuffix string) PathSegment {
	return PathSegment(strings.TrimSuffix(string(ps), TrimSuffix))
}
