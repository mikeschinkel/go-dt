package dt

import (
	"strings"
)

// SplitSegments splits a string by sep separator and returns a slice of the specified segment type.
// It counts separators first to allocate exact capacity, optimizing for GC and heavy usage.
func SplitSegments[S ~string](s, sep string) []S {
	var before string
	var found bool
	var rest string
	var count int

	if len(s) == 0 {
		return make([]S, 0)
	}

	// Count separators to determine exact segment count
	count = 1
	for i := 0; i < len(s); i++ {
		if string(s[i]) == sep {
			count++
		}
	}

	out := make([]S, 0, count)
	rest = s
	for {
		before, rest, found = strings.Cut(rest, sep)
		out = append(out, S(before))
		if !found {
			goto end
		}
	}
end:
	return out
}

// IndexSegments returns the segment at the given index in a string split by sep.
// Returns empty string (cast to S) if index is out of bounds or negative.
func IndexSegments[S ~string](s, sep string, index int) S {
	var before string
	var found bool
	var rest string
	var curIdx int

	if len(s) == 0 {
		goto end
	}
	if index < 0 {
		goto end
	}

	rest = s
	curIdx = 0

	for {
		before, rest, found = strings.Cut(rest, sep)
		if curIdx == index {
			return S(before)
		}
		if !found {
			goto end
		}
		curIdx++
	}
end:
	return S("")
}

// SliceSegments returns a slice of segments from start (inclusive) to end (exclusive)
// in a string split by sep, similar to string[start:end] slicing.
// Supports -1 for end to mean "to the last segment".
// Returns empty slice if indices are invalid (negative start, start >= end, or out of bounds).
func SliceSegments[S ~string](s, sep string, start, end int) []S {
	var before string
	var found bool
	var rest string
	var curIdx int
	var capacity int

	out := make([]S, 0)

	if len(s) == 0 {
		goto end
	}
	if start < 0 {
		goto end
	}
	if end < 0 && end != -1 {
		goto end
	}
	if end != -1 && start >= end {
		goto end
	}

	// Count separators to estimate segment count for capacity
	capacity = 1
	for i := 0; i < len(s); i++ {
		if string(s[i]) == sep {
			capacity++
		}
	}

	// Set capacity - use end-start if end is specified, otherwise use remaining capacity
	if end != -1 {
		out = make([]S, 0, end-start)
	} else {
		out = make([]S, 0, capacity-start)
	}

	rest = s
	curIdx = 0

	for {
		before, rest, found = strings.Cut(rest, sep)
		if curIdx >= start && (end == -1 || curIdx < end) {
			out = append(out, S(before))
		}
		if end != -1 && curIdx >= end {
			goto end
		}
		if !found {
			goto end
		}
		curIdx++
	}
end:
	return out
}

// SliceSegmentsScalar returns a scalar (joined) value of segments from start (inclusive) to end (exclusive)
// in a string split by sep, similar to string[start:end] slicing but returns a joined string.
// Supports -1 for end to mean "to the last segment".
// Returns empty string if indices are invalid.
// Uses a single pass through the string to find byte positions, avoiding allocations.
func SliceSegmentsScalar[S ~string](s, sep string, start, end int) S {
	var curIdx int
	var startPos int
	var endPos int

	if len(s) == 0 {
		return S("")
	}
	if start < 0 {
		return S("")
	}
	if end < 0 && end != -1 {
		return S("")
	}
	if end != -1 && start >= end {
		return S("")
	}

	// Single pass through string to find start and end byte positions
	rest := s
	pos := 0
	startPos = -1
	endPos = len(s)

	for {
		before, after, found := strings.Cut(rest, sep)

		// Record position of this segment
		if curIdx == start {
			startPos = pos
		}
		if curIdx == end && end != -1 {
			endPos = pos
			break
		}

		if !found {
			// We've reached the last segment
			if end == -1 || curIdx < end {
				endPos = len(s)
			}
			break
		}

		// Move to next segment
		pos += len(before) + len(sep)
		rest = after
		curIdx++
	}

	// Return the substring from startPos to endPos
	if startPos == -1 {
		return S("")
	}
	if startPos >= len(s) {
		return S("")
	}
	if endPos > len(s) {
		endPos = len(s)
	}
	if startPos >= endPos {
		return S("")
	}

	// Handle the case where we need to remove the trailing separator
	result := s[startPos:endPos]
	if len(result) > 0 && strings.HasSuffix(result, sep) {
		result = result[:len(result)-len(sep)]
	}

	return S(result)
}

// JoinSegments returns joined slice of segments
func JoinSegments[S ~string](ss []S, sep string) (s S) {
	var size int
	var sb strings.Builder
	var tmp string

	if len(ss) == 0 {
		goto end
	}
	for _, s := range ss {
		size += len(s) + 1
	}
	size--
	sb.Grow(size)
	for _, s := range ss {
		sb.Write([]byte(s))
		sb.WriteString(sep)
	}
	tmp = sb.String()
	s = S(tmp[:len(tmp)-1])
end:
	return s
}
