package dtx

import (
	"github.com/mikeschinkel/go-dt"
)

type EntryPaths []dt.EntryPath

// Unique returns a new EntryPaths slice containing only the first occurrence
// of each EntryPath in p, preserving the original order.
//
// Uniqueness is determined using EntryPath equality; EntryPath values are assumed
// to already satisfy their normalization invariants.
func (eps EntryPaths) Unique() EntryPaths {
	return UniqueSlice(eps)
}

func (eps EntryPaths) DirPaths() (dps []dt.DirPath) {
	dps = make([]dt.DirPath, len(eps))
	for i, ep := range eps {
		dps[i] = dt.DirPath(ep)
	}
	return dps
}

func (eps EntryPaths) Filepaths() (fps []dt.Filepath) {
	fps = make([]dt.Filepath, len(eps))
	for i, ep := range eps {
		fps[i] = dt.Filepath(ep)
	}
	return fps
}
