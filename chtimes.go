package dt

import (
	"os"
	"time"
)

func Chtimes[T ~string](path T, atime time.Time, mtime time.Time) error {
	return os.Chtimes(string(path), atime, mtime)
}
func ChangeFileTimes(fp Filepath, atime time.Time, mtime time.Time) error {
	return Chtimes(fp, atime, mtime)
}
func ChangeDirTimes(dp DirPath, atime time.Time, mtime time.Time) error {
	return Chtimes(dp, atime, mtime)
}
