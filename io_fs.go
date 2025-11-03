package dt

import (
	"io/fs"
)

func FSReadFile(fileSys fs.FS, fp Filepath) ([]byte, error) {
	return fs.ReadFile(fileSys, string(fp))
}

func FSStat[T ~string](fileSys fs.FS, entry T) (fs.FileInfo, error) {
	return fs.Stat(fileSys, string(entry))
}
func FSStatFile(fileSys fs.FS, fp Filepath) (fs.FileInfo, error) {
	return FSStat(fileSys, fp)
}
func FSStatDir(fileSys fs.FS, dp DirPath) (fs.FileInfo, error) {
	return FSStat(fileSys, dp)
}
