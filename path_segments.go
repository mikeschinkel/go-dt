package dt

import (
	"path/filepath"
	"strings"
)

// PathSegments is one or more path segments for a filesystem path
type PathSegments string

// Split returns all segments of the path separated by the OS separator as a slice of PathSegment.
func (pss PathSegments) Split() []PathSegment {
	return SplitSegments[PathSegment](string(pss), string(filepath.Separator))
}

// Segments returns all segments of the path, same as Split.
func (pss PathSegments) Segments() []PathSegment {
	return pss.Split()
}

// Segment returns the segment at the given index in the path.
// Returns empty PathSegment if index is out of bounds or negative.
func (pss PathSegments) Segment(index int) PathSegment {
	return IndexSegments[PathSegment](string(pss), string(filepath.Separator), index)
}

// Slice returns segments from start (inclusive) to end (exclusive),
// similar to string[start:end] slicing.
func (pss PathSegments) Slice(start, end int) []PathSegment {
	return SliceSegments[PathSegment](string(pss), string(filepath.Separator), start, end)
}

// SliceScalar returns a scalar (joined) value of segments from start (inclusive) to end (exclusive),
// similar to string[start:end] slicing but returns a joined string using the OS separator.
// Supports -1 for end to mean "to the last segment".
func (pss PathSegments) SliceScalar(start, end int) PathSegments {
	return PathSegments(SliceSegmentsScalar[PathSegment](string(pss), string(filepath.Separator), start, end))
}

func (pss PathSegments) LastIndex(sep string) int {
	return strings.LastIndex(string(pss), sep)
}
