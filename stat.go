package dt

import (
	"os"
)

func Stat[T ~string](entry T) (os.FileInfo, error) {
	return os.Stat(string(entry))
}
func StatFile(fp Filepath) (os.FileInfo, error) {
	return Stat(fp)
}
func StatDir(dp DirPath) (os.FileInfo, error) {
	return Stat(dp)
}
