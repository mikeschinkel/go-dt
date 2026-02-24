package dtx

import (
	"maps"
	"slices"
)

func UniqueSlice[T comparable](slice []T) (out []T) {
	var seen map[T]struct{}

	if len(slice) < 2 {
		out = slice
		goto end
	}

	seen = make(map[T]struct{}, len(slice))
	out = make([]T, 0, len(slice))

	for _, ele := range slice {
		_, ok := seen[ele]
		if ok {
			continue
		}

		seen[ele] = struct{}{}
		out = append(out, ele)
	}
end:
	return out
}

// ConvertStringSlice takes a slice whose element values are derived from string
// (or are of string) and then returns a slice whose element values are of
// another type also derived from string, or are of string itself.
func ConvertStringSlice[OUT, IN ~string](slice []IN) (out []OUT) {
	out = make([]OUT, len(slice))
	for i, e := range slice {
		out[i] = OUT(e)
	}
	return out
}

// StringKeys takes a map with keys whose type is derived from string and returns
// a slice of string for those keys. This is useful for passing to functions like
// strings.Join(), etc.
func StringKeys[K ~string, V any](m map[K]V) (out []string) {
	return ConvertStringSlice[string](slices.Collect(maps.Keys(m)))
}
