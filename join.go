package dt

import (
	"path/filepath"
)

func DirPathJoin[T1, T2 ~string](a T1, b T2) DirPath {
	return DirPath(filepath.Join(string(a), string(b)))
}

func DirPathJoin3[T1, T2, T3 ~string](a T1, b T2, c T3) DirPath {
	return DirPath(filepath.Join(string(a), string(b), string(c)))
}
func FilepathJoin[T1, T2 ~string](a T1, b T2) Filepath {
	return Filepath(filepath.Join(string(a), string(b)))
}
