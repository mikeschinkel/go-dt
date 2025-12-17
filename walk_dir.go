package dt

import (
	"io/fs"
	"iter"
	"os"
	"path/filepath"
)

// walkDirEntry adapts an fs.FileInfo to fs.DirEntry.
// Used for the synthetic root entry in WalkDir.
type walkDirEntry struct {
	info fs.FileInfo
}

func (e walkDirEntry) Name() string               { return e.info.Name() }
func (e walkDirEntry) IsDir() bool                { return e.info.IsDir() }
func (e walkDirEntry) Type() fs.FileMode          { return e.info.Mode().Type() }
func (e walkDirEntry) Info() (fs.FileInfo, error) { return e.info, nil }

func WalkDir(root DirPath) iter.Seq2[DirEntry, error] {
	return func(yield func(DirEntry, error) bool) {
		type dirState struct {
			dir     string // absolute or as given
			entries []os.DirEntry
			i       int
		}

		var skipDir bool
		var stack []dirState

		// 1. Stat the root, like filepath.WalkDir does.
		rootPath := string(root)

		info, err := os.Lstat(rootPath)

		rootEntry := DirEntry{
			Root:    root,
			Rel:     ".", // relative *within* this walk
			Entry:   nil, // we’ll wrap info below if ok
			skipDir: &skipDir,
		}

		if err == nil {
			rootEntry.Entry = walkDirEntry{info: info}
		}

		skipDir = false
		if !yield(rootEntry, err) {
			goto end
		}

		if err != nil || skipDir {
			goto end
		}

		// 2. Non-recursive walk of children with os.ReadDir.
		stack = []dirState{{dir: rootPath}}

		for len(stack) > 0 {
			s := &stack[len(stack)-1]

			if s.entries == nil {
				ents, readErr := os.ReadDir(s.dir)
				if readErr != nil {
					rel := relPathWithinRoot(rootPath, s.dir) // "." or "sub/dir"

					entry := DirEntry{
						Root:    root,
						Rel:     RelPath(rel),
						Entry:   nil,
						skipDir: &skipDir,
					}

					skipDir = false
					if !yield(entry, readErr) {
						goto end
					}

					stack = stack[:len(stack)-1]
					continue
				}

				s.entries = ents
				s.i = 0
			}

			if s.i >= len(s.entries) {
				stack = stack[:len(stack)-1]
				continue
			}

			de := s.entries[s.i]
			s.i++

			childPath := filepath.Join(s.dir, de.Name())
			rel := relPathWithinRoot(rootPath, childPath)

			entry := DirEntry{
				Root:    root,
				Rel:     RelPath(rel),
				Entry:   de,
				skipDir: &skipDir,
			}

			skipDir = false

			if !yield(entry, nil) {
				goto end
			}

			if !de.IsDir() {
				continue
			}

			if skipDir {
				continue
			}

			stack = append(stack, dirState{dir: childPath})
		}

	end:
		return
	}
}

// relPathWithinRoot returns a path relative to rootPath.
//
//   - If path == rootPath, it returns ".".
//   - Otherwise it returns a relative path such that
//     filepath.Join(rootPath, rel) == path, when possible.
//   - If filepath.Rel fails (shouldn't in normal use), it falls back
//     to the cleaned child path.
func relPathWithinRoot(rootPath, path string) string {
	rootClean := filepath.Clean(rootPath)
	pathClean := filepath.Clean(path)

	if pathClean == rootClean {
		return "."
	}

	rel, err := filepath.Rel(rootClean, pathClean)
	if err != nil || rel == "" {
		// Fallback: in practice this shouldn't happen if path is under root.
		return pathClean
	}

	return rel
}
