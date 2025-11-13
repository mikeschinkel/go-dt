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

func DirPathJoin4[T1, T2, T3, T4 ~string](a T1, b T2, c T3, d T4) DirPath {
	return DirPath(filepath.Join(string(a), string(b), string(c), string(d)))
}

func DirPathJoin5[T1, T2, T3, T4, T5 ~string](a T1, b T2, c T3, d T4, e T5) DirPath {
	return DirPath(filepath.Join(string(a), string(b), string(c), string(d), string(e)))
}

func FilepathJoin[T1, T2 ~string](a T1, b T2) Filepath {
	return Filepath(filepath.Join(string(a), string(b)))
}
func FilepathJoin3[T1, T2, T3 ~string](a T1, b T2, c T3) Filepath {
	return Filepath(filepath.Join(string(a), string(b), string(c)))
}

func RelFilepathJoin[T1, T2 ~string](a T1, b T2) RelFilepath {
	return RelFilepath(filepath.Join(string(a), string(b)))
}

func RelFilepathJoin3[T1, T2, T3 ~string](a T1, b T2, c T3) RelFilepath {
	return RelFilepath(filepath.Join(string(a), string(b), string(c)))
}

func RelFilepathJoin4[T1, T2, T3, T4 ~string](a T1, b T2, c T3, d T4) RelFilepath {
	return RelFilepath(filepath.Join(string(a), string(b), string(c), string(d)))
}

func RelFilepathJoin5[T1, T2, T3, T4, T5 ~string](a T1, b T2, c T3, d T4, e T5) RelFilepath {
	return RelFilepath(filepath.Join(string(a), string(b), string(c), string(d), string(e)))
}
