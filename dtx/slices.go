package dtx

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
	return slice
}
