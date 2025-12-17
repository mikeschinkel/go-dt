package dtx

import (
	"strings"

	"github.com/mikeschinkel/go-dt"
)

type DirPaths []dt.DirPath

// Unique returns a new DirPaths slice containing only the first occurrence
// of each DirPath in p, preserving the original order.
//
// Uniqueness is determined using DirPath equality; DirPath values are assumed
// to already satisfy their normalization invariants.
func (dps DirPaths) Unique() DirPaths {
	return UniqueSlice(dps)
}

func (dps DirPaths) Join(sep string) string {
	sb := strings.Builder{}
	last := len(dps) - 1
	for i, d := range dps {
		sb.WriteString(string(d))
		if i == last {
			break
		}
		sb.WriteString(sep)
	}
	return sb.String()
}
