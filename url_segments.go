package dt

import (
	"path/filepath"
	"strings"
)

func ParseURLSegments(s string) (uss URLSegments, err error) {
	if s == "" {
		err = NewErr(
			ErrInvalidURLSegments,
			ErrEmpty,
		)
		goto end
	}
	uss = URLSegments(s)
end:
	return uss, err
}

// URLSegments is one or more path segments for a filepath, dir path, or URL
type URLSegments string

func (uss URLSegments) Base() URLSegment {
	return URLSegment(filepath.Base(string(uss)))
}

//func (pss URLSegments) Split() (out []URLSegment) {
//	out = make([]URLSegment, strings.Count(string(pss), "/"))
//	for i, ps := range strings.Split(string(pss), "/") {
//		out[i] = URLSegment(ps)
//	}
//	return out
//}

func (uss URLSegments) Segments() []URLSegment {
	return uss.Split()
}

func (uss URLSegments) LastIndex(sep string) int {
	return strings.LastIndex(string(uss), sep)
}

func (uss URLSegments) Split() []URLSegment {
	return SplitSegments[URLSegment](string(uss), "/")
}

func (uss URLSegments) Segment(index int) URLSegment {
	return IndexSegments[URLSegment](string(uss), "/", index)
}

func (uss URLSegments) Slice(start, end int) []URLSegment {
	return SliceSegments[URLSegment](string(uss), "/", start, end)
}

func (uss URLSegments) SliceScalar(start, end int, sep string) URLSegment {
	return SliceSegmentsScalar[URLSegment](string(uss), sep, start, end)
}
