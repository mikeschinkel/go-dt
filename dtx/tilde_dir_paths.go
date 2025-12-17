package dtx

import (
	"github.com/mikeschinkel/go-dt"
)

type TildeDirPaths []dt.TildeDirPath

// Unique returns a new TildeDirPaths slice containing only the first occurrence
// of each TildeDirPath in p, preserving the original order.
//
// Uniqueness is determined using TildeDirPath equality; TildeDirPath values are assumed
// to already satisfy their normalization invariants.
func (dps TildeDirPaths) Unique() TildeDirPaths {
	return UniqueSlice(dps)
}
