package dt

import (
	"path/filepath"
)

func Dir(fp Filepath) DirPath {
	return DirPath(filepath.Dir(string(fp)))
}
func FileBase(fp Filepath) Filename {
	return Filename(filepath.Base(string(fp)))
}
func DirBase(dp DirPath) PathSegments {
	return PathSegments(filepath.Base(string(dp)))
}
