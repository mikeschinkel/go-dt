package dt

import (
	"strings"

	"github.com/mikeschinkel/go-dt/de"
)

// PathSegments is one or more path segments for a filepath, dir path, or URL
type PathSegments string

func (pss PathSegments) Segments() (out []PathSegment) {
	out = make([]PathSegment, strings.Count(string(pss), "/"))
	for i, ps := range strings.Split(string(pss), "/") {
		out[i] = PathSegment(ps)
	}
	return out
}

func (pss PathSegments) CanWrite() (bool, error) {
	return CanWrite(EntryPath(pss))
}

type PathSegment string

func (ps PathSegment) Contains(part any) bool {
	return EntryPath(ps).Contains(part)
}

func ParsePathSegment(s string) (ps PathSegment, err error) {
	if s == "" {
		err = NewErr(
			ErrInvalidPathSegment,
			de.ErrEmpty,
		)
		goto end
	}
	ps = PathSegment(s)
end:
	return ps, err
}
