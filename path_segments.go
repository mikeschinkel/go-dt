package dt

import (
	"path/filepath"
	"strings"
)

func ParsePathSegment(s string) (ps PathSegment, err error) {
	if s == "" {
		err = NewErr(
			ErrInvalidPathSegment,
			ErrEmpty,
		)
		goto end
	}
	ps = PathSegment(s)
end:
	return ps, err
}

// PathSegments is one or more path segments for a filepath, dir path, or URL
type PathSegments string

func (pss PathSegments) Base() PathSegment {
	return PathSegment(filepath.Base(string(pss)))
}

func (pss PathSegments) Split() (out []PathSegment) {
	out = make([]PathSegment, strings.Count(string(pss), "/"))
	for i, ps := range strings.Split(string(pss), "/") {
		out[i] = PathSegment(ps)
	}
	return out
}

func (pss PathSegments) CanWrite() (bool, error) {
	return CanWrite(EntryPath(pss))
}

func (pss PathSegments) HasDotDotPrefix() bool {
	return EntryPath(pss).HasDotDotPrefix()
}

func (pss PathSegments) Segments() []PathSegment {
	return pss.Split()
}
